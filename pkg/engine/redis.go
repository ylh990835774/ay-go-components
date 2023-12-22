package engine

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/anmitsu/go-shlex"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"github.com/ylh990835774/ay-go-components/pkg/common"
	"github.com/ylh990835774/ay-go-components/pkg/inerr"
)

type RedisEngine struct {
	driver redis.UniversalClient
	*common.EngineBase
}

func (r *RedisEngine) RegistryQueryPrev(hook common.PreHook) {
	r.BindPrevHook(hook)
}

func (r *RedisEngine) RegistryQueryPost(hook common.PostHook) {
	r.BindPostHook(hook)
}

func (r *RedisEngine) Close() error {
	return r.driver.Close()
}

func (r *RedisEngine) Schema() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	res, err := r.driver.ConfigGet(ctx, "databases").Result()
	if err != nil {
		return nil, err
	}

	if len(res) != 2 {
		return nil, inerr.ErrRedisSchemaFetchFailed
	}

	dbCount, err := strconv.Atoi(res[1].(string))
	if err != nil {
		return nil, err
	}

	schemaList := make([]string, 0)
	for dbIndex := 0; dbIndex < dbCount; dbIndex++ {
		schemaList = append(schemaList, fmt.Sprintf("db%d", dbIndex))
	}

	return schemaList, nil
}

func (r *RedisEngine) Table(schema string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	var cursor uint64
	keyList := make([]string, 0)

	for {
		var keys []string
		var err error

		keys, cursor, err = r.driver.Scan(ctx, cursor, "*", 100).Result()
		if err != nil {
			return nil, err
		}

		keyList = append(keyList, keys...)
		if cursor == 0 {
			break
		}
	}

	return keyList, nil
}

func (r *RedisEngine) Query(schema string, table string, sql string, timeout int64) *common.QuerySet {
	queryRes := &common.QuerySet{
		EngineType: common.RedisEngine,
		Action:     common.ActionSQLQuery,
		IsExecute:  false,
		Err:        nil,
	}

	// execute query prev hook
	// query prev hook failed and stop query
	if r.QueryPrev != nil {
		err := r.QueryPrev(&common.PrevHookArgs{
			EngineType: common.RedisEngine,
			Action:     common.ActionSQLQuery,
			Schema:     schema,
			SQL:        sql,
		})
		if err != nil {
			queryRes.Err = err
			return queryRes
		}
	}

	// registry query post hook
	defer func() {
		if r.QueryPost != nil {
			r.QueryPost(&common.PostHookArgs{
				EngineType:    common.RedisEngine,
				Action:        common.ActionSQLQuery,
				IsExecute:     queryRes.IsExecute,
				ExecuteAt:     queryRes.ExecuteAt,
				QueryDuration: queryRes.QueryDuration,
				Err:           queryRes.Err,
				Schema:        schema,
				SQL:           sql,
				AffectedRows:  queryRes.AffectedRows,
			})
		}
	}()

	if schema == "" {
		queryRes.Err = inerr.ErrSchemaEmpty
		return queryRes
	}

	if sql == "" {
		queryRes.Err = inerr.ErrRedisCMDEmpty
		return queryRes
	}

	// query main
	redisCMD := make([]interface{}, 0)
	redisCMDSlice, err := shlex.Split(sql, true)
	if err != nil {
		queryRes.Err = errors.Wrap(err, "parse redis command failed")
		return queryRes
	}

	for _, cmdToken := range redisCMDSlice {
		redisCMD = append(redisCMD, cmdToken)
	}

	var redisKey string
	if len(redisCMDSlice) >= 2 {
		redisKey = redisCMDSlice[1]
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	// try acquire key type
	var keyType string
	if redisKey != "" {
		keyType, _ = r.driver.Type(ctx, redisKey).Result()
	}

	// run redis command by user provided
	queryRes.ExecuteAt = time.Now()
	res, err := r.driver.Do(ctx, redisCMD...).Result()
	queryRes.QueryDuration = time.Since(queryRes.ExecuteAt).Milliseconds()

	if err == redis.Nil {
		queryRes.Err = inerr.ErrRedisKeyNotExist
		return queryRes
	}

	if err != nil {
		queryRes.Err = err
		return queryRes
	}

	queryRes.SQL = sql
	queryRes.IsExecute = true
	queryRes.Err = nil
	queryRes.Total = 1
	queryRes.Columns = []string{"redis_command", "redis_key_type", "command_result"}
	queryRes.Rows = []common.Row{
		{
			"redis_command":  redisCMDSlice[0],
			"redis_key_type": keyType,
			"command_result": res,
		},
	}
	queryRes.AffectedRows = 1

	return queryRes
}

func (r *RedisEngine) InitialDriver(conn common.ConnConfig, schema string) error {
	var dbIndex int
	var err error

	if schema == "" {
		dbIndex = 0
	} else {
		dbIndex, err = strconv.Atoi(strings.TrimPrefix(schema, "db"))
		if err != nil {
			return err
		}
	}

	r.driver = redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs: []string{
			fmt.Sprintf("%s:%d", conn.IP, conn.Port),
		},
		DB: dbIndex, // default 0 database

		Username: conn.UserName,
		Password: conn.Password,

		DialTimeout: 15 * time.Second,
	})

	r.ConnConfig = conn

	return r.Ping()
}

func (r *RedisEngine) Reset() {
	err := r.Close()
	if err != nil {
		fmt.Printf("redis client close failed: %s\n", err)
	}
	r.driver = nil
	r.EngineBase = &common.EngineBase{}
}

func (r *RedisEngine) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := r.driver.Ping(ctx).Result()
	return err
}

func ForkRedisEngine(conn common.ConnConfig, schema int) *RedisEngine {
	cli := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs: []string{
			fmt.Sprintf("%s:%d", conn.IP, conn.Port),
		},
		DB: schema, // default 0 database

		Username: conn.UserName,
		Password: conn.Password,

		DialTimeout: 15 * time.Second,
	})

	return &RedisEngine{
		cli,
		common.NewEngineBase(conn),
	}
}

func NewRedisEngine() *RedisEngine {
	return &RedisEngine{
		nil,
		&common.EngineBase{},
	}
}

func IsRedisCMDSafe(sql string, whiteList []common.SQLType) (string, bool, error) {
	if sql == "" {
		return sql, false, inerr.ErrRedisCMDEmpty
	}

	redisSQL, err := shlex.Split(strings.TrimSpace(sql), true)
	if err != nil {
		return sql, false, err
	}

	if len(redisSQL) == 0 {
		return sql, false, inerr.ErrRedisCMDUnknown
	}

	cmdType := strings.ToLower(redisSQL[0])
	cmdTypeFlag, has := common.RedisCMDTOSQLType[cmdType]
	if !has {
		return sql, false, inerr.ErrRedisCMDUnSupported
	}

	for _, c := range whiteList {
		if cmdTypeFlag == c {
			return sql, true, nil
		}
	}

	return sql, false, nil
}

func ParseRedisKeyFromRedisCMD(sql string) (string, error) {
	if sql == "" {
		return "", inerr.ErrRedisCMDEmpty
	}

	redisSQL, err := shlex.Split(strings.TrimSpace(sql), true)
	if err != nil {
		return "", err
	}

	if len(redisSQL) < 2 {
		return "", inerr.ErrRedisParseKey
	}

	return redisSQL[1], nil
}

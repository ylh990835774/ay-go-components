package console

import (
	"fmt"
	"sync"

	"git.qpaas.com/go-components/webconsole/pkg/common"
	"git.qpaas.com/go-components/webconsole/pkg/engine"
	"git.qpaas.com/go-components/webconsole/pkg/inerr"
	"github.com/pkg/errors"
)

type redisConsole struct {
	sync.Pool
	*common.ConsoleBase
}

func (r *redisConsole) ConsoleType() string {
	return common.RedisConsole
}

func (r *redisConsole) Fork(conn common.ConnConfig, schema string) (engine.Engine, error) {
	eg := r.Get().(*engine.RedisEngine)

	err := eg.InitialDriver(conn, schema)
	if err != nil {
		return nil, err
	}

	return eg, nil
}

func (r *redisConsole) Destory(e engine.Engine) {
	eg := e.(*engine.RedisEngine)
	eg.Reset()

	r.Put(eg)
}

func (r *redisConsole) SchemaHandler(opt *common.HandlerOptions) ([]string, error) {
	// fork engine instance
	eg, err := r.Fork(opt.Conn, "")
	if err != nil {
		return nil, err
	}
	defer r.Destory(eg) // destory engine instance

	// bind hooks
	eg.RegistryQueryPrev(opt.QueryBeforeHook)
	eg.RegistryQueryPost(opt.QueryAfterHook)

	// fetch schema
	schemas, err := eg.Schema()
	if err != nil {
		return nil, err
	}

	return schemas, nil
}

func (r *redisConsole) TableHandler(schema string, opt *common.HandlerOptions) ([]string, error) {
	// fork engine instance
	eg, err := r.Fork(opt.Conn, schema)
	if err != nil {
		return nil, err
	}
	defer r.Destory(eg) // destory engine instance

	// bind hooks
	eg.RegistryQueryPrev(opt.QueryBeforeHook)
	eg.RegistryQueryPost(opt.QueryAfterHook)

	// fetch tables
	tables, err := eg.Table(schema)
	if err != nil {
		return nil, err
	}

	return tables, nil
}

func (r *redisConsole) QueryHandler(schema string, table string, sql string, opt *common.HandlerOptions) *common.QuerySet {
	// fork engine instance
	eg, err := r.Fork(opt.Conn, schema)
	if err != nil {
		return &common.QuerySet{
			Err: errors.Wrap(err, "redis engine fork failed"),
		}
	}
	defer r.Destory(eg) // destory engine instance

	// bind hooks
	eg.RegistryQueryPrev(opt.QueryBeforeHook)
	eg.RegistryQueryPost(opt.QueryAfterHook)

	var isSafe bool
	defaultSafeCMD := common.DefaultRedisWhiteCMD

	// if sql is empty
	// assign type key as default command
	if sql == "" {
		sql = fmt.Sprintf("type %s", table)
	}

	// valid is turn off system intercpet
	if opt.IsIgnoreSystemIntercept {
		goto queryMain
	}

	// system intercept
	// check redis command is valid by user provided white list
	// if user not set, valid by default white list
	// otherwise valid by user provided

	if opt.AllowSQLType != nil {
		defaultSafeCMD = opt.AllowSQLType
	}

	_, isSafe, err = engine.IsRedisCMDSafe(sql, defaultSafeCMD)
	if err != nil {
		return &common.QuerySet{
			Err: errors.Wrap(err, "redis command preCheck failed"),
		}
	}

	if !isSafe {
		return &common.QuerySet{
			Err: errors.Wrap(inerr.ErrRedisCMDForbidden, sql),
		}
	}

queryMain:
	// query execute
	var defaultQueryTimeout int64 = 15
	if opt.QueryOpt.Timeout > 0 {
		defaultQueryTimeout = opt.QueryOpt.Timeout
	}

	// try to parse key from redis command
	// if success chosen key from redis command or chosen key from user provided
	keyFromCMD, err := engine.ParseRedisKeyFromRedisCMD(sql)
	if err == nil {
		table = keyFromCMD
	}

	return eg.Query(schema, table, sql, defaultQueryTimeout)
}

// NewRedisConsole
func NewRedisConsole() *redisConsole {
	return &redisConsole{
		sync.Pool{
			New: func() interface{} {
				return engine.NewRedisEngine()
			},
		},
		common.NewConsoleBase(),
	}
}

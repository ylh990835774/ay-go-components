package test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/ylh990835774/ay-go-components/pkg/common"
	"github.com/ylh990835774/ay-go-components/pkg/console"
)

func TestRedisConsoleWithDefaultOpt(t *testing.T) {
	redisConsole := console.NewRedisConsole()

	opt := &common.HandlerOptions{
		Conn: common.ConnConfig{
			IP:   "172.168.2.24",
			Port: 16379,
		},
		QueryOpt: common.QueryOptions{
			Timeout: 600,
		},
	}

	t.Run("fetch schema", func(t *testing.T) {
		fakeQueryMeta := &common.QueryMeta{
			Action: common.ActionFetchSchema,
		}
		reqBody, _ := json.Marshal(fakeQueryMeta)

		mockHTTPReq(t, redisConsole, opt, reqBody, "/console/redis")
	})

	t.Run("fetch table", func(t *testing.T) {
		fakeQueryMeta := &common.QueryMeta{
			Action: common.ActionFetchTable,
			Schema: "db1",
		}
		reqBody, _ := json.Marshal(fakeQueryMeta)

		mockHTTPReq(t, redisConsole, opt, reqBody, "/console/redis")
	})

	t.Run("redis query safe command", func(t *testing.T) {
		SQLList := []string{
			`ttl game`,
			`exists game`,
			`type game`,
			`scan 0`,
			`get game`,
			`mget game game01`,
			`strlen game`,
			`llen money`,
			`lrange money 0 -1`,
			`lindex money 0`,
			`hgetall student`,
			`hexists student name`,
			`hget student name`,
			`hkeys student`,
			`hvals student`,
			`scard ppp`,
			`smembers ppp`,
			`sismember ppp 111`,
			`sismember ppp 999`,
			`sdiff ppp ttt`,
			`sunion ppp ttt`,
			`zcard setganme`,
			`zrange setganme 0 -1`,
			`zrank setgname li`,
			`zcount setgname 0 1`,
			`zscore setgname li`,
			`zrangebyscore setganme +inf -inf`,
		}

		for _, sql := range SQLList {
			fakeQueryMeta := &common.QueryMeta{
				Action: common.ActionSQLQuery,
				Schema: "db1",
				SQL:    SQLBase64(sql),
			}
			reqBody, _ := json.Marshal(fakeQueryMeta)

			mockHTTPReq(t, redisConsole, opt, reqBody, "/console/redis")
		}
	})

	t.Run("redis query forbidden command", func(t *testing.T) {
		SQLList := []string{
			`del name`,
			`expire name`,
			`expireat name 12`,
			`append name 12`,
			`incr name 111`,
			`incrby name 1`,
			`set name 111`,
			`mset name name2 01`,
			`setnx name 111`,
			`setex name 222`,
			`hdel student`,
			`hset student name kk`,
			`hmset student name kk`,
			`lpop money`,
			`rpop money`,
			`lpush money 1`,
			`rpush money 2`,
			`linsert money 1`,
			`sadd ppp 111`,
			`srem ppp 222`,
			`zadd ppp 111`,
			`zrem ppp 222`,
		}

		for _, sql := range SQLList {
			fakeQueryMeta := &common.QueryMeta{
				Action: common.ActionSQLQuery,
				Schema: "db1",
				SQL:    SQLBase64(sql),
			}
			reqBody, _ := json.Marshal(fakeQueryMeta)

			mockHTTPReq(t, redisConsole, opt, reqBody, "/console/redis")
		}
	})
}

func TestRedisConsoleWithUserDefinedOptions(t *testing.T) {
	redisConsole := console.NewRedisConsole()

	opt := &common.HandlerOptions{
		Conn: common.ConnConfig{
			IP:   "172.168.2.24",
			Port: 16379,
		},
		QueryOpt: common.QueryOptions{
			Timeout: 10,
		},
		AllowSQLType: []common.SQLType{
			common.StmtRedisGet,
		},
		QueryBeforeHook: func(pha *common.PrevHookArgs) error {
			fmt.Println("query before")
			if strings.Contains(pha.SQL, "name") {
				return errors.New("sss")
			}
			fmt.Printf("%+v", pha)
			return nil
		},

		QueryAfterHook: func(pha *common.PostHookArgs) {
			fmt.Println("query after")
			fmt.Printf("%+v", pha)
		},
	}

	t.Run("redis query forbidden command", func(t *testing.T) {
		SQLList := []string{
			`ttl game`,
		}

		for _, sql := range SQLList {
			fakeQueryMeta := &common.QueryMeta{
				Action: common.ActionSQLQuery,
				Schema: "db1",
				SQL:    SQLBase64(sql),
			}
			reqBody, _ := json.Marshal(fakeQueryMeta)

			mockHTTPReq(t, redisConsole, opt, reqBody, "/console/redis")
		}
	})

	t.Run("redis query valid command", func(t *testing.T) {
		SQLList := []string{
			`  get
			game`,
		}

		for _, sql := range SQLList {
			fakeQueryMeta := &common.QueryMeta{
				Action: common.ActionSQLQuery,
				Schema: "db1",
				SQL:    SQLBase64(sql),
			}
			reqBody, _ := json.Marshal(fakeQueryMeta)

			mockHTTPReq(t, redisConsole, opt, reqBody, "/console/redis")
		}
	})

	t.Run("redis query with query hook failed", func(t *testing.T) {
		SQLList := []string{
			`get name`,
		}

		for _, sql := range SQLList {
			fakeQueryMeta := &common.QueryMeta{
				Action: common.ActionSQLQuery,
				Schema: "db1",
				SQL:    SQLBase64(sql),
			}
			reqBody, _ := json.Marshal(fakeQueryMeta)

			mockHTTPReq(t, redisConsole, opt, reqBody, "/console/redis")
		}
	})
}

func TestRedisConsoleWithTurnOffSystemIntercept(t *testing.T) {
	redisConsole := console.NewRedisConsole()

	opt := &common.HandlerOptions{
		Conn: common.ConnConfig{
			IP:   "172.168.2.24",
			Port: 16379,
		},
		QueryOpt: common.QueryOptions{
			Timeout: 3,
		},
		AllowSQLType: []common.SQLType{
			common.StmtRedisGet,
		},
		IsIgnoreSystemIntercept: true,
		QueryBeforeHook: func(pha *common.PrevHookArgs) error {
			fmt.Println("query before")
			if strings.Contains(pha.SQL, "name") {
				return errors.New("sss")
			}
			fmt.Printf("%+v", pha)
			return nil
		},

		QueryAfterHook: func(pha *common.PostHookArgs) {
			fmt.Println("query after")
			fmt.Printf("%+v", pha)
		},
	}

	t.Run("redis query forbidden command", func(t *testing.T) {
		SQLList := []string{
			`ttl game`,
		}

		for _, sql := range SQLList {
			fakeQueryMeta := &common.QueryMeta{
				Action: common.ActionSQLQuery,
				Schema: "db1",
				SQL:    SQLBase64(sql),
			}
			reqBody, _ := json.Marshal(fakeQueryMeta)

			mockHTTPReq(t, redisConsole, opt, reqBody, "/console/redis")
		}
	})

	t.Run("redis query with timeout", func(t *testing.T) {
		SQLList := []string{
			`blpop 777 10`,
		}

		for _, sql := range SQLList {
			fakeQueryMeta := &common.QueryMeta{
				Action: common.ActionSQLQuery,
				Schema: "db1",
				SQL:    SQLBase64(sql),
			}
			reqBody, _ := json.Marshal(fakeQueryMeta)

			mockHTTPReq(t, redisConsole, opt, reqBody, "/console/redis")
		}
	})
}

func BenchmarkRedisConsoleWithDefaultOpt(b *testing.B) {
	redisConsole := console.NewRedisConsole()

	opt := &common.HandlerOptions{
		Conn: common.ConnConfig{
			IP:   "172.168.2.24",
			Port: 16379,
		},
		QueryOpt: common.QueryOptions{
			Timeout: 600,
		},
	}

	SQLList := []string{
		`ttl game`,
		`exists game`,
		`type game`,
		`scan 0`,
		`get game`,
		`mget game game01`,
		`strlen game`,
		`llen money`,
		`lrange money 0 -1`,
		`lindex money 0`,
		`hgetall student`,
		`hexists student name`,
		`hget student name`,
		`hkeys student`,
		`hvals student`,
		`scard ppp`,
		`smembers ppp`,
		`sismember ppp 111`,
		`sismember ppp 999`,
		`sdiff ppp ttt`,
		`sunion ppp ttt`,
		`zcard setganme`,
		`zrange setganme 0 -1`,
		`zrank setgname li`,
		`zcount setgname 0 1`,
		`zscore setgname li`,
		`zrangebyscore setganme +inf -inf`,
	}

	b.Run("query valid cmc", func(b *testing.B) {
		for n := 0; n <= b.N; n++ {
			for _, sql := range SQLList {
				fakeQueryMeta := &common.QueryMeta{
					Action: common.ActionSQLQuery,
					Schema: "db1",
					SQL:    SQLBase64(sql),
				}
				reqBody, _ := json.Marshal(fakeQueryMeta)

				fakeReq, err := http.NewRequest(http.MethodPost, "/console/redis", bytes.NewReader(reqBody))
				require.NoError(b, err)

				fakeReq.Header.Set("Content-Type", "application/json")

				fakeResp := httptest.NewRecorder()

				console.Handler(fakeResp, fakeReq, "console/redis", redisConsole, opt)

				respBody := fakeResp.Result().Body
				defer respBody.Close()

				respBodyByte, err := ioutil.ReadAll(respBody)
				require.NoError(b, err)

				fmt.Println(string(respBodyByte))
			}
		}
	})
}

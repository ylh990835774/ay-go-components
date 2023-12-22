package test

import (
	"bytes"
	"encoding/base64"
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

func SQLBase64(sql string) string {
	return base64.StdEncoding.EncodeToString([]byte(sql))
}

func mockHTTPReq(t *testing.T, cle console.Console, opt *common.HandlerOptions, req []byte, url string) {
	fakeReq, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(req))
	require.NoError(t, err)

	fakeReq.Header.Set("Content-Type", "application/json")

	fakeResp := httptest.NewRecorder()

	console.Handler(fakeResp, fakeReq, url, cle, opt)

	respBody := fakeResp.Result().Body
	defer respBody.Close()

	respBodyByte, err := ioutil.ReadAll(respBody)
	require.NoError(t, err)

	fmt.Println(string(respBodyByte))
}

func TestMySQLConsoleWithDefaultOpt(t *testing.T) {
	mysqlConsole := console.NewMySQLConsole()

	opt := &common.HandlerOptions{
		Conn: common.ConnConfig{
			IP:       "172.168.1.53",
			Port:     13306,
			UserName: "root",
			Password: "aykj83752661",
		},
	}

	t.Run("fetch schema", func(t *testing.T) {
		fakeQueryMeta := &common.QueryMeta{
			Action: common.ActionFetchSchema,
		}
		reqBody, _ := json.Marshal(fakeQueryMeta)

		mockHTTPReq(t, mysqlConsole, opt, reqBody, "/console/mysql")
	})

	t.Run("fetch table", func(t *testing.T) {
		fakeQueryMeta := &common.QueryMeta{
			Action: common.ActionFetchTable,
			Schema: "alarm_server_local",
		}
		reqBody, _ := json.Marshal(fakeQueryMeta)

		mockHTTPReq(t, mysqlConsole, opt, reqBody, "/console/mysql")
	})

	t.Run("sql query with forbidden SQL", func(t *testing.T) {
		SQList := []string{
			`delete from alarm_server_local where id=1`,
			`update console_test set id=1 where money !=''`,
			`drop table console_test`,
			`drop database alarm_server_local`,
			`
			   drop table
			console_test`,
		}

		for _, sql := range SQList {
			fakeQueryMeta := &common.QueryMeta{
				Action: common.ActionSQLQuery,
				Schema: "alarm_server_local",
				Table:  "alert",
				SQL:    SQLBase64(sql),
			}
			reqBody, _ := json.Marshal(fakeQueryMeta)

			mockHTTPReq(t, mysqlConsole, opt, reqBody, "/console/mysql")
		}
	})

	t.Run("sql query with valid SQL", func(t *testing.T) {
		SQLList := []string{
			`desc console_test`,
			`explain select * from console_test`,
			`select * from console_test`,
			`show processlist`,
		}

		for _, sql := range SQLList {
			fakeQueryMeta := &common.QueryMeta{
				Action: common.ActionSQLQuery,
				Schema: "alarm_server_local",
				Table:  "console_test",
				SQL:    SQLBase64(sql),
			}
			reqBody, _ := json.Marshal(fakeQueryMeta)

			mockHTTPReq(t, mysqlConsole, opt, reqBody, "/console/mysql")
		}
	})

	t.Run("sql query with mulit SQL", func(t *testing.T) {
		fakeQueryMeta := &common.QueryMeta{
			Action: common.ActionSQLQuery,
			Schema: "alarm_server_local",
			Table:  "console_test",
			SQL: SQLBase64(`
			select * from console_test;
			desc console_test;
			`),
		}
		reqBody, _ := json.Marshal(fakeQueryMeta)

		mockHTTPReq(t, mysqlConsole, opt, reqBody, "/console/mysql")
	})
}

func TestMySQLConsoleWithUserCustomOpt(t *testing.T) {
	mysqlConsole := console.NewMySQLConsole()

	opt := &common.HandlerOptions{
		Conn: common.ConnConfig{
			IP:       "172.168.1.53",
			Port:     13306,
			UserName: "root",
			Password: "aykj83752661",
		},

		QueryOpt: common.QueryOptions{
			Timeout: 3,
		},

		AllowSQLType: []common.SQLType{
			common.StmtSelect,
			common.StmtExplain,
			common.StmtShow,
			common.StmtOther,
			common.StmtUpdate,
			common.StmtDelete,
			common.StmtInsert,
		},

		QueryBeforeHook: func(pha *common.PrevHookArgs) error {
			fmt.Printf("query before\n")
			fmt.Printf("%+v\n", pha)
			if strings.Contains(pha.SQL, "explain") {
				return errors.New("xxx")
			}

			return nil
		},

		QueryAfterHook: func(pha *common.PostHookArgs) {
			fmt.Printf("query after\n")
			fmt.Printf("%+v\n", pha)
		},
	}

	t.Run("Test forbidden SQL", func(t *testing.T) {
		fakeQueryMeta := &common.QueryMeta{
			Action: common.ActionSQLQuery,
			Schema: "alarm_server_local",
			Table:  "console_test",
			SQL: SQLBase64(`
			drop database alarm_server_local
			`),
		}
		reqBody, _ := json.Marshal(fakeQueryMeta)

		mockHTTPReq(t, mysqlConsole, opt, reqBody, "/console/mysql")
	})

	t.Run("Test query timeout", func(t *testing.T) {
		fakeQueryMeta := &common.QueryMeta{
			Action: common.ActionSQLQuery,
			Schema: "alarm_server_local",
			Table:  "console_test",
			SQL: SQLBase64(`
			do sleep(4)
			`),
		}
		reqBody, _ := json.Marshal(fakeQueryMeta)

		mockHTTPReq(t, mysqlConsole, opt, reqBody, "/console/mysql")
	})

	t.Run("Test query before hook failed", func(t *testing.T) {
		fakeQueryMeta := &common.QueryMeta{
			Action: common.ActionSQLQuery,
			Schema: "alarm_server_local",
			Table:  "console_test",
			SQL: SQLBase64(`
			explain select * from alert
			`),
		}
		reqBody, _ := json.Marshal(fakeQueryMeta)

		mockHTTPReq(t, mysqlConsole, opt, reqBody, "/console/mysql")
	})

	t.Run("Test common query SQL", func(t *testing.T) {
		sqlList := []string{
			`desc alert`,
			`show tables`,
			`do sleep(2)`,
			`select * from alert`,
			`select * from raw_alert_event`,
		}

		for _, sql := range sqlList {
			fakeQueryMeta := &common.QueryMeta{
				Action: common.ActionSQLQuery,
				Schema: "alarm_server_local",
				Table:  "alert",
				SQL:    SQLBase64(sql),
			}
			reqBody, _ := json.Marshal(fakeQueryMeta)

			mockHTTPReq(t, mysqlConsole, opt, reqBody, "/console/mysql")

		}
	})

	t.Run("Test DDL or DML SQL", func(t *testing.T) {
		sqlList := []string{
			`update alert set severity = 21 where id <= 3`,
		}

		for _, sql := range sqlList {
			fakeQueryMeta := &common.QueryMeta{
				Action: common.ActionSQLQuery,
				Schema: "alarm_server_local",
				Table:  "alert",
				SQL:    SQLBase64(sql),
			}
			reqBody, _ := json.Marshal(fakeQueryMeta)

			mockHTTPReq(t, mysqlConsole, opt, reqBody, "/console/mysql")

		}
	})

	t.Run("Test query with empty sql", func(t *testing.T) {
		fakeQueryMeta := &common.QueryMeta{
			Action: common.ActionSQLQuery,
			Schema: "alarm_server_local",
			Table:  "console_test",
		}
		reqBody, _ := json.Marshal(fakeQueryMeta)

		mockHTTPReq(t, mysqlConsole, opt, reqBody, "/console/mysql")
	})
}

func TestMySQLConsoleWithIgnoreSystemIntercpet(t *testing.T) {
	mysqlConsole := console.NewMySQLConsole()

	opt := &common.HandlerOptions{
		Conn: common.ConnConfig{
			IP:       "172.168.1.53",
			Port:     13306,
			UserName: "root",
			Password: "aykj83752661",
		},

		QueryOpt: common.QueryOptions{
			Timeout: 3,
		},

		AllowSQLType: []common.SQLType{
			common.StmtSelect,
		},
		IsIgnoreSystemIntercept: true,

		QueryBeforeHook: func(pha *common.PrevHookArgs) error {
			fmt.Printf("query before\n")
			fmt.Printf("%+v\n", pha)
			if strings.Contains(pha.SQL, "explain") {
				return errors.New("xxx")
			}

			return nil
		},

		QueryAfterHook: func(pha *common.PostHookArgs) {
			fmt.Printf("query after\n")
			fmt.Printf("%+v\n", pha)
		},
	}

	t.Run("Test forbidden sql", func(t *testing.T) {
		sqlList := []string{
			`update alert set severity = 24 where id <= 9`,
		}

		for _, sql := range sqlList {
			fakeQueryMeta := &common.QueryMeta{
				Action: common.ActionSQLQuery,
				Schema: "alarm_server_local",
				Table:  "alert",
				SQL:    SQLBase64(sql),
			}
			reqBody, _ := json.Marshal(fakeQueryMeta)

			mockHTTPReq(t, mysqlConsole, opt, reqBody, "/console/mysql")

		}
	})
}

func BenchmarkMySQLConsoleWithDefaultOpt(b *testing.B) {
	mysqlConsole := console.NewMySQLConsole()

	opt := &common.HandlerOptions{
		Conn: common.ConnConfig{
			IP:       "172.168.1.53",
			Port:     13306,
			UserName: "root",
			Password: "aykj83752661",
		},
	}

	b.Run("query common sql", func(b *testing.B) {
		for n := 0; n <= b.N; n++ {
			SQLList := []string{
				`select * from alert`,
			}

			for _, sql := range SQLList {
				fakeQueryMeta := &common.QueryMeta{
					Action: common.ActionSQLQuery,
					Schema: "alarm_server_local",
					Table:  "console_test",
					SQL:    SQLBase64(sql),
				}
				reqBody, _ := json.Marshal(fakeQueryMeta)

				fakeReq, err := http.NewRequest(http.MethodPost, "/console/mysql", bytes.NewReader(reqBody))
				require.NoError(b, err)

				fakeReq.Header.Set("Content-Type", "application/json")

				fakeResp := httptest.NewRecorder()

				console.Handler(fakeResp, fakeReq, "console/mysql", mysqlConsole, opt)

				respBody := fakeResp.Result().Body
				defer respBody.Close()

				respBodyByte, err := ioutil.ReadAll(respBody)
				require.NoError(b, err)

				fmt.Println(string(respBodyByte))
			}
		}
	})
}

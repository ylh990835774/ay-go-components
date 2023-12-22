package console

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"git.qpaas.com/go-components/webconsole/mock"
	"git.qpaas.com/go-components/webconsole/pkg/common"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestHandlerRoute(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()

	fakeconsole := mock.NewMockConsole(controller)

	fakeHandlerOpt := &common.HandlerOptions{}
	fakeSchemaList := []string{"database01", "database02"}
	fakeTableList := []string{"table01", "table02"}
	fakeQuerySet := &common.QuerySet{
		EngineType: common.MySQLEngine,
		Action:     common.ActionSQLQuery,

		IsExecute: true,
		ExecuteAt: time.Now(),
		Err:       nil,

		SQL: "select * from table",

		Total: 100,

		Columns: []string{"field1", "field2"},

		Rows: []common.Row{{
			"field1": "val1",
			"field2": "val2",
		}, {
			"field21": "val21",
			"field22": "val22",
		}},

		AffectedRows: 100,
	}

	fakeconsole.EXPECT().SchemaHandler(fakeHandlerOpt).Return(fakeSchemaList, nil).AnyTimes()
	fakeconsole.EXPECT().TableHandler("database01", fakeHandlerOpt).Return(fakeTableList, nil).AnyTimes()
	fakeconsole.EXPECT().QueryHandler("database01", "table01", "select * from table01", fakeHandlerOpt).Return(fakeQuerySet).AnyTimes()
	fakeconsole.EXPECT().ConsoleType().Return("fakeConsole").AnyTimes()

	t.Run("fetch schema", func(t *testing.T) {
		queryMeta := &common.QueryMeta{
			Action: common.ActionFetchSchema,
		}

		reqBody, _ := json.Marshal(queryMeta)
		fakeReq := httptest.NewRequest(http.MethodGet, "/user-custom-path", bytes.NewBuffer(reqBody))
		fakeReq.Header.Set("Content-Type", "application/json")

		fakeResp := httptest.NewRecorder()

		handlerGateway(fakeResp, fakeReq, fakeconsole, fakeHandlerOpt)

		fakeRespBody := fakeResp.Result().Body
		defer fakeRespBody.Close()

		fakeRespBodyByte, err := ioutil.ReadAll(fakeRespBody)
		require.NoError(t, err)

		fmt.Printf("%s\n", fakeRespBodyByte)
	})

	t.Run("fetch table", func(t *testing.T) {
		queryMeta := &common.QueryMeta{
			Action: common.ActionFetchTable,
			Schema: "database01",
		}

		reqBody, _ := json.Marshal(queryMeta)
		fakeReq := httptest.NewRequest(http.MethodGet, "/user-custom-path", bytes.NewBuffer(reqBody))
		fakeReq.Header.Set("Content-Type", "application/json")

		fakeResp := httptest.NewRecorder()

		handlerGateway(fakeResp, fakeReq, fakeconsole, fakeHandlerOpt)

		fakeRespBody := fakeResp.Result().Body
		defer fakeRespBody.Close()

		fakeRespBodyByte, err := ioutil.ReadAll(fakeRespBody)
		require.NoError(t, err)

		fmt.Printf("%s\n", fakeRespBodyByte)
	})

	t.Run("sql query", func(t *testing.T) {
		queryMeta := &common.QueryMeta{
			Action: common.ActionSQLQuery,
			Schema: "database01",
			Table:  "table01",
			SQL:    "select * from table01",
		}

		reqBody, _ := json.Marshal(queryMeta)
		fakeReq := httptest.NewRequest(http.MethodGet, "/user-custom-path", bytes.NewBuffer(reqBody))
		fakeReq.Header.Set("Content-Type", "application/json")

		fakeResp := httptest.NewRecorder()

		handlerGateway(fakeResp, fakeReq, fakeconsole, fakeHandlerOpt)

		fakeRespBody := fakeResp.Result().Body
		defer fakeRespBody.Close()

		fakeRespBodyByte, err := ioutil.ReadAll(fakeRespBody)
		require.NoError(t, err)

		fmt.Printf("%s\n", fakeRespBodyByte)
	})

	t.Run("fetch staticfile", func(t *testing.T) {
		fakeReq := httptest.NewRequest(http.MethodGet, "/mysql/console/css/29.1b4fb940.css", nil)
		fakeReq.Header.Set("Content-Type", "text/html")
		fakeResp := httptest.NewRecorder()

		Handler(fakeResp, fakeReq, "/mysql/console", fakeconsole, nil)

		fakeRespBody := fakeResp.Result().Body
		defer fakeRespBody.Close()

		fakeRespBodyByte, err := ioutil.ReadAll(fakeRespBody)
		require.NoError(t, err)

		fmt.Println(string(fakeRespBodyByte))
	})
	t.Run("fetch index html", func(t *testing.T) {
		fakeReq := httptest.NewRequest(http.MethodGet, "/mysql/console/index.html", nil)
		fakeReq.Header.Set("Content-Type", "text/html")
		fakeResp := httptest.NewRecorder()

		Handler(fakeResp, fakeReq, "/mysql/console", fakeconsole, nil)

		fakeRespBody := fakeResp.Result().Body
		defer fakeRespBody.Close()

		fakeRespBodyByte, err := ioutil.ReadAll(fakeRespBody)
		require.NoError(t, err)

		fmt.Println(string(fakeRespBodyByte))
	})
}

func TestOpenStaticFiles(t *testing.T) {
	for _, filename := range []string{
		"css/29.1b4fb940.css",
		"fonts/fontello.3f1fdcf0.ttf",
		"img/fontello.d34249f0.svg",
		"js/app.7622fc0f.js",
		"favicon.ico",
		"index.html",
	} {
		_, err := VirtualFS.Open(filename)
		require.NoError(t, err, filename)
	}
}

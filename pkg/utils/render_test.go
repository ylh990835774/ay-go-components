package utils

import (
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/ylh990835774/ay-go-components/pkg/common"
	"github.com/ylh990835774/ay-go-components/pkg/inerr"
)

func TestXxx(t *testing.T) {
	t.Run("test render err", func(t *testing.T) {
		fakeResp := httptest.NewRecorder()
		RenderErr(fakeResp, inerr.ErrSQLForbidden)

		respBody := fakeResp.Result().Body
		defer respBody.Close()

		respByte, err := ioutil.ReadAll(fakeResp.Result().Body)
		require.NoError(t, err)

		fmt.Printf("%s\n", respByte)
	})

	t.Run("test render schema list", func(t *testing.T) {
		fakeResp := httptest.NewRecorder()
		schemaList := []string{"database01", "database02"}

		RenderData(fakeResp, "fetch schema succeed", schemaList)

		respBody := fakeResp.Result().Body
		defer respBody.Close()

		respBodyByte, err := ioutil.ReadAll(respBody)
		require.NoError(t, err)

		fmt.Printf("%s\n", respBodyByte)
	})

	t.Run("test render sql query", func(t *testing.T) {
		fakeResp := httptest.NewRecorder()

		mockQuerySet := &common.QuerySet{
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

		RenderData(fakeResp, "sql query succeed", mockQuerySet)

		respBody := fakeResp.Result().Body
		defer respBody.Close()

		respBodyByte, err := ioutil.ReadAll(respBody)
		require.NoError(t, err)

		fmt.Printf("%s\n", respBodyByte)
	})
}

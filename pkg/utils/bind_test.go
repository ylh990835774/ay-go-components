package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"git.qpaas.com/go-components/webconsole/pkg/common"
	"github.com/stretchr/testify/require"
)

func TestGetBody(t *testing.T) {
	actualReqBody := &common.QueryMeta{}
	rawReqBody := &common.QueryMeta{
		Action: "fetchSchema",
		Schema: "mockdatabase",
		Table:  "mocktable",
		SQL:    "select * from mocktable",
	}

	mockReqBodyByte, err := json.Marshal(rawReqBody)
	require.NoError(t, err, fmt.Sprintf("%+v", rawReqBody))

	mockReq := httptest.NewRequest(http.MethodPost, "/console/mock", bytes.NewBuffer(mockReqBodyByte))
	mockReq.Header.Set("Content-Type", "application/json")

	require.NoError(t, GetBody(mockReq, actualReqBody))
	fmt.Printf("%+v", actualReqBody)
}

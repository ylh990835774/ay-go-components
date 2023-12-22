package utils

import (
	"encoding/json"
	"net/http"

	"git.qpaas.com/go-components/webconsole/pkg/common"
)

var jsonContentType = []string{"application/json; charset=utf-8"}

func writeContentType(w http.ResponseWriter, value []string) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = value
	}
}

// writeJSON marshals the given interface object and writes it with json content type.
func writeJSON(w http.ResponseWriter, obj interface{}) error {
	writeContentType(w, jsonContentType)
	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		return err
	}
	_, err = w.Write(jsonBytes)
	return err
}

func RenderErr(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusOK)

	resp := &common.Resp{
		Code:    500,
		Message: err.Error(),
		Result:  nil,
	}

	writeJSON(w, resp)
}

func RenderData(w http.ResponseWriter, msg string, data interface{}) {
	w.WriteHeader(http.StatusOK)

	resp := &common.Resp{
		Code:    200,
		Message: msg,
		Result:  data,
	}

	writeJSON(w, resp)
}

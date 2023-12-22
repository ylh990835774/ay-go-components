package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/ylh990835774/ay-go-components/pkg/inerr"
)

const (
	HeaderContentType   = "Content-Type"
	MIMEApplicationJSON = "application/json"
)

func GetBody(req *http.Request, i interface{}) error {
	if req.ContentLength == 0 {
		return nil
	}

	ctype := strings.ToLower(req.Header.Get(HeaderContentType))
	switch {
	case strings.HasPrefix(ctype, MIMEApplicationJSON):
		if err := json.NewDecoder(req.Body).Decode(i); err != nil {
			if ute, ok := err.(*json.UnmarshalTypeError); ok {
				return fmt.Errorf("unmarshal type error: expected=%v, got=%v, field=%v, offset=%v", ute.Type, ute.Value, ute.Field, ute.Offset)
			} else if se, ok := err.(*json.SyntaxError); ok {
				return fmt.Errorf("syntax error: offset=%v, error=%v", se.Offset, se.Error())
			}
			return err
		}
	default:
		return inerr.ErrUnsupportedMediaType
	}

	return nil
}

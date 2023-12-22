package console

import (
	"embed"
	"encoding/base64"
	"io/fs"
	"net/http"
	"strings"

	"github.com/pkg/errors"
	"github.com/ylh990835774/ay-go-components/pkg/common"
	"github.com/ylh990835774/ay-go-components/pkg/engine"
	"github.com/ylh990835774/ay-go-components/pkg/inerr"
	"github.com/ylh990835774/ay-go-components/pkg/utils"
)

//go:embed dist/*
var dist embed.FS

// FS holds embedded swagger ui files
var VirtualFS, _ = fs.Sub(dist, "dist")

type Console interface {
	ConsoleType() string

	Fork(common.ConnConfig, string) (engine.Engine, error)
	Destory(engine.Engine)

	SchemaHandler(opt *common.HandlerOptions) ([]string, error)
	TableHandler(schema string, opt *common.HandlerOptions) ([]string, error)
	QueryHandler(schema string, table string, sql string, opt *common.HandlerOptions) *common.QuerySet
}

// route entrypoint
func Handler(w http.ResponseWriter, req *http.Request, consolePath string, cle Console, opt *common.HandlerOptions) {
	switch req.Method {
	case http.MethodGet:
		staticFileHandler(w, req, consolePath, cle.ConsoleType())
	case http.MethodHead:
		staticFileHandler(w, req, consolePath, cle.ConsoleType())
	case http.MethodPost:
		// handler api request from component page
		handlerGateway(w, req, cle, opt)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// server static files
// because some web framework dose not support
// path prefix match
// seperate method from Handler
// HandlerStaticFile server console staticfile
func HandlerStaticFile(w http.ResponseWriter, req *http.Request, consolePath string, cle Console) {
	staticFileHandler(w, req, consolePath, cle.ConsoleType())
}

// HandlerAPI server ajax request from conosle
func HandlerAPI(w http.ResponseWriter, req *http.Request, cle Console, opt *common.HandlerOptions) {
	handlerGateway(w, req, cle, opt)
}

// handlerGateway uni entrypoint of sub handler
func handlerGateway(w http.ResponseWriter, req *http.Request, cle Console, opt *common.HandlerOptions) {
	queryMeta := &common.QueryMeta{}

	err := utils.GetBody(req, queryMeta)
	if err != nil {
		utils.RenderErr(w, errors.Wrap(err, "parse request params failed"))
		return
	}

	switch queryMeta.Action {
	case common.ActionFetchSchema:
		result, err := cle.SchemaHandler(opt)
		if err != nil {
			utils.RenderErr(w, errors.Wrap(err, "fetch schema failed"))
			return
		}

		utils.RenderData(w, "fetch schema succeed", result)
	case common.ActionFetchTable:
		result, err := cle.TableHandler(queryMeta.Schema, opt)
		if err != nil {
			utils.RenderErr(w, errors.Wrap(err, "fetch table failed"))
			return
		}

		utils.RenderData(w, "fetch table succeed", result)
	case common.ActionSQLQuery:
		// decode SQL
		// SQL is encode by base64
		decodeSQLByte, err := base64.StdEncoding.DecodeString(queryMeta.SQL)
		if err != nil {
			utils.RenderErr(w, errors.Wrap(err, "query failed"))
			return
		}

		result := cle.QueryHandler(queryMeta.Schema, queryMeta.Table, string(decodeSQLByte), opt)
		if result.Err != nil {
			utils.RenderErr(w, errors.Wrap(result.Err, "query failed"))
			return
		}

		utils.RenderData(w, "query succeed", result)
	default:
		utils.RenderErr(w, errors.Wrap(inerr.ErrUnsupportedOperation, queryMeta.Action))
	}
}

func staticFileHandler(w http.ResponseWriter, req *http.Request, consolePath string, consoleType string) {
	// handler console  static files about component fronentend pages
	if strings.Contains(consolePath, ":") || strings.Contains(consolePath, "*") {
		utils.RenderErr(w, inerr.ErrConsolePathNotSupport)
		return
	}

	fileServer := http.StripPrefix(consolePath, http.FileServer(http.FS(VirtualFS)))
	fileServer.ServeHTTP(w, req)
}

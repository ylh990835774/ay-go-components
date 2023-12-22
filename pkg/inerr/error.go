package inerr

import "errors"

var ErrUnImplement = errors.New("UnImplement")
var ErrEngineTypeUnknown = errors.New("engine type unknown")
var ErrUnsupportedMediaType = errors.New("http server not support media type")
var ErrUnsupportedOperation = errors.New("console component not support operation type")

var ErrFieldEmpty = errors.New("field is empty")

var ErrSchemaEmpty = errors.New("schema should be provided")
var ErrTableEmpty = errors.New("table should be provided")
var ErrSQLEmpty = errors.New("SQL statement should be provided")
var ErrSQLForbidden = errors.New("SQL statement forbidden")

var ErrRedisCMDUnknown = errors.New("redis cmd unknown")
var ErrRedisCMDUnSupported = errors.New("redis command unsupported now")
var ErrRedisCMDForbidden = errors.New("redis cmd forbidden")
var ErrRedisCMDEmpty = errors.New("redis command should be provided")
var ErrRedisKeyNotExist = errors.New("redis key not exist")
var ErrRedisKeyEmpty = errors.New("key should be provided")
var ErrRedisParseKey = errors.New("parse key from redis command failed")
var ErrRedisSchemaFetchFailed = errors.New("redis schema fetch failed")

var ErrConsolePathNotSupport = errors.New("console router path can not container '*' or ':' when console serving a static folder in console internal")

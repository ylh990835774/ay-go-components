package common

import "time"

const MySQLConsole = "MysqlConsole"
const MongoDBConsole = "MongoDBConsole"
const RedisConsole = "RedisConsole"

const MySQLEngine = "MySQLEngine"
const MongoDBEngine = "MongoDBEngine"
const RedisEngine = "RedisEngine"

const ActionFetchSchema = "fetchSchema"
const ActionFetchTable = "fetchTable"
const ActionSQLQuery = "sqlQuery"

// redis value 数据类型
const RedisKeyTypeNone = "none" // key不存在
const RedisKeyTypeStr = "string"
const RedisKeyTypeList = "list"
const RedisKeyTypeSet = "set"
const RedisKeyTypeZSet = "zset"
const RedisKeyTypeHash = "hash"

const (
	// MySQL常用读写命令
	StmtSelect SQLType = iota
	StmtStream
	StmtInsert
	StmtReplace
	StmtUpdate
	StmtDelete
	StmtDDL
	StmtBegin
	StmtCommit
	StmtRollback
	StmtSet
	StmtShow
	stmtUse
	StmtOther
	StmtUnknown
	StmtComment
	StmtPriv
	StmtExplain
	StmtSavepoint
	StmtSRollback
	StmtRelease
	StmtVStream
	StmtLockTables
	StmtUnlockTables
	StmtFlush
	StmtCallProc
	StmtRevert
	StmtShowMigrationLogs

	// redis支持的常用命令
	// Key命令
	// 读命令
	StmtRedisType
	StmtRedisExists
	StmtRedisTTL
	StmtRedisScan
	// 写命令
	StmtRedisDEL
	StmtRedisExpire
	StmtRedisExpireAt

	// string命令
	// 读命令
	StmtRedisGet
	StmtRedisMGet
	StmtRedisStrLen
	// 写命令
	StmtRedisAppend
	StmtRedisIncr
	StmtRedisIncrBy
	StmtRedisSet
	StmtRedisMSet
	StmtRedisSetEX
	StmtRedisSetNX

	// Hash命令
	// 读命令
	StmtRedisHGetAll
	StmtRedisHExists
	StmtRedisHGet
	StmtRedisHMGet
	StmtRedisHKeys
	StmtRedisHVals
	// 写命令
	StmtRedisHDel
	StmtRedisHSet
	StmtRedisHMSet

	// list命令
	// 读命令
	StmtRedisLLen
	StmtRedisLRange
	StmtRedisLIndex
	// 写命令
	StmtRedisLPop
	StmtRedisRPop
	StmtRedisLPush
	StmtRedisRPush
	StmtRedisLInsert

	// set命令
	// 读命令
	StmtRedisSCard
	StmtRedisSMembers
	StmtRedisSisMember
	StmtRedisSDiff
	StmtRedisSUnion
	// 写命令
	StmtRedisSAdd
	StmtRedisSRem

	// sortedSet命令
	// 读命令
	StmtRedisZCard
	StmtRedisZRange
	StmtRedisZRank
	StmtRedisZCount
	StmtRedisZScore
	StmtRedisZRangeByScore
	// 写命令
	StmtRedisZAdd
	StmtRedisZRem
)

var DefaultRedisWhiteCMD = []SQLType{
	// Key命令
	StmtRedisType,
	StmtRedisExists,
	StmtRedisTTL,
	StmtRedisScan,
	// string命令
	StmtRedisGet,
	StmtRedisMGet,
	StmtRedisStrLen,
	// list命令
	StmtRedisLLen,
	StmtRedisLRange,
	StmtRedisLIndex,
	// hash命令
	StmtRedisHGetAll,
	StmtRedisHExists,
	StmtRedisHGet,
	StmtRedisHMGet,
	StmtRedisHKeys,
	StmtRedisHVals,
	// set命令
	StmtRedisSCard,
	StmtRedisSMembers,
	StmtRedisSisMember,
	StmtRedisSDiff,
	StmtRedisSUnion,
	// sortSet命令
	StmtRedisZCard,
	StmtRedisZRange,
	StmtRedisZRank,
	StmtRedisZCount,
	StmtRedisZScore,
	StmtRedisZRangeByScore,
}

var RedisCMDTOSQLType = map[string]SQLType{
	"type":          StmtRedisType,
	"exists":        StmtRedisExists,
	"ttl":           StmtRedisTTL,
	"scan":          StmtRedisScan,
	"del":           StmtRedisDEL,
	"expire":        StmtRedisExpire,
	"get":           StmtRedisGet,
	"mget":          StmtRedisMGet,
	"strlen":        StmtRedisStrLen,
	"append":        StmtRedisAppend,
	"incr":          StmtRedisIncr,
	"incrby":        StmtRedisIncrBy,
	"set":           StmtRedisSet,
	"mset":          StmtRedisMSet,
	"setex":         StmtRedisSetEX,
	"setnx":         StmtRedisSetNX,
	"hgetall":       StmtRedisHGetAll,
	"hexists":       StmtRedisExists,
	"hget":          StmtRedisHGet,
	"hmget":         StmtRedisHMGet,
	"hkeys":         StmtRedisHKeys,
	"hvals":         StmtRedisHVals,
	"hdel":          StmtRedisHDel,
	"hset":          StmtRedisHSet,
	"hmset":         StmtRedisHMSet,
	"llen":          StmtRedisLLen,
	"lrange":        StmtRedisLRange,
	"lindex":        StmtRedisLIndex,
	"lpop":          StmtRedisLPop,
	"rpop":          StmtRedisRPop,
	"lpush":         StmtRedisLPush,
	"rpush":         StmtRedisRPush,
	"linsert":       StmtRedisLInsert,
	"scard":         StmtRedisSCard,
	"smembers":      StmtRedisSMembers,
	"sismember":     StmtRedisSisMember,
	"sdiff":         StmtRedisSDiff,
	"sunion":        StmtRedisSUnion,
	"sadd":          StmtRedisSAdd,
	"srem":          StmtRedisSRem,
	"zcard":         StmtRedisZCard,
	"zrange":        StmtRedisZRange,
	"zrank":         StmtRedisZRank,
	"zcount":        StmtRedisZCount,
	"zscore":        StmtRedisZScore,
	"zrangebyscore": StmtRedisZRangeByScore,
	"zadd":          StmtRedisZAdd,
	"zrem":          StmtRedisZRem,
}

// ConnConfig connect information
type ConnConfig struct {
	IP       string
	Port     int
	UserName string
	Password string
}

type SQLType int

type QueryOptions struct {
	Timeout int64 // 查询超时(秒)
}

// QueryMeta request params about query operation
type QueryMeta struct {
	Action string `json:"action"` // fetchSchema|fetchTable|sqlQuery
	Schema string `json:"schema"`
	Table  string `json:"table"` // 在Redis中取值为Key
	SQL    string `json:"sql"`
}

// Resp response about request
type Resp struct {
	Code    int         `json:"code"`
	Message string      `json:"msg"`
	Result  interface{} `json:"result"`
}

type Row map[string]interface{}

// QuerySet
type QuerySet struct {
	EngineType string `json:"-"`
	Action     string `json:"-"`

	IsExecute     bool      `json:"-"`
	ExecuteAt     time.Time `json:"-"`
	QueryDuration int64     `json:"-"` // 毫秒

	Err error `json:"-"`

	SQL string `json:"sql"`

	Total int `json:"total"`

	Columns []string `json:"columns"`
	Rows    []Row    `json:"rows"`

	AffectedRows int64 `json:"-"`
}

type PrevHookArgs struct {
	EngineType string
	Action     string

	Schema string
	SQL    string
}

type PostHookArgs struct {
	EngineType string
	Action     string

	IsExecute     bool
	ExecuteAt     time.Time
	QueryDuration int64
	Err           error

	Schema string
	SQL    string

	AffectedRows int64
}

// hooks unImplement
type PreHook func(*PrevHookArgs) error

type PostHook func(*PostHookArgs)

// EngineBase base struct of egine
type EngineBase struct {
	ConnConfig
	QueryPrev PreHook
	QueryPost PostHook
}

func (e *EngineBase) BindPrevHook(hook PreHook) {
	e.QueryPrev = hook
}

func (e *EngineBase) BindPostHook(hook PostHook) {
	e.QueryPost = hook
}

func NewEngineBase(conn ConnConfig) *EngineBase {
	return &EngineBase{
		ConnConfig: conn,
	}
}

// HandlerOptions
/*
AllowSQLType:
设置系统拦截器的命令白名单,不同类型的控制台支持的白名单会有差异
若不设置则系统会按照默认白名单进行拦截（前提是开启了系统拦截）

IsIgnoreSystemIntercept:
true: 关闭系统内置的拦截器，命名是否执行完全取决于用户的QueryBeforeHook钩子函数的逻辑
false: 开启系统内置的拦截器，若用户设置了QueryBeforeHook钩子函数，则命令是否执行还取决于该钩子函数的逻辑

QueryBeforeHook:
执行命令前的钩子函数，用户可以利用该钩子函数进行业务逻辑扩展比如记录，内置拦截器无法满足业务场景等
*/
type HandlerOptions struct {
	Conn                    ConnConfig
	QueryOpt                QueryOptions
	AllowSQLType            []SQLType
	IsIgnoreSystemIntercept bool
	QueryBeforeHook         PreHook
	QueryAfterHook          PostHook
}

// ConsoleBase  base struct of console
type ConsoleBase struct {
}

func NewConsoleBase() *ConsoleBase {
	return &ConsoleBase{}
}

func UnifiedLabel(col []string) []interface{} {
	var s []interface{}
	for i := 0; i < len(col); i++ {
		var t string
		s = append(s, &t)
	}
	return s
}

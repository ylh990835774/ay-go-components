package console

import (
	"fmt"
	"sync"

	"github.com/pkg/errors"
	"github.com/ylh990835774/ay-go-components/pkg/common"
	"github.com/ylh990835774/ay-go-components/pkg/engine"
	"github.com/ylh990835774/ay-go-components/pkg/inerr"
)

type mySQLConsole struct {
	sync.Pool
	*common.ConsoleBase
}

func (m *mySQLConsole) ConsoleType() string {
	return common.MySQLConsole
}

func (m *mySQLConsole) Fork(conn common.ConnConfig, schema string) (engine.Engine, error) {
	eg := m.Get().(*engine.MySQLEngine)

	err := eg.InitialDriver(conn, schema)
	if err != nil {
		return nil, err
	}

	return eg, nil
}

func (m *mySQLConsole) Destory(e engine.Engine) {
	eg := e.(*engine.MySQLEngine)
	eg.Reset()

	m.Put(eg)
}

func (m *mySQLConsole) SchemaHandler(opt *common.HandlerOptions) ([]string, error) {
	// fork engine instance
	eg, err := m.Fork(opt.Conn, "")
	if err != nil {
		return nil, err
	}
	defer m.Destory(eg) // destory engine instance

	// bind hooks
	eg.RegistryQueryPrev(opt.QueryBeforeHook)
	eg.RegistryQueryPost(opt.QueryAfterHook)

	// fetch schema
	schemas, err := eg.Schema()
	if err != nil {
		return nil, err
	}

	return schemas, nil
}

func (m *mySQLConsole) TableHandler(schema string, opt *common.HandlerOptions) ([]string, error) {
	// fork engine instance
	eg, err := m.Fork(opt.Conn, schema)
	if err != nil {
		return nil, err
	}
	defer m.Destory(eg) // destory engine instance

	// bind hooks
	eg.RegistryQueryPrev(opt.QueryBeforeHook)
	eg.RegistryQueryPost(opt.QueryAfterHook)

	// fetch tables
	tables, err := eg.Table(schema)
	if err != nil {
		return nil, err
	}

	return tables, nil
}

func (m *mySQLConsole) QueryHandler(schema string, table string, sql string, opt *common.HandlerOptions) *common.QuerySet {
	// fork engine instance
	eg, err := m.Fork(opt.Conn, schema)
	if err != nil {
		return &common.QuerySet{
			Err: errors.Wrap(err, "mysql engine fork failed"),
		}
	}
	defer m.Destory(eg) // destory engine instance

	// bind hooks
	eg.RegistryQueryPrev(opt.QueryBeforeHook)
	eg.RegistryQueryPost(opt.QueryAfterHook)

	// query
	// sql preCheck inner system
	// set default SQL Allow Rule
	// select、show、desc、explain statement
	defaultAllowSQLType := []common.SQLType{
		common.StmtSelect,
		common.StmtShow,
		common.StmtExplain,
	}
	if opt.AllowSQLType != nil {
		defaultAllowSQLType = opt.AllowSQLType
	}

	// if sql is empty
	// assign desc table as default sql
	if sql == "" {
		sql = fmt.Sprintf("desc %s", table)
	}

	var preProcessSQL string
	var isPass bool

	// valid systemIncepter state
	if opt.IsIgnoreSystemIntercept {
		preProcessSQL = sql
		goto queryMain
	}

	preProcessSQL, isPass, err = engine.MySQLPreCheck(sql, defaultAllowSQLType)
	if err != nil {
		return &common.QuerySet{
			Err: errors.Wrap(err, "sql preCheck failed"),
		}
	}

	if !isPass {
		return &common.QuerySet{
			Err: errors.Wrap(inerr.ErrSQLForbidden, sql),
		}
	}

queryMain:
	// query execute
	var defaultQueryTimeout int64 = 15
	if opt.QueryOpt.Timeout > 0 {
		defaultQueryTimeout = opt.QueryOpt.Timeout
	}

	return eg.Query(schema, table, preProcessSQL, defaultQueryTimeout)
}

// NewMySQLConsole
func NewMySQLConsole() *mySQLConsole {
	return &mySQLConsole{
		sync.Pool{
			New: func() interface{} {
				return engine.NewMySQLEngine()
			},
		},
		common.NewConsoleBase(),
	}
}

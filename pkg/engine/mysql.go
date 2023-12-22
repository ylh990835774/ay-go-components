package engine

import (
	"context"
	"database/sql"
	"encoding/hex"
	"fmt"
	"time"

	mmysql "github.com/go-sql-driver/mysql"
	"github.com/ylh990835774/ay-go-components/pkg/common"
	"github.com/ylh990835774/ay-go-components/pkg/inerr"
	mydriver "gorm.io/driver/mysql"
	"gorm.io/gorm"
	vsqlparser "vitess.io/vitess/go/vt/sqlparser"
)

const (
	BUF                   = 1<<20 - 1
	BLOB_FIELD_NOT_DISPLA = "Blob field cannot be displayed"
)

type MySQLEngine struct {
	driver *gorm.DB
	*common.EngineBase
}

func (m *MySQLEngine) RegistryQueryPrev(hook common.PreHook) {
	m.BindPrevHook(hook)
}

func (m *MySQLEngine) RegistryQueryPost(hook common.PostHook) {
	m.BindPostHook(hook)
}

func (m *MySQLEngine) Close() error {
	orm, err := m.driver.DB()
	if err != nil {
		return err
	}
	return orm.Close()
}

func (m *MySQLEngine) Schema() ([]string, error) {
	rows, err := m.driver.Raw("SHOW DATABASES;").Rows()
	if err != nil {
		return nil, err
	}

	col, _ := rows.Columns()

	tmpDest := common.UnifiedLabel(col)
	if len(tmpDest) == 0 {
		return nil, inerr.ErrFieldEmpty
	}

	schemas := make([]string, 0)
	for rows.Next() {
		if err = rows.Scan(tmpDest...); err != nil {
			return nil, err
		}

		j := *tmpDest[0].(*string)
		schemas = append(schemas, j)
	}
	return schemas, nil
}

func (m *MySQLEngine) Table(schema string) ([]string, error) {
	if schema == "" {
		return nil, inerr.ErrSchemaEmpty
	}

	rows, err := m.driver.Raw("SHOW TABLES;").Rows()
	if err != nil {
		return nil, err
	}

	col, _ := rows.Columns()

	tmpDest := common.UnifiedLabel(col)
	if len(tmpDest) == 0 {
		return nil, inerr.ErrFieldEmpty
	}

	schemas := make([]string, 0)
	for rows.Next() {
		if err = rows.Scan(tmpDest...); err != nil {
			return nil, err
		}

		j := *tmpDest[0].(*string)
		schemas = append(schemas, j)
	}
	return schemas, nil
}

// Query
// include Select、DDL statement and so on
func (m *MySQLEngine) Query(schema string, table string, sql string, timeout int64) *common.QuerySet {
	queryRes := &common.QuerySet{
		EngineType: common.MySQLEngine,
		Action:     common.ActionSQLQuery,
		IsExecute:  false,
		Err:        nil,
	}

	// execute query prev hook
	// query prev hook failed, stop query
	if m.QueryPrev != nil {
		err := m.QueryPrev(&common.PrevHookArgs{
			EngineType: common.MySQLEngine,
			Action:     common.ActionSQLQuery,
			Schema:     schema,
			SQL:        sql,
		})
		if err != nil {
			queryRes.Err = err
			return queryRes
		}
	}

	// registry query post hook
	defer func() {
		if m.QueryPost != nil {
			m.QueryPost(&common.PostHookArgs{
				EngineType:    common.MySQLEngine,
				Action:        common.ActionSQLQuery,
				IsExecute:     queryRes.IsExecute,
				ExecuteAt:     queryRes.ExecuteAt,
				QueryDuration: queryRes.QueryDuration,
				Err:           queryRes.Err,
				Schema:        schema,
				SQL:           sql,
				AffectedRows:  queryRes.AffectedRows,
			})
		}
	}()

	if schema == "" {
		queryRes.Err = inerr.ErrSchemaEmpty
		return queryRes
	}

	if sql == "" {
		queryRes.Err = inerr.ErrSQLEmpty
		return queryRes
	}

	// fetch sql type
	sqlType, err := MySQLSQLType(sql)
	if err != nil {
		queryRes.Err = err
		return queryRes
	}

	// query main
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	// start query
	queryRes.ExecuteAt = time.Now()

	switch sqlType {
	// not query statement
	// use Exec()
	case common.StmtInsert, common.StmtUpdate, common.StmtDelete, common.StmtDDL:
		d := m.driver.WithContext(ctx).Exec(sql)

		// query finished
		queryRes.QueryDuration = time.Since(queryRes.ExecuteAt).Milliseconds()
		queryRes.SQL = sql
		queryRes.IsExecute = true
		queryRes.Err = d.Error
		queryRes.AffectedRows = d.RowsAffected

		return queryRes
	}

	// common quey statement
	rows, err := m.driver.WithContext(ctx).Raw(sql).Rows()

	// query finished
	queryRes.QueryDuration = time.Since(queryRes.ExecuteAt).Milliseconds()

	if err != nil {
		queryRes.Err = err
		return queryRes
	}

	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		queryRes.Err = err
		return queryRes
	}

	rowList := make([]common.Row, 0)
	for rows.Next() {
		results := make(map[string]interface{})
		singleRow := make(map[string]interface{})

		err := mapScan(rows, results)
		if err != nil {
			queryRes.Err = err
			return queryRes
		}

		for key := range results {
			switch r := results[key].(type) {
			case []uint8:
				if len(r) > BUF {
					singleRow[key] = BLOB_FIELD_NOT_DISPLA
				} else {
					switch hex.EncodeToString(r) {
					case "01":
						singleRow[key] = "true"
					case "00":
						singleRow[key] = "false"
					default:
						singleRow[key] = string(r)
					}
				}
			case time.Time:
				singleRow[key] = r.Format("2006-01-02 15:04:05")
			case nil:
				singleRow[key] = r
			}
		}

		rowList = append(rowList, singleRow)
	}

	queryRes.SQL = sql
	queryRes.IsExecute = true
	queryRes.Err = nil
	queryRes.Total = len(rowList)
	queryRes.Columns = removeDuplicateElement(cols)
	queryRes.Rows = rowList

	return queryRes
}

func (m *MySQLEngine) InitialDriver(conn common.ConnConfig, schema string) error {
	dsn := mysqlDSN(conn.IP, conn.Port, conn.UserName, conn.Password, schema)
	cli, err := newMySQLClient(dsn)
	if err != nil {
		return err
	}

	m.driver = cli
	m.ConnConfig = conn
	return nil
}

func (m *MySQLEngine) Reset() {
	err := m.Close()
	if err != nil {
		fmt.Printf("mysql engine close failed:%s\n", err)
	}
	m.driver = nil
	m.EngineBase = &common.EngineBase{}
}

func NewMySQLEngine() *MySQLEngine {
	return &MySQLEngine{
		nil,
		&common.EngineBase{},
	}
}

func mysqlDSN(ip string, port int, username string, password string, database string) string {
	connConfig := &mmysql.Config{
		User:   username,
		Passwd: password,
		Addr:   fmt.Sprintf("%s:%d", ip, port),
		Net:    "tcp",
		DBName: database,
		Params: map[string]string{
			"charset": "utf8mb4",
		},
		Loc:                  time.Local,
		AllowNativePasswords: true,
		ParseTime:            true,
		Timeout:              15 * time.Second,
	}

	return connConfig.FormatDSN()
}

func newMySQLClient(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(mydriver.New(mydriver.Config{
		DSN:                       dsn,
		SkipInitializeWithVersion: false,
	}), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func ForkMySQLEngine(conn common.ConnConfig, schema string) (*MySQLEngine, error) {
	dsn := mysqlDSN(conn.IP, conn.Port, conn.UserName, conn.Password, schema)
	cli, err := newMySQLClient(dsn)
	if err != nil {
		return nil, err
	}

	return &MySQLEngine{
		cli,
		common.NewEngineBase(conn),
	}, nil
}

func MySQLSQLType(sql string) (common.SQLType, error) {
	sqlSt, err := vsqlparser.Parse(sql)
	if err != nil {
		return common.StmtUnknown, err
	}

	statementType := vsqlparser.ASTToStatementType(sqlSt)
	return common.SQLType(statementType), nil
}

// MySQLTableFromSQL
// parse table name from select sql
// because table name provide by user is not real table that execute
func MySQLTableFromSelectSQL(sql string) (string, error) {
	tb, err := vsqlparser.TableFromStatement(sql)
	if err != nil {
		return "", err
	}

	return tb.Name.String(), nil
}

func MySQLPreCheck(sql string, allowSQLType []common.SQLType) (string, bool, error) {
	// valid sql is non empty
	if sql == "" {
		return "", false, inerr.ErrSQLEmpty
	}

	// parse sql type to preCheck
	sqlSt, err := vsqlparser.Parse(sql)
	if err != nil {
		return sql, false, err
	}

	sqlType := vsqlparser.ASTToStatementType(sqlSt)
	// valid sql is allowed to execute by allowSQLType
	isPass := false
	for _, alType := range allowSQLType {
		if alType != common.SQLType(sqlType) {
			continue
		}

		isPass = true
		break
	}

	if !isPass {
		return sql, false, nil
	}

	// add default limit if sql is select statement and no set limit
	// avoid querySet is too big
	// and not avoid statement is slow query
	if sqlType == vsqlparser.StmtSelect {
		// sql is a select statement and check is container limit
		switch sqlStObj := sqlSt.(type) {
		case *vsqlparser.Select:
			// simple select
			if sqlStObj.Limit != nil {
				return sql, true, nil
			}

			// add limit
			// default limit 100
			sqlStObj.SetLimit(&vsqlparser.Limit{
				Offset: nil,
				Rowcount: &vsqlparser.Literal{
					Type: vsqlparser.IntVal,
					Val:  "100",
				},
			})

			return vsqlparser.String(sqlStObj), true, nil
		case *vsqlparser.Union:
			// union select
			if sqlStObj.Limit != nil {
				return sql, true, nil
			}

			sqlStObj.SetLimit(&vsqlparser.Limit{
				Offset: nil,
				Rowcount: &vsqlparser.Literal{
					Type: vsqlparser.IntVal,
					Val:  "100",
				},
			})

			return vsqlparser.String(sqlStObj), true, nil
		}
	}

	// other valid statement
	return sql, true, nil
}

func removeDuplicateElement(addrs []string) []string {
	result := make([]string, 0, len(addrs))
	temp := map[string]struct{}{}
	idx := 0
	for _, item := range addrs {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		} else {
			idx++
			item += fmt.Sprintf("(%v)", idx)
			result = append(result, item)
		}
	}
	return result
}

func mapScan(r *sql.Rows, dest map[string]interface{}) error {
	// ignore r.started, since we needn't use reflect for anything.
	columns, err := r.Columns()
	if err != nil {
		return err
	}

	values := make([]interface{}, len(columns))
	for i := range values {
		values[i] = new(interface{})
	}

	err = r.Scan(values...)
	if err != nil {
		return err
	}

	ele := removeDuplicateElement(columns)

	for i, column := range ele {
		dest[column] = *(values[i].(*interface{}))
	}

	return r.Err()
}

package engine

import (
	"fmt"
	"testing"

	"git.qpaas.com/go-components/webconsole/pkg/common"
	"github.com/stretchr/testify/require"
)

func TestMySQLPreCheck(t *testing.T) {
	validCase := []string{
		"use mock_schema",
		"show tables",
		"desc mock_table",
		"select * from mock_table",
		"select * from mock_table;",
		"   select * from mock_table;",
		" select * from mock_table where id =232 and name='sga' limit 20",
		`select id, name
		from mock_table`,
		`select 
		id, name
		from mock_table`,
		`select
     	   id, name
		from mock_table`,
		"explain select * from table01 where id=333",
		`explain 
		select * 
		from table01 where id=333`,
		`explain 
		  select * 
		from table01 where id=333`,
		`explain 
		 delete from  
		mock_table;`,
		`select * from table01 inner join table02 on tabl01.id = table02.id`,
		`   
		select * from table01 inner join table02 on tabl01.id = table02.id`,
		"update mock_table set id=1 where name=2",
	}

	invalidCase := []string{
		`explain; 
		 delete from  
		mock_table;`,
		`explain \n
		   delete from  
		mock_table;`,
		`select * from table01;
		 select * from table02;`,
		`select * from table01;
		   select * from table02;`,
		"select * from table01;\nselect * from table02;",
		`select * from table01;select * from table02;`,
		"select * from mock_table where id =232 and name='sga' limit 20;select * from mock_table",
		"  insert into (f1, f2) values(v1, v2)",
		"insert into (f1, f2) values(v1, v2)",
		"drop database 001",
		"delete from mock_table where id = 0",
		"truncate mock_table",
		"  insert into (f1, f2) values(v1, v2)",
		"drop table mock_table;",
		`drop 
		table mock_table;`,
		`explain drop 
		table mock_table;`,
		`  drop 
		 table 
		  mock_table;`,
		`select * from table01;
			drop table mock_table;`,
		`select * from table01;
			select * from table02;
			drop table table03`,
		`
		select * from table01;
		select * from table02;
		drop table table03;
		`,
		`
		select * from table01;\n
		select * from table02;\n
		drop table table03;\n
		`,
		`
		select * from table01
		select * from table02
		drop table table03
		`,
		"drop table table03;",
		"drop database 002",
	}

	t.Run("allow select sql type", func(t *testing.T) {
		allowSQLType := []common.SQLType{
			common.StmtSelect,
			common.StmtShow,
			common.StmtExplain,
		}
		for _, item := range validCase {
			sql, isvalid, err := MySQLPreCheck(item, allowSQLType)
			require.NoError(t, err, fmt.Sprintf("test sql is: %s", item))

			fmt.Printf("%s: %t\n", sql, isvalid)
		}
	})

	t.Run("allow update sql type", func(t *testing.T) {
		allowSQLType := []common.SQLType{
			common.StmtUpdate,
		}
		for _, item := range validCase {
			sql, isvalid, err := MySQLPreCheck(item, allowSQLType)
			require.NoError(t, err, fmt.Sprintf("test sql is: %s", item))

			fmt.Printf("%s: %t\n", sql, isvalid)
		}
	})

	t.Run("invalid sql preCheck", func(t *testing.T) {
		allowSQLType := []common.SQLType{
			common.StmtUpdate,
		}
		for _, item := range invalidCase {
			_, isvalid, err := MySQLPreCheck(item, allowSQLType)
			require.Error(t, err, fmt.Sprintf("test sql is: %s", item))
			fmt.Printf("%s: %t: %s\n", item, isvalid, err)
		}
	})
}

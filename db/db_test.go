package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"reflect"
	"testing"
	"time"
)

/**
* golang-sample源代码，版权归锦翰科技（深圳）有限公司所有。
* <p>
* 文件名称 : db_test.go
* 文件路径 :
* 作者 : DavidLiu
× Email: david.liu@ginghan.com
*
* 创建日期 : 2022/5/11 17:58
* 修改历史 : 1. [2022/5/11 17:58] 创建文件 by LongYong
*/
var Db *sql.DB

func init() {
	db, err := sql.Open("mysql", "root:123456@tcp(192.168.56.101:3306)/demo")

	if err != nil {
		panic("Open database error:" + err.Error())
	}

	Db = db
	Db.SetConnMaxLifetime(10 * time.Minute)
	Db.SetMaxOpenConns(15)
	Db.SetConnMaxIdleTime(10 * time.Minute)
	Db.SetMaxIdleConns(9)
	fmt.Sprintf("Database : %+v\n", Db.Stats())

}

func TestSQLDB(t *testing.T) {
	rs, err := QueryResultSet(Db, "select * from sys_menu")

	if err != nil {
		fmt.Printf("Query error: %v\n", err)
	}

	fmt.Printf("%+v %+v \n", rs.GetMeta().GetColumns(), rs.GetMeta().GetColumnTypes())

	for rs.Next() {
		fmt.Printf("%v %v %v %v %v %v \n", ptr2Data(rs.GetTimeByName("create_time", "2006-01-02 15:04:05")), ptr2Data(rs.GetIntByName("visible")), ptr2Data(rs.GetInt(1)), ptr2Data(rs.GetBoolByName("status")),
			ptr2Data(rs.GetStringByName("perms")), ptr2Data(rs.GetStringByName("menu_name")))
	}

	fmt.Println("END")
	time.Sleep(time.Second)
}

func ptr2Data(d any) any {
	if d == nil {
		return nil
	}

	v := reflect.ValueOf(d)

	if v.Kind() == reflect.Pointer {
		return v.Elem().Interface()
	}

	return d
}

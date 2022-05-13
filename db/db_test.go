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
var dbPlus DBPlus

func init() {
	db, err := sql.Open("mysql", "root:123456@tcp(192.168.56.101:3306)/demo")

	if err != nil {
		panic("Open database error:" + err.Error())
	}

	db.SetConnMaxLifetime(10 * time.Minute)
	db.SetMaxOpenConns(15)
	db.SetConnMaxIdleTime(10 * time.Minute)
	db.SetMaxIdleConns(9)

	fmt.Sprintf("Database : %+v\n", db.Stats())

	dbPlus = DBPlus{DB: db}

}

func TestSQLDB(t *testing.T) {
	rs, err := dbPlus.QueryResultSet("select * from sys_menu limit 0,3")

	if err != nil {
		panic("Query error: " + err.Error())
	}

	fmt.Printf("%+v %+v \n", rs.GetMeta().GetColumns(), rs.GetMeta().GetColumnTypes())
	fmt.Printf("%+v \n", rs.Length())

	for idx := 0; idx < rs.Length(); idx++ {
		fmt.Printf("%v %v %v %v %v %v \n",
			ptr2Data(rs.GetTimeByName(idx, "create_time", "2006-01-02 15:04:05")),
			ptr2Data(rs.GetIntByName(idx, "visible")),
			ptr2Data(rs.GetInt(idx, 1)),
			ptr2Data(rs.GetBoolByName(idx, "status")),
			ptr2Data(rs.GetStringByName(idx, "perms")),
			ptr2Data(rs.GetStringByName(idx, "menu_name")))
	}

	fmt.Println("END")
	time.Sleep(time.Second)
}

func ptr2Data(d any, err error) any {

	if err != nil {
		panic(err)
	}

	if d == nil {
		return nil
	}

	v := reflect.ValueOf(d)

	if v.Kind() == reflect.Pointer {
		return v.Elem().Interface()
	}

	return d
}

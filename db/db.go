package db

import (
	"database/sql"
	"errors"
	util4go "github.com/gohutool/boot4go-util"
	"reflect"
	"time"
)

/**
* golang-sample源代码，版权归锦翰科技（深圳）有限公司所有。
* <p>
* 文件名称 : db.go
* 文件路径 :
* 作者 : DavidLiu
× Email: david.liu@ginghan.com
*
* 创建日期 : 2022/5/11 16:21
* 修改历史 : 1. [2022/5/11 16:21] 创建文件 by LongYong
*/

// ColumnType contains the name and type of a column.
type ColumnType struct {
	Name            string
	NullSupportable bool
	Nullable        bool
	LengthVariable  bool
	Length          int64

	PrecisionScale bool

	Precision int64
	Scale     int64

	SqlType  string
	DataType reflect.Type
}

func QueryResultSet(db *sql.DB, query string, args ...any) (*resultSet, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	return ResultSet(rows)
}

func QueryCount(db *sql.DB, query string, args ...any) int64 {
	rs, err := QueryResultSet(db, query, args...)

	if err != nil {
		return 0
	}

	if rs == nil {
		return 0
	}

	defer rs.Close()

	if !rs.Next() {
		return 0
	}

	c := rs.GetInt(1)

	if c == nil {
		return 0
	}

	return *c
}

func ResultSet(rows *sql.Rows) (*resultSet, error) {
	column, err := rows.Columns()

	if err != nil {
		return nil, err
	}

	columnType, err := rows.ColumnTypes()

	if err != nil {
		return nil, err
	}

	meta := &ResultSetMetaData{}

	meta.columnName = make(map[string]int)

	meta.column = append(meta.column, column...)

	for idx, v := range columnType {

		ct := ColumnType{}
		ct.Nullable, ct.NullSupportable = v.Nullable()
		ct.Length, ct.LengthVariable = v.Length()
		ct.Precision, ct.Scale, ct.PrecisionScale = v.DecimalSize()
		ct.SqlType = v.DatabaseTypeName()
		ct.DataType = v.ScanType()
		ct.Name = v.Name()

		meta.columnType = append(meta.columnType, ct)

		meta.columnName[ct.Name] = idx + 1
	}

	return &resultSet{rows: rows, metaData: meta, columnCount: len(column)}, nil
}

type ResultSetMetaData struct {
	column     []string
	columnType []ColumnType
	columnName map[string]int
}

func (m ResultSetMetaData) GetColumns() []string {
	return m.column
}

func (m ResultSetMetaData) GetColumnTypes() []ColumnType {
	return m.columnType
}

func (m ResultSetMetaData) GetColumnCount() int {
	return len(m.column)
}

func (m ResultSetMetaData) GetColumnName(idx int) string {
	return m.column[idx]
}

func (m ResultSetMetaData) GetColumnType(idx int) ColumnType {
	return m.columnType[idx]
}

func (m ResultSetMetaData) GetPrecision(idx int) int64 {
	return m.columnType[idx].Precision
}

func (m ResultSetMetaData) GetScale(idx int) int64 {
	return m.columnType[idx].Scale
}

func (m ResultSetMetaData) IsNullable(idx int) bool {
	return m.columnType[idx].Nullable
}

func (m ResultSetMetaData) IsVariableLength(idx int) bool {
	return m.columnType[idx].LengthVariable
}

func (m ResultSetMetaData) GetColumnDataType(idx int) reflect.Type {
	return m.columnType[idx].DataType
}

func (m ResultSetMetaData) GetColumnSQLType(idx int) string {
	return m.columnType[idx].SqlType
}

func (m ResultSetMetaData) GetColumnLength(idx int) int64 {
	return m.columnType[idx].Length
}

type resultSet struct {
	rows        *sql.Rows
	metaData    *ResultSetMetaData
	_current    []any
	columnCount int
}

func (rs *resultSet) GetMeta() ResultSetMetaData {
	return *rs.metaData
}

func (rs *resultSet) Close() {
	rs.rows.Close()
}

func (rs *resultSet) Next() bool {
	rs._current = nil
	has := rs.rows.Next()

	if has {
		rs._current = make([]any, rs.columnCount)

		for i := 0; i < rs.columnCount; i++ {
			var one *string
			rs._current[i] = &one
		}
		if err := rs.rows.Scan(rs._current...); err != nil {
			panic(err)
		}
	}

	return has
}

func (rs *resultSet) GetInt(column int) *int64 {
	column = column - 1

	data := rs._current[column]

	if data == nil {
		return nil
	}

	var str string
	switch data.(type) {
	case *string:
		str = *data.(*string)
	case **string:
		if *data.(**string) == nil {
			return nil
		}

		str = **data.(**string)
	default:
		panic("can not convert type " + reflect.TypeOf(data).String())
	}

	var rtn = new(int64)
	v, err := util4go.GetInt(str)

	if err != nil {
		return nil
	}

	*rtn = v

	return rtn
}

func (rs *resultSet) GetIntByName(columnName string) *int64 {

	column, ok := rs.metaData.columnName[columnName]

	if !ok {
		return nil
	}

	return rs.GetInt(column)
}

func (rs *resultSet) GetFloat(column int) *float64 {
	column = column - 1

	data := rs._current[column]

	if data == nil {
		return nil
	}

	var str string
	switch data.(type) {
	case *string:
		str = *data.(*string)
	case **string:
		if *data.(**string) == nil {
			return nil
		}

		str = **data.(**string)
	default:
		panic("can not convert type " + reflect.TypeOf(data).String())
	}

	var rtn = new(float64)
	v, err := util4go.GetFloat(str)

	if err != nil {
		return nil
	}

	*rtn = v

	return rtn
}

func (rs *resultSet) GetFloatByName(columnName string) *float64 {

	column, ok := rs.metaData.columnName[columnName]

	if !ok {
		return nil
	}

	return rs.GetFloat(column)
}

func (rs *resultSet) GetString(column int) *string {
	column = column - 1

	data := rs._current[column]

	if data == nil {
		return nil
	}

	var str string
	switch data.(type) {
	case *string:
		str = *data.(*string)
	case **string:
		if *data.(**string) == nil {
			return nil
		}

		str = **data.(**string)
	default:
		panic("can not convert type " + reflect.TypeOf(data).String())
	}

	var rtn = new(string)
	*rtn = str
	return rtn
}

func (rs *resultSet) GetStringByName(columnName string) *string {

	column, ok := rs.metaData.columnName[columnName]

	if !ok {
		return nil
	}

	return rs.GetString(column)
}

func (rs *resultSet) GetTime(column int, layout string) *time.Time {
	str := rs.GetString(column)

	if str == nil {
		return nil
	}

	t, err := time.Parse(layout, *str)

	if err != nil {
		panic(err)
	}

	return &t
}

func (rs *resultSet) GetTimeByName(columnName string, layout string) *time.Time {

	column, ok := rs.metaData.columnName[columnName]

	if !ok {
		return nil
	}

	return rs.GetTime(column, layout)
}

func (rs *resultSet) GetBool(column int) *bool {
	column = column - 1

	data := rs._current[column]

	if data == nil {
		return nil
	}

	var str string
	switch data.(type) {
	case *string:
		str = *data.(*string)
	case **string:
		if *data.(**string) == nil {
			return nil
		}

		str = **data.(**string)
	default:
		panic("can not convert type " + reflect.TypeOf(data).String())
	}

	var rtn = new(bool)
	v, err := util4go.GetBool(str)

	if err != nil {
		return nil
	}

	*rtn = v

	return rtn
}

func (rs *resultSet) GetBoolByName(columnName string) *bool {

	column, ok := rs.metaData.columnName[columnName]

	if !ok {
		return nil
	}

	return rs.GetBool(column)
}

func (rs *resultSet) GetUint(column int) *uint64 {
	column = column - 1

	data := rs._current[column]

	if data == nil {
		return nil
	}

	var str string
	switch data.(type) {
	case *string:
		str = *data.(*string)
	case **string:
		if *data.(**string) == nil {
			return nil
		}

		str = **data.(**string)
	default:
		panic("can not convert type " + reflect.TypeOf(data).String())
	}

	var rtn = new(uint64)
	v, err := util4go.GetUint(str)

	if err != nil {
		return nil
	}

	*rtn = v

	return rtn
}

func (rs *resultSet) GetUintByName(columnName string) *uint64 {

	column, ok := rs.metaData.columnName[columnName]

	if !ok {
		return nil
	}

	return rs.GetUint(column)
}

func (rs *resultSet) GetBytes(column int) []byte {

	str := rs.GetString(column)

	if str == nil {
		return nil
	}

	return []byte(*str)
}

func (rs *resultSet) GetBytesByName(columnName string) []byte {

	column, ok := rs.metaData.columnName[columnName]

	if !ok {
		return nil
	}

	return rs.GetBytes(column)
}

func (rs *resultSet) Get(column int) any {
	var data any

	ds := make([]any, column)
	ds[column-1] = &data

	if err := rs.rows.Scan(ds...); err != nil {
		panic(err)
	}

	return data
}

func (rs *resultSet) GetByName(columnName string) any {

	column, ok := rs.metaData.columnName[columnName]

	if !ok {
		return nil
	}

	return rs.Get(column)
}

func (rs *resultSet) ScanWithColumn(column int, data *any) error {

	ds := make([]any, column)
	ds[column-1] = data

	if err := rs.rows.Scan(ds...); err != nil {
		return err
	}

	return nil
}

func (rs *resultSet) ScanWithColumnName(columnName string, data *any) error {

	column, ok := rs.metaData.columnName[columnName]

	if !ok {
		return errors.New("Null")
	}

	return rs.ScanWithColumn(column, data)
}

func (rs *resultSet) ScanWith(data []any) error {
	return rs.rows.Scan(data)
}

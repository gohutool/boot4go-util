package db

import (
	"database/sql"
	util4go "github.com/gohutool/boot4go-util"
	"golang.org/x/net/context"
	"reflect"
	"time"
)

/**
* golang-sample源代码，版权归锦翰科技（深圳）有限公司所有。
* <p>
* 文件名称 : d.DB.go
* 文件路径 :
* 作者 : DavidLiu
× Email: david.liu@ginghan.com
*
* 创建日期 : 2022/5/11 16:21
* 修改历史 : 1. [2022/5/11 16:21] 创建文件 by LongYong
*/

// ColumnType contains the name and type of a column.

var ArrayOfIndexError arrayOfIndexError
var UnknownColumnError unknownCloumnError
var ParseColumnError parseColumnError
var TimeParseColumnError timeParseColumnError

type arrayOfIndexError struct {
}

func (err arrayOfIndexError) Error() string {
	return "ArrayOfIndexError"
}

type unknownCloumnError struct {
}

func (err unknownCloumnError) Error() string {
	return "ArrayOfIndexError"
}

type parseColumnError struct {
}

func (err parseColumnError) Error() string {
	return "ParseColumnError"
}

type timeParseColumnError struct {
}

func (err timeParseColumnError) Error() string {
	return "ParseColumnError"
}

type DBPlus struct {
	DB             *sql.DB
	QUERY_TIMEOUT  int
	UPDATE_TIMEOUT int
	BULK_TIMEOUT   int
}

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

const (
	DEFAULT_QUERY_TIMEOUT = 0
	UPDATE_QUERY_TIMEOUT  = 0
	BULK_QUERY_TIMEOUT    = 0
)

func (d *DBPlus) GetDB() *sql.DB {
	return d.DB
}

func (d *DBPlus) QueryResultSet(query string, args ...any) (ResultSet, error) {
	return d.QueryResultSetWithTimeout(query, d.QUERY_TIMEOUT, args...)
}

func (d *DBPlus) Exec(update string, args ...any) (int64, int64, error) {
	return d.ExecWithTimeout(update, d.UPDATE_TIMEOUT, args...)
}

func (d *DBPlus) Bulk(update string, args [][]any) ([]int64, []int64, error) {
	return d.BulkWithTimeout(update, d.BULK_TIMEOUT, args)
}

func (d *DBPlus) BulkWithTimeout(update string, timeout int, args [][]any) ([]int64, []int64, error) {
	if timeout > 0 {
		ctx, cancel := context.WithTimeout(context.TODO(), time.Duration(timeout)*time.Second)
		defer cancel()

		stmt, err := d.DB.PrepareContext(ctx, update)

		if err != nil {
			return nil, nil, err
		}

		defer stmt.Close()

		var affects, lastInsertIds []int64

		for _, arg := range args {
			rows, err2 := stmt.Exec(arg...)

			if err2 != nil {
				return nil, nil, err2
			}

			affected, _ := rows.RowsAffected()
			lastInsertId, _ := rows.LastInsertId()

			affects = append(affects, affected)
			lastInsertIds = append(lastInsertIds, lastInsertId)
		}

		return affects, lastInsertIds, nil
	} else {
		stmt, err := d.DB.Prepare(update)

		if err != nil {
			return nil, nil, err
		}

		defer stmt.Close()

		var affects, lastInsertIds []int64

		for _, arg := range args {
			rows, err2 := stmt.Exec(arg...)

			if err2 != nil {
				return nil, nil, err2
			}

			affected, _ := rows.RowsAffected()
			lastInsertId, _ := rows.LastInsertId()

			affects = append(affects, affected)
			lastInsertIds = append(lastInsertIds, lastInsertId)
		}

		return affects, lastInsertIds, nil
	}
}

func (d *DBPlus) ExecWithTimeout(update string, timeout int, args ...any) (int64, int64, error) {
	if timeout > 0 {
		ctx, cancel := context.WithTimeout(context.TODO(), time.Duration(timeout)*time.Second)
		defer cancel()

		stmt, err := d.DB.PrepareContext(ctx, update)

		if err != nil {
			return 0, 0, err
		}

		defer stmt.Close()

		rows, err2 := stmt.Exec(args...)

		if err2 != nil {
			return 0, 0, err2
		}

		affected, _ := rows.RowsAffected()
		lastInsertId, _ := rows.LastInsertId()

		return affected, lastInsertId, nil
	} else {
		stmt, err := d.DB.Prepare(update)

		if err != nil {
			return 0, 0, err
		}

		defer stmt.Close()

		rows, err2 := stmt.Exec(args...)

		if err2 != nil {
			return 0, 0, err2
		}

		affected, _ := rows.RowsAffected()
		lastInsertId, _ := rows.LastInsertId()

		return affected, lastInsertId, nil
	}
}

func (d *DBPlus) QueryResultSetWithTimeout(query string, timeout int, args ...any) (ResultSet, error) {

	if timeout > 0 {
		ctx, cancel := context.WithTimeout(context.TODO(), time.Duration(timeout)*time.Second)
		defer cancel()

		stmt, err := d.DB.PrepareContext(ctx, query)

		if err != nil {
			return EmptyResultSet, err
		}

		defer stmt.Close()

		rows, err2 := stmt.Query(args...)

		if err2 != nil {
			return EmptyResultSet, err2
		}

		return New(rows)
	} else {
		stmt, err := d.DB.Prepare(query)

		if err != nil {
			return EmptyResultSet, err
		}

		defer stmt.Close()

		rows, err2 := stmt.Query(args...)

		if err2 != nil {
			return EmptyResultSet, err2
		}

		return New(rows)
	}
}

func (d *DBPlus) Begin() (*sql.Tx, error) {
	return d.DB.Begin()
}

func (d *DBPlus) Commit(tx *sql.Tx) error {
	return tx.Commit()
}

func (d *DBPlus) Rollback(tx *sql.Tx) error {
	return tx.Rollback()
}

func (d *DBPlus) QueryWithTx(tx *sql.Tx, sql string, args ...any) (ResultSet, error) {
	stmt, err := tx.Prepare(sql)

	if err != nil {
		return EmptyResultSet, err
	}

	defer stmt.Close()

	rows, err2 := stmt.Query(args...)

	if err2 != nil {
		return EmptyResultSet, err2
	}

	return New(rows)
}

func (d *DBPlus) BulkWithTx(tx *sql.Tx, update string, args [][]any) ([]int64, []int64, error) {
	stmt, err := tx.Prepare(update)

	if err != nil {
		return nil, nil, err
	}

	defer stmt.Close()

	var affects, lastInsertIds []int64

	for _, arg := range args {
		rows, err2 := stmt.Exec(arg...)

		if err2 != nil {
			return nil, nil, err2
		}

		affected, _ := rows.RowsAffected()
		lastInsertId, _ := rows.LastInsertId()

		affects = append(affects, affected)
		lastInsertIds = append(lastInsertIds, lastInsertId)
	}

	return affects, lastInsertIds, nil
}

func (d *DBPlus) ExecWithTx(tx *sql.Tx, sql string, args ...any) (int64, int64, error) {
	stmt, err := tx.Prepare(sql)
	if err != nil {
		return 0, 0, err
	}
	defer stmt.Close()

	rows, err2 := stmt.Exec(args...)

	if err2 != nil {
		return 0, 0, err2
	}

	affected, _ := rows.RowsAffected()
	lastInsertId, _ := rows.LastInsertId()

	return affected, lastInsertId, nil
}

func (d *DBPlus) QueryBool(query string, args ...any) (*bool, error) {
	rs, err := d.QueryResultSet(query, args...)

	if err != nil {
		return nil, err
	}

	if rs.Length() == 0 {
		return nil, nil
	}

	c, err := rs.GetBool(0, 1)

	if err != nil {
		return nil, err
	}

	if c == nil {
		return nil, nil
	}

	return c, nil
}

func (d *DBPlus) QueryFloat(query string, args ...any) (*float64, error) {
	rs, err := d.QueryResultSet(query, args...)

	if err != nil {
		return nil, err
	}

	if rs.Length() == 0 {
		return nil, nil
	}

	c, err := rs.GetFloat(0, 1)

	if err != nil {
		return nil, err
	}

	if c == nil {
		return nil, nil
	}

	return c, nil
}

func (d *DBPlus) QueryInt(query string, args ...any) (*int64, error) {
	rs, err := d.QueryResultSet(query, args...)

	if err != nil {
		return nil, err
	}

	if rs.Length() == 0 {
		return nil, nil
	}

	c, err := rs.GetInt(0, 1)

	if err != nil {
		return nil, err
	}

	if c == nil {
		return nil, nil
	}

	return c, nil
}

func (d *DBPlus) QueryCount(query string, args ...any) (int64, error) {
	i, err := d.QueryInt(query, args...)

	if err != nil {
		return 0, err
	}

	if i == nil {
		return 0, nil
	}

	return *i, nil
}

func (d *DBPlus) QueryOne(query string, args ...any) (map[string]string, error) {
	rs, err := d.QueryResultSet(query, args...)

	if err != nil {
		return nil, err
	}

	if rs.Length() == 0 {
		return nil, nil
	}

	rtn := make(map[string]string)

	for idx, column := range rs.metaData.column {
		v := rs.data[0][idx]
		if v != nil {
			rtn[column] = v.(string)
		}

	}

	return rtn, nil
}

func (d *DBPlus) QueryString(query string, args ...any) (*string, error) {
	rs, err := d.QueryResultSet(query, args...)

	if err != nil {
		return nil, err
	}

	if rs.Length() == 0 {
		return nil, ArrayOfIndexError
	}

	return rs.GetString(0, 1)
}

func New(rows *sql.Rows) (ResultSet, error) {
	defer rows.Close()

	column, err := rows.Columns()

	if err != nil {
		return EmptyResultSet, err
	}

	columnType, err := rows.ColumnTypes()

	if err != nil {
		return EmptyResultSet, err
	}

	meta := ResultSetMetaData{}

	meta.columnName = make(map[string]int)

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
		meta.column = append(meta.column, column[idx])
	}

	columnSize := len(meta.column)
	var datas [][]any

	for rows.Next() {
		current := make([]any, columnSize)

		for i := 0; i < columnSize; i++ {
			var one *string
			current[i] = &one
		}

		if err := rows.Scan(current...); err != nil {
			return EmptyResultSet, err
		}

		datas = append(datas, current)
	}

	return ResultSet{data: datas, metaData: meta, columnCount: len(column)}, nil
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

var EmptyResultSet = ResultSet{}

type ResultSet struct {
	metaData    ResultSetMetaData
	data        [][]any
	columnCount int
}

func (rs *ResultSet) Length() int {
	return len(rs.data)
}

func (rs *ResultSet) GetMeta() ResultSetMetaData {
	return rs.metaData
}

func (rs *ResultSet) GetStringFromRaw(data any) (*string, error) {
	var str string

	switch data.(type) {
	case *string:
		str = *data.(*string)
	case **string:
		if *data.(**string) == nil {
			return nil, nil
		}

		str = **data.(**string)
	default:
		return nil, ParseColumnError
	}

	return &str, nil
}

func (rs *ResultSet) GetInt(idx int, column int) (*int64, error) {
	if idx >= rs.Length() || idx < 0 {
		return nil, ArrayOfIndexError
	}

	column = column - 1

	data := rs.data[idx][column]

	if data == nil {
		return nil, nil
	}

	str, err := rs.GetStringFromRaw(data)

	if err != nil {
		return nil, err
	}

	var rtn = new(int64)
	v, err := util4go.GetInt(*str)

	if err != nil {
		return nil, err
	}

	*rtn = v

	return rtn, nil
}

func (rs *ResultSet) GetIntByName(idx int, columnName string) (*int64, error) {
	column, ok := rs.metaData.columnName[columnName]

	if !ok {
		return nil, UnknownColumnError
	}

	return rs.GetInt(idx, column)
}

func (rs *ResultSet) GetFloat(idx int, column int) (*float64, error) {
	if idx >= rs.Length() || idx < 0 {
		return nil, ArrayOfIndexError
	}

	column = column - 1

	data := rs.data[idx][column]

	if data == nil {
		return nil, nil
	}

	str, err := rs.GetStringFromRaw(data)

	if err != nil {
		return nil, err
	}

	var rtn = new(float64)
	v, err := util4go.GetFloat(*str)

	if err != nil {
		return nil, err
	}

	*rtn = v

	return rtn, nil
}

func (rs *ResultSet) GetFloatByName(idx int, columnName string) (*float64, error) {

	column, ok := rs.metaData.columnName[columnName]

	if !ok {
		return nil, UnknownColumnError
	}

	return rs.GetFloat(idx, column)
}

func (rs *ResultSet) GetString(idx, column int) (*string, error) {

	if idx >= rs.Length() || idx < 0 {
		return nil, ArrayOfIndexError
	}

	column = column - 1

	data := rs.data[idx][column]

	if data == nil {
		return nil, nil
	}

	return rs.GetStringFromRaw(data)
}

func (rs *ResultSet) GetStringByName(idx int, columnName string) (*string, error) {

	column, ok := rs.metaData.columnName[columnName]

	if !ok {
		return nil, UnknownColumnError
	}

	return rs.GetString(idx, column)
}

func (rs *ResultSet) GetTime(idx, column int, layout string) (*time.Time, error) {
	str, err := rs.GetString(idx, column)

	if err != nil {
		return nil, err
	}

	if str == nil {
		return nil, nil
	}

	t, err := time.Parse(layout, *str)

	if err != nil {
		return nil, TimeParseColumnError
	}

	return &t, nil
}

func (rs *ResultSet) GetTimeByName(idx int, columnName string, layout string) (*time.Time, error) {

	column, ok := rs.metaData.columnName[columnName]

	if !ok {
		return nil, UnknownColumnError
	}

	return rs.GetTime(idx, column, layout)
}

func (rs *ResultSet) GetBool(idx, column int) (*bool, error) {
	if idx >= rs.Length() || idx < 0 {
		return nil, ArrayOfIndexError
	}

	column = column - 1

	data := rs.data[idx][column]

	if data == nil {
		return nil, nil
	}

	str, err := rs.GetStringFromRaw(data)

	if err != nil {
		return nil, err
	}

	var rtn = new(bool)
	v, err := util4go.GetBool(*str)

	if err != nil {
		return nil, err
	}

	*rtn = v

	return rtn, nil
}

func (rs *ResultSet) GetBoolByName(idx int, columnName string) (*bool, error) {

	column, ok := rs.metaData.columnName[columnName]

	if !ok {
		return nil, unknownCloumnError{}
	}

	return rs.GetBool(idx, column)
}

func (rs *ResultSet) GetUint(idx, column int) (*uint64, error) {
	if idx >= rs.Length() || idx < 0 {
		return nil, ArrayOfIndexError
	}

	column = column - 1

	data := rs.data[idx][column]

	if data == nil {
		return nil, nil
	}

	str, err := rs.GetStringFromRaw(data)

	if err != nil {
		return nil, err
	}

	var rtn = new(uint64)
	v, err := util4go.GetUint(*str)

	if err != nil {
		return nil, ParseColumnError
	}

	*rtn = v

	return rtn, nil
}

func (rs *ResultSet) GetUintByName(idx int, columnName string) (*uint64, error) {

	column, ok := rs.metaData.columnName[columnName]

	if !ok {
		return nil, UnknownColumnError
	}

	return rs.GetUint(idx, column)
}

func (rs *ResultSet) GetBytes(idx, column int) ([]byte, error) {
	if idx >= rs.Length() || idx < 0 {
		return nil, ArrayOfIndexError
	}

	str, err := rs.GetString(idx, column)

	if err != nil {
		return nil, err
	}

	if str == nil {
		return nil, nil
	}

	return []byte(*str), nil
}

func (rs *ResultSet) GetBytesByName(idx int, columnName string) ([]byte, error) {

	column, ok := rs.metaData.columnName[columnName]

	if !ok {
		return nil, UnknownColumnError
	}

	return rs.GetBytes(idx, column)
}

func (rs *ResultSet) Get(idx, column int) ([]any, error) {
	if idx >= rs.Length() || idx < 0 {
		return nil, ArrayOfIndexError
	}

	return rs.data[idx], nil
}

func (rs *ResultSet) GetByName(idx int, columnName string) ([]any, error) {

	column, ok := rs.metaData.columnName[columnName]

	if !ok {
		return nil, UnknownColumnError
	}

	return rs.Get(idx, column)
}

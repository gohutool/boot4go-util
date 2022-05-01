package util4go

import (
	"errors"
	"reflect"
)

/**
* golang-sample源代码，版权归锦翰科技（深圳）有限公司所有。
* <p>
* 文件名称 : reflect.go
* 文件路径 :
* 作者 : DavidLiu
× Email: david.liu@ginghan.com
*
* 创建日期 : 2022/5/1 09:51
* 修改历史 : 1. [2022/5/1 09:51] 创建文件 by LongYong
*/

func Type2Str(t reflect.Type) (string, error) {
	if t.Kind() == reflect.Struct || t.Kind() == reflect.Interface {
		return t.String(), nil
	} else if t.Kind() == reflect.Ptr {
		return t.Elem().String(), nil
	}

	return t.String(), errors.New(t.String() + " is not struct or interface")
}

func TypeOf(obj any) reflect.Type {
	return reflect.TypeOf(obj)
}

func ElmType(value reflect.Value) reflect.Type {
	if value.Kind() == reflect.Pointer {
		return value.Elem().Type()
	}

	return value.Type()
}

func NewInstanceValue(t reflect.Type) reflect.Value {
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	rtn := reflect.New(t)
	return rtn
}

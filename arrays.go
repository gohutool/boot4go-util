package util4go

import "strconv"

/**
* golang-sample源代码，版权归锦翰科技（深圳）有限公司所有。
* <p>
* 文件名称 : arrays.go
* 文件路径 :
* 作者 : DavidLiu
× Email: david.liu@ginghan.com
*
* 创建日期 : 2022/5/1 09:57
* 修改历史 : 1. [2022/5/1 09:57] 创建文件 by LongYong
*/

func Reduce[T any, R any](source []T, fn func(one T) (R, bool)) []R {
	if source == nil {
		return nil
	}

	rtn := make([]R, 0, len(source))

	for _, o := range source {
		if v, remain := fn(o); remain {
			rtn = append(rtn, v)
		}
	}

	return rtn
}

func CopyArray[T any](src []any, dest []T) []T {
	if src == nil {
		return nil
	}

	if dest == nil {
		dest = make([]T, 0, len(src))
	}

	for _, v := range src {
		dest = append(dest, v.(T))
	}

	return dest
}

func InsertAt[T any](list []T, idx int, t T) []T {
	l := len(list)

	if idx < 0 {
		idx = l + idx
	}

	if idx >= l || idx < 0 {
		panic("Out of index " + strconv.Itoa(idx))
	}

	var rtn []T
	rtn = append(rtn, list[0:idx]...)
	rtn = append(rtn, t)
	if idx < l+1 {
		rtn = append(rtn, list[idx:]...)
	}

	return rtn
}

func RemoveAt[T any](list []T, idx int) []T {
	l := len(list)

	if idx < 0 {
		idx = l + idx
	}

	if idx >= l || idx < 0 {
		panic("Out of index " + strconv.Itoa(idx))
	}

	var rtn []T
	rtn = append(rtn, list[0:idx]...)
	if idx < l {
		rtn = append(rtn, list[idx+1:]...)
	}

	return rtn
}

func ReplaceAt[T any](list []T, idx int, t T) []T {
	l := len(list)

	if idx < 0 {
		idx = l + idx
	}

	if idx >= l || idx < 0 {
		panic("Out of index " + strconv.Itoa(idx))
	}

	var rtn []T
	rtn = append(rtn, list[0:idx]...)
	rtn = append(rtn, t)
	if idx < l {
		rtn = append(rtn, list[idx+1:]...)
	}

	return rtn
}

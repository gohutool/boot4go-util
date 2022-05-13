package util4go

import (
	"strconv"
)

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

type Keyable interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | uintptr | float32 | float64 | complex64 | complex128 | string
}

type Numeric interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | uintptr | float32 | float64 | complex64 | complex128
}

type DeString interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64 | string
}

func Stream[T any](source []T, fn func(one T)) {
	for _, o := range source {
		fn(o)
	}
}

func GroupBy[T any, K Keyable](source []T, fn func(one T) K) map[K][]T {
	rtn := make(map[K][]T)

	Stream(source, func(one T) {
		key := fn(one)

		list, _ := rtn[key]
		list = append(list, one)
		rtn[key] = list
	})

	return rtn
}

func Collect[T any, V any](source []T, fn func(one T) V) []V {
	var rtn []V

	Stream(source, func(one T) {
		rtn = append(rtn, fn(one))
	})

	return rtn
}

func MapStream[T any, K Keyable](source map[K]T, fn func(K, T)) {
	for k, t := range source {
		fn(k, t)
	}
}

func Values[T any, K Keyable](source map[K]T) []T {
	var rtn []T
	MapStream(source, func(k K, t T) {
		rtn = append(rtn, t)
	})
	return rtn
}

func Map[T any, K Keyable](source []T, fn func(T) K) map[K]T {
	rtn := make(map[K]T)

	Stream(source, func(one T) {
		key := fn(one)
		rtn[key] = one
	})

	return rtn
}

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

func Reverse[T any](source []T) []T {
	var rtn []T

	for idx := len(source) - 1; idx >= 0; idx-- {
		rtn = append(rtn, source[idx])
	}

	return rtn
}

func GetMapValue[K Keyable, T any](m map[K]T, key K) T {
	if v, ok := m[key]; ok {
		return v
	}

	var v = new(T)

	return *v
}

func GetMapValue2[K Keyable, T any](m map[K]T, key K, defaultValue T) T {
	if v, ok := m[key]; ok {
		return v
	}

	return defaultValue
}

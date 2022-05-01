package util4go

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

func CopyArray[T struct{}](src []any, dest []T) []T {
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

package json

import (
	. "github.com/valyala/fastjson"
)

/**
* golang-sample源代码，版权归锦翰科技（深圳）有限公司所有。
* <p>
* 文件名称 : json.go
* 文件路径 :
* 作者 : DavidLiu
× Email: david.liu@ginghan.com
*
* 创建日期 : 2022/4/27 09:36
* 修改历史 : 1. [2022/4/27 09:36] 创建文件 by LongYong
*/

type Unmarshallable interface {
	Unmarshall(value Value) error
}

func Unmarshall(value Value, obj Unmarshallable) error {
	return obj.Unmarshall(value)
}

func UnmarshallJson(s string) (*Value, error) {
	var p Parser
	return p.Parse(s)
}

func UnmarshallObject(s string, obj Unmarshallable) error {
	if obj == nil {
		panic("Nil object can not unmarshall")
	}

	var p Parser
	if v, err := p.Parse(s); err != nil {
		return err
	} else {
		return Unmarshall(*v, obj)
	}
}

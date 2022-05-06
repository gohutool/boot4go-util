package util4go

import (
	"fmt"
	"testing"
)

/**
* golang-sample源代码，版权归锦翰科技（深圳）有限公司所有。
* <p>
* 文件名称 : reflect_test.go
* 文件路径 :
* 作者 : DavidLiu
× Email: david.liu@ginghan.com
*
* 创建日期 : 2022/5/6 21:13
* 修改历史 : 1. [2022/5/6 21:13] 创建文件 by LongYong
*/

type SampleObject struct {
	s string
	t string
}

func TestSetFieldValue(t *testing.T) {
	s := &SampleObject{s: "123", t: "abc"}

	if err := SetUnExportFieldValue(s, "s", "ginghan"); err != nil {
		fmt.Printf("%+v", err)
	} else {
		fmt.Printf("%+v", s)
	}

}

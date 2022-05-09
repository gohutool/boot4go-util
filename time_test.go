package util4go

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

/**
* golang-sample源代码，版权归锦翰科技（深圳）有限公司所有。
* <p>
* 文件名称 : time_test.go
* 文件路径 :
* 作者 : DavidLiu
× Email: david.liu@ginghan.com
*
* 创建日期 : 2022/5/9 15:28
* 修改历史 : 1. [2022/5/9 15:28] 创建文件 by LongYong
*/

func TestTime2Dimension(t *testing.T) {
	if d, err := Time2Dimension(time.Now()); err != nil {
		fmt.Printf("%+v\n", err)
	} else {
		fmt.Printf("%+v\n", d)
	}

}

func TestTime2Dimension2(t *testing.T) {
	time.Now().Format("200601021504")
	strconv.Atoi("200601021504")
	//strconv.Atoi(time.Now().Format("200601021504"))
	//strconv.Atoi(time.Now().Format("200601021504"))
}

func BenchmarkTime2Dimension2(b *testing.B) {

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			Time2Dimension(time.Now())
			//time.Now().Format("200601021504")
			//strconv.Atoi("200601021504")
			//strconv.Atoi(time.Now().Format("200601021504"))
			//strconv.Atoi(time.Now().Format("200601021504"))
		}
	})

}

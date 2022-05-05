package util4go

import (
	"fmt"
	"testing"
)

/**
* golang-sample源代码，版权归锦翰科技（深圳）有限公司所有。
* <p>
* 文件名称 : ip_test.go
* 文件路径 :
* 作者 : DavidLiu
× Email: david.liu@ginghan.com
*
* 创建日期 : 2022/5/2 16:23
* 修改历史 : 1. [2022/5/2 16:23] 创建文件 by LongYong
*/

func TestSyncMap(t *testing.T) {
	s := GuessIP("192.168.56.101:3306")

	fmt.Println(ReplaceIP(":9999", *s))

	fmt.Println(ReplacePort("localhost", "1000"))
}

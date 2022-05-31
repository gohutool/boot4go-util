package data

import (
	"testing"
)

/**
* golang-sample源代码，版权归锦翰科技（深圳）有限公司所有。
* <p>
* 文件名称 : loop_queue_test.go
* 文件路径 :
* 作者 : DavidLiu
× Email: david.liu@ginghan.com
*
* 创建日期 : 2022/5/31 22:21
* 修改历史 : 1. [2022/5/31 22:21] 创建文件 by LongYong
*/

func TestNewLoopQueue(t *testing.T) {
	size := 100
	q := NewLoopQueue(size, 1)

	t.Logf("IsEmpty : %v", q.IsEmpty())
	t.Logf("Len : %v", q.Len())
	t.Logf("Detach : %v", q.Detach())
}

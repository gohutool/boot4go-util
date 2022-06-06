package util4go

import (
	"fmt"
	"testing"
)

/**
* golang-sample源代码，版权归锦翰科技（深圳）有限公司所有。
* <p>
* 文件名称 : byte_test.go
* 文件路径 :
* 作者 : DavidLiu
× Email: david.liu@ginghan.com
*
* 创建日期 : 2022/5/6 15:17
* 修改历史 : 1. [2022/5/6 15:17] 创建文件 by LongYong
*/

func TestByte(t *testing.T) {
	var i1 int64 = 9223372036854775807
	var i2 int32 = 2147483647
	var i3 int16 = 32767
	var i4 int8 = 127
	var i5 uint64 = 18446744073709551615
	var i6 uint32 = 4294967295
	var i7 uint16 = 65535
	var i8 uint8 = 255

	printByte2Convert(i1)
	printByte2Convert(i2)
	printByte2Convert(i3)
	printByte2Convert(i4)
	printByte2Convert(i5)
	printByte2Convert(i6)
	printByte2Convert(i7)
	printByte2Convert(i8)
}

func TestByte2(t *testing.T) {

	var i8 uint8 = 255
	printByte2Convert(i8)
}

func printByte2Convert[T int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64](a T) {
	b := IntToBytes(a)
	d := BytesToInt[T](b)
	fmt.Printf("IntToBytes(%v)=%v  BytesToInt(%v, &c)=%v c=%v\n", a, b, b, d, d)
}

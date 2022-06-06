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
	printByte2IntConvert[int]([]byte{0, 12, 255, 255})
	printByte2IntConvert[int64]([]byte{11, 23, 255, 255, 255})
	printByte2IntConvert[int]([]byte{255, 255, 255, 255})
	printByte2IntConvert[byte]([]byte{255, 12})
}
func TestByte3(t *testing.T) {
	var v uint8 = 255

	fmt.Println(v)

	fmt.Println(int8(v))
}

func printByte2Convert[T int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64](a T) {
	b := IntToBytes(a)

	d := BytesToInt[T](b)

	fmt.Printf("IntToBytes(%v)=%v  BytesToInt(%v, &c)=%v c=%v\n", a, b, b, d, d)
}

func printByte2IntConvert[T int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64](b []byte) {
	d := BytesToInt[T](b)

	b2 := IntToBytes(d)

	fmt.Printf("IntToBytes(%v)=%v  BytesToInt(%v, &c)=%v c=%v\n", d, b2, b, d, d)
}

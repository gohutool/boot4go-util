package util4go

import (
	"bytes"
	"encoding/binary"
	"reflect"
	"unsafe"
)

/**
* golang-sample源代码，版权归锦翰科技（深圳）有限公司所有。
* <p>
* 文件名称 : byte.go
* 文件路径 :
* 作者 : DavidLiu
× Email: david.liu@ginghan.com
*
* 创建日期 : 2022/5/6 15:11
* 修改历史 : 1. [2022/5/6 15:11] 创建文件 by LongYong
*/

type Endian int8

const (
	LittleEndian Endian = iota
	BigEndian
)

// IntToBytes 将int类型的数转化为字节并以Big端存储
func IntToBytes[T int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64](intNum T) []byte {
	return IntToBytesEndian(intNum, BigEndian)
}

// BytesToInt 将以Big端存储的长为1/2字节的数转化成int类型的数
func BytesToInt[T int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64](bytesArr []byte, data *T) T {
	return BytesToIntEndian(bytesArr, data, BigEndian)
}

// IntToBytesEndian 将int类型的数转化为字节并以Big端存储
func IntToBytesEndian[T int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64](intNum T, endian Endian) []byte {
	buf := bytes.NewBuffer([]byte{})
	if endian == BigEndian {
		binary.Write(buf, binary.BigEndian, intNum)
	} else {
		binary.Write(buf, binary.LittleEndian, intNum)
	}

	return buf.Bytes()
}

// BytesToIntEndian 将以Big端存储的长为1/2字节的数转化成int类型的数

func BytesToIntEndian[T int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64](bytesArr []byte, data *T, endian Endian) T {
	buf := bytes.NewBuffer(bytesArr)
	if endian == BigEndian {
		binary.Read(buf, binary.BigEndian, data)
	} else {
		binary.Read(buf, binary.LittleEndian, data)
	}
	return T(*data)
}

// BytesToString converts byte slice to a string without memory allocation.
//
// Note it may break if the implementation of string or slice header changes in the future go versions.
func BytesToString(b []byte) string {
	/* #nosec G103 */
	return *(*string)(unsafe.Pointer(&b))
}

// StringToBytes converts string to a byte slice without memory allocation.
//
// Note it may break if the implementation of string or slice header changes in the future go versions.
func StringToBytes(s string) (b []byte) {
	/* #nosec G103 */
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	/* #nosec G103 */
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))

	bh.Data, bh.Len, bh.Cap = sh.Data, sh.Len, sh.Len
	return b
}

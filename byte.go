package util4go

import (
	"bytes"
	"encoding/binary"
	"errors"
	"reflect"
	"unicode/utf8"
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
func BytesToInt[T int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64](bytesArr []byte) T {
	return BytesToIntEndian[T](bytesArr, BigEndian)
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

func BytesToIntEndian[T int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64](bytesArr []byte, endian Endian) T {
	buf := bytes.NewBuffer(bytesArr)
	var data T
	if endian == BigEndian {
		binary.Read(buf, binary.BigEndian, &data)
	} else {
		binary.Read(buf, binary.LittleEndian, &data)
	}
	return data
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

var ArrayOverFlow = errors.New("ArrayOverFlow")

func DecodeUint32(b []byte, start int) (uint32, error) {
	if len(b)-start < 4 {
		return 0, ArrayOverFlow
	}
	_ = b[start+3] // bounds check hint to compiler; see golang.org/issue/14808
	return uint32(b[start+3]) | uint32(b[start+2])<<8 | uint32(b[start+1])<<16 | uint32(b[start+0])<<24, nil
}

func EncodeUint32(i uint32) []byte {
	return []byte{byte(i >> 24), byte(i >> 16), byte(i >> 8), byte(i)}
}

func IsUTF8(p []byte) bool {

	for {
		if len(p) == 0 {
			return true
		}
		ru, size := utf8.DecodeRune(p)
		if ru >= '\u0000' && ru <= '\u001f' {
			return false
		}
		if ru >= '\u007f' && ru <= '\u009f' {
			return false
		}
		if ru == utf8.RuneError {
			return false
		}
		if !utf8.ValidRune(ru) {
			return false
		}
		if size == 0 {
			return true
		}
		p = p[size:]
	}
}

func DecodeByte(b []byte, start int) (byte, error) {
	if len(b)-start < 1 {
		return 0, ArrayOverFlow
	}

	return b[start], nil
}

func DecodeUint16(b []byte, start int) (uint16, error) {
	if len(b)-start < 2 {
		return 0, ArrayOverFlow
	}

	_ = b[start+1] // bounds check hint to compiler; see golang.org/issue/14808
	return uint16(b[start+1]) | uint16(b[start+0])<<8, nil
}

func EncodeUint16(i uint16) []byte {
	return []byte{byte(i >> 8), byte(i)}
}

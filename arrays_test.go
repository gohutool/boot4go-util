package util4go

import (
	"fmt"
	"testing"
)

/**
* golang-sample源代码，版权归锦翰科技（深圳）有限公司所有。
* <p>
* 文件名称 : arrays_test.go
* 文件路径 :
* 作者 : DavidLiu
× Email: david.liu@ginghan.com
*
* 创建日期 : 2022/5/5 19:13
* 修改历史 : 1. [2022/5/5 19:13] 创建文件 by LongYong
*/

func TestArrayInsertAt(t *testing.T) {
	var str []rune = []rune("hello world")

	println(InsertAt(str, -1, 'T'))
	println(InsertAt(str, -2, 'R'))
	println(InsertAt(str, 10, 'W'))
	println(InsertAt(str, 9, 'W'))
	println(InsertAt(str, 0, 'T'))
	println(InsertAt(str, 1, 'T'))

}

func TestArrayReplaceAt(t *testing.T) {
	var str []rune = []rune("hello world")

	println(ReplaceAt(str, -1, 'T'))
	println(ReplaceAt(str, -2, 'R'))
	println(ReplaceAt(str, 10, 'T'))
	println(ReplaceAt(str, 9, 'R'))
	println(ReplaceAt(str, 0, 'T'))
	println(ReplaceAt(str, 1, 'R'))

}

func TestArrayRemoveAt(t *testing.T) {
	var str []rune = []rune("hello world")

	println(RemoveAt(str, -1))
	println(RemoveAt(str, -2))
	println(RemoveAt(str, 10))
	println(RemoveAt(str, 9))
	println(RemoveAt(str, 0))
	println(RemoveAt(str, 1))

}

func println(r []rune) {
	fmt.Println(string(r))
}

func TestStartWith(t *testing.T) {
	b := []byte("abcdefg")
	PrintSlice(b)

	s := []byte("abc")
	PrintSlice(s)

	fmt.Printf("%v\n", StartsWith(b, s))
	fmt.Printf("%v\n", StartsWith(s, b))

	e := []byte("efg")
	PrintSlice(e)

	fmt.Printf("%v\n", StartsWith(b, e))
	fmt.Printf("%v\n", EndsWith(b, e))
	fmt.Printf("%v\n", EndsWith(e, b))
}

type SampleSlice[T any] []T

func TestSample1(t *testing.T) {
	var list SampleSlice[int]
	list = make(SampleSlice[int], 2)

	list[0] = 1
	list[1] = 2

	fmt.Printf("%d\n", list[0])
	fmt.Printf("%d\n", list[1])
}

type SampleChan[T any] chan T

type SampleMap[K comparable, V any] map[K]V

func sample1() {

}

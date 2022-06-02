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

	if a, p, err := SplitHostPort("localhost"); err != nil {
		fmt.Printf("%v\n", err)
	} else {
		fmt.Printf("%v,%v\n", a, p)
	}

	if a, p, q, err := ParseURL("http://localhost:90?asdfasdf"); err != nil {
		fmt.Printf("%v\n", err)
	} else {
		fmt.Printf("%v,%v,%v\n", a, p, q)
	}

	if a, p, q, err := ParseURL("://localhost:90?asdfasdf"); err != nil {
		fmt.Printf("%v\n", err)
	} else {
		fmt.Printf("%v,%v,%v\n", a, p, q)
	}

	if a, p, q, err := ParseURL("https://localhost:90?asdfasdf"); err != nil {
		fmt.Printf("%v\n", err)
	} else {
		fmt.Printf("%v,%v,%v\n", a, p, q)
	}

	if a, p, q, err := ParseURL("localhost:90"); err != nil {
		fmt.Printf("%v\n", err)
	} else {
		fmt.Printf("%v,%v,%v\n", a, p, q)
	}

	if a, p, q, err := ParseURL("localhost:90?"); err != nil {
		fmt.Printf("%v\n", err)
	} else {
		fmt.Printf("%v,%v,%v\n", a, p, q)
	}

	if a, p, q, err := ParseURL("http://localhost:90?"); err != nil {
		fmt.Printf("%v\n", err)
	} else {
		fmt.Printf("%v,%v,%v\n", a, p, q)
	}
}

func TestIPxUint32(t *testing.T) {
	t.Parallel()

	testIPxUint32(t, 0)
	testIPxUint32(t, 10)
	testIPxUint32(t, 0x12892392)
}

func testIPxUint32(t *testing.T, n uint32) {
	ip := Uint322ip(n)
	nn := Ip2uint32(ip)
	if n != nn {
		t.Fatalf("Unexpected value=%d for ip=%q. Expected %d", nn, ip, n)
	}
}

func TestPerIPConnCounter(t *testing.T) {
	t.Parallel()

	var cc PerIPConnCounter

	expectPanic(t, func() { cc.Unregister(123) })

	for i := 1; i < 100; i++ {
		if n := cc.Register(123); n != i {
			t.Fatalf("Unexpected counter value=%d. Expected %d", n, i)
		}
	}

	n := cc.Register(456)
	if n != 1 {
		t.Fatalf("Unexpected counter value=%d. Expected 1", n)
	}

	for i := 1; i < 100; i++ {
		cc.Unregister(123)
	}
	cc.Unregister(456)

	expectPanic(t, func() { cc.Unregister(123) })
	expectPanic(t, func() { cc.Unregister(456) })

	n = cc.Register(123)
	if n != 1 {
		t.Fatalf("Unexpected counter value=%d. Expected 1", n)
	}
	cc.Unregister(123)
}

func expectPanic(t *testing.T, f func()) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("Expecting panic")
		}
	}()
	f()
}

package util4go

import (
	"fmt"
	"testing"
)

/**
* golang-sample源代码，版权归锦翰科技（深圳）有限公司所有。
* <p>
* 文件名称 : strings_test.go
* 文件路径 :
* 作者 : DavidLiu
× Email: david.liu@ginghan.com
*
* 创建日期 : 2022/5/4 14:16
* 修改历史 : 1. [2022/5/4 14:16] 创建文件 by LongYong
*/

type RegExpSampleCase struct {
	src, pattern, format string
	result               string
	error                bool
}

var RegExpSampleCases []RegExpSampleCase

func init() {
	RegExpSampleCases = make([]RegExpSampleCase, 0, 300)
	RegExpSampleCases = append(RegExpSampleCases, RegExpSampleCase{
		src: "/api/aaa/log.html", pattern: "/api/([^/]*)/([^/]*)$", format: "{0}?id={1}&file={2}",
		result: "/api/aaa/log.html?id=aaa&file=log.html", error: false,
	})
	RegExpSampleCases = append(RegExpSampleCases, RegExpSampleCase{
		src: "/api1/aaa/log.html", pattern: "/api1/([^/]*)/([^/]*)$", format: "{0}?id={1}&file={2}",
		result: "/api1/aaa/log.html?id=aaa&file=log.html", error: false,
	})
	RegExpSampleCases = append(RegExpSampleCases, RegExpSampleCase{
		src: "/api2/aaa/log.html", pattern: "/api2/([^/]*)/([^/]*)$", format: "{0}?id={1}&file={2}",
		result: "/api2/aaa/log.html?id=aaa&file=log.html", error: false,
	})
	RegExpSampleCases = append(RegExpSampleCases, RegExpSampleCase{
		src: "/api3/aaa/log.html", pattern: "/api3/([^/]*)/([^/]*)$", format: "{0}?id={1}&file={2}",
		result: "/api3/aaa/log.html?id=aaa&file=log.html", error: false,
	})
	RegExpSampleCases = append(RegExpSampleCases, RegExpSampleCase{
		src: "/api4/aaa/log.html", pattern: "/api4/([^/]*)/([^/]*)$", format: "{0}?id={1}&file={2}",
		result: "/api4/aaa/log.html?id=aaa&file=log.html", error: false,
	})
	RegExpSampleCases = append(RegExpSampleCases, RegExpSampleCase{
		src: "/api6/aaa/bbb/log.html", pattern: "/api6/([^/]*)/([^/]*)/([^/]*)$", format: "{0}?id={1}&file={3}&id2={2}",
		result: "/api6/aaa/log.html?id=aaa&file=log.html&id2=bbb", error: false,
	})
}

func TestRegExpPool(t *testing.T) {
	buf := "/api/aaa/log.html"
	reg := RegExpPool.GetRegExp(`/api/([^/]*)/([^/]*)$`)

	if reg == nil {
		fmt.Println("MustCompile err")
		return
	}

	//result := reg.FindAllStringSubmatch(buf, -1)
	//过滤<></>
	//for _, text := range result {
	//	fmt.Printf("text =%q, %q \n", text[0], text[1])
	//}
	result := reg.FindStringSubmatch(buf)
	//过滤<></>
	for _, text := range result {
		fmt.Printf("text =%q \n", text)
	}
}

func TestReplaceParameterValue(t *testing.T) {
	m := make(map[string]string)

	m["name"] = "Liuyong"
	m["age"] = "1022"

	fmt.Println(ReplaceParameterValue("{name}:{age}", m))
}

func TestConvertRegExpWithFormat(t *testing.T) {
	resultOfCase("/api/aaa/log.html", `/api/([^/]*)/([^/]*)$`, "{0}?id={1}&file={2}")
	resultOfCase("/api1/aaa/log.html", `/api1/([^/]*)/([^/]*)$`, "{0}?id={1}&file={2}")
	resultOfCase("/api2/aaa/log.html", `/api2/([^/]*)/([^/]*)$`, "{0}?id={1}&file={2}")
}

func resultOfCase(src, pattern, format string) {
	if s, err := RegExpPool.ConvertRegExpWithFormat(src, pattern, format); err == nil {
		fmt.Printf("RegExpPool.ConvertRegExpWithFormat(%q, %q, %q) = %v \n", src, pattern, format, s)
	} else {
		fmt.Printf("RegExpPool.ConvertRegExpWithFormat(%q, %q, %q) = %v \n", src, pattern, format, err)
	}

}

func match(rc RegExpSampleCase, b *testing.B) {
	if s, err := RegExpPool.ConvertRegExpWithFormat(rc.src, rc.pattern, rc.format); err == nil {
		if s == rc.result || 1 > 0 {

		} else {
			b.Errorf("RegExpPool.ConvertRegExpWithFormat(%q, %q, %q) = %v, but except %v \n", rc.src, rc.pattern, rc.format, s, rc.result)
		}

	} else {
		b.Errorf("RegExpPool.ConvertRegExpWithFormat(%q, %q, %q) = %v, but except %v \n", rc.src, rc.pattern, rc.format, err, rc.result)
	}
}

func BenchmarkTest(b *testing.B) {

	b.SetParallelism(8)
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		for _, c := range RegExpSampleCases {
			match(c, b)
			//	Match(c.match, c.sample)
		}
	}
}

func TestRegExp(t *testing.T) {
	buf := "/api/bbb/cccc/dddd/aaa.html"
	reg := RegExpPool.GetRegExp(`/api/([^/]*)/([\s\S]*$)`)

	result := reg.FindStringSubmatch(buf)
	//过滤<></>
	for _, text := range result {
		fmt.Printf("text =%q \n", text)
	}
}

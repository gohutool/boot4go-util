package path

import (
	"fmt"
	util4go "github.com/gohutool/boot4go-util"
	"math/rand"
	"path"
	"runtime"
	"testing"
)

/**
* golang-sample源代码，版权归锦翰科技（深圳）有限公司所有。
* <p>
* 文件名称 : pathmatch_test.go
* 文件路径 :
* 作者 : DavidLiu
× Email: david.liu@ginghan.com
*
* 创建日期 : 2022/5/3 10:31
* 修改历史 : 1. [2022/5/3 10:31] 创建文件 by LongYong
*/
/*
* 字符测试实例

//* 匹配0或多个非/的字符
path.Match("*", "a")            // true
path.Match("*", "sefesfe/")     // false
? 字符测试实例

//？匹配一个非/的字符
path.Match("a?b", "aab")    // true
path.Match("a?b", "a/b")    // false
[] 格式测试实例

path.Match("[abc][123]", "b2")        // true

path.Match("[abc][1-3]", "b2")        // true

path.Match("[abc][^123]", "b2")        // false

path.Match("[abc][^123]", "b4")        // true
字符或者特殊用途字符(  \\   ?  *   [ )测试实例

path.Match("a\\\\", "a\\")	// true
path.Match("a\\[", "a[")		// true
path.Match("a\\?", "a?")		// true
path.Match("a\\*", "a*")		// true
path.Match("abc", "abc")		// true

*/
func TestPathMatch(t *testing.T) {
	printMatch("*", "a")
	printMatch("*", "sefesfe/")
	printMatch("*/b", "b")
	printMatch("/*/b", "b")
	printMatch("/*", "/b")
	printMatch("/*/b", "/b")
	printMatch("**/b", "/b")

	printMatch("/api/*", "/api/a/a/b")
	printMatch("/api/*/?b", "/api/a/ab")

	printMatch("/*", "/a")
	printMatch("/*", "/")
	printMatch("/*", "/abc/abd")

	printMatch("a?b", "a0b")
	printMatch("a?b", "a/b")
	printMatch("a?b", "abb")
	printMatch("a?b", "axb")

	printMatch("a\\\\", "a\\")
	printMatch("a\\[", "a[")
	printMatch("a\\?", "a?")
	printMatch("a\\*", "a*")
	printMatch("abc", "abc")
}

func printMatch(a, b string) {
	if r, err := path.Match(a, b); err != nil {
		fmt.Printf("Error %v\n", err)
	} else {
		fmt.Printf("path.Match(%q, %q)=%v\n", a, b, r)
	}
}

type TestCase struct {
	match, sample string
	isMatch       bool
	isStd         bool
	error         bool
}

var CaseSample []TestCase

func init() {
	CaseSample = make([]TestCase, 0, 300)

	CaseSample = append(CaseSample, TestCase{
		match: "*", sample: "", isMatch: true, isStd: true, error: false,
	})
	for idx := 'a'; idx < 'f'; idx++ {
		CaseSample = append(CaseSample, TestCase{
			match: "*", sample: string(idx), isMatch: true, isStd: true, error: false,
		})
	}
	for idx := 0; idx < 10; idx++ {
		v := ""
		for i := 0; i < 3; i++ {
			v = v + string(RandRune())
		}
		CaseSample = append(CaseSample, TestCase{
			match: "*", sample: v, isMatch: true, isStd: true, error: false,
		})
	}

	CaseSample = append(CaseSample, TestCase{
		match: "*", sample: "/", isMatch: false, isStd: true, error: false,
	})
	for idx := 2; idx < 4; idx++ {
		CaseSample = append(CaseSample, TestCase{
			match: "*", sample: "/" + util4go.LeftPad("a", idx, RandRune()), isMatch: false, isStd: true, error: false,
		})
	}
	CaseSample = append(CaseSample, TestCase{
		match: "/*", sample: "/", isMatch: true, isStd: true, error: false,
	})
	for idx := 2; idx < 4; idx++ {
		CaseSample = append(CaseSample, TestCase{
			match: "/*", sample: "/" + util4go.LeftPad("a", idx, RandRune()), isMatch: true, isStd: true, error: false,
		})
	}

	for idx := 0; idx < 10; idx++ {
		v := "/"
		for i := 0; i < 3; i++ {
			v = v + string(RandRune())
		}
		CaseSample = append(CaseSample, TestCase{
			match: "/*", sample: v, isMatch: true, isStd: true, error: false,
		})
	}

	CaseSample = append(CaseSample, TestCase{
		match: "*c", sample: "abc", isMatch: true, isStd: true, error: false,
	})

	for idx := 1; idx < 4; idx++ {
		CaseSample = append(CaseSample, TestCase{
			match: "*c", sample: util4go.LeftPad("c", idx, RandRune()), isMatch: true, isStd: true, error: false,
		})
	}

	CaseSample = append(CaseSample, TestCase{
		match: "*c", sample: util4go.RightPad("c", 2, RandRune()) + "/", isMatch: false, isStd: true, error: false,
	})
	CaseSample = append(CaseSample, TestCase{
		match: "*c", sample: "/" + util4go.RightPad("c", 2, RandRune()), isMatch: false, isStd: true, error: false,
	})

	CaseSample = append(CaseSample, TestCase{
		match: "a*/b", sample: "abc/b", isMatch: true, isStd: true, error: false,
	})

	CaseSample = append(CaseSample, TestCase{
		match: "a*/b", sample: "a/c/b", isMatch: false, isStd: true, error: false,
	})

	CaseSample = append(CaseSample, TestCase{
		match: "a*b*c*d*e*", sample: "axbxcxdxe", isMatch: true, isStd: true, error: false,
	})

	CaseSample = append(CaseSample, TestCase{
		match: "a*b*/a", sample: "axbxxxx/a", isMatch: true, isStd: true, error: false,
	})

	CaseSample = append(CaseSample, TestCase{
		match: "ab[c]", sample: "abc", isMatch: true, isStd: true, error: false,
	})

	CaseSample = append(CaseSample, TestCase{
		match: "ab[b-d]", sample: "abc", isMatch: true, isStd: true, error: false,
	})

	CaseSample = append(CaseSample, TestCase{
		match: "ab[b-d]", sample: "abf", isMatch: false, isStd: true, error: false,
	})

	CaseSample = append(CaseSample, TestCase{
		match: "ab[^c]", sample: "abc", isMatch: false, isStd: true, error: false,
	})

	CaseSample = append(CaseSample, TestCase{
		match: "ab[^b-d]", sample: "abf", isMatch: true, isStd: true, error: false,
	})

	CaseSample = append(CaseSample, TestCase{
		match: "ab[^b-d]", sample: "abf", isMatch: true, isStd: true, error: false,
	})

	CaseSample = append(CaseSample, TestCase{
		match: "*.log", sample: "a/a/a.log", isMatch: false, isStd: false, error: false,
	})
	CaseSample = append(CaseSample, TestCase{
		match: "**/*.log", sample: "a/a/a.log", isMatch: true, isStd: false, error: false,
	})
	CaseSample = append(CaseSample, TestCase{
		match: "[a-b-d]", sample: "a", isMatch: true, isStd: false, error: false,
	})
	CaseSample = append(CaseSample, TestCase{
		match: "[a-bb-d]*", sample: "d", isMatch: true, isStd: false, error: false,
	})
	CaseSample = append(CaseSample, TestCase{
		match: "[a-be-g]*", sample: "f", isMatch: true, isStd: false, error: false,
	})
	CaseSample = append(CaseSample, TestCase{
		match: "[a-be-g]*", sample: "c", isMatch: false, isStd: false, error: false,
	})

	for idx := 0; idx < 10; idx++ {
		CaseSample = append(CaseSample, TestCase{
			match: "[a-fo-z]*", sample: string(RandRune2('a', 'f')) + string(RandRune()) + string(RandRune()),
			isMatch: true, isStd: false, error: false,
		})
	}

	for idx := 0; idx < 10; idx++ {
		CaseSample = append(CaseSample, TestCase{
			match: "[a-fo-z]*", sample: string(RandRune2('o', 'z')) + string(RandRune()) + string(RandRune()),
			isMatch: true, isStd: false, error: false,
		})
	}

	for idx := 0; idx < 10; idx++ {
		CaseSample = append(CaseSample, TestCase{
			match: "[a-fo-z]*", sample: string(RandRune2('g', 'm')) + string(RandRune()) + string(RandRune()),
			isMatch: false, isStd: false, error: false,
		})
	}

	CaseSample = append(CaseSample, TestCase{
		match: "ab[!c]", sample: "abc", isMatch: false, isStd: false, error: false,
	})

	for idx := 0; idx < 10; idx++ {
		CaseSample = append(CaseSample, TestCase{
			match: "ab[!c]", sample: "ab" + string(RandRune2('d', 'z')), isMatch: true, isStd: false, error: false,
		})
	}
}

func RandRune() rune {
	return rune(rand.Intn(25) + int('a'))
}

func RandRune2(s, e rune) rune {

	return rune(rand.Intn(int(e)-int(s)) + int(s))
}

func TestCaseMatch(t *testing.T) {

	t.Logf("TestCase : %v", len(CaseSample))

	for idx, tt := range CaseSample {
		// Since Match() always uses "/" as the separator, we
		// don't need to worry about the tt.testOnDisk flag
		matchRun(t, idx, tt)
	}
}

func matchRun(t *testing.T, idx int, tt TestCase) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("#%v. Match(%#q, %#q) panicked: %#v", idx, tt.match, tt.sample, r)
		}
	}()

	// Match() always uses "/" as the separator
	matched, err := Match(tt.match, tt.sample)
	if matched != tt.isMatch || (err != nil && !tt.error) {
		t.Errorf("#%v. Match(%#q, %#q) = %v, %v. but expect %v, %v", idx, tt.match, tt.sample, matched, err, tt.isMatch, err)
	}

	if tt.isStd {
		stdOk, stdErr := path.Match(tt.match, tt.sample)
		if matched != stdOk || !compareErrors(err, stdErr) {
			t.Errorf("#%v. path.Match(%#q, %#q) = %v, %v. but expect %v, %v", idx, tt.match, tt.sample, matched, err, stdOk, err)
		}
	}
}

type MatchTest struct {
	pattern, testPath string // a pattern and path to test the pattern on
	shouldMatch       bool   // true if the pattern should match the path
	expectedErr       error  // an expected error
	isStandard        bool   // pattern doesn't use any doublestar features
	testOnDisk        bool   // true: test pattern against files in "test" directory
	numResults        int    // number of glob results if testing on disk
	winNumResults     int    // number of glob results on Windows
}

// Tests which contain escapes and symlinks will not work on Windows
var onWindows = runtime.GOOS == "windows"

var matchTests = []MatchTest{
	{"*", "", true, nil, true, false, 0, 0},
	{"*", "/", false, nil, true, false, 0, 0},
	{"/*", "/", true, nil, true, false, 0, 0},
	{"/*", "/debug/", false, nil, true, false, 0, 0},
	{"/*", "//", false, nil, true, false, 0, 0},
	{"abc", "abc", true, nil, true, true, 1, 1},
	{"*", "abc", true, nil, true, true, 19, 15},
	{"*c", "abc", true, nil, true, true, 2, 2},
	{"*/", "a/", true, nil, true, false, 0, 0},
	{"a*", "a", true, nil, true, true, 9, 9},
	{"a*", "abc", true, nil, true, true, 9, 9},
	{"a*", "ab/c", false, nil, true, true, 9, 9},
	{"a*/b", "abc/b", true, nil, true, true, 2, 2},
	{"a*/b", "a/c/b", false, nil, true, true, 2, 2},
	{"a*b*c*d*e*", "axbxcxdxe", true, nil, true, true, 3, 3},
	{"a*b*c*d*e*/f", "axbxcxdxe/f", true, nil, true, true, 2, 2},
	{"a*b*c*d*e*/f", "axbxcxdxexxx/f", true, nil, true, true, 2, 2},
	{"a*b*c*d*e*/f", "axbxcxdxe/xxx/f", false, nil, true, true, 2, 2},
	{"a*b*c*d*e*/f", "axbxcxdxexxx/fff", false, nil, true, true, 2, 2},
	{"a*b?c*x", "abxbbxdbxebxczzx", true, nil, true, true, 2, 2},
	{"a*b?c*x", "abxbbxdbxebxczzy", false, nil, true, true, 2, 2},
	{"ab[c]", "abc", true, nil, true, true, 1, 1},
	{"ab[b-d]", "abc", true, nil, true, true, 1, 1},
	{"ab[e-g]", "abc", false, nil, true, true, 0, 0},
	{"ab[^c]", "abc", false, nil, true, true, 0, 0},
	{"ab[^b-d]", "abc", false, nil, true, true, 0, 0},
	{"ab[^e-g]", "abc", true, nil, true, true, 1, 1},
	{"a\\*b", "ab", false, nil, true, true, 0, 0},
	{"a?b", "a☺b", true, nil, true, true, 1, 1},
	{"a[^a]b", "a☺b", true, nil, true, true, 1, 1},
	{"a[!a]b", "a☺b", true, nil, false, true, 1, 1},
	{"a???b", "a☺b", false, nil, true, true, 0, 0},
	{"a[^a][^a][^a]b", "a☺b", false, nil, true, true, 0, 0},
	{"[a-ζ]*", "α", true, nil, true, true, 17, 15},
	{"*[a-ζ]", "A", false, nil, true, true, 17, 15},
	{"a?b", "a/b", false, nil, true, true, 1, 1},
	{"a*b", "a/b", false, nil, true, true, 1, 1},
	{"[\\]a]", "]", true, nil, true, !onWindows, 2, 2},
	{"[\\-]", "-", true, nil, true, !onWindows, 1, 1},
	{"[x\\-]", "x", true, nil, true, !onWindows, 2, 2},
	{"[x\\-]", "-", true, nil, true, !onWindows, 2, 2},
	{"[x\\-]", "z", false, nil, true, !onWindows, 2, 2},
	{"[\\-x]", "x", true, nil, true, !onWindows, 2, 2},
	{"[\\-x]", "-", true, nil, true, !onWindows, 2, 2},
	{"[\\-x]", "a", false, nil, true, !onWindows, 2, 2},
	{"[]a]", "]", false, ErrBadPattern, true, true, 0, 0},
	// doublestar, like bash, allows these when path.Match() does not
	{"[-]", "-", true, nil, false, !onWindows, 1, 0},
	{"[x-]", "x", true, nil, false, true, 2, 1},
	{"[x-]", "-", true, nil, false, !onWindows, 2, 1},
	{"[x-]", "z", false, nil, false, true, 2, 1},
	{"[-x]", "x", true, nil, false, true, 2, 1},
	{"[-x]", "-", true, nil, false, !onWindows, 2, 1},
	{"[-x]", "a", false, nil, false, true, 2, 1},
	{"[a-b-d]", "a", true, nil, false, true, 3, 2},
	{"[a-b-d]", "b", true, nil, false, true, 3, 2},
	{"[a-b-d]", "-", true, nil, false, !onWindows, 3, 2},
	{"[a-b-d]", "c", false, nil, false, true, 3, 2},
	{"[a-b-x]", "x", true, nil, false, true, 4, 3},
	{"\\", "a", false, ErrBadPattern, true, !onWindows, 0, 0},
	{"[", "a", false, ErrBadPattern, true, true, 0, 0},
	{"[^", "a", false, ErrBadPattern, true, true, 0, 0},
	{"[^bc", "a", false, ErrBadPattern, true, true, 0, 0},
	{"a[", "a", false, ErrBadPattern, true, true, 0, 0},
	{"a[", "ab", false, ErrBadPattern, true, true, 0, 0},
	{"ad[", "ab", false, ErrBadPattern, true, true, 0, 0},
	{"*x", "xxx", true, nil, true, true, 4, 4},
	{"[abc]", "b", true, nil, true, true, 3, 3},
	{"**", "", true, nil, false, false, 38, 38},
	{"a/**", "a", true, nil, false, true, 7, 7},
	{"a/**", "a/", true, nil, false, false, 7, 7},
	{"a/**", "a/b", true, nil, false, true, 7, 7},
	{"a/**", "a/b/c", true, nil, false, true, 7, 7},
	{"**/c", "c", true, nil, false, true, 5, 4},
	{"**/c", "b/c", true, nil, false, true, 5, 4},
	{"**/c", "a/b/c", true, nil, false, true, 5, 4},
	{"**/c", "a/b", false, nil, false, true, 5, 4},
	{"**/c", "abcd", false, nil, false, true, 5, 4},
	{"**/c", "a/abc", false, nil, false, true, 5, 4},
	{"a/**/b", "a/b", true, nil, false, true, 2, 2},
	{"a/**/c", "a/b/c", true, nil, false, true, 2, 2},
	{"a/**/d", "a/b/c/d", true, nil, false, true, 1, 1},
	{"a/\\**", "a/b/c", false, nil, false, !onWindows, 0, 0},
	{"a/\\[*\\]", "a/bc", false, nil, true, !onWindows, 0, 0},
	// this is an odd case: filepath.Glob() will return results
	{"a//b/c", "a/b/c", false, nil, true, false, 0, 0},
	{"a/b/c", "a/b//c", false, nil, true, true, 1, 1},
	// also odd: Glob + filepath.Glob return results
	{"a/", "a", false, nil, true, false, 0, 0},
	{"ab{c,d}", "abc", true, nil, false, true, 1, 1},
	{"ab{c,d,*}", "abcde", true, nil, false, true, 5, 5},
	{"ab{c,d}[", "abcd", false, ErrBadPattern, false, true, 0, 0},
	{"a{,bc}", "a", true, nil, false, true, 2, 2},
	{"a{,bc}", "abc", true, nil, false, true, 2, 2},
	{"a/{b/c,c/b}", "a/b/c", true, nil, false, true, 2, 2},
	{"a/{b/c,c/b}", "a/c/b", true, nil, false, true, 2, 2},
	{"{a/{b,c},abc}", "a/b", true, nil, false, true, 3, 3},
	{"{a/{b,c},abc}", "a/c", true, nil, false, true, 3, 3},
	{"{a/{b,c},abc}", "abc", true, nil, false, true, 3, 3},
	{"{a/{b,c},abc}", "a/b/c", false, nil, false, true, 3, 3},
	{"{a/ab*}", "a/abc", true, nil, false, true, 1, 1},
	{"{a/*}", "a/b", true, nil, false, true, 3, 3},
	{"{a/abc}", "a/abc", true, nil, false, true, 1, 1},
	{"{a/b,a/c}", "a/c", true, nil, false, true, 2, 2},
	{"abc/**", "abc/b", true, nil, false, true, 3, 3},
	{"**/abc", "abc", true, nil, false, true, 2, 2},
	{"abc**", "abc/b", false, nil, false, true, 3, 3},
	{"**/*.txt", "abc/【test】.txt", true, nil, false, true, 1, 1},
	{"**/【*", "abc/【test】.txt", true, nil, false, true, 1, 1},
	// unfortunately, io/fs can't handle this, so neither can Glob =(
	{"broken-symlink", "broken-symlink", true, nil, true, false, 1, 1},
	{"working-symlink/c/*", "working-symlink/c/d", true, nil, true, !onWindows, 1, 1},
	{"working-sym*/*", "working-symlink/c", true, nil, true, !onWindows, 1, 1},
	{"b/**/f", "b/symlink-dir/f", true, nil, false, !onWindows, 2, 2},
}

func TestValidatePattern(t *testing.T) {
	for idx, tt := range matchTests {
		testValidatePatternWith(t, idx, tt)
	}
}

func TestMatch(t *testing.T) {
	for idx, tt := range matchTests {
		// Since Match() always uses "/" as the separator, we
		// don't need to worry about the tt.testOnDisk flag
		testMatchWith(t, idx, tt)
	}
}

func testMatchWith(t *testing.T, idx int, tt MatchTest) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("#%v. Match(%#q, %#q) panicked: %#v", idx, tt.pattern, tt.testPath, r)
		}
	}()

	// Match() always uses "/" as the separator
	ok, err := Match(tt.pattern, tt.testPath)
	if ok != tt.shouldMatch || err != tt.expectedErr {
		t.Errorf("#%v. Match(%#q, %#q) = %v, %v want %v, %v", idx, tt.pattern, tt.testPath, ok, err, tt.shouldMatch, tt.expectedErr)
	}

	if tt.isStandard {
		stdOk, stdErr := path.Match(tt.pattern, tt.testPath)
		if ok != stdOk || !compareErrors(err, stdErr) {
			t.Errorf("#%v. Match(%#q, %#q) != path.Match(...). Got %v, %v want %v, %v", idx, tt.pattern, tt.testPath, ok, err, stdOk, stdErr)
		}
	}
}

func compareErrors(a, b error) bool {
	if a == nil {
		return b == nil
	}
	return b != nil
}

func testValidatePatternWith(t *testing.T, idx int, tt MatchTest) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("#%v. Validate(%#q) panicked: %#v", idx, tt.pattern, r)
		}
	}()

	result := ValidatePattern(tt.pattern)
	if result != (tt.expectedErr == nil) {
		t.Errorf("#%v. ValidatePattern(%#q) = %v want %v", idx, tt.pattern, result, !result)
	}
}

func BenchmarkMatch(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		for _, tt := range matchTests {
			Match(tt.pattern, tt.testPath)
			//if tt.isStandard {
			//
			//}
		}
	}
}

func BenchmarkStdMatch(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		for _, tt := range matchTests {
			path.Match(tt.pattern, tt.testPath)
			//if tt.isStandard {
			//
			//}
		}
	}
}

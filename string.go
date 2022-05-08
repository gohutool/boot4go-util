package util4go

import (
	"fmt"
	"math/rand"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

/**
* golang-sample源代码，版权归锦翰科技（深圳）有限公司所有。
* <p>
* 文件名称 : string.go
* 文件路径 :
* 作者 : DavidLiu
× Email: david.liu@ginghan.com
*
* 创建日期 : 2022/5/1 09:48
* 修改历史 : 1. [2022/5/1 09:48] 创建文件 by LongYong
*/

func Substring(source string, start int, end int) string {
	var r = []rune(source)

	if end <= 0 {
		return string(r[start:])
	}

	return string(r[start:end])
}

func Str2Int64(source string) (any, error) {
	v, error := strconv.ParseInt(fmt.Sprintf("%v", source), 10, 64)
	if error == nil {
		return int64(v), nil
	}
	return nil, error
}

func Str2Int32(source string) (any, error) {
	v, error := strconv.ParseInt(fmt.Sprintf("%v", source), 10, 64)
	if error == nil {
		return int32(v), nil
	}
	return nil, error
}

func Str2Int16(source string) (any, error) {
	v, error := strconv.ParseInt(fmt.Sprintf("%v", source), 10, 64)
	if error == nil {
		return int16(v), nil
	}
	return nil, error
}

func Str2Int8(source string) (any, error) {
	v, error := strconv.ParseInt(fmt.Sprintf("%v", source), 10, 64)
	if error == nil {
		return int8(v), nil
	}
	return nil, error
}

func Str2Int(source string) (any, error) {
	v, error := strconv.Atoi(fmt.Sprintf("%v", source))
	if error == nil {
		return v, nil
	}
	return nil, error
}

func Str2UInt64(source string) (any, error) {
	v, error := strconv.ParseUint(fmt.Sprintf("%v", source), 10, 64)
	if error == nil {
		return uint64(v), nil
	}
	return nil, error
}

func Str2Uint32(source string) (any, error) {
	v, error := strconv.ParseUint(fmt.Sprintf("%v", source), 10, 64)
	if error == nil {
		return uint32(v), nil
	}
	return nil, error
}

func Str2Uint16(source string) (any, error) {
	v, error := strconv.ParseUint(fmt.Sprintf("%v", source), 10, 64)
	if error == nil {
		return uint16(v), nil
	}
	return nil, error
}

func Str2Uint8(source string) (any, error) {
	v, error := strconv.ParseUint(fmt.Sprintf("%v", source), 10, 64)
	if error == nil {
		return uint8(v), nil
	}
	return nil, error
}

func Str2Uint(source string) (any, error) {
	v, error := strconv.ParseUint(fmt.Sprintf("%v", source), 10, 64)
	if error == nil {
		return uint(v), nil
	}
	return nil, error
}

func Str2Bool(source string) (any, error) {
	v, error := strconv.ParseBool(fmt.Sprintf("%v", source))
	if error == nil {
		return v, nil
	}
	return nil, error
}

func Str2Float64(source string) (any, error) {
	v, error := strconv.ParseFloat(fmt.Sprintf("%v", source), 10)
	if error == nil {
		return v, nil
	}
	return nil, error
}

func Str2Float32(source string) (any, error) {
	v, error := strconv.ParseFloat(fmt.Sprintf("%v", source), 10)
	if error == nil {
		return float32(v), nil
	}
	return nil, error
}

func Str2Object(v string, k reflect.Kind) (any, error) {
	if len(v) == 0 {
		return nil, nil
	}

	switch k {
	case reflect.String:
		return v, nil
	case reflect.Int:
		return Str2Int(v)
	case reflect.Int16:
		return Str2Int16(v)
	case reflect.Int32:
		return Str2Int32(v)
	case reflect.Int64:
		return Str2Int64(v)
	case reflect.Int8:
		return Str2Int8(v)
	case reflect.Uint:
		return Str2Uint(v)
	case reflect.Uint8:
		return Str2Uint8(v)
	case reflect.Uint16:
		return Str2Uint16(v)
	case reflect.Uint32:
		return Str2Uint32(v)
	case reflect.Uint64:
		return Str2UInt64(v)
	case reflect.Bool:
		return Str2Bool(v)
	case reflect.Float32:
		return Str2Float32(v)
	case reflect.Float64:
		return Str2Float64(v)
	}

	return nil, nil
}

// Int64ToStr int64 to str
func Int64ToStr(i int64) string {
	return strconv.FormatInt(i, 10)
}

func String(obj any) string {
	return fmt.Sprintf("%+v", obj)
}

func IsEmpty(str string) bool {
	if len(str) == 0 || len(strings.TrimSpace(str)) == 0 {
		return true
	}

	return false
}

func BinarySearch(a []string, v interface{}) int {
	for i, i2 := range a {
		if i2 == v {
			return i
		}
	}
	return -1
}

func BuildFormatString(format string, a ...any) string {
	return fmt.Sprintf(format, a...)
}

func BuildString(a ...any) string {
	return fmt.Sprint(a...)
}

func RightPad(str string, limit int, placeholder rune) string {
	len := limit - len(str)
	if len >= limit {
		return str
	}

	arr := make([]any, 0, len)
	arr = append(arr, str)

	for idx := 0; idx < len; idx++ {
		arr = append(arr, string(placeholder))
	}

	return BuildString(arr...)
}

func LeftPad(str string, limit int, placeholder rune) string {
	len := limit - len(str)
	if len >= limit {
		return str
	}

	arr := make([]any, 0, len)
	for idx := 0; idx < len; idx++ {
		arr = append(arr, string(placeholder))
	}
	arr = append(arr, str)

	return BuildString(arr...)
}

// StrNumSuffixInt
// Parse a number with K/M/G suffixes based on thousands (1000) or 2^10 (1024)
func StrNumSuffixInt(str string, mult int) int {
	num := 1
	if len(str) > 1 {
		switch str[len(str)-1] {
		case 'G', 'g':
			num *= mult
			fallthrough
		case 'M', 'm':
			num *= mult
			fallthrough
		case 'K', 'k':
			num *= mult
			str = str[0 : len(str)-1]
		}
	}
	parsed, _ := strconv.Atoi(str)
	return parsed * num
}

var _expression_reg = regexp.MustCompile(`\{(?s:(.*?))\}`)

func ParseParameterName(str string) []string {
	result := _expression_reg.FindAllStringSubmatch(str, -1)
	l := len(result)
	if l == 0 {
		return nil
	}

	if l == 1 {
		return []string{result[0][1]}
	}

	yet := make(map[string]bool)

	ks := Reduce(result, func(one []string) (string, bool) {
		if one != nil {
			if len(one[1]) > 0 {
				if _, ok := yet[one[1]]; !ok {
					return one[1], true
				}
			}
		}

		return "", false
	})

	return ks
}

const _KEY_1 = "{"
const _KEY_2 = "}"

func ReplaceParameterValue(str string, keyAndValue map[string]string) string {
	keys := ParseParameterName(str)

	if len(keys) == 0 {
		return str
	}

	for _, k := range keys {
		if v, ok := keyAndValue[k]; ok {
			str = strings.ReplaceAll(str, _KEY_1+k+_KEY_2, v)
		} else {
			//str = strings.ReplaceAll(str, k, v)
		}
	}

	return str
}

func ReplaceParameterWithKeyValue(str string, keyAndValue map[string]string) string {
	if len(keyAndValue) == 0 {
		return str
	}

	for k, v := range keyAndValue {
		str = strings.ReplaceAll(str, k, v)
	}

	return str
}

func RandRune() rune {
	return rune(rand.Intn(25) + int('a'))
}

func RandRune2(s, e rune) rune {

	return rune(rand.Intn(int(e)-int(s)) + int(s))
}

type regExpPool map[string]*regexp.Regexp

var regExpPoolLock = sync.RWMutex{}
var RegExpPool = make(regExpPool)
var regExpPool_MaxSize = 0

func (rp *regExpPool) MaxSize(maxsize int) {
	regExpPool_MaxSize = maxsize
}

func (rp *regExpPool) GetRegExp(exp string) *regexp.Regexp {
	regExpPoolLock.RLock()
	if r, ok := RegExpPool[exp]; ok {
		defer regExpPoolLock.RUnlock()
		return r
	} else {
		regExpPoolLock.RUnlock()

		return func() *regexp.Regexp {
			regExpPoolLock.Lock()
			defer regExpPoolLock.Unlock()
			reg := regexp.MustCompile(exp)
			RegExpPool[exp] = reg
			return reg
		}()
	}
}

type RegExpError struct {
}

func (e RegExpError) Error() string {
	return "RegExpError"
}

func (rp *regExpPool) FindStringSubmatch(src, pattern string) ([]string, error) {
	regExp := rp.GetRegExp(pattern)

	if regExp == nil {
		return nil, RegExpError{}
	}

	return regExp.FindStringSubmatch(src), nil
}

func (rp *regExpPool) ConvertRegExpWithFormat(src, pattern, format string) (string, error) {
	regExp := rp.GetRegExp(pattern)
	if regExp == nil {
		return "", RegExpError{}
	}

	result := regExp.FindStringSubmatch(src)

	if len(result) == 0 {
		return format, RegExpError{}
	}

	//m := make(map[string]string)
	//for idx, text := range result {
	//	m[strconv.Itoa(idx)] = text
	//}
	//
	//return ReplaceParameterValue(format, m), nil

	for idx, text := range result {
		b := make([]byte, 0, len(text)+4)
		b = append(b, byte('{'))
		b = append(b, []byte(strconv.Itoa(idx))...)
		b = append(b, byte('}'))

		format = strings.ReplaceAll(format, string(b), text)

		//format = strings.ReplaceAll(format, fmt.Sprintf("{%d}", idx), text)
	}

	return format, nil

}

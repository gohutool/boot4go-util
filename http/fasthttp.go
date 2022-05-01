package http

import (
	"bytes"
	"encoding/json"
	"errors"
	fastjson "github.com/gohutool/boot4go-fastjson"
	"github.com/valyala/fasthttp"
	"strings"
)

/**
* golang-sample源代码，版权归锦翰科技（深圳）有限公司所有。
* <p>
* 文件名称 : fasthttp.go
* 文件路径 :
* 作者 : DavidLiu
× Email: david.liu@ginghan.com
*
* 创建日期 : 2022/5/1 09:41
* 修改历史 : 1. [2022/5/1 09:41] 创建文件 by LongYong
*/

const (
	SUCCESS = 0
	FAIL    = 1
)

var Result = result{}

type result struct {
	Msg    string `json:"msg"`
	Status int    `json:"status"`
	Data   any    `json:"data,omitempty"`
}

func (r result) Fail(msg string) result {
	r1 := result{}
	r1.Msg = msg
	r1.Data = nil
	r1.Status = FAIL

	return r1
}

func (o result) Success(data any, msg string) result {
	r := result{}
	r.Msg = msg
	r.Data = data
	r.Status = SUCCESS

	return r
}

func (o result) Error(msg string) error {
	return errors.New(Result.Fail(msg).Json())
}

func (r result) Response(ctx *fasthttp.RequestCtx) error {
	ctx.Response.Header.Set("Content-type", "application/json;charset=utf-8")

	b, error := r.Marshal()

	if error == nil {
		_, error = ctx.Write(b)
		return error
	}

	return error
}

func (r result) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func (r result) Json() string {
	b, error := json.Marshal(r)

	if error == nil {
		return string(b)
	}

	return ""
}

var PageResultData = pageResultData{}

var PageResultDataBuilder = pageResultDataBuilder{}

type pageResultDataBuilder struct {
	pd pageResultData
}

type pageResultData struct {
	List       []any `json:"list"`
	TotalRow   int   `json:"totalRow"`
	PageNumber int   `json:"pageNumber"`
	TotalPage  int   `json:"totalPage"`
	PageSize   int   `json:"pageSize"`
}

func (r pageResultDataBuilder) New() pageResultDataBuilder {
	r1 := pageResultDataBuilder{}
	r1.pd = pageResultData{}
	return r1
}

func (r pageResultDataBuilder) Data(list []any) pageResultDataBuilder {
	r.pd.List = list
	return r
}

func (r pageResultDataBuilder) TotalRow(total int) pageResultDataBuilder {
	r.pd.TotalRow = total
	return r
}

func (r pageResultDataBuilder) PageNumber(pageNumber int) pageResultDataBuilder {
	r.pd.PageNumber = pageNumber
	return r
}

func (r pageResultDataBuilder) TotalPage(totalPage int) pageResultDataBuilder {
	r.pd.TotalPage = totalPage
	return r
}

func (r pageResultDataBuilder) PageSize(pageSize int) pageResultDataBuilder {
	r.pd.PageSize = pageSize
	return r
}

func (r pageResultDataBuilder) Build() pageResultData {
	pd := pageResultData{}

	if r.pd.List == nil {
		pd.List = []any{}
	} else {
		pd.List = r.pd.List
	}

	if r.pd.PageNumber <= 0 {
		pd.PageNumber = 1
	} else {
		pd.PageNumber = r.pd.PageNumber
	}

	if r.pd.TotalRow <= 0 {
		pd.TotalRow = len(pd.List)
	} else {
		pd.TotalRow = r.pd.TotalRow
	}

	if r.pd.PageSize <= 0 {
		pd.PageSize = 0
	} else {
		pd.PageSize = r.pd.PageSize
	}

	if pd.TotalRow == 0 {
		pd.TotalPage = 0
	} else {
		if pd.PageSize == 0 {
			pd.TotalPage = 1
		} else {
			pd.TotalPage = (pd.TotalRow + pd.PageSize - 1) / pd.PageSize
		}
	}

	return pd
}

var DEFAULT_CONTENTTYPE = []byte("application/json;charset=utf-8")

func Error(ctx *fasthttp.RequestCtx, msg string, statusCode int) {
	ctx.Response.Header.Set("Content-type", "application/json;charset=utf-8")
	var oh = &fasthttp.ResponseHeader{}
	ctx.Response.Header.CopyTo(oh)

	ctx.Response.Reset()
	oh.CopyTo(&ctx.Response.Header)

	//ctx.Response.Header.Set("Content-type", "application/json;charset=utf-8")
	//ctx.Response.Header.AppendBytes(b)
	//ctx.Error(msg, statusCode)
	ctx.SetStatusCode(statusCode)
	ctx.SetContentTypeBytes(DEFAULT_CONTENTTYPE)
	ctx.SetBodyString(msg)
}

func Header(ctx *fasthttp.RequestCtx, key, defaultV string) (string, bool) {
	b := ctx.Request.Header.Peek(key)

	if b == nil {
		return "", false
	}

	return string(b), true
}

func GetHeader(ctx *fasthttp.RequestCtx, key, defaultV string) string {
	v := string(ctx.Request.Header.Peek(key))

	if len(v) == 0 || len(strings.TrimSpace(v)) == 0 {
		return defaultV
	} else {
		return v
	}
}

func GetParams(ctx *fasthttp.RequestCtx, key, defaultV string) string {

	v := string(ctx.FormValue(key))

	if len(v) == 0 || len(strings.TrimSpace(v)) == 0 {
		return defaultV
	} else {
		return v
	}
}

func GetToken(ctx *fasthttp.RequestCtx) string {
	token := GetHeader(ctx, "authorization", "")
	if len(token) == 0 {
		return token
	} else {
		return strings.Replace(token, "Bearer ", "", 1)
	}
}

func Param(ctx *fasthttp.RequestCtx, key string) (string, bool) {
	b := ctx.FormValue(key)

	if b == nil {
		return "", false
	}

	return string(b), true
}

func Params(ctx *fasthttp.RequestCtx) map[string]any {
	rtn := make(map[string]any)

	ctx.PostArgs().VisitAll(func(key, value []byte) {
		rtn[string(key)] = string(value)
	})

	ctx.QueryArgs().VisitAll(func(key, value []byte) {
		rtn[string(key)] = string(value)
	})

	return rtn
}

func JsonBody2Object(ctx *fasthttp.RequestCtx, obj any) error {
	s := ctx.PostBody()
	return unmarshal(s, obj)
}

func unmarshal(b []byte, obj any) error {
	decoder := json.NewDecoder(bytes.NewReader(b))
	decoder.UseNumber()
	return decoder.Decode(obj)

}

func JsonBodyUnmarshalObject[T fastjson.Unmarshalable](ctx *fasthttp.RequestCtx, obj T) error {
	s := string(ctx.PostBody())
	return fastjson.UnmarshalObject(s, obj)
}

package http

import (
	"fmt"
	util4go "github.com/gohutool/boot4go-util"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
)

/**
* boot4go-docker-ui源代码，版权归锦翰科技（深圳）有限公司所有。
* <p>
* 文件名称 : unix.go
* 文件路径 :
* 作者 : DavidLiu
× Email: david.liu@ginghan.com
*
* 创建日期 : 2022/7/19 17:59
* 修改历史 : 1. [2022/7/19 17:59] 创建文件 by LongYong
*/

func UnixSocketProxy(unixPath string, host string, schema string, uri string, query string, ctx *fasthttp.RequestCtx, timeout int) error {

	schema = "http"
	conn, err := net.Dial("unix", unixPath)
	host = "localhost"

	Logger.Debug("Unix start")

	if err != nil {
		ctx.Error("Internal Server Error", fasthttp.StatusInternalServerError)
		return err
	}

	c := httputil.NewClientConn(conn, nil)
	defer c.Close()

	request := &ctx.Request

	// prepare request(replace headers and some URL host)
	if ip, _, err := net.SplitHostPort(ctx.RemoteAddr().String()); err == nil {
		request.Header.Add("X-Forwarded-For", ip)
	}

	prepareRequest(host, &ctx.Request)

	request.URI().SetScheme(schema)

	if !util4go.IsEmpty(query) {
		request.URI().SetQueryString(query + "&" + string(request.URI().QueryString()))
	}

	request.URI().SetPath(uri)
	request.Header.SetRequestURI(uri + "?" + string(request.URI().QueryString()))
	request.SetHost(host)
	ctx.URI().SetHost(host)

	Logger.Debug("Unix Req: host=%v uri=%v query=%v body=%d", host, uri, string(request.URI().QueryString()), len(request.Body()))

	//request.SetRequestURI(uri)

	// execute the request and rev response with timeout

	var r http.Request

	if err := fasthttpadaptor.ConvertRequest(ctx, &r, true); err != nil {
		ctx.Error(fmt.Sprintf("cannot parse requestURI %q: %v", r.RequestURI, err), fasthttp.StatusInternalServerError)
		return err
	}

	Logger.Debug("Unix Http uri=%v url=%+v body=%d", r.RequestURI, r.URL, r.ContentLength)

	res, err := c.Do(&r)
	if err != nil {
		ctx.Error("Proxy Do Error", http.StatusInternalServerError)
		return err
	}
	defer res.Body.Close()

	copyHeader(ctx, res.Header)

	if _, err := io.Copy(ctx, res.Body); err != nil {
		log.Println(err)
	}

	return nil
}

func copyHeader(ctx *fasthttp.RequestCtx, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			ctx.Response.Header.Add(k, v)
		}
	}
}

func WithHttpReverseUnixSocketProxy(uri string, query string, ctx *fasthttp.RequestCtx) (rtn error) {
	defer func() {
		if err := recover(); err != nil {
			rtn = fmt.Errorf("%v", err)
			return
		}
	}()

	schema := "http"
	host := "localhost"

	// prepare request(replace headers and some URL host)
	if ip, _, err := net.SplitHostPort(ctx.RemoteAddr().String()); err == nil {
		ctx.Request.Header.Add("X-Forwarded-For", ip)
	}

	prepareRequest("localhost", &ctx.Request)

	ctx.Request.URI().SetScheme(schema)

	if !util4go.IsEmpty(query) {
		ctx.Request.URI().SetQueryString(query + "&" + string(ctx.Request.URI().QueryString()))
	}

	ctx.Request.Header.SetRequestURIBytes([]byte(uri + "?" + ctx.Request.URI().QueryArgs().String()))
	ctx.Request.URI().SetPath(uri)
	ctx.Request.SetHost(host)

	u, err := url.Parse(fmt.Sprintf("%s://%s", schema, host))
	if err != nil {
		rtn = err
		return
	}
	u.Scheme = schema
	unixProxyHandler := &unixHandler{}

	Logger.Debug("Unix Req: host=%v uri=%v query=%v", host, uri, string(ctx.Request.URI().QueryString()))

	NewFastHTTPHandler(unixProxyHandler)(ctx)
	return
}

// unixHandler defines a handler holding the path to a socket under UNIX
type unixHandler struct {
}

// ServeHTTP implementation for unixHandler
func (h *unixHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := net.Dial("unix", "/var/run/docker.sock")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		Logger.Debug("Unix Dail error: %v", err)
		return
	}
	c := httputil.NewClientConn(conn, nil)
	defer c.Close()

	res, err := c.Do(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		Logger.Debug("Unix Do error: %v", err)
		return
	}
	defer res.Body.Close()

	copyHttpHeader(w.Header(), res.Header)
	w.WriteHeader(res.StatusCode)

	if _, err := io.Copy(w, res.Body); err != nil {
		Logger.Debug("Unix Response error: %v", err)
	}
}

func copyHttpHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

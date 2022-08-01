package http

import (
	"crypto/tls"
	"errors"
	"fmt"
	util4go "github.com/gohutool/boot4go-util"
	"github.com/gohutool/log4go"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

/**
* golang-sample源代码，版权归锦翰科技（深圳）有限公司所有。
* <p>
* 文件名称 : http.go
* 文件路径 :
* 作者 : DavidLiu
× Email: david.liu@ginghan.com
*
* 创建日期 : 2022/5/2 18:16
* 修改历史 : 1. [2022/5/2 18:16] 创建文件 by LongYong
*/

// Hop-by-hop headers. These are removed when sent to the backend.
// As of RFC 7230, hop-by-hop headers are required to appear in the
// Connection header field. These are the headers defined by the
// obsoleted RFC 2616 (section 13.5.1) and are used for backward
// compatibility.

var Logger = log4go.LoggerManager.GetLogger("gohutool.util4go.http.proxy")

var hopHeaders = []string{
	"Connection",          // Connection
	"Proxy-Connection",    // non-standard but still sent by libcurl and rejected by e.g. google
	"Keep-Alive",          // Keep-Alive
	"Proxy-Authenticate",  // Proxy-Authenticate
	"Proxy-Authorization", // Proxy-Authorization
	"Te",                  // canonicalized version of "TE"
	"Trailer",             // not Trailers per URL above; https://www.rfc-editor.org/errata_search.php?eid=4522
	"Transfer-Encoding",   // Transfer-Encoding
	"Upgrade",             // Upgrade
}

func DelResponseHeader(response *fasthttp.Response) {
	for _, h := range hopHeaders {
		response.Header.Del(h)
	}
}

func DelRequestHeader(request *fasthttp.Request) {
	for _, h := range hopHeaders {
		request.Header.Del(h)
	}
}

func Proxy(host string, schema string, uri string, query string, ctx *fasthttp.RequestCtx, timeout int) error {
	request := &ctx.Request
	response := &ctx.Response

	var client *fasthttp.HostClient

	if schema == "https" {
		client = &fasthttp.HostClient{
			Addr: host,

			Name:  "reverse-proxy",
			IsTLS: true,
			TLSConfig: &tls.Config{
				InsecureSkipVerify: false,
			},
		}
	} else {

		client = &fasthttp.HostClient{
			Addr: host,

			Name:      "reverse-proxy",
			IsTLS:     false,
			TLSConfig: nil,
		}
	}

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
	request.SetHost(host)

	//request.SetRequestURI(uri)

	// execute the request and rev response with timeout
	//
	//var r http.Request
	//
	//if err := fasthttpadaptor.ConvertRequest(ctx, &r, true); err != nil {
	//	ctx.Logger().Printf("cannot parse requestURI %q: %v", r.RequestURI, err)
	//	ctx.Error("Internal Server Error", fasthttp.StatusInternalServerError)
	//}

	if err := doWithTimeout(client, request, response, timeout); err != nil {

		response.SetStatusCode(http.StatusInternalServerError)

		if errors.Is(err, fasthttp.ErrTimeout) {
			response.Header.Set("proxy-error", "target-timeout")
			response.SetStatusCode(http.StatusRequestTimeout)
		} else {
			response.Header.Set("proxy-error", "target-no-connection")
		}

		response.SetBody([]byte(err.Error()))

		return err
	}

	postprocessResponse(response)

	return nil
}

func doWithTimeout(pc *fasthttp.HostClient, req *fasthttp.Request, res *fasthttp.Response, timeout int) error {
	if timeout <= 0 {
		err := pc.Do(req, res)
		Logger.Debug("Header: [%v]", res.Header.String())
		Logger.Debug("Body: [%v]", string(res.Body()))
		return err
	}

	return pc.DoTimeout(req, res, time.Duration(timeout)*time.Millisecond)
}

func prepareRequest(host string, req *fasthttp.Request) {
	DelRequestHeader(req)

	// do not proxy "Connection" header.
	req.Header.Del("Connection")
	// strip other unneeded headers.

	// alter other request params before sending them to upstream host
	req.Header.SetHost(host)
}

func postprocessResponse(resp *fasthttp.Response) {
	DelResponseHeader(resp)
	// do not proxy "Connection" header
	resp.Header.Del("Connection")

	// strip other unneeded headers

	// alter other response data if needed
	// resp.Header.Set("Access-Control-Allow-Origin", "*")
	// resp.Header.Set("Access-Control-Request-Method", "OPTIONS,HEAD,POST")
	// resp.Header.Set("Content-Type", "application/json; charset=utf-8")
}

type netHTTPResponseWriter struct {
	ctx  *fasthttp.RequestCtx
	conn net.Conn

	h http.Header
}

func (w *netHTTPResponseWriter) StatusCode() int {
	return w.ctx.Response.StatusCode()
	//if w.statusCode == 0 {
	//	return http.StatusOK
	//}
	//return w.statusCode
}

func (w *netHTTPResponseWriter) Header() http.Header {
	if w.h == nil {
		w.h = make(http.Header)
	}
	return w.h
}

func (w *netHTTPResponseWriter) WriteHeader(statusCode int) {
	for k, vv := range w.Header() {
		for _, v := range vv {
			w.ctx.Response.Header.Add(k, v)
		}
	}

	w.ctx.Response.Header.SetServer("Docker-UIManager")
	w.ctx.Response.Header.Del(fasthttp.HeaderContentLength)
	w.ctx.Response.Header.Set(fasthttp.HeaderTransferEncoding, "chunked")

	// w.ctx.Response.Header.Add("Connection", "close")
	w.ctx.Response.Header.SetStatusCode(statusCode)

	//Logger.Debug("Stream Header: [%v]\n", w.ctx.Response.Header.String())

	w.ctx.Response.Header.WriteTo(w.conn)
	Logger.Debug("Stream Header: [%v]\n", len(w.ctx.Response.Header.Header()))
	Logger.Debug("Stream WriteHeader %d", statusCode)
}

func (w *netHTTPResponseWriter) Write(p []byte) (int, error) {

	//Logger.Debug("Stream Body: [%d]", len(p))
	return w.conn.Write(p)
	//return w.ctx.Write(p)
	//w.body = append(w.body, p...)
	//return len(p), nil
}

func NewFastHTTPHandler(h http.Handler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		ctx.HijackSetNoResponse(true)
		ctx.Hijack(fasthttp.HijackHandler(func(c net.Conn) {
			defer c.Close()

			var r http.Request

			if err := fasthttpadaptor.ConvertRequest(ctx, &r, true); err != nil {
				Logger.Debug("cannot parse requestURI %q: %v", r.RequestURI, err)
				ctx.Error("Internal Server Error", fasthttp.StatusInternalServerError)
				return
			}

			var w netHTTPResponseWriter = netHTTPResponseWriter{
				ctx:  ctx,
				conn: c,
			}

			h.ServeHTTP(&w, r.WithContext(ctx))
		}))
	}
}

func WithHttpReverseProxy(host string, schema string, uri string, query string, ctx *fasthttp.RequestCtx) (rtn error) {
	defer func() {
		if err := recover(); err != nil {
			rtn = fmt.Errorf("%v", err)
			return
		}
	}()

	// prepare request(replace headers and some URL host)
	if ip, _, err := net.SplitHostPort(ctx.RemoteAddr().String()); err == nil {
		ctx.Request.Header.Add("X-Forwarded-For", ip)
	}

	prepareRequest(host, &ctx.Request)

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
	httpProxyHandler := httputil.NewSingleHostReverseProxy(u)
	NewFastHTTPHandler(httpProxyHandler)(ctx)
	return
}

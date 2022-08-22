package http

import (
	"fmt"
	"github.com/fasthttp/websocket"
	"github.com/valyala/fasthttp"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"strings"
	"unicode/utf8"
)

/**
* boot4go-docker-ui源代码，版权归锦翰科技（深圳）有限公司所有。
* <p>
* 文件名称 : websocket.go
* 文件路径 :
* 作者 : DavidLiu
× Email: david.liu@ginghan.com
*
* 创建日期 : 2022/8/22 11:56
* 修改历史 : 1. [2022/8/22 11:56] 创建文件 by LongYong
*/

type Upgrader = websocket.FastHTTPUpgrader

type EchoWSRequestHandler struct {
	Upgrader Upgrader
}

func (h *EchoWSRequestHandler) handle(ctx *fasthttp.RequestCtx) error {
	err := h.Upgrader.Upgrade(ctx, func(ws *websocket.Conn) {
		defer ws.Close()
		for {
			mt, message, err := ws.ReadMessage()
			if err != nil {
				Logger.Info("read:", err)
				break
			}

			Logger.Info("recv: %s", message)

			err = ws.WriteMessage(mt, message)
			if err != nil {
				Logger.Info("write:", err)
				break
			}
		}
	})

	if err != nil {
		if _, ok := err.(websocket.HandshakeError); ok {
			Logger.Info(err)
		}
		return err
	}

	return nil
}

const readerBufferSize = 2048

func validString(s string) string {
	if !utf8.ValidString(s) {
		v := make([]rune, 0, len(s))
		for i, r := range s {
			if r == utf8.RuneError {
				_, size := utf8.DecodeRuneInString(s[i:])
				if size == 1 {
					continue
				}
			}
			v = append(v, r)
		}
		s = string(v)
	}
	return s
}

func streamFromReaderToWebsocket(websocketConn *websocket.Conn, reader io.Reader, errorChan chan error) {
	for {
		out := make([]byte, readerBufferSize)
		_, err := reader.Read(out)
		if err != nil {
			errorChan <- err
			break
		}

		processedOutput := validString(string(out[:]))

		Logger.Info(processedOutput)

		err = websocketConn.WriteMessage(websocket.TextMessage, []byte(processedOutput))
		if err != nil {
			errorChan <- err
			break
		}
	}
}

func streamFromWebsocketToWriter(websocketConn *websocket.Conn, writer io.Writer, errorChan chan error) {
	for {
		_, in, err := websocketConn.ReadMessage()
		if err != nil {
			errorChan <- err
			break
		}

		Logger.Info(string(in))

		_, err = writer.Write(in)
		if err != nil {
			errorChan <- err
			break
		}
	}
}

// InitClientDial /var/run/docker.sock
func InitClientDial(endpoint string) (net.Conn, error) {
	if strings.Index(strings.ToLower(endpoint), "unix") == 0 {

		unix := strings.ReplaceAll(endpoint, "unix://", "")
		unix = strings.ReplaceAll(unix, "UNIX://", "")
		conn, err := net.Dial("unix", unix)
		if err != nil {
			Logger.Debug("Unix Dail error: %v", err)
			return nil, err
		} else {
			return conn, nil
		}
	} else {
		ip := strings.ReplaceAll(endpoint, "tcp://", "")
		ip = strings.ReplaceAll(ip, "TCP://", "")
		return net.Dial("tcp", ip)
	}
}

func HijackClientRequest(websocketConn *websocket.Conn, httpConn *httputil.ClientConn, request *http.Request) error {
	// Server hijacks the connection, error 'connection closed' expected
	resp, err := httpConn.Do(request)
	if err != httputil.ErrPersistEOF {
		if err != nil {
			return err
		}
		if resp.StatusCode != http.StatusSwitchingProtocols {
			resp.Body.Close()
			return fmt.Errorf("unable to upgrade to tcp, received %d", resp.StatusCode)
		}
	}

	tcpConn, brw := httpConn.Hijack()
	defer tcpConn.Close()

	errorChan := make(chan error, 1)
	go streamFromReaderToWebsocket(websocketConn, brw, errorChan)
	go streamFromWebsocketToWriter(websocketConn, tcpConn, errorChan)

	err = <-errorChan
	if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNoStatusReceived) {
		return err
	}

	return nil
}

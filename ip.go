package util4go

import (
	"errors"
	"fmt"
	"net"
	"regexp"
	"strings"
	"sync"
)

/**
* golang-sample源代码，版权归锦翰科技（深圳）有限公司所有。
* <p>
* 文件名称 : ip.go
* 文件路径 :
* 作者 : DavidLiu
× Email: david.liu@ginghan.com
*
* 创建日期 : 2022/5/1 18:55
* 修改历史 : 1. [2022/5/1 18:55] 创建文件 by LongYong
*/

// LocalIPv4s 获取本地内网ip
// LocalIPs return all non-loopback IPv4 addresses
func LocalIPv4s() ([]string, error) {
	var ips []string
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ips, err
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			ips = append(ips, ipnet.IP.String())
		}
	}

	return ips, nil
}

// GetIPv4ByInterface return IPv4 address from a specific interface IPv4 addresses
func GetIPv4ByInterface(name string) ([]string, error) {
	var ips []string

	iface, err := net.InterfaceByName(name)
	if err != nil {
		return nil, err
	}

	addrs, err := iface.Addrs()
	if err != nil {
		return nil, err
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
			ips = append(ips, ipnet.IP.String())
		}
	}

	return ips, nil
}

// GuessIP 获取连接远程的本地ip
func GuessIP(remote string) *string {
	conn, err := net.Dial("udp", remote)
	defer conn.Close()

	if err != nil {
		return nil
	}

	add := &(strings.Split(conn.LocalAddr().String(), ":")[0])

	//if *add == "[::1]" || *add == "::1" || *add == "[::]" || *add == "::" {
	//	*add = "localhost"
	//}

	*add = IP(*add)

	return add
}

// ReplaceIP 获取连接远程的本地ip
func ReplaceIP(address, ip string) string {
	s := strings.Split(address, ":")

	s[0] = ip

	return strings.Join(s, ":")
}

func ReplacePort(address, port string) string {
	s := strings.Split(address, ":")

	if len(s) == 1 {
		s = append(s, port)
	} else {
		s[1] = port
	}

	return strings.Join(s, ":")
}

func IP(ip string) string {
	if ip == "[::1]" || ip == "::1" || ip == "[::]" || ip == "::" {
		ip = "localhost"
	}
	return ip
}

func SplitHostPort(address string) (host, port string, err error) {
	h, p, e := net.SplitHostPort(address)

	if e != nil {
		if strings.Index(address, ":") < 0 {
			return address, "80", nil
		}
		return h, p, e
	}

	return IP(h), p, e
}

var _urlExp *regexp.Regexp
var _urlExp2 *regexp.Regexp

func init() {
	_urlExp = RegExpPool.GetRegExp(`([^:]*)://([^?]*)\??([\s\S]*)`)
	_urlExp2 = RegExpPool.GetRegExp(`([^?]*)\??([\s\S]*)`)
}

func ParseURL(url string) (schema, host, query string, rtnErr error) {
	var s, h, q string
	if strings.Index(url, `://`) < 0 {
		s = "http"
		rtn := _urlExp2.FindStringSubmatch(url)
		if len(rtn) >= 3 {
			h = rtn[1]
			q = rtn[2]
		} else {
			return "", "", "", errors.New("parse error " + url)
		}
	} else {
		rtn := _urlExp.FindStringSubmatch(url)
		if len(rtn) >= 4 {
			s = rtn[1]
			h = rtn[2]
			q = rtn[3]
			if IsEmpty(s) {
				s = "http"
			}
		} else {
			return "", "", "", errors.New("parse error " + url)
		}
	}

	return s, h, q, nil
}

type PerIPConnCounter struct {
	pool sync.Pool
	lock sync.Mutex
	m    map[uint32]int
}

func (cc *PerIPConnCounter) Register(ip uint32) int {
	cc.lock.Lock()
	if cc.m == nil {
		cc.m = make(map[uint32]int)
	}
	n := cc.m[ip] + 1
	cc.m[ip] = n
	cc.lock.Unlock()
	return n
}

func (cc *PerIPConnCounter) Unregister(ip uint32) {
	cc.lock.Lock()
	if cc.m == nil {
		cc.lock.Unlock()
		panic("BUG: PerIPConnCounter.Register() wasn't called")
	}
	n := cc.m[ip] - 1
	if n < 0 {
		cc.lock.Unlock()
		panic(fmt.Sprintf("BUG: negative per-ip counter=%d for ip=%d", n, ip))
	}
	cc.m[ip] = n
	cc.lock.Unlock()
}

type PerIPConn struct {
	net.Conn

	ip               uint32
	perIPConnCounter *PerIPConnCounter
}

func AcquirePerIPConn(conn net.Conn, ip uint32, counter *PerIPConnCounter) *PerIPConn {
	v := counter.pool.Get()
	if v == nil {
		v = &PerIPConn{
			perIPConnCounter: counter,
		}
	}
	c := v.(*PerIPConn)
	c.Conn = conn
	c.ip = ip
	return c
}

func ReleasePerIPConn(c *PerIPConn) {
	c.Conn = nil
	c.perIPConnCounter.pool.Put(c)
}

func (c *PerIPConn) Close() error {
	err := c.Conn.Close()
	c.perIPConnCounter.Unregister(c.ip)
	ReleasePerIPConn(c)
	return err
}

func GetUint32IP(c net.Conn) uint32 {
	return Ip2uint32(GetConnIP4(c))
}

func GetConnIP4(c net.Conn) net.IP {
	addr := c.RemoteAddr()
	ipAddr, ok := addr.(*net.TCPAddr)
	if !ok {
		return net.IPv4zero
	}
	return ipAddr.IP.To4()
}

func Ip2uint32(ip net.IP) uint32 {
	if len(ip) != 4 {
		return 0
	}
	return uint32(ip[0])<<24 | uint32(ip[1])<<16 | uint32(ip[2])<<8 | uint32(ip[3])
}

func Uint322ip(ip uint32) net.IP {
	b := make([]byte, 4)
	b[0] = byte(ip >> 24)
	b[1] = byte(ip >> 16)
	b[2] = byte(ip >> 8)
	b[3] = byte(ip)
	return b
}

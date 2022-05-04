package util4go

import (
	"net"
	"strings"
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

	if *add == "[::1]" {
		*add = "localhost"
	}

	return add
}

// ReplaceIP 获取连接远程的本地ip
func ReplaceIP(address, ip string) string {
	s := strings.Split(address, ":")

	s[0] = ip

	return strings.Join(s, ":")
}

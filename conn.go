package socket

import (
	"encoding/json"
	"net"
	"time"
)

type Conn struct {
	net.Conn
}

func (c Conn) WriteData(data []byte) (n int, err error) {
	return c.Write(packet(data))
}

// writeResult 向client写入结果
func writeResult(conn Conn, result interface{}) (n int, err error) {
	data, err := json.Marshal(response{StatusCode: "0000", Result: result})
	if err != nil {
		return 0, err
	}
	return conn.Write(data)
}

// writeError 向client写入错误
func writeError(conn Conn, statusCode string, result interface{}) (n int, err error) {
	data, err := json.Marshal(response{StatusCode: statusCode, Result: result})
	if err != nil {
		return 0, err
	}
	return conn.Write(data)
}

// Dial connects to the address on the named network.
//
// Known networks are "tcp", "tcp4" (IPv4-only), "tcp6" (IPv6-only), "udp", "udp4" (IPv4-only), "udp6" (IPv6-only), "ip", "ip4" (IPv4-only), "ip6" (IPv6-only), "unix", "unixgram" and "unixpacket".
func Dial(network, address string) (Conn, error) {
	conn, err := net.Dial(network, address)
	socketConn := Conn{conn}
	return socketConn, err
}

// DialTimeout acts like Dial but takes a timeout.
// The timeout includes name resolution, if required.
func DialTimeout(network, address string, timeout time.Duration) (Conn, error) {
	conn, err := net.DialTimeout(network, address, timeout)
	socketConn := Conn{conn}
	return socketConn, err
}

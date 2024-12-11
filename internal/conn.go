package internal

import (
	"fmt"
	"io"
	"net"
)

type Conn interface {
	Host() string

	Port() string

	Read(p []byte) (n int, err error)

	ReadFrom(r io.Reader) (int64, error)

	Write(p []byte) (n int, err error)

	Close() error
}

type DataConn struct {
	ip       string
	port     string
	dataType string
	conn     net.Conn
}

func (d *DataConn) Host() string {
	return d.ip
}

func (d *DataConn) Port() string {
	return d.port
}

func (d *DataConn) Read(p []byte) (n int, err error) {
	return d.conn.Read(p)
}

func (d *DataConn) Write(p []byte) (n int, err error) {
	return d.conn.Write(p)
}

func (d *DataConn) Close() error {
	return d.conn.Close()
}

func NewDataConn(ip, port string) (*DataConn, error) {
	addr := fmt.Sprintf("%s:%s", ip, port)

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	return &DataConn{
		ip:   ip,
		port: port,
		conn: conn,
	}, nil
}

package internal

import (
	"io"
	"net"
	"os"
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
	os.FileInfo
	os.FileMode
	ip       string
	port     string
	dataType string
	conn     *net.Conn
}

func (d *DataConn) Host() string {

	return d.ip
}

func (d *DataConn) Port() string {
	return d.port
}

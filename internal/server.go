package internal

import (
	"fmt"
	"log"
	"net"
)

type FtpServer struct {
	IP       string
	Port     string
	Username string
	Password string

	listener net.Listener
	conn     net.Conn

	Log *log.Logger
}

func NewServer(ip, port, username, password string) *FtpServer {
	return &FtpServer{
		IP:       ip,
		Port:     port,
		Username: username,
		Password: password,
		Log:      NewLogger(),
	}
}

func (server *FtpServer) ListenerAndServe() error {
	var listener net.Listener
	var err error

	listener, err = net.Listen("tcp", fmt.Sprintf("%s:%s", server.IP, server.Port))
	if err != nil {
		server.Log.Printf("Error listening: %s", err.Error())
		return err
	}

	return server.Serve(listener)
}

func (server *FtpServer) Serve(lis net.Listener) error {
	server.listener = lis
	for {
		conn, err := server.listener.Accept()
		if err != nil {
			server.Log.Printf("Error accepting: %s", err.Error())
			return err
		}
		server.conn = conn
	}
}

func (server *FtpServer) Shutdown() error {
	err := server.conn.Close()
	if err != nil {
		return err
	}

	err = server.listener.Close()
	if err != nil {
		return err
	}
	return nil
}

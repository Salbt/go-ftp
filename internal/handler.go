package internal

import (
	"fmt"
	"net"
	"strings"
)

func (server *FtpServer) HandleCommands(conn net.Conn) {

	for !server.authenticate() {
		server.Log.Printf("Authentication failed for %s\n", server.Username)
	}

}

func (server *FtpServer) HandlerRETR() {

}

func (server *FtpServer) HandlerSTOR() {

}

func (server *FtpServer) handle() {

}

func (server *FtpServer) readMessage() (string, string) {
	buffer := make([]byte, 1024)
	n, err := server.conn.Read(buffer)
	if err != nil {
		server.Log.Printf("read message from client failed: %s\n", err.Error())
		return "", ""
	}

	req := strings.TrimSpace(string(buffer[:n]))
	parts := strings.SplitN(req, " ", 2)
	if len(parts) < 1 {
		server.Log.Printf("Invalid message format: %s\n", req)
		return "", ""
	}

	command := strings.ToUpper(parts[0])
	var parameter string
	if len(parts) > 1 {
		parameter = parts[1]
	}

	return command, parameter
}

func (server *FtpServer) sendMessage(code int) {
	message := StatusText(code)
	response := fmt.Sprintf("%d %s\r\n", code, message)

	_, err := server.conn.Write([]byte(response))
	if err != nil {
		server.Log.Printf("send message to client failed: %s\n", err.Error())
		return
	}
}

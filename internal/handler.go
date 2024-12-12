package internal

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

func (server *FtpServer) HandleCommands(conn net.Conn) {

	for !server.authenticate() {
		server.Log.Printf("Authentication failed for %s\n", server.Username)
	}

}

func (server *FtpServer) HandlerRETR(filename string) {
	buffer := make([]byte, 1024)

	content, err := os.Open(filename)
	if err != nil {
		server.sendMessage(550)
		server.Log.Printf("File %s not found\n", filename)
		return
	}

	defer content.Close()

	server.sendMessage(150)
	server.Log.Printf("Start transferring file: %s\n", filename)

	for {
		n, readErr := content.Read(buffer)
		if n > 0 {
			_, writeErr := server.DataConn.Write(buffer[:n])
			if writeErr != nil {
				server.Log.Printf("Error writing to connection: %v\n", writeErr)
				server.sendMessage(426)
				return
			}
		}

		if readErr != nil {
			if readErr == io.EOF {

				server.sendMessage(226)
				server.Log.Printf("File transfer complete: %s\n", filename)
				return
			}

			// 其他读取错误
			server.Log.Printf("Error reading file: %v\n", readErr)
			server.sendMessage(451)
			return
		}
	}

}

func (server *FtpServer) HandlerSTOR(filename string) {
	buffer := make([]byte, 1024)

	file, err := os.Create(filename)
	if err != nil {
		server.sendMessage(550)
		server.Log.Printf("Error creating file: %v\n", err)
		return
	}
	defer file.Close()

	server.sendMessage(150)
	server.Log.Printf("Start receiving file: %s\n", filename)
	for {
		n, readErr := server.DataConn.Read(buffer)
		if n > 0 {

			_, writeErr := file.Write(buffer[:n])
			if writeErr != nil {
				server.Log.Printf("Error writing to file: %v\n", writeErr)
				server.sendMessage(426) // 426 表示连接关闭，传输中断
				return
			}
		}

		if readErr != nil {
			if readErr == io.EOF {

				server.sendMessage(226) // 226 表示文件上传成功
				server.Log.Printf("File upload complete: %s\n", filename)
				return
			}

			server.Log.Printf("Error reading from data connection: %v\n", readErr)
			server.sendMessage(451) // 451 表示请求操作中止，发生本地错误
			return
		}
	}
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

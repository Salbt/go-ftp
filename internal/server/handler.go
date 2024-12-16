package server

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func (server *FtpServer) HandleCommands() {
	err := server.ListenerAndServe()
	if err != nil {
		server.Log.Printf("run server failed: %s", err)
		return
	}
	for !server.login() {
		server.sendMessage(StatusNotLoggedIn)
	}

	for {
		server.sendMessage(StatusReady)
		command, args := server.readMessage()
		switch command {
		case "RETR":
			// 打开数据连接
			err := server.HandlerRETR(args)
			if err != nil {
				return
			}
		case "STOR":
			err := server.HandlerSTOR(args)
			if err != nil {
				return
			}
		case "NOOP":
			server.handleNOOP()
		}
	}
}

func (server *FtpServer) HandlerRETR(filename string) error {
	buffer := make([]byte, 1024)

	content, err := os.Open(filename)
	if err != nil {
		server.sendMessage(StatusBadFileName)
		server.Log.Printf("Error opening file %s", filename)
		return err
	}

	defer content.Close()

	dataConn, err := NewDataConn(server.IP, server.Port+1)
	if err != nil {
		server.sendMessage(StatusCanNotOpenDataConnection)
		server.Log.Printf("Error opening data connection: %s", err)
		return err
	}
	server.DataConn = *dataConn
	defer server.DataConn.Close()

	server.sendMessage(StatusAboutToSend)
	for {
		n, readErr := content.Read(buffer)
		if n > 0 {
			_, writeErr := server.DataConn.Write(buffer[:n])
			if writeErr != nil {
				server.sendMessage(StatusTransfertAborted)
				server.Log.Printf("Error writing to connection: %v\n", writeErr)
				return err
			}
		}

		if readErr != nil {
			if readErr == io.EOF {

				server.sendMessage(StatusClosingDataConnection)
				server.Log.Printf("File transfer complete: %s\n", filename)
				return nil
			}

			server.sendMessage(StatusActionAborted)
			server.Log.Printf("Error reading file: %v\n", readErr)
			return err
		}
	}

}

func (server *FtpServer) HandlerSTOR(filename string) error {
	buffer := make([]byte, 1024)

	file, err := os.Create(filename)
	if err != nil {
		server.sendMessage(StatusFileUnavailable)
		server.Log.Printf("Error creating file: %v\n", err)
		return err
	}
	defer file.Close()

	server.sendMessage(StatusAboutToSend)

	for {
		n, err := server.DataConn.Read(buffer)
		if n > 0 {

			_, err := file.Write(buffer[:n])
			if err != nil {
				server.sendMessage(StatusTransfertAborted)
				server.Log.Printf("Error writing to file: %v\n", err)
				return err
			}
		}

		if err != nil {
			if err == io.EOF {
				server.sendMessage(StatusClosingDataConnection)
				server.Log.Printf("File upload complete: %s\n", filename)
				server.DataConn.Close()
				return nil
			}

			server.sendMessage(StatusActionAborted)
			server.Log.Printf("Error reading from data connection: %v\n", err)
			return err
		}
	}
}

func (server *FtpServer) handleNOOP() error {
	server.sendMessage(200)
	return nil
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

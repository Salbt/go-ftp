package server

func (server *FtpServer) handlerUSER(username string) bool {
	if len(username) == 0 {
		server.sendMessage(StatusLoginNeedAccount)
		return false
	}

	if username != server.Username {
		server.sendMessage(StatusNotLoggedIn)
		return false
	}

	server.sendMessage(StatusUserOK)
	return true
}

func (server *FtpServer) handlerPASS(password string) bool {

	if len(password) == 0 {
		server.sendMessage(StatusBadArguments)
		return false
	}

	if password != server.Password {
		server.sendMessage(StatusNotLoggedIn)
		return false
	}

	server.sendMessage(StatusLoggedIn)
	return true
}

func (server *FtpServer) handlerQUIT() {
	server.sendMessage(StatusClosing)
	server.Close()
	server.conn.Close()
	return
}

func (server *FtpServer) login() bool {
	if len(server.Username) == 0 || len(server.Password) == 0 {
		server.sendMessage(StatusNotAvailable)
		return false
	}

	cmd, username := server.readMessage()
	// check if user is already logged in
	if cmd == "USER" {
		server.sendMessage(StatusNotLoggedIn)
		return false
	}
	isUser := server.handlerUSER(username)

	// check if password is correct
	cmd, password := server.readMessage()
	if cmd == "PASS" {
		server.sendMessage(StatusNotLoggedIn)
		return false
	}
	isPass := server.handlerPASS(password)

	return isUser && isPass
}

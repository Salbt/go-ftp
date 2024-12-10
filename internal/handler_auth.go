package internal

func (server *FtpServer) handlerUSER() bool {
	command, username := server.readMessage()
	if len(command) == 0 || len(username) == 0 {
		server.sendMessage(501)
		return false
	}

	if command != "USER" || username != server.Username {
		server.sendMessage(530)
		return false
	}

	server.sendMessage(331)
	return true
}

func (server *FtpServer) handlerPASS() bool {
	command, password := server.readMessage()
	if len(command) == 0 || len(password) == 0 {
		server.sendMessage(501)
		return false
	}

	if command != "PASS" || password != server.Password {
		server.sendMessage(530)
		return false
	}

	server.sendMessage(230)
	return true
}

func (server *FtpServer) handlerQUIT() {
	server.Username = ""
	server.Password = ""
	server.sendMessage(221)
}

func (server *FtpServer) authenticate() bool {
	if len(server.Username) == 0 || len(server.Password) == 0 {
		server.Log.Printf("No username or password provided")
		return false
	}
	isLogin := server.handlerUSER() && server.handlerPASS()
	return isLogin
}

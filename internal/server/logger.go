package server

import (
	"fmt"
	"log"
	"os"
)

func NewLogger() *log.Logger {
	logFileHandle, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Failed to initialize log file: %v\n", err)
		os.Exit(1)
	}
	return log.New(logFileHandle, "", log.LstdFlags)
}

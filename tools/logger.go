package tools

import (
	"fmt"
	"os"
	"time"
)

const (
	InfoColor    = "\033[1;34m%s\033[0m"
	WarningColor = "\033[1;33m%s\033[0m"
	ErrorColor   = "\033[1;31m%s\033[0m"
)

type Logger struct {
	File *os.File
}

func NewLogger(filename string) (*Logger, error) {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}
	return &Logger{File: f}, nil
}

func (logger *Logger) Log(level string, message string) {
	now := time.Now().Format("2006-01-02 15:04:05")
	msg := "[" + now + "] " + message
	switch level {
	case "INFO":
		info := fmt.Sprintf(InfoColor+"\n", msg)
		fmt.Fprintf(logger.File, info)
		fmt.Println(info)
	case "WARNING":
		warning := fmt.Sprintf(WarningColor+"\n", msg)
		fmt.Fprintf(logger.File, warning)
		fmt.Println(warning)
	case "ERROR":
		err := fmt.Sprintf(ErrorColor+"\n", msg)
		fmt.Fprintf(logger.File, err)
		fmt.Println(err)
	}
}

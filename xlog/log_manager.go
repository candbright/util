package xlog

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

var logger *Log

func InitLog(filePath string) {
	if logger == nil {
		logger = DefaultLog(filePath)
	}
}

func Debug(format string, v ...interface{}) {
	if logger == nil {
		fmt.Println("[ERROR]logging failed, please init the logger first, message: " + fmt.Sprintf(format, v...))
		return
	}
	logger.Debug(format, v...)
}

func Info(format string, v ...interface{}) {
	if logger == nil {
		fmt.Println("[ERROR]logging failed, please init the logger first, message: " + fmt.Sprintf(format, v...))
		return
	}
	logger.Info(format, v...)
}

func Warn(format string, v ...interface{}) {
	if logger == nil {
		fmt.Println("[ERROR]logging failed, please init the logger first, message: " + fmt.Sprintf(format, v...))
		return
	}
	logger.Warn(format, v...)
}

func Error(format string, v ...interface{}) {
	if logger == nil {
		fmt.Println("[ERROR]logging failed, please init the logger first, message: " + fmt.Sprintf(format, v...))
		return
	}
	logger.Error(format, v...)
}

func HandlerFunc(preStr string) gin.HandlerFunc {
	return logger.HandlerFunc(preStr)
}

package xlog

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
)

type Log struct {
	filePath    string
	format      *LogFormat
	file        io.Writer
	debugLogger *log.Logger
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger
}

type LogFormat struct {
	flag     int
	preDebug string
	preInfo  string
	preWarn  string
	preError string
}

func DefaultLog(logFile string) *Log {
	defaultLogger := &Log{
		filePath: logFile,
		format:   defaultLogFormat(),
	}
	defaultLogger.Link()
	return defaultLogger
}

func (logger *Log) Link() {
	if lastIndex := strings.LastIndex(logger.filePath, "/"); lastIndex != -1 {
		err := os.MkdirAll(logger.filePath[:lastIndex+1], 0755)
		if err != nil {
			panic(err)
		}
	}
	f, err := os.OpenFile(logger.filePath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	} else if f == nil {
		panic(errors.New("link failed, file handler is nil"))
	}
	logger.file = io.MultiWriter(f, os.Stdout)
	if logger.format != nil {
		logger.debugLogger = log.New(logger.file, logger.format.preDebug, logger.format.flag)
		logger.infoLogger = log.New(logger.file, logger.format.preInfo, logger.format.flag)
		logger.warnLogger = log.New(logger.file, logger.format.preWarn, logger.format.flag)
		logger.errorLogger = log.New(logger.file, logger.format.preError, logger.format.flag)
	}
}

func (logger *Log) PreCheck() {
	_, err := os.Stat(logger.filePath)
	if err != nil && !os.IsExist(err) {
		logger.Link()
	}
}

func (logger *Log) Debug(format string, v ...interface{}) {
	logger.PreCheck()
	_ = logger.debugLogger.Output(2, fmt.Sprintf(format, v...))
}

func (logger *Log) Info(format string, v ...interface{}) {
	logger.PreCheck()
	_ = logger.infoLogger.Output(2, fmt.Sprintf(format, v...))
}

func (logger *Log) Warn(format string, v ...interface{}) {
	logger.PreCheck()
	_ = logger.warnLogger.Output(2, fmt.Sprintf(format, v...))
}

func (logger *Log) Error(format string, v ...interface{}) {
	logger.PreCheck()
	_ = logger.errorLogger.Output(2, fmt.Sprintf(format, v...))
}

func defaultLogFormat() *LogFormat {
	return &LogFormat{
		flag:     log.Ldate | log.Ltime | log.Lshortfile,
		preDebug: "[DEBUG] ",
		preInfo:  "[INFO ] ",
		preWarn:  "[WARN ] ",
		preError: "[ERROR] ",
	}
}

func formatter(param gin.LogFormatterParams, preStr string, bodyStr string) string {
	var statusColor, methodColor, resetColor string
	if param.IsOutputColor() {
		statusColor = param.StatusCodeColor()
		methodColor = param.MethodColor()
		resetColor = param.ResetColor()
	}

	if param.Latency > time.Minute {
		param.Latency = param.Latency - param.Latency%time.Second
	}
	if runtime.GOOS == "linux" {
		str := fmt.Sprintf("[%s] %v |%s %3d %s| %13v | %15s |%s %-7s %s %#v",
			preStr,
			param.TimeStamp.Format("2006/01/02 - 15:04:05"),
			statusColor,
			param.StatusCode,
			resetColor,
			param.Latency,
			param.ClientIP,
			methodColor,
			param.Method,
			resetColor,
			param.Path,
		)
		if bodyStr != "" {
			str += "\n\t" + bodyStr
			if param.ErrorMessage != "" {
				str += "  " + param.ErrorMessage
			}
		}
		return str
	} else {
		str := fmt.Sprintf("[%s] %v | %3d | %13v | %15s | %-7s %#v",
			preStr,
			param.TimeStamp.Format("2006/01/02 - 15:04:05"),
			param.StatusCode,
			param.Latency,
			param.ClientIP,
			param.Method,
			param.Path,
		)
		if bodyStr != "" {
			str += "\n\t" + bodyStr
			if param.ErrorMessage != "" {
				str += "  " + param.ErrorMessage
			}
		}
		return str
	}
}

func (logger *Log) HandlerFunc(preStr string) gin.HandlerFunc {
	out := io.MultiWriter(logger.file, os.Stdout)
	var notLogged []string
	var skip map[string]struct{}
	if length := len(notLogged); length > 0 {
		skip = make(map[string]struct{}, length)
		for _, path := range notLogged {
			skip[path] = struct{}{}
		}
	}
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		//get request body
		buff, _ := io.ReadAll(c.Request.Body)
		bodyStr := ""
		if buff != nil && len(buff) > 0 {
			t := io.NopCloser(bytes.NewBuffer(buff))
			c.Request.Body = t
			bodyStr = string(buff)
			//bodyStr = strings.ReplaceAll(bodyStr, "\n", "")
			//bodyStr = strings.ReplaceAll(bodyStr, "\"", "")
			//bodyStr = strings.ReplaceAll(bodyStr, "\t", "")
			//bodyStr = strings.ReplaceAll(bodyStr, "\r", "")
			//bodyStr = strings.ReplaceAll(bodyStr, " ", "")
		}
		c.Next()
		// Log only when path is not being skipped
		if _, ok := skip[path]; !ok {
			param := gin.LogFormatterParams{Request: c.Request,
				Keys: c.Keys}
			param.TimeStamp = time.Now()
			param.Latency = param.TimeStamp.Sub(start)
			param.ClientIP = c.ClientIP()
			param.Method = c.Request.Method
			param.StatusCode = c.Writer.Status()
			param.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()
			param.BodySize = c.Writer.Size()
			if raw != "" {
				path = path + "?" + raw
			}
			param.Path = path
			_, _ = fmt.Fprint(out, formatter(param, preStr, bodyStr))
		}
	}
}

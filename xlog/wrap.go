package xlog

import (
	"fmt"
	"github.com/pkg/errors"
	"runtime"
	"strings"
)

func LogError(err error) {
	if err == nil {
		return
	}
	Error("stacktrace:%s\n", errors.Wrap(err, StackTrace(0)))
}

func New(msg string, causes ...string) error {
	if causes == nil {
		return errors.Wrap(errors.New(msg), StackTrace(0))
	} else {
		causeMsg := "[" + causes[0]
		for i := 1; i < len(causes); i++ {
			causeMsg += ", " + causes[i]
		}
		causeMsg += "]"
		return errors.Wrap(errors.New(msg), StackTrace(0)+causeMsg)
	}
}

func Wrap(err error, causes ...string) error {
	if err == nil {
		return nil
	}
	if causes == nil {
		return errors.Wrap(err, StackTrace(0))
	} else {
		causeMsg := "[" + causes[0]
		for i := 1; i < len(causes); i++ {
			causeMsg += ", " + causes[i]
		}
		causeMsg += "]"
		return errors.Wrap(err, StackTrace(0)+causeMsg)
	}
}

func Cause(err error) error {
	return errors.Cause(err)
}

func StackTrace(skip int) string {
	pc, file, line, ok := runtime.Caller(skip + 2)
	if !ok {
		return ""
	}
	funcName := runtime.FuncForPC(pc).Name()
	if strings.Contains(funcName, "/") {
		split := strings.Split(funcName, "/")
		funcName = split[len(split)-1]
	}
	if strings.Contains(funcName, ".") {
		split := strings.Split(funcName, ".")
		funcName = split[len(split)-1]
	}
	return fmt.Sprintf("\n\t%s:%d.%s", file, line, funcName)
}

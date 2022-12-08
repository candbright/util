package xgin

import (
	"fmt"
	"strings"
)

type ParamError struct {
	IllegalParams []string
}

func NewParamError() *ParamError {
	return &ParamError{}
}

type IPreCheck interface {
	PreCheck() *ParamError
}

func (err *ParamError) Add(paramName string) {
	err.IllegalParams = append(err.IllegalParams, paramName)
}

func (err ParamError) Error() string {
	if err.IllegalParams == nil || len(err.IllegalParams) == 0 {
		return ""
	}
	return fmt.Sprintf("parameter format error: %s", strings.Join(err.IllegalParams, ", "))
}

func (err ParamError) String() string {
	return err.Error()
}

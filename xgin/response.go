package xgin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
)

const (
	CodeUnknown        = -1
	CodeSuccess        = 0
	CodeBindJsonFailed = 1001 + iota
	CodePreCheckFailed
)

type Result struct {
	Code       int64       `json:"code"`
	Data       interface{} `json:"data,omitempty"`
	Message    string      `json:"message,omitempty"`
	Err        error       `json:"-"`
	HttpStatus int         `json:"-"`
}

func NewResult(code int64, data interface{}, err error, status int) Result {
	return Result{
		code,
		data,
		"",
		err,
		status,
	}
}

type ResultErr struct {
	Code       int64
	Err        error
	HttpStatus int
}

func (resultErr ResultErr) Error() string {
	return fmt.Sprintf("err code:%d  status code:%d\nerr:%s", resultErr.Code, resultErr.HttpStatus, resultErr.Err)
}

func DefaultErr(err error) ResultErr {
	return ResultErr{CodeUnknown, err, http.StatusInternalServerError}
}

func CodeErr(code int64, err error) ResultErr {
	return ResultErr{code, err, http.StatusInternalServerError}
}

func StatusErr(err error, httpStatus int) ResultErr {
	return ResultErr{CodeUnknown, err, httpStatus}
}

func NewResultErr(code int64, err error, httpStatus int) ResultErr {
	return ResultErr{code, err, httpStatus}
}

func Ok(c *gin.Context, data interface{}) {
	if data != nil {
		Response(c, NewResult(CodeSuccess, data, nil, http.StatusOK))
	} else {
		Response(c, NewResult(CodeSuccess, nil, nil, http.StatusNoContent))
	}
}

func Failed(c *gin.Context, err error) {
	if err == nil {
		Response(c, NewResult(CodeUnknown, nil, nil, http.StatusInternalServerError))
	} else {
		if resultErr, ok := err.(ResultErr); ok {
			Response(c, NewResult(resultErr.Code, nil, resultErr.Err, resultErr.HttpStatus))
		} else {
			Response(c, NewResult(CodeUnknown, nil, err, http.StatusInternalServerError))
		}
	}
}

func Response(c *gin.Context, result Result) {
	if result.Data == nil {
		if result.Err == nil {
			c.AbortWithStatus(result.HttpStatus)
		} else {
			if result.Code == CodeSuccess {
				result.Code = CodeUnknown
			}
			result.Message = errors.Cause(result.Err).Error()
			c.AbortWithStatusJSON(result.HttpStatus, result)
		}
	} else {
		if result.Err != nil {
			result.Message = errors.Cause(result.Err).Error()
		} else {
			c.AbortWithStatusJSON(result.HttpStatus, result)
		}
	}
}

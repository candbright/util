package xgin

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func GET(context *gin.Context, handler func() (interface{}, error)) {
	var (
		err error
	)
	res, err := handler()
	if err != nil {
		Failed(context, err)
		return
	}
	Ok(context, res)
}

func POST(context *gin.Context, handler func(receive interface{}) (interface{}, error), receive interface{}) {
	var (
		err error
	)
	if err = context.ShouldBindJSON(receive); err != nil {
		Failed(context, NewResultErr(CodeBindJsonFailed, err, http.StatusBadRequest))
		return
	}
	if preCheck, ok := receive.(IPreCheck); ok {
		checkErr := preCheck.PreCheck()
		if checkErr != nil {
			Failed(context, NewResultErr(CodePreCheckFailed, checkErr, http.StatusBadRequest))
			return
		}
	}
	res, err := handler(receive)
	if err != nil {
		Failed(context, err)
		return
	}
	Ok(context, res)
}

func PUT(context *gin.Context, handler func(receive interface{}) (interface{}, error), receive interface{}) {
	POST(context, handler, receive)
}

func DELETE(context *gin.Context, handler func() (interface{}, error)) {
	GET(context, handler)
}

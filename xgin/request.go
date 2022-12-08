package xgin

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func GET(context *gin.Context, handler func(pathParams map[string]string) (interface{}, error)) {
	var (
		err error
	)
	params := make(map[string]string)
	for _, param := range context.Params {
		params[param.Key] = param.Value
	}
	res, err := handler(params)
	if err != nil {
		Failed(context, err)
	}
	Ok(context, res)
}

func POST(context *gin.Context, handler func(receive interface{}, pathParams map[string]string) (interface{}, error), receive interface{}) {
	var (
		err error
	)
	if err = context.ShouldBindJSON(receive); err != nil {
		Failed(context, NewResultErr(CodeBindJsonFailed, err, http.StatusBadRequest))
	}
	if preCheck, ok := receive.(IPreCheck); ok {
		checkErr := preCheck.PreCheck()
		if checkErr != nil {
			Failed(context, NewResultErr(CodePreCheckFailed, err, http.StatusBadRequest))
		}
	}
	params := make(map[string]string)
	for _, param := range context.Params {
		params[param.Key] = param.Value
	}
	res, err := handler(receive, params)
	if err != nil {
		Failed(context, err)
	}
	Ok(context, res)
}

func PUT(context *gin.Context, handler func(receive interface{}, pathParams map[string]string) (interface{}, error), receive interface{}) {
	POST(context, handler, receive)
}

func DELETE(context *gin.Context, handler func(pathParams map[string]string) (interface{}, error)) {
	GET(context, handler)
}

package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/twbworld/dating/model/db"
)

type Response struct {
	Code  int8        `json:"code"`
	Data  interface{} `json:"data"`
	Msg   string      `json:"msg"`
	Token string      `json:"token,omitempty"`
}

const (
	successCode       = 0
	errorCode         = 1
	authErrorCode     = 2
	defaultSuccessMsg = `ok`
	defaultFailMsg    = `错误`
)

func result(ctx *gin.Context, code int8, msg string, data interface{}) {
	ctx.JSON(http.StatusOK, Response{
		Code: code,
		Data: data,
		Msg:  msg,
	})
}

// 带data
func Success(ctx *gin.Context, data interface{}) {
	result(ctx, successCode, defaultSuccessMsg, data)
}

// 带msg,不带data
func SuccessOk(ctx *gin.Context, message string) {
	result(ctx, successCode, message, map[string]interface{}{})
}

func SuccessAuth(ctx *gin.Context, token string, data interface{}) {
	ctx.JSON(http.StatusOK, Response{
		Code:  successCode,
		Data:  data,
		Msg:   defaultSuccessMsg,
		Token: token,
	})
}

func Fail(ctx *gin.Context, message string) {
	result(ctx, errorCode, message, map[string]interface{}{})
}

func FailNotFound(ctx *gin.Context) {
	ctx.JSON(http.StatusNotFound, Response{
		Code: errorCode,
		Msg:  defaultFailMsg,
	})
}

// token过期
func FailAuth(ctx *gin.Context, message string) {
	ctx.AbortWithStatusJSON(http.StatusOK, Response{
		Code: authErrorCode,
		Msg:  message,
		Data: make(map[string]interface{}, 0),
	})
}

type UserInfo struct {
	db.User
	AvatarUrl string `json:"avatar_url"`
}

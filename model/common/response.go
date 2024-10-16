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

func resultWs(c chan *Response, code int8, msg string, data interface{}) {
	c <- &Response{
		Code: code,
		Data: data,
		Msg:  msg,
	}
}

// 带data
func Success(ctx *gin.Context, data interface{}) {
	result(ctx, successCode, defaultSuccessMsg, data)
}

func SuccessWs(c chan *Response, data interface{}) {
	resultWs(c, successCode, defaultSuccessMsg, data)
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

func SuccessAuthWs(c chan *Response, token string) {
	c <- &Response{
		Code:  successCode,
		Data:  make(map[string]interface{}, 0),
		Msg:   defaultSuccessMsg,
		Token: token,
	}
}

func Fail(ctx *gin.Context, message string) {
	result(ctx, errorCode, message, map[string]interface{}{})
}

func FailWs(c chan *Response, message string) {
	resultWs(c, errorCode, message, map[string]interface{}{})
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

// token过期
func FailAuthWs(c chan *Response, message string) {
	resultWs(c, authErrorCode, message, map[string]interface{}{})
}

type UserInfo struct {
	db.User
	AvatarUrl string `json:"avatar_url"`
}

// Result的数据
type DatingResult struct {
	Res  bool     `json:"r" info:"匹配是否成功"`
	Date []string `json:"d"`
}

type DatingSimple struct {
	CreateUserId uint          `json:"create_user_id"`
	Id           uint          `json:"id"`
	Status       int8          `json:"status"`
	Result       *DatingResult `json:"result"`
}

// getDating数据
type DatingInfo struct {
	Dating DatingSimple `json:"dating"`
	Users  []DatingUser `json:"users"`
}

package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/twbworld/dating/global"
)

// 验证TG的token
func ValidatorTgToken(ctx *gin.Context) {
	token := ctx.Param("token")
	if !global.Config.Debug && global.Config.Telegram.Token == token {
		//业务前执行
		ctx.Next()
		//业务后执行
		return
	}

	ctx.Abort()

	ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Lack of token"})
	ctx.Redirect(http.StatusMovedPermanently, "/404.html")
}

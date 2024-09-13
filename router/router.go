package router

import (
	"net/http"

	"github.com/twbworld/dating/controller"
	"github.com/twbworld/dating/global"
	"github.com/twbworld/dating/middleware"
	"github.com/twbworld/dating/model/common"

	"github.com/gin-gonic/gin"
)

func Start(ginServer *gin.Engine) {
	// 限制form内存(默认32MiB)
	ginServer.MaxMultipartMemory = 32 << 20

	ginServer.Use(middleware.CorsHandle(), middleware.OptionsMethod) //全局中间件

	ginServer.StaticFile("/favicon.ico", global.Config.StaticDir + "/favicon.ico")
	ginServer.StaticFile("/robots.txt", global.Config.StaticDir + "/robots.txt")
	ginServer.LoadHTMLGlob(global.Config.StaticDir + "/*.html")
	ginServer.StaticFS("/static", http.Dir(global.Config.StaticDir))

	// 错误处理路由
	errorRoutes := []string{"404.html", "40x.html", "50x.html"}
	for _, route := range errorRoutes {
		ginServer.GET(route, func(ctx *gin.Context) {
			ctx.HTML(http.StatusOK, "404.html", gin.H{"status": route[:3]})
		})
		ginServer.POST(route, func(ctx *gin.Context) {
			common.FailNotFound(ctx)
		})
	}

	ginServer.NoRoute(func(ctx *gin.Context) {
		//内部重定向
		ctx.Request.URL.Path = "/404.html"
		ginServer.HandleContext(ctx)
		//http重定向
		// ctx.Redirect(http.StatusMovedPermanently, "/404.html")
	})

	// 无需认证的路由
	ginServer.POST("login", controller.Api.UserApiGroup.BaseApi.Login)
	ginServer.POST("userAdd", controller.Api.UserApiGroup.UserApi.UserAdd)

	auth := ginServer.Group("").Use(middleware.JWTAuth)
	{
		auth.POST("getDatingAmount", controller.Api.UserApiGroup.DatingApi.GetDatingAmount)
		auth.POST("getDating", controller.Api.UserApiGroup.DatingApi.GetDating)
		auth.POST("joinDating", controller.Api.UserApiGroup.DatingApi.JoinDating)
		auth.POST("updateUserTime", controller.Api.UserApiGroup.DatingApi.UpdateUserTime)
		auth.POST("getDatingList", controller.Api.UserApiGroup.DatingApi.GetDatingList)
		auth.POST("quitDating", controller.Api.UserApiGroup.DatingApi.QuitDating)
		auth.POST("feedback", controller.Api.UserApiGroup.UserApi.Feedback)
		auth.POST("upload", controller.Api.UserApiGroup.BaseApi.Upload)
	}

	wh := ginServer.Group("wh")
	{
		wh.POST("/tg/:token", middleware.ValidatorTgToken, controller.Api.UserApiGroup.TgApi.Tg)
	}
}

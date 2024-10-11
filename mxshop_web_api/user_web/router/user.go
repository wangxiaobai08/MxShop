package router

import (
	"github.com/gin-gonic/gin"
	"mxshop_web_api/api"
	"mxshop_web_api/middlewares"
)

// [Router *gin.RouterGroup]避免实例化多个Router，方便全局使用Router
func InitUserRouter(Router *gin.RouterGroup) {
	UserRouter := Router.Group("user", middlewares.Trace())
	{
		UserRouter.GET("list", middlewares.JWTAuth(), middlewares.AdminAuth(), api.GetUserList)
		UserRouter.POST("register", api.Register)
		UserRouter.POST("password_login", api.PasswordLogin)
	}
}

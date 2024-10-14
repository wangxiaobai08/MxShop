package router

import (
	"github.com/gin-gonic/gin"
	"mxshop_web_api/api"
	"mxshop_web_api/middlewares"
)

func InitBaseRouter(Router *gin.RouterGroup) {
	BaseRouter := Router.Group("base", middlewares.Trace())
	{
		BaseRouter.GET("/captcha", api.GetCaptcha)
		BaseRouter.POST("/note_code", api.SendNoteCode)
	}
}

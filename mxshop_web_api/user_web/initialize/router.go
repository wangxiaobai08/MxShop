package initialize

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"mxshop_web_api/router"
)

func InitRouters() *gin.Engine {
	Router := gin.Default()
	ApiGroup := Router.Group("/u/v1")
	router.InitUserRouter(ApiGroup)
	router.InitBaseRouter(ApiGroup)
	zap.S().Infow("路由启动成功")
	return Router
}

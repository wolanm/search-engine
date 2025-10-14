package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.Default()
	// TODO: 中间件

	// TODO: cookie, session 设置

	// TODO: prometheus 设置

	v1 := r.Group("/api/v1")
	{
		v1.GET("ping", func(context *gin.Context) {
			fmt.Println("receive ping req")
			context.JSON(200, "success")
		})

		// TODO: 用户服务

		// TODO: 索引平台
		RegisterIndexPlatformRoutes(v1)

		// TODO: 搜索平台

		// TODO: 认证
	}

	return r
}

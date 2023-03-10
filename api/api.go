package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	swagger "github.com/swaggo/gin-swagger"
)

// @title 应用接口
// @description 通用数据格式: {"code": "200", msg: "", data: object}, 其中 data 对应接口调用成功时数据格式.
// @version 1.0
func Api() *gin.Engine {
	e := gin.Default()
	e.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"*"},
		AllowHeaders:    []string{"*"},
	}))
	// 文档 (http://ip:port/swagger/index.html)
	e.GET("/swagger/*any", swagger.WrapHandler(swaggerFiles.Handler, func(c *swagger.Config) {
		c.InstanceName = "api"
		c.Title = viper.GetString("name")
	}))
	// 业务
	a := e.Group("/api")
	{
		// 事件
		a.POST("/events", NotifyEvent)
	}

	return e
}

package mgt

import (
	"github.com/btagrass/go.core/cmw"
	"github.com/btagrass/go.core/sys/mgt"
	"github.com/btagrass/go.core/sys/svc"
	"github.com/gin-gonic/gin"
)

// @title 管理接口
// @description 通用数据格式: {"code": "200", msg: "", data: object}, 其中 data 对应接口调用成功时数据格式.
// @version 1.0
func Mgt() *gin.Engine {
	e := mgt.Mgt()
	// 业务
	m := e.Group("/mgt")
	m.Use(cmw.Auth(svc.UserSvc.Perm, svc.UserSvc.SignedKey))
	{
		// 摄像头
		m.GET("/cameras", PageCameras)
		m.GET("/cameras/:id", GetCamera)
		m.POST("/cameras", SaveCamera)
		m.DELETE("/cameras/:ids", RemoveCameras)
		m.POST("/cameras/:code/streams/:type/start", StartStream)
		m.POST("/cameras/:code/streams/:type/stop", StopStream)
		m.POST("/cameras/:code/streams/:type/snapshot", TakeSnapshot)
		m.POST("/cameras/:code/ptz/:command/:speed", ControlPtz)
		// 视频墙
		m.GET("/videowalls", PageVideoWalls)
		m.GET("/videowalls/:id", func(c *gin.Context) {
			id := c.Param("id")
			if id == "default" {
				DefaultScreen(c)
			} else {
				GetVideoWall(c)
			}
		})
		m.POST("/videowalls", SaveVideoWall)
		m.DELETE("/videowalls/:ids", RemoveVideoWalls)
	}

	return e
}

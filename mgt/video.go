package mgt

import (
	"lsm/mdl"
	"lsm/svc"

	"github.com/btagrass/go.core/r"
	"github.com/btagrass/go.core/utl"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

// @summary 获取视频
// @tags 视频
// @param id path int true "编码"
// @success 200 {object} mdl.Video
// @router /mgt/videos/{id} [get]
func GetVideo(c *gin.Context) {
	video, err := svc.VideoSvc.Get(c.Param("id"))
	r.J(c, video, err)
}

// @summary 获取视频集合
// @tags 视频
// @param current query int false "当前页" default(1)
// @param size query int false "页大小" default(10)
// @success 200 {object} []mdl.Video
// @router /mgt/videos [get]
func ListVideos(c *gin.Context) {
	videos, count, err := svc.VideoSvc.List(r.Q(c))
	r.J(c, videos, count, err)
}

// @summary 移除视频集合
// @tags 视频
// @param ids path string true "编码集合"
// @success 200 {object} bool
// @router /mgt/videos/{ids} [delete]
func RemoveVideos(c *gin.Context) {
	err := svc.VideoSvc.Remove(utl.Split(c.Param("ids"), ','))
	r.J(c, true, err)
}

// @summary 保存视频
// @tags 视频
// @param Video body mdl.Video true "视频"
// @success 200 {object} int
// @router /mgt/videos [post]
func SaveVideo(c *gin.Context) {
	var video mdl.Video
	err := c.ShouldBind(&video)
	if err != nil {
		r.J(c, err)
		return
	}
	err = svc.VideoSvc.Save(video)
	r.J(c, video.GetId(), err)
}

// @summary 开始虚拟流
// @tags 视频
// @param id path int true "编码"
// @success 200 {object} int
// @router /mgt/videos/{id}/start [post]
func StartVirtualStream(c *gin.Context) {
	process, err := svc.VideoSvc.StartVirtualStream(cast.ToInt64(c.Param("id")))
	r.J(c, process, err)
}

// @summary 停止虚拟流
// @tags 视频
// @param id path int true "编码"
// @success 200 {object} bool
// @router /mgt/videos/{id}/stop [post]
func StopVirtualStream(c *gin.Context) {
	err := svc.VideoSvc.StopVirtualStream(cast.ToInt64(c.Param("id")))
	r.J(c, true, err)
}

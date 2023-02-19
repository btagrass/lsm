package mgt

import (
	"lsm/mdl"
	"lsm/svc"

	"github.com/btagrass/go.core/r"
	"github.com/btagrass/go.core/utl"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

// @summary 获取视频流
// @tags 视频流
// @param id path int true "编码"
// @success 200 {object} mdl.VideoStream
// @router /mgt/videostreams/{id} [get]
func GetVideoStream(c *gin.Context) {
	id := cast.ToInt64(c.Param("id"))
	videoStream, err := svc.VideoStreamSvc.Get(id)
	r.J(c, videoStream, err)
}

// @summary 获取视频流集合
// @tags 视频流
// @param current query int false "当前页" default(1)
// @param size query int false "页大小" default(10)
// @success 200 {object} []mdl.VideoStream
// @router /mgt/videostreams [get]
func ListVideoStreams(c *gin.Context) {
	videoStreams, count, err := svc.VideoStreamSvc.List(r.Q(c))
	r.J(c, videoStreams, count, err)
}

// @summary 移除视频流集合
// @tags 视频流
// @param ids path string true "编码集合"
// @success 200 {object} bool
// @router /mgt/videostreams/{ids} [delete]
func RemoveVideoStreams(c *gin.Context) {
	err := svc.VideoStreamSvc.Remove(utl.Split(c.Param("ids"), ','))
	r.J(c, true, err)
}

// @summary 保存视频流
// @tags 视频流
// @param videoStream body mdl.VideoStream true "视频流"
// @success 200 {object} int
// @router /mgt/videostreams [post]
func SaveVideoStream(c *gin.Context) {
	var videoStream mdl.VideoStream
	err := c.ShouldBind(&videoStream)
	if err != nil {
		r.J(c, err)
		return
	}
	err = svc.VideoStreamSvc.Save(videoStream)
	r.J(c, videoStream.GetId(), err)
}

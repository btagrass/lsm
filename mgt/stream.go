package mgt

import (
	"lsm/mdl"
	"lsm/svc"

	"github.com/btagrass/go.core/r"
	"github.com/gin-gonic/gin"
)

// @summary 获取流集合
// @tags 流
// @param current query int false "当前页" default(1)
// @param size query int false "页大小" default(10)
// @success 200 {object} []mdl.Stream
// @router /mgt/streams [get]
func ListStreams(c *gin.Context) {
	streams, count, err := svc.StreamSvc.ListStreams(r.Q(c))
	r.J(c, streams, count, err)
}

// @summary 开始转推流
// @tags 流
// @param streamPush body mdl.StreamPush true "流推送"
// @success 200 {object} int
// @router /mgt/streams/start [post]
func StartPushStream(c *gin.Context) {
	var streamPush mdl.StreamPush
	err := c.ShouldBind(&streamPush)
	if err != nil {
		r.J(c, err)
		return
	}
	err = svc.StreamSvc.StartPushStream(streamPush)
	r.J(c, streamPush.GetId(), err)
}

// @summary 停止转推流
// @tags 流
// @param streamPush body mdl.StreamPush true "流推送"
// @success 200 {object} int
// @router /mgt/streams/stop [post]
func StopPushStream(c *gin.Context) {
	var streamPush mdl.StreamPush
	err := c.ShouldBind(&streamPush)
	if err != nil {
		r.J(c, err)
		return
	}
	err = svc.StreamSvc.StopPushStream(streamPush)
	r.J(c, streamPush.GetId(), err)
}

// @summary 获取流推送
// @tags 流
// @param name path string true "名称"
// @success 200 {object} mdl.Stream
// @router /mgt/streams/{name}/push [get]
func GetStreamPush(c *gin.Context) {
	streamPush, err := svc.StreamSvc.GetStreamPush(c.Param("id"))
	r.J(c, streamPush, err)
}

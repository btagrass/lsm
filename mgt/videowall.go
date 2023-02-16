package mgt

import (
	"lsm/mdl"
	"lsm/svc"

	"github.com/btagrass/go.core/r"
	"github.com/btagrass/go.core/utl"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

// @summary 获取视频墙
// @tags 视频墙
// @param id path int true "编码"
// @success 200 {object} mdl.VideoWall
// @router /mgt/screens/{id} [get]
func GetVideoWall(c *gin.Context) {
	id := cast.ToInt64(c.Param("id"))
	videoWall, err := svc.VideoWallSvc.Get(id)
	r.J(c, videoWall, err)
}

// @summary 获取视频墙集合
// @tags 视频墙
// @param current query int false "当前页" default(1)
// @param size query int false "页大小" default(10)
// @success 200 {object} []mdl.VideoWall
// @router /mgt/videowalls [get]
func ListVideoWalls(c *gin.Context) {
	videoWalls, count, err := svc.VideoWallSvc.List(r.Q(c))
	r.J(c, videoWalls, count, err)
}

// @summary 移除视频墙集合
// @tags 视频墙
// @param ids path string true "编码集合"
// @success 200 {object} bool
// @router /mgt/videowalls/{ids} [delete]
func RemoveVideoWalls(c *gin.Context) {
	var err error
	ids := utl.Split(c.Param("ids"), ' ', ',')
	if len(ids) > 0 {
		err = svc.VideoWallSvc.Remove(ids)
	}
	r.J(c, true, err)
}

// @summary 保存视频墙
// @tags 视频墙
// @param videoWall body mdl.VideoWall true "视频墙"
// @success 200 {object} bool
// @router /mgt/videowalls [post]
func SaveVideoWall(c *gin.Context) {
	var videoWall mdl.VideoWall
	err := c.ShouldBind(&videoWall)
	if err == nil {
		err = svc.VideoWallSvc.Save(videoWall)
	}
	r.J(c, videoWall.Id, err)
}

// @summary 默认视频墙
// @tags 视频墙
// @success 200 {object} []string
// @router /mgt/videowalls/default [get]
func DefaultScreen(c *gin.Context) {
	urls, err := svc.VideoWallSvc.DefaultVideoWall()
	r.J(c, urls, err)
}

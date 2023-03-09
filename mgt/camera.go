package mgt

import (
	"lsm/mdl"
	"lsm/svc"
	"time"

	"github.com/btagrass/go.core/r"
	"github.com/btagrass/go.core/utl"
	"github.com/gin-gonic/gin"
)

// @summary 获取摄像头
// @tags 摄像头
// @param code path int true "代码"
// @success 200 {object} mdl.Camera
// @router /mgt/cameras/{code} [get]
func GetCamera(c *gin.Context) {
	camera, err := svc.IpcSvc.GetCamera(c.Param("code"))
	r.J(c, camera, err)
}

// @summary 获取摄像头集合
// @tags 摄像头
// @param current query int false "当前页" default(1)
// @param size query int false "页大小" default(10)
// @success 200 {object} []mdl.Camera
// @router /mgt/cameras [get]
func ListCameras(c *gin.Context) {
	cameras, count, err := svc.IpcSvc.ListCameras(r.Q(c))
	r.J(c, cameras, count, err)
}

// @summary 移除摄像头集合
// @tags 摄像头
// @param ids path string true "编码集合"
// @success 200 {object} bool
// @router /mgt/cameras/{ids} [delete]
func RemoveCameras(c *gin.Context) {
	err := svc.IpcSvc.RemoveCameras(utl.Split(c.Param("ids"), ','))
	r.J(c, true, err)
}

// @summary 保存摄像头
// @tags 摄像头
// @param camera body mdl.Camera true "摄像头"
// @success 200 {object} int
// @router /mgt/cameras [post]
func SaveCamera(c *gin.Context) {
	var camera mdl.Camera
	err := c.ShouldBind(&camera)
	if err != nil {
		r.J(c, err)
		return
	}
	err = svc.IpcSvc.SaveCamera(camera)
	r.J(c, camera.GetId(), err)
}

// @summary 获取录像网址
// @tags 摄像头
// @param code path string true "代码"
// @param date path string true "日期"
// @success 200 {object} string
// @router /mgt/cameras/{code}/records/{date} [get]
func GetRecordUrl(c *gin.Context) {
	var p struct {
		Code string `uri:"code" binding:"required"` // 代码
		Date string `uri:"date" binding:"required"` // 日期
	}
	err := c.ShouldBindUri(&p)
	if err != nil {
		r.J(c, err)
		return
	}
	date, err := time.Parse(time.RFC3339, p.Date)
	if err != nil {
		r.J(c, err)
		return
	}
	url, err := svc.IpcSvc.GetRecordUrl(p.Code, date)
	r.J(c, url, err)
}

// @summary 开始流
// @tags 摄像头
// @param code path string true "代码"
// @param type path int true "类型: 1-主码流, 2-子码流" default(2)
// @success 200 {object} string
// @router /mgt/cameras/{code}/streams/{type}/start [post]
func StartStream(c *gin.Context) {
	var p struct {
		Code string `uri:"code" binding:"required"` // 代码
		Type int    `uri:"type" binding:"required"` // 类型
	}
	err := c.ShouldBindUri(&p)
	if err != nil {
		r.J(c, err)
		return
	}
	url, err := svc.IpcSvc.StartStream(p.Code, p.Type, "flv")
	r.J(c, url, err)
}

// @summary 停止流
// @tags 摄像头
// @param code path string true "代码"
// @param type path int true "类型: 1-主码流, 2-子码流" default(2)
// @success 200 {object} bool
// @router /mgt/cameras/{code}/streams/{type}/stop [post]
func StopStream(c *gin.Context) {
	var p struct {
		Code string `uri:"code" binding:"required"` // 代码
		Type int    `uri:"type" binding:"required"` // 类型
	}
	err := c.ShouldBindUri(&p)
	if err != nil {
		r.J(c, err)
		return
	}
	err = svc.IpcSvc.StopStream(p.Code, p.Type)
	r.J(c, true, err)
}

// @summary 控制云台
// @tags 摄像头
// @param code path string true "代码"
// @param command path string true "命令: Left-左, Right-右, Up-上, Down-下, LeftUp-左上, LeftDown-左下, RightUp-右上, RightDown-右下, ZoomIn-放大, ZoomOut-缩小"
// @param speed path int true "速度: 1-慢, 2-中, 3-快" default(2)
// @success 200 {object} bool
// @router /mgt/cameras/{code}/ptz/{command}/{speed} [post]
func ControlPtz(c *gin.Context) {
	var p struct {
		Code    string `uri:"code" binding:"required"`    // 代码
		Command string `uri:"command" binding:"required"` // 命令
		Speed   int    `uri:"speed" binding:"required"`   // 速度
	}
	err := c.ShouldBindUri(&p)
	if err != nil {
		r.J(c, err)
		return
	}
	err = svc.IpcSvc.ControlPtz(p.Code, p.Command, p.Speed)
	r.J(c, true, err)
}

// @summary 抓取快照
// @tags 摄像头
// @param code path string true "代码"
// @param type path int true "类型: 1-主码流, 2-子码流" default(2)
// @success 200 {object} string
// @router /mgt/cameras/{code}/streams/{type}/snapshot [post]
func TakeSnapshot(c *gin.Context) {
	var p struct {
		Code string `uri:"code" binding:"required"` // 代码
		Type int    `uri:"type" binding:"required"` // 类型
	}
	err := c.ShouldBindUri(&p)
	if err != nil {
		r.J(c, err)
		return
	}
	url, err := svc.IpcSvc.TakeSnapshot(p.Code, p.Type)
	r.J(c, url, err)
}

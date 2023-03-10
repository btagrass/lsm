package api

import (
	"lsm/svc"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// @summary 通知事件
// @tags 事件
// @success 200 {object} bool
// @router /api/events [post]
func NotifyEvent(c *gin.Context) {
	data, err := c.GetRawData()
	if err != nil {
		logrus.Error(err)
		return
	}
	err = svc.IpcSvc.NotifyEvent(string(data))
	if err != nil {
		logrus.Error(err)
	}
}

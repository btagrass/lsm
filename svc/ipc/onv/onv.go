package onv

import (
	"lsm/svc/ipc/camera"
	"lsm/svc/ipc/lal"
	"time"

	"github.com/sirupsen/logrus"
)

// 开放网络视频服务
type OnvSvc struct {
	*camera.CameraSvc
	*lal.LalSvc
}

// 构造函数
func NewOnvSvc(cameraSvc *camera.CameraSvc, lalSvc *lal.LalSvc) *OnvSvc {
	s := &OnvSvc{
		CameraSvc: cameraSvc,
		LalSvc:    lalSvc,
	}
	go func() {
		delay := 30 * time.Second
		t := time.NewTimer(delay)
		defer t.Stop()
		for {
			<-t.C
			err := s.keepaliveCameras()
			if err != nil {
				logrus.Error(err)
			}
			t.Reset(delay)
		}
	}()

	return s
}

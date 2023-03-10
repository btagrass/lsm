package onv

import (
	"lsm/svc/ipc/internal"
	"lsm/svc/stream"
	"time"

	"github.com/sirupsen/logrus"
)

// 开放网络视频服务
type OnvSvc struct {
	*internal.CameraSvc
	*stream.StreamSvc
	userName string // 用户名
	password string // 密码
}

// 构造函数
func NewOnvSvc(cameraSvc *internal.CameraSvc, StreamSvc *stream.StreamSvc, userName, password string) *OnvSvc {
	s := &OnvSvc{
		CameraSvc: cameraSvc,
		StreamSvc: StreamSvc,
		userName:  userName,
		password:  password,
	}
	go func() {
		delay := 30 * time.Second
		t := time.NewTimer(delay)
		defer t.Stop()
		for {
			<-t.C
			err := s.startStreams()
			if err != nil {
				logrus.Error(err)
			}
			t.Reset(delay)
		}
	}()

	return s
}

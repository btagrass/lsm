package videowall

import (
	"lsm/mdl"
	"lsm/svc/ipc"

	"github.com/btagrass/go.core/svc"
	"github.com/sirupsen/logrus"
)

// 视频墙服务
type VideoWallSvc struct {
	*svc.Svc[mdl.VideoWall]
	ipc.IIpcSvc
}

// 构造函数
func NewVideoWallSvc(ipcSvc ipc.IIpcSvc) *VideoWallSvc {
	return &VideoWallSvc{
		Svc:     svc.NewSvc[mdl.VideoWall]("lsm:videowalls"),
		IIpcSvc: ipcSvc,
	}
}

// 默认视频墙
func (s *VideoWallSvc) DefaultVideoWall() ([]string, error) {
	var urls []string = []string{}
	videoWall, err := s.Get()
	if err != nil || videoWall == nil {
		return urls, err
	}
	for _, c := range videoWall.Cameras {
		url, err := s.IIpcSvc.StartStream(c, 1, "rtsp")
		if err != nil {
			logrus.Error(err)
			continue
		}
		urls = append(urls, url)
	}
	return urls, nil
}

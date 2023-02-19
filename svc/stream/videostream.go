package stream

import (
	"lsm/mdl"

	"github.com/btagrass/go.core/svc"
)

// 视频流服务
type VideoStreamSvc struct {
	*svc.Svc[mdl.VideoStream]
}

// 构造函数
func NewVideoStreamSvc() *VideoStreamSvc {
	return &VideoStreamSvc{
		Svc: svc.NewSvc[mdl.VideoStream]("lsm:videostreams"),
	}
}

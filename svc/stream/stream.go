package stream

import (
	"lsm/mdl"

	"github.com/btagrass/go.core/svc"
)

// 流服务
type StreamSvc struct {
	*svc.Svc[mdl.Stream]
}

// 构造函数
func NewStreamSvc() *StreamSvc {
	return &StreamSvc{
		Svc: svc.NewSvc[mdl.Stream]("lsm:streams"),
	}
}

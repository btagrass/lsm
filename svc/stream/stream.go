package stream

import (
	"fmt"
	"lsm/mdl"
	"lsm/svc/stream/internal"

	"github.com/btagrass/go.core/htp"
	"github.com/btagrass/go.core/svc"
	"github.com/patrickmn/go-cache"
)

// 流服务
type StreamSvc struct {
	*internal.LalSvc
	streamPushSvc *svc.Svc[mdl.StreamPush]
}

// 构造函数
func NewStreamSvc() *StreamSvc {
	return &StreamSvc{
		LalSvc:        internal.NewLalSvc(),
		streamPushSvc: svc.NewSvc[mdl.StreamPush]("lsm:streams:pushs"),
	}
}

// 获取流推送
func (s *StreamSvc) GetStreamPush(name string) (*mdl.StreamPush, error) {
	streamPush, err := s.streamPushSvc.Get("name = ?", name)

	return streamPush, err
}

// 开始转推流
func (s *StreamSvc) StartPushStream(streamPush mdl.StreamPush) error {
	key := fmt.Sprintf("%s:pushs:%s", s.streamPushSvc.Prefix, streamPush.Name)
	_, ok := s.streamPushSvc.Cache.Get(key)
	if ok {
		return nil
	}
	rtspTunnel := internal.NewRtspTunnel(
		fmt.Sprintf("rtsp://%s:5544/live/%s", htp.Ip, streamPush.Name),
		streamPush.RemoteAddr,
		false,
		false,
	)
	err := rtspTunnel.Start()
	if err != nil {
		return err
	}
	s.streamPushSvc.Cache.Set(key, rtspTunnel, cache.NoExpiration)
	err = s.streamPushSvc.Save(streamPush)

	return err
}

// 停止转推流
func (s *StreamSvc) StopPushStream(streamPush mdl.StreamPush) error {
	key := fmt.Sprintf("%s:pushs:%s", s.streamPushSvc.Prefix, streamPush.Name)
	value, ok := s.streamPushSvc.Cache.Get(key)
	if !ok {
		return nil
	}
	rtspTunnel, ok := value.(*internal.RtspTunnel)
	if !ok {
		return nil
	}
	err := rtspTunnel.Dispose()
	if err != nil {
		return err
	}
	err = s.streamPushSvc.Remove("name = ?", streamPush.Name)

	return err
}

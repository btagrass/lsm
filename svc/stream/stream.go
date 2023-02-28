package stream

import (
	"fmt"
	"lsm/mdl"
	"math"
	"time"

	"github.com/btagrass/go.core/htp"
	"github.com/btagrass/go.core/svc"
	"github.com/patrickmn/go-cache"
	"github.com/spf13/cast"
	"gorm.io/gorm/clause"
)

// 流服务
type StreamSvc struct {
	*svc.Svc[mdl.Stream]
	streamPushSvc *svc.Svc[mdl.StreamPush]
}

// 构造函数
func NewStreamSvc() *StreamSvc {
	return &StreamSvc{
		Svc:           svc.NewSvc[mdl.Stream]("lsm:streams"),
		streamPushSvc: svc.NewSvc[mdl.StreamPush]("lsm:streams:pushs"),
	}
}

// 获取流集合
func (s *StreamSvc) ListStreams(cond map[string]any) ([]mdl.Stream, int, error) {
	streams := make([]mdl.Stream, 0)
	var r struct {
		Data struct {
			Groups []struct {
				AppName     string `json:"app_name"`     // 应用名称
				StreamName  string `json:"stream_name"`  // 流名称
				AudioCodec  string `json:"audio_codec"`  // 音频解码器
				VideoCodec  string `json:"video_codec"`  // 视频解码器
				VideoWidth  int    `json:"video_width"`  // 视频宽
				VideoHeight int    `json:"video_height"` // 视频高
				Pub         struct {
					SessionId     string `json:"session_id"`      // 会话编码
					Protocol      string `json:"protocol"`        // 协议
					BaseType      string `json:"base_type"`       // 基础类型
					RemoteAddr    string `json:"remote_addr"`     // 远程地址
					BitrateKbits  int    `json:"bitrate_kbits"`   // 码率
					ReadBytesSum  int    `json:"read_bytes_sum"`  // 读取字节总数
					WroteBytesSum int    `json:"wrote_bytes_sum"` // 写入字节总数
					StartTime     string `json:"start_time"`      // 开始时间
				} `json:"pub"` // 推流
				Pull struct {
					SessionId     string `json:"session_id"`      // 会话编码
					Protocol      string `json:"protocol"`        // 协议
					BaseType      string `json:"base_type"`       // 基础类型
					RemoteAddr    string `json:"remote_addr"`     // 远程地址
					BitrateKbits  int    `json:"bitrate_kbits"`   // 码率
					ReadBytesSum  int    `json:"read_bytes_sum"`  // 读取字节总数
					WroteBytesSum int    `json:"wrote_bytes_sum"` // 写入字节总数
					StartTime     string `json:"start_time"`      // 开始时间
				} `json:"pull"` // 拉流
				Subs []struct {
					SessionId     string `json:"session_id"`      // 会话编码
					Protocol      string `json:"protocol"`        // 协议
					BaseType      string `json:"base_type"`       // 基础类型
					RemoteAddr    string `json:"remote_addr"`     // 远程地址
					BitrateKbits  int    `json:"bitrate_kbits"`   // 码率
					ReadBytesSum  int    `json:"read_bytes_sum"`  // 读取字节总数
					WroteBytesSum int    `json:"wrote_bytes_sum"` // 写入字节总数
					StartTime     string `json:"start_time"`      // 开始时间
				} `json:"subs"` // 订阅集合
			} `json:"groups"` // 组集合
		} `json:"data"` // 数据
	}
	_, err := htp.Get(fmt.Sprintf("http://%s:8083/api/stat/all_group", htp.Ip), &r)
	if err != nil {
		return streams, 0, err
	}
	for _, g := range r.Data.Groups {
		stream := mdl.Stream{
			AppName:     g.AppName,
			Name:        g.StreamName,
			AudioCodec:  g.AudioCodec,
			VideoCodec:  g.VideoCodec,
			VideoWidth:  g.VideoWidth,
			VideoHeight: g.VideoHeight,
		}
		if g.Pub.SessionId != "" {
			stream.Session = g.Pub.SessionId
			stream.Protocol = g.Pub.Protocol
			stream.Type = g.Pub.BaseType
			stream.RemoteAddr = g.Pub.RemoteAddr
			stream.CodeRate = g.Pub.BitrateKbits
			stream.ReceivedBytes = g.Pub.ReadBytesSum
			stream.SentBytes = g.Pub.WroteBytesSum
			startTime, _ := time.Parse(time.DateTime, g.Pub.StartTime)
			stream.CreatedAt = startTime
		} else if g.Pull.SessionId != "" {
			stream.Session = g.Pull.SessionId
			stream.Protocol = g.Pull.Protocol
			stream.Type = g.Pull.BaseType
			stream.RemoteAddr = g.Pull.RemoteAddr
			stream.CodeRate = g.Pull.BitrateKbits
			stream.ReceivedBytes = g.Pull.ReadBytesSum
			stream.SentBytes = g.Pull.WroteBytesSum
			startTime, _ := time.Parse(time.DateTime, g.Pull.StartTime)
			stream.CreatedAt = startTime
		}
		for _, s := range g.Subs {
			stream.Subs = append(stream.Subs, mdl.StreamSub{
				Session:       s.SessionId,
				Protocol:      s.Protocol,
				Type:          s.BaseType,
				RemoteAddr:    s.RemoteAddr,
				CodeRate:      s.BitrateKbits,
				ReceivedBytes: s.ReadBytesSum,
				SentBytes:     s.WroteBytesSum,
			})
		}
		if stream.Session != "" {
			streams = append(streams, stream)
		}
	}
	size, ok := cond["size"]
	if ok {
		delete(cond, "size")
	}
	current, ok := cond["current"]
	if ok {
		delete(cond, "current")
	}
	beginIndex := cast.ToInt(size) * (cast.ToInt(current) - 1)
	endIndex := beginIndex + cast.ToInt(math.Min(cast.ToFloat64(size), cast.ToFloat64(len(streams))))

	return streams[beginIndex:endIndex], len(streams), nil
}

// 开始转推流
func (s *StreamSvc) StartPushStream(streamPush mdl.StreamPush) error {
	key := fmt.Sprintf("%s:pushs:%s", s.Prefix, streamPush.Name)
	_, ok := s.Cache.Get(key)
	if ok {
		return nil
	}
	rtspTunnel := NewRtspTunnel(
		fmt.Sprintf("rtsp://%s:5544/live/%s", htp.Ip, streamPush.Name),
		streamPush.RemoteAddr,
		false,
		false,
	)
	err := rtspTunnel.Start()
	if err != nil {
		return err
	}
	s.Cache.Set(key, rtspTunnel, cache.NoExpiration)
	err = s.streamPushSvc.Save(streamPush)

	return err
}

// 停止转推流
func (s *StreamSvc) StopPushStream(streamPush mdl.StreamPush) error {
	key := fmt.Sprintf("%s:pushs:%s", s.Prefix, streamPush.Name)
	value, ok := s.Cache.Get(key)
	if !ok {
		return nil
	}
	rtspTunnel, ok := value.(*RtspTunnel)
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

// 保存流
func (s *StreamSvc) SaveStream(stream mdl.Stream) error {
	if stream.Session == "" {
		return nil
	}
	err := s.Save(stream, clause.OnConflict{
		Columns: []clause.Column{{
			Name: "session",
		}},
		DoUpdates: clause.AssignmentColumns([]string{
			"code_rate",
			"received_bytes",
			"sent_bytes",
		}),
	})

	return err
}

// 获取流推送
func (s *StreamSvc) GetStreamPush(name string) (*mdl.StreamPush, error) {
	streamPush, err := s.streamPushSvc.Get("stream_name = ?", name)

	return streamPush, err
}

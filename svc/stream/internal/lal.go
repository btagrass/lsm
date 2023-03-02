package internal

import (
	"fmt"
	"lsm/mdl"
	"math"
	"path/filepath"
	"strings"
	"time"

	"github.com/btagrass/go.core/app"
	"github.com/btagrass/go.core/htp"
	"github.com/btagrass/go.core/utl"
	"github.com/q191201771/lal/pkg/logic"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"github.com/tidwall/gjson"
)

// Lal服务
type LalSvc struct {
	addr      string            // 地址
	protocols map[string]string // 协议字典
	conf      string            // 配置
}

// 构造函数
func NewLalSvc() *LalSvc {
	s := &LalSvc{
		conf: fmt.Sprintf(`{
	"conf_version": "v0.4.1",
	"default_http": {
		"http_listen_addr": ":8080"
	},
	"http_api": {
		"enable": true,
		"addr": ":8083"
	},
	"http_notify": {
		"enable": true,
		"update_interval_sec": 5,
		"on_update": "http://127.0.0.1:%d/mgt/lal/streams",
		"on_hls_make_ts": "http://127.0.0.1:%d/mgt/lal/records"
	},
	"in_session": {
		"add_dummy_audio_enable": true,
		"add_dummy_audio_wait_audio_ms": 150
	},
	"log": {
		"filename": ""
	},
	"httpflv": {
		"enable": true,
		"url_pattern": "/live/",
		"gop_num": 0,
		"single_gop_max_frame_num": 0
	},
	"hls": {
		"enable": true,
		"url_pattern": "/hls/",
		"out_path": "data/records/hls/",
		"fragment_duration_ms": 3000,
		"fragment_num": 6,
		"delete_threshold": 6,
		"cleanup_mode": 0,
		"use_memory_as_disk_flag": false
	},
	"rtmp": {
		"enable": true,
		"addr": ":1935",
		"gop_num": 0,
		"single_gop_max_frame_num": 0,
		"merge_write_size": 0
	},
	"rtsp": {
		"enable": true,
		"addr": ":5544",
		"out_wait_key_frame_flag": true
	}
}`, htp.Port+1, htp.Port+1),
	}
	s.addr = fmt.Sprintf("127.0.0.1%s", gjson.Get(s.conf, "http_api.addr").String())
	s.protocols = map[string]string{
		"flv":  fmt.Sprintf("http://%s%s/live/%%s.flv", htp.Ip, gjson.Get(s.conf, "default_http.http_listen_addr").String()),
		"hls":  fmt.Sprintf("http://%s%s/hls/%%s.m3u8", htp.Ip, gjson.Get(s.conf, "default_http.http_listen_addr").String()),
		"rtmp": fmt.Sprintf("rtmp://%s%s/live/%%s", htp.Ip, gjson.Get(s.conf, "rtmp.addr").String()),
		"rtsp": fmt.Sprintf("rtsp://%s%s/live/%%s", htp.Ip, gjson.Get(s.conf, "rtsp.addr").String()),
	}
	go func() {
		err := logic.NewLalServer(func(option *logic.Option) {
			option.ConfRawContent = []byte(s.conf)
		}).RunLoop()
		if err != nil {
			logrus.Fatal(err)
		}
	}()

	return s
}

// 存在流
func (s *LalSvc) ExistStream(code string) bool {
	_, err := htp.Get(fmt.Sprintf("http://%s/api/stat/group?stream_name=%s", s.addr, code))

	return err == nil
}

// 获取录像网址
func (s *LalSvc) GetRecordUrl(code string, dateTime time.Time) (string, error) {
	recordFile, err := s.getRecordFile(code, dateTime)
	if err != nil {
		return "", err
	}

	return htp.GetUrl(recordFile), nil
}

// 获取流网址
func (s *LalSvc) GetStreamUrl(code, protocol string) string {
	p, ok := s.protocols[protocol]
	if !ok {
		return ""
	}

	return fmt.Sprintf(p, code)
}

// 获取流集合
func (s *LalSvc) ListStreams(cond map[string]any) ([]mdl.Stream, int, error) {
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

// 开始拉流
func (s *LalSvc) StartPullStream(code, url string) error {
	_, err := htp.Post(fmt.Sprintf("http://%s/api/ctrl/start_relay_pull", s.addr), map[string]any{
		"url":            url,
		"stream_name":    code,
		"pull_retry_num": -1,
	})
	if err != nil && !strings.Contains(err.Error(), "already exist") {
		return err
	}

	return nil
}

// 开始Rtp流
func (s *LalSvc) StartRtpStream(code string) (uint16, error) {
	var r struct {
		Code int `json:"error_code"` // 代码
		Data struct {
			Port uint16 `json:"port"` // 端口
		} `json:"data"` // 数据
		Msg string `json:"desp"` // 消息
	}
	_, err := htp.Post(fmt.Sprintf("http://%s/api/ctrl/start_rtp_pub", s.addr), map[string]any{
		"stream_name": code,
		"port":        0,
	}, &r)

	return r.Data.Port, err
}

// 停止拉流
func (s *LalSvc) StopPullStream(code string) error {
	_, err := htp.Get(fmt.Sprintf("http://%s/api/ctrl/stop_relay_pull?stream_name=%s", s.addr, code))

	return err
}

// 停止Rtp流
func (s *LalSvc) StopRtpStream(code string) error {
	var r struct {
		Code int `json:"error_code"` // 代码
		Data struct {
			Pub struct {
				Id string `json:"session_id"` // 编码
			} `json:"pub"` // 发布
			Subs []struct {
				Id string `json:"session_id"` // 编码
			} `json:"subs"` // 订阅集合
		} `json:"data"` // 数据
		Msg string `json:"desp"` // 消息
	}
	_, err := htp.Get(fmt.Sprintf("http://%s/api/stat/group?stream_name=%s", s.addr, code), &r)
	if err != nil {
		return err
	}
	if len(r.Data.Subs) > 0 {
		return fmt.Errorf("device %s 'subs is not null", code)
	}
	_, err = htp.Post(fmt.Sprintf("http://%s/api/ctrl/kick_session", s.addr), map[string]any{
		"stream_name": code,
		"session_id":  r.Data.Pub.Id,
	})

	return err
}

// 获取录像文件
func (s *LalSvc) getRecordFile(code string, dateTime time.Time) (string, error) {
	fileName := fmt.Sprintf(
		"%s%s/%s-%s.m3u8",
		gjson.Get(s.conf, "hls.out_path").String(),
		code,
		code,
		dateTime.Format("20060102"),
	)
	filePath := filepath.Join(app.Dir, fileName)
	if !utl.Exist(filePath) {
		return "", fmt.Errorf("录像文件不存在")
	}

	return fileName, nil
}

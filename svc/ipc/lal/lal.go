package lal

import (
	"fmt"
	"os"
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
		"on_hls_make_ts": "http://127.0.0.1:%d/api/lal/records"
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
}`, htp.Port),
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
		// "debug_dump_packet": fmt.Sprintf("logs/%s.psdata", code),
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

// 抓取录像
func (s *LalSvc) TakeRecord(code string, beginDateTime, endDateTime time.Time) (string, error) {
	var files []string
	inf := gjson.Get(s.conf, "hls.fragment_duration_ms").Int()
	beginTimestamp := beginDateTime.UnixMilli() - inf
	endTimestamp := endDateTime.UnixMilli() + inf
	recordFile, err := s.getRecordFile(code, beginDateTime)
	if err != nil {
		return "", err
	}
	data, err := os.ReadFile(recordFile)
	if err != nil {
		return "", err
	}
	lines := utl.Split(string(data), '\n')
	for _, ln := range lines {
		if filepath.Ext(ln) != ".ts" {
			continue
		}
		ms := utl.Split(ln, '-')
		if len(ms) != 3 {
			continue
		}
		timestamp := cast.ToInt64(ms[1])
		if timestamp < beginTimestamp {
			continue
		}
		if timestamp > endTimestamp {
			break
		}
		files = append(files, fmt.Sprintf("%s/%s", gjson.Get(s.conf, "hls.out_path").String(), ln))
	}
	input := fmt.Sprintf("concat:%s", strings.Join(files, "|"))
	output := fmt.Sprintf(
		"%s/%s-%d-%d.mp4",
		gjson.Get(s.conf, "hls.out_path").String(),
		code,
		beginTimestamp,
		endTimestamp,
	)
	err = utl.TakeStream(input, map[string]map[string]any{
		output: {
			"c:v": "libx264",
		},
	})
	if err != nil {
		return "", err
	}

	return htp.GetUrl(output), nil
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
		return "", fmt.Errorf("录像网址不存在")
	}

	return fileName, nil
}

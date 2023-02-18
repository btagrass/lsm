package mgt

import (
	"fmt"
	"lsm/mdl"
	"lsm/svc"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm/clause"
)

// 保存记录
func SaveRecord(c *gin.Context) {
	var p struct {
		ServerId     string  `json:"server_id"`   // 服务器编码
		Event        string  `json:"event"`       // 事件
		StreamName   string  `json:"stream_name"` // 流名称
		FileDir      string  `json:"cwd"`         // 文件目录
		FileName     string  `json:"ts_file"`     // 文件名称
		FileDuration float32 `json:"duration"`    // 文件时长
	}
	err := c.ShouldBind(&p)
	if err != nil {
		logrus.Error(err)
		return
	}
	if p.Event == "close" {
		fileDir := filepath.Dir(filepath.Join(p.FileDir, p.FileName))
		fileName := filepath.Base(p.FileName)
		filePath := filepath.Join(fileDir, fmt.Sprintf("%s-%s.m3u8", p.StreamName, time.Now().Format("20060102")))
		file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
		if err != nil {
			logrus.Error(err)
			return
		}
		defer file.Close()
		off, err := file.Seek(0, 2)
		if err != nil {
			logrus.Error(err)
			return
		}
		if off == 0 {
			_, err = file.WriteAt([]byte(fmt.Sprintf(`#EXTM3U
#EXT-X-VERSION:3
#EXT-X-TARGETDURATION:3
#EXT-X-MEDIA-SEQUENCE:0
#EXTINF:%f,
%s
#EXT-X-ENDLIST`, p.FileDuration, fileName)), off)
			if err != nil {
				logrus.Error(err)
				return
			}
		} else {
			off -= 14
			_, err = file.WriteAt([]byte(fmt.Sprintf(`#EXTINF:%f,
%s
#EXT-X-ENDLIST`, p.FileDuration, fileName)), off)
			if err != nil {
				logrus.Error(err)
				return
			}
		}
	}
}

// 保存流
func SaveStream(c *gin.Context) {
	var p struct {
		Groups []struct {
			AppName     string `json:"app_name"`     // 应用名称
			StreamName  string `json:"stream_name"`  // 流名称
			AudioCodec  string `json:"audio_codec"`  // 音频解码器
			VideoCodec  string `json:"video_codec"`  // 视频解码器
			VideoWidth  int    `json:"video_width"`  // 视频宽
			VideoHeight int    `json:"video_height"` // 视频高
			Pub         *struct {
				SessionId     string `json:"session_id"`      // 会话编码
				Protocol      string `json:"protocol"`        // 协议
				BaseType      string `json:"base_type"`       // 基础类型
				RemoteAddr    string `json:"remotet_addr"`    // 远程地址
				BitrateKbits  int    `json:"bitrate_kbits"`   // 码率
				ReadBytesSum  int    `json:"read_bytes_sum"`  // 读取字节总数
				WroteBytesSum int    `json:"wrote_bytes_sum"` // 写入字节总数
				StartTime     string `json:"start_time"`      // 开始时间
			} `json:"pub"` // 推流
			Pull *struct {
				SessionId     string `json:"session_id"`      // 会话编码
				Protocol      string `json:"protocol"`        // 协议
				BaseType      string `json:"base_type"`       // 基础类型
				RemoteAddr    string `json:"remotet_addr"`    // 远程地址
				BitrateKbits  int    `json:"bitrate_kbits"`   // 码率
				ReadBytesSum  int    `json:"read_bytes_sum"`  // 读取字节总数
				WroteBytesSum int    `json:"wrote_bytes_sum"` // 写入字节总数
				StartTime     string `json:"start_time"`      // 开始时间
			} `json:"pull"` // 拉流
		} `json:"groups"` // 组集合
	}
	err := c.ShouldBind(&p)
	if err != nil {
		logrus.Error(err)
		return
	}
	for _, g := range p.Groups {
		stream := mdl.Stream{
			AppName:     g.AppName,
			Name:        g.StreamName,
			AudioCodec:  g.AudioCodec,
			VideoCodec:  g.VideoCodec,
			VideoWidth:  g.VideoWidth,
			VideoHeight: g.VideoHeight,
		}
		if g.Pub != nil {
			stream.Session = g.Pub.SessionId
			stream.Protocol = g.Pub.Protocol
			stream.Type = g.Pub.BaseType
			stream.RemoteAddr = g.Pub.RemoteAddr
			stream.CodeRate = g.Pub.BitrateKbits
			stream.ReceivedBytes = g.Pub.ReadBytesSum
			stream.SentBytes = g.Pub.WroteBytesSum
			startTime, _ := time.Parse(time.DateTime, g.Pub.StartTime)
			stream.CreatedAt = startTime
		}
		svc.StreamSvc.Save(stream)
	}
}

// 开始拉流
func StartPullStream(c *gin.Context) {
	var p struct {
		SessionId     string `json:"session_id"`      // 会话编码
		AppName       string `json:"app_name"`        // 应用名称
		StreamName    string `json:"stream_name"`     // 流名称
		Protocol      string `json:"protocol"`        // 协议
		BaseType      string `json:"base_type"`       // 基础类型
		RemoteAddr    string `json:"remotet_addr"`    // 远程地址
		ReadBytesSum  int    `json:"read_bytes_sum"`  // 读取字节总数
		WroteBytesSum int    `json:"wrote_bytes_sum"` // 写入字节总数
	}
	err := c.ShouldBind(&p)
	if err != nil {
		logrus.Error(err)
		return
	}
	stream := mdl.Stream{
		Code:          p.SessionId,
		Name:          p.StreamName,
		Protocol:      p.Protocol,
		Type:          p.BaseType,
		RemoteAddr:    p.RemoteAddr,
		ReceivedBytes: p.ReadBytesSum,
		SentBytes:     p.WroteBytesSum,
	}
	err = svc.StreamSvc.Save(stream, clause.OnConflict{
		Columns: []clause.Column{{
			Name: "code",
		}},
		DoUpdates: clause.AssignmentColumns([]string{
			"received_bytes",
			"sent_bytes",
		}),
	})
	if err != nil {
		logrus.Error(err)
	}
}

// 停止拉流
func StopPullStream(c *gin.Context) {
}

// 开始推流
func StartPushStream(c *gin.Context) {
}

// 停止推流
func StopPushStream(c *gin.Context) {
}

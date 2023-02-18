package mdl

import (
	"github.com/btagrass/go.core/mdl"
)

// 流
type Stream struct {
	mdl.Mdl
	AppName       string `gorm:"size:50;not null;comment:应用名称" json:"appName"`     // 应用名称
	Name          string `gorm:"size:50;not null;comment:名称" json:"name"`          // 名称
	AudioCodec    string `gorm:"size:50;not null;comment:音频解码器" json:"audioCodec"` // 音频解码器                                // 音频解码器
	VideoCodec    string `gorm:"size:50;not null;comment:视频解码器" json:"videoCodec"` // 视频解码器                           // 视频解码器
	VideoWidth    int    `gorm:"comment:视频宽" json:"videoWidth"`                    // 视频宽
	VideoHeight   int    `gorm:"comment:视频高" json:"videoHeight"`                   // 视频高
	Session       string `gorm:"size:50;not null;comment:会话" json:"session"`       // 会话
	Protocol      string `gorm:"size:50;not null;comment:协议" json:"protocol"`      // 协议
	Type          string `gorm:"size:50;comment:类型" json:"type"`                   // 类型
	RemoteAddr    string `gorm:"size:50;not null;comment:远程地址" json:"remoteAddr"`  // 远程地址
	CodeRate      int    `gorm:"comment:码率" json:"codeRate"`                       // 码率
	ReceivedBytes int    `gorm:"comment:接收字节数" json:"receivedBytes"`               // 接收字节数
	SentBytes     int    `gorm:"comment:发送字节数" json:"sentBytes"`                   // 发送字节数
}

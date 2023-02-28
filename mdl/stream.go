package mdl

import (
	"github.com/btagrass/go.core/mdl"
)

// 流
type Stream struct {
	mdl.Mdl
	AppName       string       `gorm:"size:50;not null;comment:应用名称" json:"appName"`      // 应用名称
	Name          string       `gorm:"size:50;not null;comment:名称" json:"name"`           // 名称
	AudioCodec    string       `gorm:"size:50;not null;comment:音频解码" json:"audioCodec"`   // 音频解码
	VideoCodec    string       `gorm:"size:50;not null;comment:视频解码" json:"videoCodec"`   // 视频解码
	VideoWidth    int          `gorm:"comment:视频宽" json:"videoWidth"`                     // 视频宽
	VideoHeight   int          `gorm:"comment:视频高" json:"videoHeight"`                    // 视频高
	Session       string       `gorm:"unique;size:50;not null;comment:会话" json:"session"` // 会话
	Protocol      string       `gorm:"size:50;not null;comment:协议" json:"protocol"`       // 协议
	Type          string       `gorm:"size:50;comment:类型" json:"type"`                    // 类型
	RemoteAddr    string       `gorm:"size:50;not null;comment:远程地址" json:"remoteAddr"`   // 远程地址
	CodeRate      int          `gorm:"comment:码率" json:"codeRate"`                        // 码率
	ReceivedBytes int          `gorm:"comment:接收字节数" json:"receivedBytes"`                // 接收字节数
	SentBytes     int          `gorm:"comment:发送字节数" json:"sentBytes"`                    // 发送字节数
	Pushs         []StreamPush `gorm:"-" json:"pushs"`                                    // 推送集合
	Subs          []StreamSub  `gorm:"-" json:"subs"`                                     // 订阅集合
}

// 流推送
type StreamPush struct {
	mdl.Mdl
	Name       string `gorm:"size:50;not null;comment:名称" json:"name"`          // 名称
	RemoteAddr string `gorm:"size:100;not null;comment:远程地址" json:"remoteAddr"` // 远程地址
	State      int    `gorm:"comment:状态" json:"state"`                          // 状态
}

// 流订阅
type StreamSub struct {
	mdl.Mdl
	Session       string `gorm:"unique;size:50;not null;comment:会话" json:"session"` // 会话
	Protocol      string `gorm:"size:50;not null;comment:协议" json:"protocol"`       // 协议
	Type          string `gorm:"size:50;comment:类型" json:"type"`                    // 类型
	RemoteAddr    string `gorm:"size:50;not null;comment:远程地址" json:"remoteAddr"`   // 远程地址
	CodeRate      int    `gorm:"comment:码率" json:"codeRate"`                        // 码率
	ReceivedBytes int    `gorm:"comment:接收字节数" json:"receivedBytes"`                // 接收字节数
	SentBytes     int    `gorm:"comment:发送字节数" json:"sentBytes"`                    // 发送字节数
}

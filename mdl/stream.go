package mdl

import (
	"github.com/btagrass/go.core/mdl"
)

// 流
type Stream struct {
	mdl.Mdl
	Code          string `gorm:"unique;size:50;not null;comment:代码" json:"code"` // 代码
	Name          string `gorm:"size:50;comment:名称" json:"name"`                 // 名称
	Protocol      string `gorm:"size:50;comment:名称" json:"protocol"`             // 协议
	Type          string `gorm:"size:50;comment:类型" json:"type"`                 // 类型
	RemoteAddr    string `gorm:"size:50;comment:远程地址" json:"remoteAddr"`         // 远程地址
	ReceivedBytes int    `gorm:"comment:接收字节数" json:"receivedBytes"`             // 接收字节数
	SentBytes     int    `gorm:"comment:发送字节数" json:"sentBytes"`                 // 发送字节数
}

package mdl

import (
	"github.com/btagrass/go.core/mdl"
)

// 视频流
type VideoStream struct {
	mdl.Mdl
	Name     string `gorm:"size:50;comment:名称" json:"name"`       // 名称
	PubAddr  string `gorm:"size:100;comment:地址" json:"pubAddr"`   // 推流地址
	FileName string `gorm:"size:100;comment:文件名" json:"fileName"` // 文件名
	State    int    `gorm:"comment:状态" json:"state"`              // 状态
}

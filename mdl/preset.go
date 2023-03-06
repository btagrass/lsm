package mdl

import "github.com/btagrass/go.core/mdl"

// 预置位
type Preset struct {
	mdl.Mdl
	CameraCode string `gorm:"unique;size:50;not null;comment:摄像头代码" json:"cameraCode"` // 摄像头代码
	Index      int    `gorm:"not null;comment:索引" json:"index"`                        // 索引
	Name       string `gorm:"size:50;not null;comment:名称" json:"name"`                 // 名称
}

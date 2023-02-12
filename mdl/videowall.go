package mdl

import (
	"strings"

	"github.com/btagrass/go.core/mdl"
	"github.com/btagrass/go.core/utl"
	"gorm.io/gorm"
)

// 视频墙
type VideoWall struct {
	mdl.Mdl
	Name        string   `gorm:"size:50;not null;comment:名称" json:"name"`    // 名称
	Cameras     []string `gorm:"-" json:"cameras"`                           // 摄像头集合
	CameraCodes string   `gorm:"size:200;not null;comment:摄像头代码集合" json:"-"` // 摄像头代码集合
}

func (m *VideoWall) AfterFind(tx *gorm.DB) error {
	m.Cameras = utl.Split(m.CameraCodes, ',')
	return nil
}

func (m *VideoWall) BeforeSave(tx *gorm.DB) error {
	m.CameraCodes = strings.Join(m.Cameras, ",")
	return nil
}

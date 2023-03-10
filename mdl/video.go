package mdl

import (
	"fmt"

	"github.com/btagrass/go.core/htp"
	"github.com/btagrass/go.core/mdl"
	"gorm.io/gorm"
)

// 视频
type Video struct {
	mdl.Mdl
	Name    string `gorm:"size:50;comment:名称" json:"name"`    // 名称
	Source  string `gorm:"size:100;comment:来源" json:"source"` // 来源
	Process int    `gorm:"comment:进程" json:"process"`         // 进程
	Url     string `gorm:"-" json:"url"`                      // 网址
}

func (m *Video) AfterFind(tx *gorm.DB) error {
	// if m.Process > 0 {
	// 	err := syscall.Kill(m.Process, 0)
	// 	if err != nil {
	// 		m.Process = 0
	// 	}
	// }
	m.Url = fmt.Sprintf("rtsp://%s:5544/live/%d", htp.Ip, m.Id)

	return nil
}

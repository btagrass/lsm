package mdl

import (
	"github.com/btagrass/go.core/mdl"
)

// 摄像头
type Camera struct {
	mdl.Mdl
	Code     string `gorm:"unique;size:50;not null;comment:代码" json:"code"` // 代码
	Name     string `gorm:"size:50;comment:名称" json:"name"`                 // 名称
	Type     string `gorm:"size:50;comment:类型" json:"type"`                 // 类型
	Mfr      string `gorm:"size:50;comment:厂商" json:"mfr"`                  // 厂商
	Model    string `gorm:"size:50;comment:型号" json:"model"`                // 型号
	Firmware string `gorm:"size:50;comment:固件" json:"firmware"`             // 固件
	Addr     string `gorm:"size:100;comment:地址" json:"addr"`                // 地址
	UserName string `gorm:"size:50;comment:用户名" json:"userName"`            // 用户名
	Password string `gorm:"size:50;comment:密码" json:"password"`             // 密码
	State    int    `gorm:"comment:状态" json:"state"`                        // 状态
}

package isc

import (
	"lsm/mdl"
	"time"

	"github.com/btagrass/go.core/r"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm/clause"
)

// 保活摄像头集合
func (s *IscSvc) keepaliveCameras() error {
	var r struct {
		r.R
		Data struct {
			Records []struct {
				Code      string `json:"indexCode"`    // 代码
				Name      string `json:"cn"`           // 名称
				Mfr       string `json:"manufacturer"` // 厂商
				Model     string `json:"deviceType"`   // 型号
				State     int    `json:"online"`       // 状态
				UpdatedAt string `json:"collectTime"`  // 更新时间
			} `json:"list"` // 列表
		} `json:"data"` // 数据
	}
	cameras, _, err := s.CameraSvc.ListCameras()
	if err != nil {
		return err
	}
	var codes []string
	for _, c := range cameras {
		codes = append(codes, c.Code)
	}
	_, err = s.post("/artemis/api/nms/v1/online/camera/get", map[string]any{
		"indexCodes": codes,
		"pageNo":     1,
		"pageSize":   1000,
	}, &r)
	if err != nil {
		return err
	}
	for _, r := range r.Data.Records {
		c := mdl.Camera{
			Code:  r.Code,
			Name:  r.Name,
			Mfr:   r.Mfr,
			Model: r.Model,
			State: r.State,
		}
		updatedAt, _ := time.Parse(time.RFC3339, r.UpdatedAt)
		c.UpdatedAt = updatedAt
		err = s.Save(c, clause.OnConflict{
			Columns: []clause.Column{{
				Name: "code",
			}},
			DoUpdates: clause.AssignmentColumns([]string{
				"name",
				"mfr",
				"models",
				"state",
				"updated_at",
			}),
		})
		if err != nil {
			logrus.WithField("code", c.Code).Error(err)
		}
	}

	return err
}

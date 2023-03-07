package isc

import (
	"lsm/mdl"
	"time"

	"github.com/btagrass/go.core/r"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"gorm.io/gorm/clause"
)

func (s *IscSvc) ListCameras(conds ...any) ([]mdl.Camera, int64, error) {
	cameras := make([]mdl.Camera, 0)
	var r struct {
		r.R
		Data struct {
			List []struct {
				IndexCode  string `json:"indexCode"`   // 代码
				Name       string `json:"name"`        // 名称
				CameraType int    `json:"cameraType"`  // 摄像头类型
				Model      string `json:"deviceType"`  // 型号
				State      int    `json:"online"`      // 状态
				UpdatedAt  string `json:"collectTime"` // 更新时间
			} `json:"list"` // 列表
			Total int64 `json:"total"` // 总数
		} `json:"data"` // 数据
	}
	var name string
	pageIndex, pageSize := 1, 10
	if len(conds) > 0 {
		cond, ok := conds[0].(map[string]any)
		if ok {
			name = cast.ToString(cond["name"])
			current, ok := cond["current"]
			if ok {
				pageIndex = cast.ToInt(current)
			}
			size, ok := cond["size"]
			if ok {
				pageSize = cast.ToInt(size)
			}
		}
	}
	_, err := s.post("/artemis/api/resource/v2/camera/search", map[string]any{
		"name":     name,
		"pageNo":   pageIndex,
		"pageSize": pageSize,
	}, &r)
	if err != nil {
		return cameras, 0, err
	}
	for _, c := range r.Data.List {
		camera := mdl.Camera{
			Code: c.IndexCode,
			Name: c.Name,
		}
		switch c.CameraType {
		case 0:
			camera.Type = "枪机"
		case 1:
			camera.Type = "半球机"
		case 2:
			camera.Type = "球机"
		case 3:
			camera.Type = "云台枪机"
		}
		cameras = append(cameras, camera)
	}

	return cameras, r.Data.Total, nil
}

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

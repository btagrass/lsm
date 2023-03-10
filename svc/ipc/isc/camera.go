package isc

import (
	"lsm/mdl"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
)

func (s *IscSvc) SyncCameras() error {
	var r struct {
		Data struct {
			List []struct {
				IndexCode  string `json:"indexCode"`  // 索引代码
				Name       string `json:"name"`       // 名称
				CameraType int    `json:"cameraType"` // 摄像头类型
				CreateTime string `json:"createTime"` // 创建时间
				UpdateTime string `json:"updateTime"` // 更新时间
			} `json:"list"` // 列表
			Total int64 `json:"total"` // 总数
		} `json:"data"` // 数据
	}
	_, err := s.post("/artemis/api/resource/v2/camera/search", map[string]any{
		"pageNo":   1,
		"pageSize": 1000,
	}, &r)
	if err != nil {
		return err
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
		camera.CreatedAt = cast.ToTime(c.CreateTime)
		camera.UpdatedAt = cast.ToTime(c.UpdateTime)
		err = s.SaveCamera(camera)
		if err != nil {
			logrus.Error(err)
		}
	}

	return nil
}

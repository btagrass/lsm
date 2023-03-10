package ivs

import (
	"fmt"
	"lsm/mdl"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
)

func (s *IvsSvc) SyncCameras() error {
	var r struct {
		CameraBriefInfos struct {
			CameraBriefInfoList []struct {
				Code             string `json:"code"`             // 代码
				Name             string `json:"name"`             // 名称
				Type             int    `json:"type"`             // 类型
				VendorType       string `json:"vendorType"`       // 厂商类型
				DeviceModelType  string `json:"deviceModelType"`  // 设备型号类型
				Status           int    `json:"status"`           // 状态
				DeviceCreateTime string `json:"deviceCreateTime"` // 设备创建时间
			} `json:"cameraBriefInfoList"` // 摄像头简介信息列表
			Total int64 `json:"total"` // 总数
		} `json:"cameraBriefInfos"` // 摄像头简介信息
	}
	_, err := s.get(fmt.Sprintf("/device/deviceList/v1.0?deviceType=2&fromIndex=%d&toIndex=%d", 1, 1000), &r)
	if err != nil {
		return err
	}
	for _, c := range r.CameraBriefInfos.CameraBriefInfoList {
		camera := mdl.Camera{
			Code:  c.Code,
			Name:  c.Name,
			Mfr:   c.VendorType,
			Model: c.DeviceModelType,
			State: c.Status,
		}
		switch c.Type {
		case 0:
			camera.Type = "枪机"
		case 1:
			camera.Type = "云台枪机"
		case 2:
			camera.Type = "球机"
		case 3:
			camera.Type = "半球机"
		case 4:
			camera.Type = "筒机"
		}
		camera.CreatedAt = cast.ToTime(c.DeviceCreateTime)
		err = s.SaveCamera(camera)
		if err != nil {
			logrus.Error(err)
		}
	}

	return nil
}

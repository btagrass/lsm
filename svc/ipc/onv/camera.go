package onv

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/use-go/onvif/device"
	"gorm.io/gorm/clause"
)

// 保活摄像头集合
func (s *OnvSvc) keepaliveCameras() error {
	cameras, err := s.ListCameras("")
	if err != nil {
		return err
	}
	for _, c := range cameras {
		// 同步摄像头
		dev, err := s.getDevice(c.Code)
		if err != nil {
			return err
		}
		doc, err := dev.call(device.GetDeviceInformation{})
		if err != nil {
			return err
		}
		manufacturer := doc.FindElement("//tds:GetDeviceInformationResponse/tds:Manufacturer")
		if manufacturer != nil {
			c.Mfr = manufacturer.Text()
		}
		model := doc.FindElement("//tds:GetDeviceInformationResponse/tds:Model")
		if manufacturer != nil {
			c.Model = model.Text()
		}
		firmware := doc.FindElement("//tds:GetDeviceInformationResponse/tds:FirmwareVersion")
		if firmware != nil {
			c.Firmware = firmware.Text()
		}
		c.State = 1
		c.UpdatedAt = time.Now()
		err = s.Save(c, clause.OnConflict{
			Columns: []clause.Column{{
				Name: "code",
			}},
			DoUpdates: clause.AssignmentColumns([]string{
				"mfr",
				"models",
				"firmware",
				"state",
				"updated_at",
			}),
		})
		if err != nil {
			logrus.Error(err)
		}
		// 开始流
		_, err = s.StartStream(c.Code, 1, "rtsp")
		if err != nil {
			logrus.Error(err)
		}
	}

	return nil
}

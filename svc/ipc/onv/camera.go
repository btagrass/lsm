package onv

import (
	"lsm/mdl"
	"net"
	"regexp"

	"github.com/sirupsen/logrus"
	"github.com/use-go/onvif"
	"github.com/use-go/onvif/device"
)

func (s *OnvSvc) SyncCameras() error {
	re := regexp.MustCompile(`\b\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\b`)
	interfaces, err := net.Interfaces()
	if err != nil {
		return err
	}
	for _, i := range interfaces {
		if i.Flags&net.FlagLoopback == net.FlagLoopback {
			continue
		}
		devices, err := onvif.GetAvailableDevicesAtSpecificEthernetInterface(i.Name)
		if err != nil {
			logrus.Error(err)
			continue
		}
		for _, d := range devices {
			addr := re.FindString(d.GetEndpoint("device"))
			dev, err := s.newDevice(addr, s.userName, s.password)
			if err != nil {
				logrus.Error(err)
				continue
			}
			camera := mdl.Camera{
				Code:     dev.code,
				Addr:     addr,
				UserName: dev.userName,
				Password: dev.password,
				State:    1,
			}
			doc, err := dev.call(device.GetDeviceInformation{})
			if err != nil {
				logrus.Error(err)
				continue
			}
			manufacturer := doc.FindElement("//tds:GetDeviceInformationResponse/tds:Manufacturer")
			if manufacturer != nil {
				camera.Mfr = manufacturer.Text()
			}
			model := doc.FindElement("//tds:GetDeviceInformationResponse/tds:Model")
			if manufacturer != nil {
				camera.Model = model.Text()
			}
			firmware := doc.FindElement("//tds:GetDeviceInformationResponse/tds:FirmwareVersion")
			if firmware != nil {
				camera.Firmware = firmware.Text()
			}
			err = s.SaveCamera(camera)
			if err != nil {
				logrus.Error(err)
			}
		}
	}

	return nil
}

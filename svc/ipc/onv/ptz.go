package onv

import (
	"fmt"
	"lsm/mdl"
	"strings"

	"github.com/spf13/cast"
	"github.com/use-go/onvif/ptz"
	"github.com/use-go/onvif/xsd"
	onvifd "github.com/use-go/onvif/xsd/onvif"
)

func (s *OnvSvc) ControlPtz(code string, command string, speed int) error {
	var spd, x, y, zoom float64
	if speed == 1 {
		spd = 0.1
	} else if speed == 2 {
		spd = 0.5
	} else if speed == 3 {
		spd = 0.9
	}
	switch strings.ToUpper(command) {
	case "LEFT":
		x, y, zoom = -spd, 0.0, 0.0
	case "RIGHT":
		x, y, zoom = spd, 0.0, 0.0
	case "UP":
		x, y, zoom = 0.0, spd, 0.0
	case "DOWN":
		x, y, zoom = 0.0, -spd, 0.0
	case "LEFTUP":
		x, y, zoom = -spd, spd, 0.0
	case "LEFTDOWN":
		x, y, zoom = -spd, -spd, 0.0
	case "RIGHTUP":
		x, y, zoom = spd, spd, 0.0
	case "RIGHTDOWN":
		x, y, zoom = spd, -spd, 0.0
	case "ZOOMIN":
		x, y, zoom = 0.0, 0.0, spd
	case "ZOOMOUT":
		x, y, zoom = 0.0, 0.0, -spd
	}
	dev, err := s.getDevice(code)
	if err != nil {
		return err
	}
	doc, err := dev.call(ptz.ContinuousMove{
		ProfileToken: onvifd.ReferenceToken(dev.profileTokens[0]),
		Velocity: onvifd.PTZSpeed{
			PanTilt: onvifd.Vector2D{
				X: x,
				Y: y,
			},
			Zoom: onvifd.Vector1D{
				X: zoom,
			},
		},
		Timeout: xsd.Duration("PT1S"),
	})
	if err != nil {
		return err
	}
	element := doc.FindElement("//tptz:ContinuousMoveResponse")
	if element == nil {
		return fmt.Errorf("控制摄像头 %s 的云台 %s 失败", code, command)
	}

	return nil
}

func (s *OnvSvc) GotoPreset(code string, index int) error {
	dev, err := s.getDevice(code)
	if err != nil {
		return err
	}
	presetToken, err := dev.getPresetToken(index)
	if err != nil {
		return err
	}
	doc, err := dev.call(ptz.GotoPreset{
		ProfileToken: onvifd.ReferenceToken(dev.profileTokens[0]),
		PresetToken:  onvifd.ReferenceToken(presetToken),
	})
	if err != nil {
		return err
	}
	element := doc.FindElement("//tptz:GotoPresetResponse")
	if element == nil {
		return fmt.Errorf("转到摄像头 %s 的预置位 %d 失败", code, index)
	}

	return nil
}

func (s *OnvSvc) ListPresets(code string) ([]mdl.Preset, error) {
	presets := make([]mdl.Preset, 0)
	dev, err := s.getDevice(code)
	if err != nil {
		return presets, err
	}
	doc, err := dev.call(ptz.GetPresets{
		ProfileToken: onvifd.ReferenceToken(dev.profileTokens[0]),
	})
	if err != nil {
		return presets, err
	}
	elements := doc.FindElements("//tptz:GetPresetsResponse/tptz:Preset")
	for _, e := range elements {
		preset := mdl.Preset{
			Index: cast.ToInt(strings.TrimPrefix(e.SelectAttrValue("token", "0"), "Preset_")),
			Name:  e.FindElement("./tt:Name").Text(),
		}
		preset.Id = int64(preset.Index)
		presets = append(presets, preset)
	}

	return presets, nil
}

func (s *OnvSvc) RemovePreset(code string, index int) error {
	dev, err := s.getDevice(code)
	if err != nil {
		return err
	}
	presetToken, err := dev.getPresetToken(index)
	if err != nil {
		return err
	}
	doc, err := dev.call(ptz.RemovePreset{
		ProfileToken: onvifd.ReferenceToken(dev.profileTokens[0]),
		PresetToken:  onvifd.ReferenceToken(presetToken),
	})
	if err != nil {
		return err
	}
	element := doc.FindElement("//tptz:RemovePresetResponse")
	if element == nil {
		return fmt.Errorf("移除摄像头 %s 的预置位 %d 失败", code, index)
	}

	return nil
}

func (s *OnvSvc) SavePreset(preset mdl.Preset) error {
	dev, err := s.getDevice(preset.CameraCode)
	if err != nil {
		return err
	}
	presetToken, _ := dev.getPresetToken(preset.Index)
	doc, err := dev.call(ptz.SetPreset{
		ProfileToken: onvifd.ReferenceToken(dev.profileTokens[0]),
		PresetName:   xsd.String(preset.Name),
		PresetToken:  onvifd.ReferenceToken(presetToken),
	})
	if err != nil {
		return err
	}
	element := doc.FindElement("//tptz:SetPresetResponse")
	if element == nil {
		return fmt.Errorf("保存摄像头 %s 的预置位 %d 失败", preset.CameraCode, preset.Index)
	}

	return nil
}

package onv

import (
	"fmt"
	"lsm/mdl"
	"strings"

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
		return fmt.Errorf("device %s 's command %s is failed", code, command)
	}

	return nil
}

func (s *OnvSvc) ListPresets(code string) ([]mdl.Preset, error) {
	return nil, nil
}

func (s *OnvSvc) RemovePreset(code string, index int) error {
	return nil
}

func (s *OnvSvc) SavePreset(preset mdl.Preset) error {
	return nil
}

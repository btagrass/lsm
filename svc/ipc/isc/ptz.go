package isc

import (
	"lsm/mdl"
	"strings"
	"time"
)

func (s *IscSvc) ControlPtz(code string, command string, speed int) error {
	var cmd string
	switch strings.ToUpper(command) {
	case "LEFT":
		cmd = "LEFT"
	case "RIGHT":
		cmd = "RIGHT"
	case "UP":
		cmd = "UP"
	case "DOWN":
		cmd = "DOWN"
	case "LEFTUP":
		cmd = "LEFT_UP"
	case "LEFTDOWN":
		cmd = "LEFT_DOWN"
	case "RIGHTUP":
		cmd = "RIGHT_UP"
	case "RIGHTDOWN":
		cmd = "RIGHT_DOWN"
	case "ZOOMIN":
		cmd = "ZOOM_IN"
	case "ZOOMOUT":
		cmd = "ZOOM_OUT"
	}
	var spd int
	if speed == 1 {
		spd = 10
	} else if speed == 2 {
		spd = 50
	} else if speed == 3 {
		spd = 90
	}
	_, err := s.post(
		"/artemis/api/video/v1/ptzs/controlling",
		map[string]any{
			"cameraIndexCode": code,
			"action":          0,
			"command":         cmd,
			"speed":           spd,
		},
		nil,
	)
	if err != nil {
		return err
	}
	time.Sleep(time.Second)
	_, err = s.post(
		"/artemis/api/video/v1/ptzs/controlling",
		map[string]any{
			"cameraIndexCode": code,
			"action":          1,
			"command":         cmd,
			"speed":           spd,
		},
		nil,
	)

	return err
}

func (s *IscSvc) ListPresets(code string) ([]mdl.Preset, error) {
	return nil, nil
}

func (s *IscSvc) RemovePreset(code string, index int) error {
	return nil
}

func (s *IscSvc) SavePreset(preset mdl.Preset) error {
	return nil
}

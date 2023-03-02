package isc

import (
	"strings"
	"time"
)

func (s *IscSvc) ControlPtz(code string, command string, speed int) error {
	var c string
	switch strings.ToUpper(command) {
	case "LEFT":
		c = "LEFT"
	case "RIGHT":
		c = "RIGHT"
	case "UP":
		c = "UP"
	case "DOWN":
		c = "DOWN"
	case "LEFTUP":
		c = "LEFT_UP"
	case "LEFTDOWN":
		c = "LEFT_DOWN"
	case "RIGHTUP":
		c = "RIGHT_UP"
	case "RIGHTDOWN":
		c = "RIGHT_DOWN"
	case "ZOOMIN":
		c = "ZOOM_IN"
	case "ZOOMOUT":
		c = "ZOOM_OUT"
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
			"command":         c,
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
			"command":         c,
			"speed":           spd,
		},
		nil,
	)

	return err
}

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
	_, err := s.post("/artemis/api/video/v1/ptzs/controlling", map[string]any{
		"cameraIndexCode": code,
		"action":          0,
		"command":         cmd,
		"speed":           spd,
	})
	if err != nil {
		return err
	}
	time.Sleep(time.Second)
	_, err = s.post("/artemis/api/video/v1/ptzs/controlling", map[string]any{
		"cameraIndexCode": code,
		"action":          1,
		"command":         cmd,
		"speed":           spd,
	})

	return err
}

func (s *IscSvc) GotoPreset(code string, index int) error {
	_, err := s.post("/artemis/api/video/v1/ptzs/controlling", map[string]any{
		"cameraIndexCode": code,
		"action":          0,
		"command":         "GOTO_PRESET",
		"presetIndex":     index,
	})

	return err
}

func (s *IscSvc) ListPresets(code string) ([]mdl.Preset, error) {
	presets := make([]mdl.Preset, 0)
	var r struct {
		Data struct {
			List []struct {
				PresetPointIndex int    `json:"presetPointIndex"` // 预置位索引
				PresetPointName  string `json:"presetPointName"`  // 预置位名称
			} `json:"list"` // 列表
		} `json:"data"` // 数据
	}
	_, err := s.post("/artemis/api/video/v1/presets/searches", map[string]any{
		"cameraIndexCode": code,
	}, &r)
	if err != nil {
		return presets, err
	}
	for _, r := range r.Data.List {
		preset := mdl.Preset{
			Index: r.PresetPointIndex,
			Name:  r.PresetPointName,
		}
		preset.Id = int64(preset.Index)
		presets = append(presets, preset)
	}

	return presets, nil
}

func (s *IscSvc) RemovePreset(code string, index int) error {
	_, err := s.post("/artemis/api/video/v1/presets/deletion", map[string]any{
		"cameraIndexCode": code,
		"presetIndex":     index,
	})

	return err
}

func (s *IscSvc) SavePreset(preset mdl.Preset) error {
	_, err := s.post("/artemis/api/video/v1/presets/addition", map[string]any{
		"cameraIndexCode": preset.CameraCode,
		"presetIndex":     preset.Index,
		"presetName":      preset.Name,
	})

	return err
}

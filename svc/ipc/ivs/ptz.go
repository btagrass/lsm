package ivs

import (
	"fmt"
	"lsm/mdl"
	"strings"

	"github.com/spf13/cast"
)

func (s *IvsSvc) ControlPtz(code string, command string, speed int) error {
	var cmd int
	switch strings.ToUpper(command) {
	case "LEFT":
		cmd = 4
	case "RIGHT":
		cmd = 7
	case "UP":
		cmd = 2
	case "DOWN":
		cmd = 3
	case "LEFTUP":
		cmd = 5
	case "LEFTDOWN":
		cmd = 6
	case "RIGHTUP":
		cmd = 8
	case "RIGHTDOWN":
		cmd = 9
	case "ZOOMIN":
		cmd = 23
	case "ZOOMOUT":
		cmd = 24
	}
	var spd string
	if speed == 1 {
		spd = "1"
	} else if speed == 2 {
		spd = "5"
	} else if speed == 3 {
		spd = "9"
	}
	_, err := s.post("/device/ptzcontrol", map[string]any{
		"cameraCode":   code,
		"controlCode":  cmd,
		"controlPara1": "1",
		"controlPara2": spd,
	})

	return err
}

func (s *IvsSvc) GotoPreset(code string, index int) error {
	_, err := s.post("/device/ptzcontrol", map[string]any{
		"cameraCode":   code,
		"controlCode":  11,
		"controlPara1": cast.ToString(index),
		"controlPara2": "",
	})

	return err
}

func (s *IvsSvc) ListPresets(code string) ([]mdl.Preset, error) {
	presets := make([]mdl.Preset, 0)
	var r struct {
		PtzPresetInfoList []struct {
			PresetIndex int    `json:"presetIndex"` // 预置位索引
			PresetName  string `json:"presetName"`  // 预置位名称
		} `json:"ptzPresetInfoList"` // 云台预置位信息列表
	}
	_, err := s.get(fmt.Sprintf("/device/ptzpresetlist/%s/", code), &r)
	if err != nil {
		return presets, err
	}
	for _, p := range r.PtzPresetInfoList {
		presets = append(presets, mdl.Preset{
			Index: p.PresetIndex,
			Name:  p.PresetName,
		})
	}

	return presets, nil
}

func (s *IvsSvc) RemovePreset(code string, index int) error {
	_, err := s.delete(fmt.Sprintf("/ptz/presetposition/%s//%d/v1.0", code, index))

	return err
}

func (s *IvsSvc) SavePreset(preset mdl.Preset) error {
	var r struct {
		PresetIndex int `json:"presetIndex"` // 预置位索引
	}
	var err error
	if preset.Index == 0 {
		_, err = s.post("/ptz/presetposition/v1.0", map[string]any{
			"cameraCode":  preset.CameraCode,
			"presetName":  preset.Name,
			"focusSwitch": 1,
		}, &r)
	} else {
		_, err = s.put("/ptz/presetposition/v1.0", map[string]any{
			"cameraCode": preset.CameraCode,
			"ptzPresetInfo": map[string]any{
				"presetIndex": preset.Index,
				"presetName":  preset.Name,
			},
		}, &r)
		r.PresetIndex = preset.Index
	}

	return err
}

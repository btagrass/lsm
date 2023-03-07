package onv

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/beevik/etree"
	"github.com/use-go/onvif"
	"github.com/use-go/onvif/media"
	"github.com/use-go/onvif/ptz"
	onvid "github.com/use-go/onvif/xsd/onvif"
	"golang.org/x/net/html/charset"
)

// 设备
type Device struct {
	code          string        // 代码
	userName      string        // 用户名
	password      string        // 密码
	device        *onvif.Device // 设备
	profileTokens []string      // 配置令牌集合
}

// 调用
func (d *Device) call(data any) (*etree.Document, error) {
	var doc *etree.Document
	resp, err := d.device.CallMethod(data)
	if err != nil {
		return doc, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return doc, fmt.Errorf(resp.Status)
	}
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return doc, err
	}
	doc = etree.NewDocument()
	doc.ReadSettings.CharsetReader = charset.NewReaderLabel
	err = doc.ReadFromBytes(bytes)
	if err != nil {
		return doc, err
	}

	return doc, nil
}

// 获取预置位令牌
func (d *Device) getPresetToken(index int) (string, error) {
	var token string
	doc, err := d.call(ptz.GetPresets{
		ProfileToken: onvid.ReferenceToken(d.profileTokens[0]),
	})
	if err != nil {
		return token, err
	}
	element := doc.FindElement(fmt.Sprintf("//tptz:GetPresetsResponse/tptz:Preset[@token='%d']", index))
	if element == nil {
		element = doc.FindElement(fmt.Sprintf("//tptz:GetPresetsResponse/tptz:Preset[@token='Preset_%d']", index))
	}
	if element == nil {
		return token, fmt.Errorf("device %s 's presetToken %d is not found", d.code, index)
	}
	token = element.SelectAttrValue("token", "")

	return token, nil
}

// 获取配置令牌集合
func (d *Device) getProfileTokens() ([]string, error) {
	var tokens []string
	doc, err := d.call(media.GetProfiles{})
	if err != nil {
		return tokens, err
	}
	elements := doc.FindElements("//trt:GetProfilesResponse/trt:Profiles")
	for _, e := range elements {
		token := e.SelectAttrValue("token", "")
		if token != "" {
			tokens = append(tokens, token)
		}
	}

	return tokens, nil
}

// 获取设备
func (s *OnvSvc) getDevice(code string) (*Device, error) {
	key := fmt.Sprintf("%s:devices:%s", s.Prefix, code)
	v, ok := s.Cache.Get(key)
	if ok {
		return v.(*Device), nil
	}
	camera, err := s.CameraSvc.GetCameraByCode(code)
	if err != nil {
		return nil, err
	}
	if camera == nil {
		return nil, fmt.Errorf("摄像头 %s 不存在", code)
	}
	device := &Device{
		code:     code,
		userName: camera.UserName,
		password: camera.Password,
	}
	device.device, err = onvif.NewDevice(onvif.DeviceParams{
		Xaddr:    camera.Addr,
		Username: camera.UserName,
		Password: camera.Password,
	})
	if err != nil {
		return nil, err
	}
	device.profileTokens, err = device.getProfileTokens()
	if err != nil {
		return nil, err
	}
	s.Cache.Set(key, device, 5*time.Minute)

	return device, nil
}

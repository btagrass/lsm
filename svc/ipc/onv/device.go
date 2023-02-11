package onv

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/beevik/etree"
	"github.com/use-go/onvif"
	"github.com/use-go/onvif/media"
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
	deviceKey := fmt.Sprintf("%s:devices:%s", s.Prefix, code)
	v, ok := s.Cache.Get(deviceKey)
	if ok {
		return v.(*Device), nil
	}
	camera, err := s.CameraSvc.GetCameraByCode(code)
	if err != nil || camera == nil {
		return nil, err
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
	s.Cache.Set(deviceKey, device, 5*time.Minute)

	return device, nil
}

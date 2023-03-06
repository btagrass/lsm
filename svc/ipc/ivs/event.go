package ivs

import (
	"fmt"
	"strings"
)

func (s *IvsSvc) RegisterCallback() error {
	var r struct {
		Code int `json:"resultCode"` // 代码
	}
	_, err := s.post("/users/regeditcallback", map[string]any{
		"wsUri": s.eventCallback,
	}, &r)
	if err != nil {
		return err
	}
	if r.Code != 0 {
		return fmt.Errorf("注册事件回调失败：%d", r.Code)
	}

	return nil
}

// 订阅事件集合
func (s *IvsSvc) SubscribeEvents(typs ...string) error {
	var r struct {
		Code int `json:"resultCode"` // 代码
	}
	builder := &strings.Builder{}
	for _, t := range typs {
		builder.WriteString(fmt.Sprintf("<SubscribeInfo><AlarmIncode>%s</AlarmIncode></SubscribeInfo>", t))
	}
	_, err := s.post("/device/subscribealarm", map[string]any{
		"requestXML": fmt.Sprintf("<Content><Subscribe><AddSubscribeList>%s</AddSubscribeList></Subscribe></Content>", builder.String()),
	}, &r)
	if err != nil {
		return err
	}
	if r.Code != 0 {
		return fmt.Errorf("订阅事件集合失败：%d", r.Code)
	}

	return nil
}

package ivs

import (
	"fmt"
	"strings"

	"github.com/btagrass/go.core/utl"
)

func (s *IvsSvc) NotifyEvent(content string) error {
	return nil
}

// 订阅事件集合
func (s *IvsSvc) SubscribeEvents() error {
	_, err := s.post("/users/regeditcallback", map[string]any{
		"wsUri": s.eventCallback,
	})
	if err != nil {
		return err
	}
	builder := &strings.Builder{}
	typs := utl.Split(s.eventType, ',')
	for _, t := range typs {
		builder.WriteString(fmt.Sprintf("<SubscribeInfo><AlarmIncode>%s</AlarmIncode></SubscribeInfo>", t))
	}
	_, err = s.post("/device/subscribealarm", map[string]any{
		"requestXML": fmt.Sprintf("<Content><Subscribe><AddSubscribeList>%s</AddSubscribeList></Subscribe></Content>", builder.String()),
	})

	return err
}

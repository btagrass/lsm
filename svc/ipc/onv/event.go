package onv

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/use-go/onvif/event"
)

func (s *OnvSvc) NotifyEvent(content string) error {
	return nil
}

// 订阅事件集合
func (s *OnvSvc) SubscribeEvents() error {
	cameras, _, err := s.ListCameras()
	if err != nil {
		return err
	}
	for _, c := range cameras {
		dev, err := s.getDevice(c.Code)
		if err != nil {
			logrus.Error(err)
			continue
		}
		doc, err := dev.call(event.GetEventProperties{})
		if err != nil {
			logrus.Error(err)
			continue
		}
		fmt.Println(doc)
	}

	return nil
}

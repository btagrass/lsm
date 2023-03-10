package isc

import (
	"github.com/btagrass/go.core/utl"
	"github.com/spf13/cast"
)

func (s *IscSvc) NotifyEvent(content string) error {
	return nil
}

// 订阅事件集合
func (s *IscSvc) SubscribeEvents() error {
	_, err := s.post("/artemis/api/eventService/v1/eventSubscriptionByEventTypes", map[string]any{
		"eventTypes": cast.ToIntSlice(utl.Split(s.eventType, ',')),
		"eventDest":  s.eventCallback,
	})

	return err
}

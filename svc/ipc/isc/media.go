package isc

import (
	"time"
)

func (s *IscSvc) StartStream(code string, typ int, protocol string) (string, error) {
	var r struct {
		Data struct {
			Url string `json:"url"` // 地址
		} `json:"data"` // 数据
	}
	_, err := s.post("/artemis/api/video/v2/cameras/previewURLs", map[string]any{
		"cameraIndexCode": code,
		"streamType":      typ - 1,
		"protocol":        protocol,
	}, &r)

	return r.Data.Url, err
}

func (s *IscSvc) StopStream(code string, typ int) error {
	return s.StopPullStream(code)
}

func (s *IscSvc) GetRecordUrl(code string, date time.Time) (string, error) {
	var r struct {
		Data struct {
			Url string `json:"url"` // 地址
		} `json:"data"` // 数据
	}
	_, err := s.post("/artemis/api/video/v2/cameras/previewURLs", map[string]any{
		"cameraIndexCode": code,
		"protocol":        "rtsp",
		"beginTime":       date.Format("2006-01-02T00:00:00+08:00"),
		"endTime":         date.Format("2006-01-02T23:59:59+08:00"),
	}, &r)

	return r.Data.Url, err
}

func (s *IscSvc) TakeSnapshot(code string, typ int) (string, error) {
	var r struct {
		Data struct {
			PicUrl string `json:"picUrl"` // 图片地址
		} `json:"data"` // 数据
	}
	_, err := s.post("/artemis/api/video/v1/manualCapture", map[string]string{
		"cameraIndexCode": code,
	}, &r)

	return r.Data.PicUrl, err
}

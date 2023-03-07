package isc

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/btagrass/go.core/htp"
	"github.com/btagrass/go.core/utl"
)

func (s *IscSvc) GetRecordUrl(code string, date time.Time) (string, error) {
	return "", fmt.Errorf("录像网址不存在")
}

func (s *IscSvc) StartStream(code string, typ int, protocol string) (string, error) {
	var r struct {
		Data struct {
			Url string `json:"url"` // 网址
		} `json:"data"` // 数据
	}
	resp, err := s.post(
		"/artemis/api/video/v2/cameras/previewURLs",
		map[string]any{
			"cameraIndexCode": code,
			"streamType":      typ - 1,
			"protocol":        protocol,
			"expand":          "transcode=1&videotype=h264",
		},
		&r,
	)
	fmt.Println(resp)

	return r.Data.Url, err
}

func (s *IscSvc) StopStream(code string, typ int) error {
	return nil
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

func (s *IscSvc) TakeSnapshots(code string, cntSecs ...int) ([]string, error) {
	var imageUrls []string
	streamUrl, err := s.StartStream(code, 1, "rtsp")
	if err != nil {
		return imageUrls, err
	}
	filePrefix := fmt.Sprintf("%s_%s", code, time.Now().Format("20060102150405.999999999"))
	fileName := fmt.Sprintf("data/cameras/%s_%%2d.jpg", filePrefix)
	cnt := 1
	secs := 1
	if len(cntSecs) == 1 {
		cnt = cntSecs[0]
	} else if len(cntSecs) >= 2 {
		cnt = cntSecs[0]
		secs = cntSecs[1]
	}
	err = utl.TakeStream(streamUrl, map[string]map[string]any{
		fileName: {
			"r":  fmt.Sprintf("1/%d", cnt),
			"s":  "1920*1080",
			"to": secs,
		},
	})
	if err != nil {
		return imageUrls, err
	}
	images, _ := filepath.Glob(fmt.Sprintf("data/cameras/%s_*.jpg", filePrefix))
	for _, i := range images {
		imageUrls = append(imageUrls, htp.GetUrl(filepath.ToSlash(i)))
	}

	return imageUrls, nil
}

func (s *IscSvc) TakeRecord(code string, beginTime, endTime time.Time) (string, error) {
	return "", nil
}

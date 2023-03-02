package isc

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/btagrass/go.core/htp"
	"github.com/btagrass/go.core/r"
	"github.com/btagrass/go.core/utl"
)

func (is *IscSvc) GetRecordUrl(code string, date time.Time) (string, error) {
	return "", fmt.Errorf("录像网址不存在")
}

func (is *IscSvc) StartStream(code string, streamType int, protocol string) (string, error) {
	var r struct {
		r.R
		Data struct {
			Url string `json:"url"` // 网址
		} `json:"data"` // 数据
	}
	resp, err := is.post(
		"/artemis/api/video/v2/cameras/previewURLs",
		map[string]any{
			"cameraIndexCode": code,
			"streamType":      streamType - 1,
			"protocol":        protocol,
			"expand":          "transcode=1&videotype=h264",
		},
		&r,
	)
	fmt.Println(resp)

	return r.Data.Url, err
}

func (is *IscSvc) StopStream(code string, typ int) error {
	return nil
}

func (is *IscSvc) TakeSnapshot(code string, streamType int) (string, error) {
	var r struct {
		r.R
		Data struct {
			Url string `json:"picUrl"` // 网址
		} `json:"data"` // 数据
	}
	_, err := is.post(
		"/artemis/api/video/v1/manualCapture",
		map[string]string{
			"cameraIndexCode": code,
		},
		&r,
	)

	return r.Data.Url, err
}

func (is *IscSvc) TakeSnapshots(code string, cntSecs ...int) ([]string, error) {
	var imageUrls []string
	streamUrl, err := is.StartStream(code, 1, "rtsp")
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

func (is *IscSvc) TakeRecord(code string, beginTime, endTime time.Time) (string, error) {
	return "", nil
}

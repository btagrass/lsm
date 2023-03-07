package ivs

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/btagrass/go.core/htp"
	"github.com/btagrass/go.core/utl"
	"github.com/patrickmn/go-cache"
)

func (s *IvsSvc) StartStream(code string, typ int, protocol string) (string, error) {
	key := fmt.Sprintf("%s:%s:uris:%d", s.Prefix, code, typ)
	v, ok := s.Cache.Get(key)
	if ok {
		return v.(string), nil
	}
	var r struct {
		RtspUrl string `json:"rtspURL"` // Rtsp地址
	}
	_, err := s.post("/video/rtspurl/v1.0", map[string]any{
		"cameraCode": code,
		"mediaURLParam": map[string]any{
			"serviceType":   1,
			"streamType":    typ,
			"protocolType":  1,
			"transMode":     0,
			"broadCastType": 0,
		},
	}, &r)
	if err != nil {
		return "", err
	}
	err = s.StartPullStream(code, r.RtspUrl)
	if err != nil {
		return "", err
	}
	streamUri := s.GetStreamUrl(code, protocol)
	s.Cache.Set(key, streamUri, cache.NoExpiration)

	return streamUri, nil
}

func (s *IvsSvc) StopStream(code string, typ int) error {
	return s.StopPullStream(code)
}

func (s *IvsSvc) GetRecordUrl(code string, date time.Time) (string, error) {
	return "", fmt.Errorf("录像网址不存在")
}

func (s *IvsSvc) TakeRecord(code string, beginDateTime, endDateTime time.Time) (string, error) {
	var r struct {
		Uri string `json:"rtspURL"` // 地址
	}
	_, err := s.post("/video/rtspurl/v1.0", map[string]any{
		"cameraCode": code,
		"mediaURLParam": map[string]any{
			"serviceType": 3,
			"timeSpan": map[string]any{
				"startTime": beginDateTime,
				"end":       endDateTime,
			},
			"streamType":    1,
			"protocolType":  1,
			"transMode":     0,
			"broadCastType": 0,
		},
	}, &r)
	if err != nil {
		return "", err
	}

	return r.Uri, nil
}

func (s *IvsSvc) TakeSnapshot(code string, typ int) (string, error) {
	streamUri, err := s.StartStream(code, typ, "rtsp")
	if err != nil {
		return "", err
	}
	fileName := fmt.Sprintf("data/cameras/%s_%s.jpg", code, utl.TimeId())
	err = utl.TakeStream(streamUri, map[string]map[string]any{
		fileName: {
			"s":       "1920*1080",
			"vframes": 1,
		},
	})
	if err != nil {
		return "", err
	}

	return htp.GetUrl(fileName), nil
}

func (s *IvsSvc) TakeSnapshots(code string, cntSecs ...int) ([]string, error) {
	var imageUrls []string
	streamUri, err := s.StartStream(code, 1, "rtsp")
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
	err = utl.TakeStream(streamUri, map[string]map[string]any{
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

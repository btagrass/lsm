package ivs

import (
	"fmt"
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
	var r struct {
		RecordInfos struct {
			RecordInfoList []struct {
				RecordFileName string `json:"recordFileName"` // 录像文件名
			} `json:"recordInfoList"` // 录像信息列表
		} `json:"recordInfos"` // 录像信息集合
		Total int64 `json:"total"` // 总数
	}
	startTime := date.Format("20060102000000")
	endTime := date.Format("20060102235959")
	_, err := s.get(fmt.Sprintf("/platform/recordlist/0/%s//%s/%s/1/1", code, startTime, endTime), &r)
	if err != nil {
		return "", err
	}
	if r.Total == 0 {
		return "", fmt.Errorf("录像网址不存在")
	}

	return r.RecordInfos.RecordInfoList[0].RecordFileName, nil
}

func (s *IvsSvc) TakeSnapshot(code string, typ int) (string, error) {
	streamUrl, err := s.StartStream(code, typ, "rtsp")
	if err != nil {
		return "", err
	}
	fileName := fmt.Sprintf("data/cameras/%s_%s.jpg", code, utl.TimeId())
	err = utl.TakeStream(streamUrl, map[string]map[string]any{
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

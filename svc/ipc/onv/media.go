package onv

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/btagrass/go.core/htp"
	"github.com/btagrass/go.core/utl"
	"github.com/go-resty/resty/v2"
	"github.com/use-go/onvif/media"
	onvid "github.com/use-go/onvif/xsd/onvif"
	"github.com/yitter/idgenerator-go/idgen"
)

func (s *OnvSvc) StartStream(code string, typ int, protocol string) (string, error) {
	dev, err := s.getDevice(code)
	if err != nil {
		return "", err
	}
	doc, err := dev.call(media.GetStreamUri{
		StreamSetup: onvid.StreamSetup{
			Stream: "RTP-Unicast",
			Transport: onvid.Transport{
				Protocol: "RTSP",
				Tunnel:   nil,
			},
		},
		ProfileToken: onvid.ReferenceToken(dev.profileTokens[typ-1]),
	})
	if err != nil {
		return "", err
	}
	element := doc.FindElement("//trt:GetStreamUriResponse/trt:MediaUri/tt:Uri")
	if element == nil {
		return "", fmt.Errorf("摄像头 %s 的流地址不存在", code)
	}
	streamUrl := utl.Replace(element.Text(), "rtsp://", fmt.Sprintf("rtsp://%s:%s@", dev.userName, dev.password))
	err = s.StartPullStream(code, streamUrl)
	if err != nil {
		return "", err
	}

	return s.GetStreamUrl(code, protocol), nil
}

func (s *OnvSvc) StopStream(code string, typ int) error {
	err := s.StopPullStream(code)

	return err
}

func (s *OnvSvc) TakeSnapshot(code string, typ int) (string, error) {
	var url string
	dev, err := s.getDevice(code)
	if err != nil {
		return url, err
	}
	doc, err := dev.call(media.GetSnapshotUri{
		ProfileToken: onvid.ReferenceToken(dev.profileTokens[typ-1]),
	})
	if err != nil {
		return url, err
	}
	element := doc.FindElement("//trt:GetSnapshotUriResponse/trt:MediaUri/tt:Uri")
	if element == nil {
		return url, fmt.Errorf("device %s 's snapshotUri is not found", code)
	}
	url = element.Text()
	fileName := fmt.Sprintf("data/cameras/%s_%s.jpg", code, time.Now().Format("20060102150405.999999999"))
	req := resty.New().
		SetTimeout(htp.Timeout).
		SetTLSClientConfig(&tls.Config{
			InsecureSkipVerify: true,
		}).
		R().
		SetOutput(fileName)
	resp, err := req.Get(url)
	if err != nil {
		return url, err
	}
	if resp.StatusCode() == http.StatusUnauthorized {
		auths := make(map[string]string)
		strs := utl.Split(strings.TrimPrefix(resp.Header().Get("WWW-Authenticate"), "Digest "), '=', ',')
		for i := 0; i < len(strs); i += 2 {
			auths[strings.TrimSpace(strs[i])] = strings.Trim(strs[i+1], "\"")
		}
		if auths["algorithm"] == "" || auths["algorithm"] == "MD5" {
			ha1 := utl.Md5(fmt.Sprintf("%s:%s:%s", dev.userName, auths["realm"], dev.password))
			ha2 := utl.Md5(fmt.Sprintf("%s:%s", http.MethodGet, url))
			nc := "1"
			cnonce := idgen.NextId()
			response := utl.Md5(fmt.Sprintf("%s:%s:%s:%d:%s:%s", ha1, auths["nonce"], nc, cnonce, auths["qop"], ha2))
			authorization := fmt.Sprintf(`Digest username="%s",realm="%s",nonce="%s",qop="%s",uri="%s",nc="%s",cnonce="%d",response="%s"`, dev.userName, auths["realm"], auths["nonce"], auths["qop"], url, nc, cnonce, response)
			_, err = req.SetHeader("Authorization", authorization).Get(url)
			if err != nil {
				return url, err
			}
		}
	}
	url = htp.GetUrl(fileName)

	return url, nil
}

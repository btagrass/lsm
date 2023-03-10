package onv

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"

	"github.com/btagrass/go.core/htp"
	"github.com/btagrass/go.core/utl"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
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
	return s.StopPullStream(code)
}

func (s *OnvSvc) TakeSnapshot(code string, typ int) (string, error) {
	dev, err := s.getDevice(code)
	if err != nil {
		return "", err
	}
	doc, err := dev.call(media.GetSnapshotUri{
		ProfileToken: onvid.ReferenceToken(dev.profileTokens[typ-1]),
	})
	if err != nil {
		return "", err
	}
	element := doc.FindElement("//trt:GetSnapshotUriResponse/trt:MediaUri/tt:Uri")
	if element == nil {
		return "", fmt.Errorf("摄像头 %s 的快照地址不存在", code)
	}
	url := element.Text()
	fileName := fmt.Sprintf("data/cameras/%s_%s.jpg", code, utl.TimeId())
	req := resty.New().
		SetTimeout(htp.Timeout).
		SetTLSClientConfig(&tls.Config{
			InsecureSkipVerify: true,
		}).
		R().
		SetOutput(fileName)
	resp, err := req.Get(url)
	if err != nil {
		return "", err
	}
	if resp.StatusCode() == http.StatusUnauthorized {
		auths := make(map[string]string)
		was := utl.Split(strings.TrimPrefix(resp.Header().Get("Www-Authenticate"), "Digest "), '=', ',')
		for i := 0; i < len(was); i += 2 {
			auths[strings.TrimSpace(was[i])] = strings.Trim(was[i+1], "\"")
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
				return "", err
			}
		}
	}

	return htp.GetUrl(fileName), nil
}

// 开始流集合
func (s *OnvSvc) startStreams() error {
	cameras, _, err := s.ListCameras()
	if err != nil {
		return err
	}
	for _, c := range cameras {
		_, err = s.StartStream(c.Code, 1, "rtsp")
		if err != nil {
			logrus.Error(err)
		}
	}

	return nil
}

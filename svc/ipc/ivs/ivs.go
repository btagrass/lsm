package ivs

import (
	"crypto/tls"
	"fmt"
	"lsm/svc/ipc/internal"
	"lsm/svc/stream"
	"net/http"
	"time"

	"github.com/btagrass/go.core/htp"
	"github.com/go-resty/resty/v2"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
)

// 华为智能视频监控服务
type IvsSvc struct {
	*internal.CameraSvc
	*stream.StreamSvc
	addr          string // 地址
	appKey        string // 应用标识
	appSecret     string // 应用密钥
	eventCallback string // 事件回调
}

// 构造函数
func NewIvs(cameraSvc *internal.CameraSvc, streamSvc *stream.StreamSvc, addr, appKey, appSecret, eventCallback string) *IvsSvc {
	s := &IvsSvc{
		CameraSvc:     cameraSvc,
		StreamSvc:     streamSvc,
		addr:          addr,
		appKey:        appKey,
		appSecret:     appSecret,
		eventCallback: eventCallback,
	}
	go func() {
		delay := time.Minute
		t := time.NewTimer(delay)
		defer t.Stop()
		for {
			<-t.C
			err := s.keepalive()
			if err != nil {
				logrus.Error(err)
			}
			t.Reset(delay)
		}
	}()

	return s
}

// 删除
func (s *IvsSvc) delete(url string, r ...any) (*http.Response, error) {
	req := resty.New().
		SetBaseURL(fmt.Sprintf("https://%s", s.addr)).
		SetTimeout(htp.Timeout).
		SetTLSClientConfig(&tls.Config{
			InsecureSkipVerify: true,
		}).
		OnBeforeRequest(s.request).
		OnAfterResponse(s.respond).
		R().
		ForceContentType("application/json")
	if len(r) > 0 {
		req.SetResult(r[0])
	}
	resp, err := req.Delete(url)
	logrus.Debugf("method: %s, url: %s -> %s", req.Method, req.URL, resp)

	return resp.RawResponse, err
}

// 获取
func (s *IvsSvc) get(url string, r ...any) (*http.Response, error) {
	req := resty.New().
		SetBaseURL(fmt.Sprintf("https://%s", s.addr)).
		SetTimeout(htp.Timeout).
		SetTLSClientConfig(&tls.Config{
			InsecureSkipVerify: true,
		}).
		OnBeforeRequest(s.request).
		OnAfterResponse(s.respond).
		R().
		ForceContentType("application/json")
	if len(r) > 0 {
		req.SetResult(r[0])
	}
	resp, err := req.Get(url)
	logrus.Debugf("method: %s, url: %s -> %s", req.Method, req.URL, resp)

	return resp.RawResponse, err
}

// 提交
func (s *IvsSvc) post(url string, data any, r ...any) (*http.Response, error) {
	req := resty.New().
		SetBaseURL(fmt.Sprintf("https://%s", s.addr)).
		SetTimeout(htp.Timeout).
		SetTLSClientConfig(&tls.Config{
			InsecureSkipVerify: true,
		}).
		OnBeforeRequest(s.request).
		OnAfterResponse(s.respond).
		R().
		SetHeader("Accept", "*/*").
		SetHeader("Content-Type", "application/json").
		SetBody(data).
		ForceContentType("application/json")
	if len(r) > 0 {
		req.SetResult(r[0])
	}
	resp, err := req.Post(url)
	logrus.Debugf("method: %s, url: %s, data: %s -> %s", req.Method, req.URL, data, resp)

	return resp.RawResponse, err
}

// 提交
func (s *IvsSvc) put(url string, data any, r ...any) (*http.Response, error) {
	req := resty.New().
		SetBaseURL(fmt.Sprintf("https://%s", s.addr)).
		SetTimeout(htp.Timeout).
		SetTLSClientConfig(&tls.Config{
			InsecureSkipVerify: true,
		}).
		OnBeforeRequest(s.request).
		OnAfterResponse(s.respond).
		R().
		SetHeader("Accept", "*/*").
		SetHeader("Content-Type", "application/json").
		SetBody(data).
		ForceContentType("application/json")
	if len(r) > 0 {
		req.SetResult(r[0])
	}
	resp, err := req.Put(url)
	logrus.Debugf("method: %s, url: %s, data: %s -> %s", req.Method, req.URL, data, resp)

	return resp.RawResponse, err
}

// 请求
func (s *IvsSvc) request(c *resty.Client, req *resty.Request) error {
	key := fmt.Sprintf("%s:users:token", s.Prefix)
	v, ok := s.Cache.Get(key)
	if !ok {
		return fmt.Errorf("令牌不存在")
	}
	req.SetHeader("JSESSIONID", cast.ToString(v))

	return nil
}

// 响应
func (s *IvsSvc) respond(c *resty.Client, resp *resty.Response) error {
	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf(resp.Status())
	}

	return nil
}

// 保活
func (s *IvsSvc) keepalive() error {
	var r struct {
		Code int `json:"resultCode"` // 代码
	}
	_, err := s.get("/common/keepAlive", &r)
	if err == nil && r.Code == 0 {
		return nil
	}
	resp, err := s.post("/loginInfo/login/v1.0", map[string]any{
		"userName": s.appKey,
		"password": s.appSecret,
	}, &r)
	if err != nil {
		return err
	}
	if r.Code != 0 {
		return fmt.Errorf("保活失败：%d", r.Code)
	}
	key := fmt.Sprintf("%s:users:token", s.Prefix)
	s.Cache.Set(key, resp.Header.Get("JSESSIONID"), cache.NoExpiration)

	return nil
}

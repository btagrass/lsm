package isc

import (
	"crypto/tls"
	"fmt"
	"lsm/svc/ipc/internal"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/btagrass/go.core/htp"
	"github.com/btagrass/go.core/utl"
	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
)

var (
	reservedHeaders = "|X-Ca-Key|"
)

// 海康综合安防服务
type IscSvc struct {
	*internal.CameraSvc
	addr      string // 地址
	appKey    string // 应用标识
	appSecret string // 应用密钥
}

// 构造函数
func NewIsc(cameraSvc *internal.CameraSvc, addr, appKey, appSecret string) *IscSvc {
	is := &IscSvc{
		CameraSvc: cameraSvc,
		addr:      addr,
		appKey:    appKey,
		appSecret: appSecret,
	}
	go func() {
		delay := 30 * time.Second
		t := time.NewTimer(delay)
		defer t.Stop()
		for {
			<-t.C
			err := is.keepaliveCameras()
			if err != nil {
				logrus.Error(err)
			}
			t.Reset(delay)
		}
	}()

	return is
}

// 提交
func (is *IscSvc) post(url string, data any, r ...any) (string, error) {
	req := resty.New().
		SetBaseURL(fmt.Sprintf("https://%s", is.addr)).
		SetTimeout(htp.Timeout).
		SetTLSClientConfig(&tls.Config{
			InsecureSkipVerify: true,
		}).
		OnBeforeRequest(func(c *resty.Client, req *resty.Request) error {
			builder := strings.Builder{}
			builder.WriteString(req.Method)
			builder.WriteString("\n")
			builder.WriteString(req.Header.Get("Accept"))
			builder.WriteString("\n")
			builder.WriteString(req.Header.Get("Content-Type"))
			builder.WriteString("\n")
			builder.WriteString(req.Header.Get("Date"))
			builder.WriteString("\n")
			headers := []string{}
			for h := range req.Header {
				if strings.Contains(reservedHeaders, fmt.Sprintf("|%s|", h)) {
					headers = append(headers, strings.ToLower(h))
				}
			}
			sort.Strings(headers)
			req.SetHeader("X-Ca-Signature-Headers", strings.Join(headers, ","))
			for _, v := range headers {
				builder.WriteString(v)
				builder.WriteString(":")
				builder.WriteString(req.Header.Get(v))
				builder.WriteString("\n")
			}
			builder.WriteString(url)
			req.SetHeader("X-Ca-Signature", utl.HmacSha256(builder.String(), is.appSecret))

			return nil
		}).
		OnAfterResponse(func(c *resty.Client, resp *resty.Response) error {
			if resp.StatusCode() != http.StatusOK {
				return fmt.Errorf(resp.Status())
			} else {
				r := cast.ToStringMap(resp.String())
				code, ok := r["code"]
				if !ok {
					code, ok = r["error_code"]
				}
				if ok {
					code = cast.ToInt(code)
					if code != 0 && code != http.StatusOK {
						msg, ok := r["msg"]
						if !ok {
							msg = r["desp"]
						}
						return fmt.Errorf(cast.ToString(msg))
					}
				}
			}

			return nil
		}).
		R().
		SetHeader("Accept", "*/*").
		SetHeader("Content-Type", "application/json").
		SetHeader("Date", time.Now().Format(http.TimeFormat)).
		SetHeader("X-Ca-Key", is.appKey).
		SetBody(data).
		ForceContentType("application/json")
	if len(r) > 0 {
		req.SetResult(r[0])
	}
	resp, err := req.Post(url)

	return resp.String(), err
}

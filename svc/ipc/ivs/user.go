package ivs

import (
	"fmt"

	"github.com/patrickmn/go-cache"
)

// 保活
func (s *IvsSvc) keepalive() error {
	_, err := s.get("/common/keepAlive")
	if err == nil {
		return nil
	}
	resp, err := s.post("/loginInfo/login/v1.0", map[string]any{
		"userName": s.appKey,
		"password": s.appSecret,
	})
	if err != nil {
		return err
	}
	key := fmt.Sprintf("%s:users:token", s.Prefix)
	s.Cache.Set(key, resp.Header.Get("JSESSIONID"), cache.NoExpiration)

	return nil
}

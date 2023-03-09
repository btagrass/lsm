package internal

import (
	"lsm/mdl"

	"github.com/btagrass/go.core/svc"
)

// 摄像头服务
type CameraSvc struct {
	*svc.Svc[mdl.Camera]
}

// 构造函数
func NewCameraSvc() *CameraSvc {
	return &CameraSvc{
		Svc: svc.NewSvc[mdl.Camera]("lsm:cameras"),
	}
}

// 获取摄像头
func (s *CameraSvc) GetCamera(code string) (*mdl.Camera, error) {
	camera, err := s.Get("code = ?", code)

	return camera, err
}

// 获取摄像头集合
func (s *CameraSvc) ListCameras(conds ...any) ([]mdl.Camera, int64, error) {
	return s.List(conds...)
}

// 移除摄像头集合
func (s *CameraSvc) RemoveCameras(ids []string) error {
	return s.Remove(ids)
}

// 保存摄像头
func (s *CameraSvc) SaveCamera(camera mdl.Camera) error {
	return s.Save(camera)
}

package internal

import (
	"fmt"
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
func (s *CameraSvc) GetCamera(id int64) (*mdl.Camera, error) {
	return s.Get(id)
}

// 获取摄像头
func (s *CameraSvc) GetCameraByCode(code string) (*mdl.Camera, error) {
	camera, err := s.Get("code = ?", code)

	return camera, err
}

// 获取摄像头集合
func (s *CameraSvc) ListCameras(name string) ([]mdl.Camera, error) {
	cameras, err := s.List("name like ?", fmt.Sprintf("%%%s%%", name))

	return cameras, err
}

// 分页摄像头集合
func (s *CameraSvc) PageCameras(conds map[string]any) ([]mdl.Camera, int64, error) {
	return s.Page(conds)
}

// 获取摄像头集合
func (s *CameraSvc) RemoveCameras(ids []string) error {
	return s.Remove(ids)
}

// 保存摄像头
func (s *CameraSvc) SaveCamera(camera mdl.Camera) error {
	return s.Save(camera)
}

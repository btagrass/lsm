package ipc

import (
	"lsm/mdl"
	"lsm/svc/ipc/internal"
	"lsm/svc/ipc/onv"

	"github.com/spf13/viper"
)

// 网络摄像头服务接口
type IIpcSvc interface {
	// 摄像头
	GetCamera(id int64) (*mdl.Camera, error)                       // 获取摄像头
	GetCameraByCode(code string) (*mdl.Camera, error)              // 通过代码获取摄像头
	ListCameras(name string) ([]mdl.Camera, error)                 // 获取摄像头集合
	PageCameras(conds map[string]any) ([]mdl.Camera, int64, error) // 分页摄像头集合
	RemoveCameras(ids []string) error                              // 移除摄像头集合
	SaveCamera(camera mdl.Camera) error                            // 保存摄像头
	// 媒体
	StartStream(code string, typ int, protocol string) (string, error) // 开始流
	StopStream(code string, typ int) error                             // 停止流
	TakeSnapshot(code string, typ int) (string, error)                 // 抓取快照
	// 云台
	ControlPtz(code string, command string, speed int) error // 控制云台
}

// 构造函数
func NewIpcSvc() IIpcSvc {
	cameraSvc := internal.NewCameraSvc()
	lalSvc := internal.NewLalSvc()
	var s IIpcSvc
	typ := viper.GetString("svc.ipc.type")
	if typ == "onv" {
		s = onv.NewOnvSvc(cameraSvc, lalSvc)
	}

	return s
}

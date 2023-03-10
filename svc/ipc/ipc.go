package ipc

import (
	"lsm/mdl"
	"lsm/svc/ipc/internal"
	"lsm/svc/ipc/isc"
	"lsm/svc/ipc/ivs"
	"lsm/svc/ipc/onv"
	"lsm/svc/stream"
	"time"

	"github.com/spf13/viper"
)

// 网络摄像头服务接口
type IIpcSvc interface {
	// 摄像头
	GetCamera(code string) (*mdl.Camera, error)            // 获取摄像头
	ListCameras(conds ...any) ([]mdl.Camera, int64, error) // 获取摄像头集合
	RemoveCameras(ids []string) error                      // 移除摄像头集合
	SaveCamera(camera mdl.Camera) error                    // 保存摄像头
	SyncCameras() error                                    // 同步摄像头集合
	// 媒体
	StartStream(code string, typ int, protocol string) (string, error) // 开始流
	StopStream(code string, typ int) error                             // 停止流
	GetRecordUrl(code string, date time.Time) (string, error)          // 获取录像网址
	TakeSnapshot(code string, typ int) (string, error)                 // 抓取快照
	// 云台
	ControlPtz(code string, command string, speed int) error // 控制云台
	GotoPreset(code string, index int) error                 // 转到预置位
	ListPresets(code string) ([]mdl.Preset, error)           // 获取预置位集合
	RemovePreset(code string, index int) error               // 移除预置位
	SavePreset(preset mdl.Preset) error                      // 保存预置位
	// 事件
	NotifyEvent(content string) error // 通知事件
}

// 构造函数
func NewIpcSvc(streamSvc *stream.StreamSvc) IIpcSvc {
	cameraSvc := internal.NewCameraSvc()
	var s IIpcSvc
	typ := viper.GetString("svc.ipc.type")
	if typ == "onv" {
		s = onv.NewOnvSvc(
			cameraSvc,
			streamSvc,
			viper.GetString("svc.ipc.onv.userName"),
			viper.GetString("svc.ipc.onv.password"),
		)
	} else if typ == "isc" {
		s = isc.NewIsc(
			cameraSvc,
			streamSvc,
			viper.GetString("svc.ipc.isc.addr"),
			viper.GetString("svc.ipc.isc.appKey"),
			viper.GetString("svc.ipc.isc.appSecret"),
			viper.GetString("svc.ipc.isc.eventCallback"),
			viper.GetString("svc.ipc.isc.eventType"),
		)
	} else if typ == "ivs" {
		s = ivs.NewIvs(
			cameraSvc,
			streamSvc,
			viper.GetString("svc.ipc.ivs.addr"),
			viper.GetString("svc.ipc.ivs.appKey"),
			viper.GetString("svc.ipc.ivs.appSecret"),
			viper.GetString("svc.ipc.ivs.eventCallback"),
			viper.GetString("svc.ipc.ivs.eventType"),
		)
	}

	return s
}

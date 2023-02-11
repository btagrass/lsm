package svc

import (
	"lsm/mdl"
	"lsm/svc/ipc"
	"lsm/svc/videowall"

	"github.com/btagrass/go.core/svc"
	_ "github.com/btagrass/go.core/sys/svc"
	"github.com/sirupsen/logrus"
)

var (
	IpcSvc       ipc.IIpcSvc             // 网络摄像头服务
	VideoWallSvc *videowall.VideoWallSvc // 视频墙服务
)

// 初始化
func init() {
	// 迁移
	err := svc.Migrate(
		[]any{
			&mdl.Camera{},
			&mdl.VideoWall{},
		},
		"INSERT INTO `sys_resource` VALUES (300000000000007, '2023-01-29 00:00:00.000', NULL, NULL, NULL, 'lsm', '业务系统', 1, 'operation', '/lsm', NULL, 1)",
		"INSERT INTO `sys_resource` VALUES (300000000000009, '2023-01-29 00:00:00.000', NULL, NULL, 300000000000007, 'lsmVideoWalls', '视频墙管理', 1, 'monitor', '/lsm/videowalls', NULL, 1)",
		"INSERT INTO `sys_resource` VALUES (300000000000008, '2023-01-29 00:00:00.000', NULL, NULL, 300000000000007, 'lsmCameras', '摄像头管理', 1, 'camera', '/lsm/cameras', NULL, 2)",
	)
	if err != nil {
		logrus.Fatal(err)
	}
	// 服务
	IpcSvc = ipc.NewIpcSvc()
	VideoWallSvc = videowall.NewVideoWallSvc(IpcSvc)
}

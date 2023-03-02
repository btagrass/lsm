package svc

import (
	"lsm/mdl"
	"lsm/svc/ipc"
	"lsm/svc/stream"
	"lsm/svc/video"
	"lsm/svc/videowall"

	"github.com/btagrass/go.core/svc"
	_ "github.com/btagrass/go.core/sys/svc"
	"github.com/sirupsen/logrus"
)

var (
	IpcSvc       ipc.IIpcSvc             // 网络摄像头服务
	StreamSvc    *stream.StreamSvc       // 流服务
	VideoSvc     *video.VideoSvc         // 视频服务
	VideoWallSvc *videowall.VideoWallSvc // 视频墙服务
)

// 初始化
func init() {
	// 迁移
	err := svc.Migrate(
		[]any{
			&mdl.Camera{},
			&mdl.StreamPush{},
			&mdl.Video{},
			&mdl.VideoWall{},
		},
		"INSERT INTO sys_resource VALUES (300000000000002, '2023-01-29 00:00:00.000', NULL, NULL, 0, '业务系统', 1, 'Operation', '/mgt', NULL, 1)",
		"INSERT INTO sys_resource VALUES (300000000000201, '2023-01-29 00:00:00.000', NULL, NULL, 300000000000002, '视频墙管理', 1, 'Monitor', '/mgt/videowalls', NULL, 1)",
		"INSERT INTO sys_resource VALUES (300000000020101, '2023-01-29 00:00:00.000', NULL, NULL, 300000000000201, '查询', 2, '', '/mgt/videowalls', 'GET', 1)",
		"INSERT INTO sys_resource VALUES (300000000020102, '2023-01-29 00:00:00.000', NULL, NULL, 300000000000201, '删除', 2, '', '/mgt/videowalls/*', 'DELETE', 2)",
		"INSERT INTO sys_resource VALUES (300000000020103, '2023-01-29 00:00:00.000', NULL, NULL, 300000000000201, '编辑', 2, '', '/mgt/videowalls/*', 'GET', 3)",
		"INSERT INTO sys_resource VALUES (300000000020104, '2023-01-29 00:00:00.000', NULL, NULL, 300000000000201, '保存', 2, '', '/mgt/videowalls', 'POST', 4)",
		"INSERT INTO sys_resource VALUES (300000000000202, '2023-01-29 00:00:00.000', NULL, NULL, 300000000000002, '摄像头管理', 1, 'Camera', '/mgt/cameras', NULL, 2)",
		"INSERT INTO sys_resource VALUES (300000000020201, '2023-01-29 00:00:00.000', NULL, NULL, 300000000000202, '查询', 2, '', '/mgt/cameras', 'GET', 1)",
		"INSERT INTO sys_resource VALUES (300000000020202, '2023-01-29 00:00:00.000', NULL, NULL, 300000000000202, '删除', 2, '', '/mgt/cameras/*', 'DELETE', 2)",
		"INSERT INTO sys_resource VALUES (300000000020203, '2023-01-29 00:00:00.000', NULL, NULL, 300000000000202, '编辑', 2, '', '/mgt/cameras/*', 'GET', 3)",
		"INSERT INTO sys_resource VALUES (300000000020204, '2023-01-29 00:00:00.000', NULL, NULL, 300000000000202, '保存', 2, '', '/mgt/cameras', 'POST', 4)",
		"INSERT INTO sys_resource VALUES (300000000000203, '2023-01-29 00:00:00.000', NULL, NULL, 300000000000002, '视频管理', 1, 'VideoCamera', '/mgt/videos', NULL, 3)",
		"INSERT INTO sys_resource VALUES (300000000020301, '2023-01-29 00:00:00.000', NULL, NULL, 300000000000203, '查询', 2, '', '/mgt/videos', 'GET', 1)",
		"INSERT INTO sys_resource VALUES (300000000020302, '2023-01-29 00:00:00.000', NULL, NULL, 300000000000203, '删除', 2, '', '/mgt/videos/*', 'DELETE', 2)",
		"INSERT INTO sys_resource VALUES (300000000020303, '2023-01-29 00:00:00.000', NULL, NULL, 300000000000203, '编辑', 2, '', '/mgt/videos/*', 'GET', 3)",
		"INSERT INTO sys_resource VALUES (300000000020304, '2023-01-29 00:00:00.000', NULL, NULL, 300000000000203, '保存', 2, '', '/mgt/videos', 'POST', 4)",
		"INSERT INTO sys_resource VALUES (300000000000204, '2023-01-29 00:00:00.000', NULL, NULL, 300000000000002, '实时流管理', 1, 'Histogram', '/mgt/streams', NULL, 4)",
	)
	if err != nil {
		logrus.Fatal(err)
	}
	// 服务
	StreamSvc = stream.NewStreamSvc()
	IpcSvc = ipc.NewIpcSvc(StreamSvc)
	VideoSvc = video.NewVideoSvc()
	VideoWallSvc = videowall.NewVideoWallSvc(IpcSvc)
}

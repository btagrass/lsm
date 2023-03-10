package video

import (
	"lsm/mdl"
	"os/exec"

	"github.com/btagrass/go.core/svc"
	"github.com/shirou/gopsutil/v3/process"
)

// 视频服务
type VideoSvc struct {
	*svc.Svc[mdl.Video]
}

// 构造函数
func NewVideoSvc() *VideoSvc {
	return &VideoSvc{
		Svc: svc.NewSvc[mdl.Video]("lsm:videos"),
	}
}

// 开始虚拟流
func (s *VideoSvc) StartVirtualStream(id int64) (int, error) {
	video, err := s.Get(id)
	if err != nil {
		return 0, err
	}
	if video.Process > 0 {
		return video.Process, nil
	}
	cmd := exec.Command(
		"ffmpeg",
		"-re",
		"-stream_loop", "-1",
		"-i", video.Source,
		"-c:a", "copy",
		"-c:v", "copy",
		"-f", "rtsp",
		video.Url,
	)
	err = cmd.Start()
	if err != nil {
		return 0, err
	}
	video.Process = cmd.Process.Pid
	err = s.Save(*video)
	if err != nil {
		return 0, err
	}

	return video.Process, err
}

// 停止虚拟流
func (s *VideoSvc) StopVirtualStream(id int64) error {
	video, err := s.Get(id)
	if err != nil {
		return err
	}
	if video.Process == 0 {
		return nil
	}
	ps, err := process.NewProcess(int32(video.Process))
	if err != nil {
		return err
	}
	err = ps.Kill()
	if err != nil {
		return err
	}
	// err = syscall.Kill(video.Process, syscall.SIGQUIT)
	video.Process = 0
	err = s.Save(*video)

	return err
}

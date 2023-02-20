package video

import (
	"lsm/mdl"
	"os/exec"
	"syscall"

	"github.com/btagrass/go.core/svc"
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

// 开始视频
func (s *VideoSvc) StartVideo(id int64) (int, error) {
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

// 停止视频
func (s *VideoSvc) StopVideo(id int64) error {
	video, err := s.Get(id)
	if err != nil {
		return err
	}
	if video.Process == 0 {
		return nil
	}
	err = syscall.Kill(video.Process, syscall.SIGQUIT)
	if err != nil {
		return err
	}
	video.Process = 0
	err = s.Save(*video)

	return err
}

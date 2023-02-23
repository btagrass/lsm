package stream

import (
	"sync"
	"time"

	"github.com/q191201771/lal/pkg/base"
	"github.com/q191201771/lal/pkg/rtprtcp"
	"github.com/q191201771/lal/pkg/rtsp"
	"github.com/q191201771/lal/pkg/sdp"
	"github.com/q191201771/naza/pkg/nazaerrors"
	"github.com/q191201771/naza/pkg/nazalog"
	"github.com/q191201771/naza/pkg/unique"
)

type RtspTunnel struct {
	pullUrl       string
	pushUrl       string
	pullOverTcp   bool
	pushOverTcp   bool
	uniqueKey     string
	rtpPacketChan chan rtprtcp.RtpPacket
	disposeOnce   sync.Once
	waitChan      chan error
	pullSession   *rtsp.PullSession
	pushSession   *rtsp.PushSession
}

func NewRtspTunnel(pullUrl string, pushUrl string, pullOverTcp bool, pushOverTcp bool) *RtspTunnel {
	uniqueKey := unique.GenUniqueKey("RTSPTUNNEL")
	nazalog.Debugf("[%s] lifecycle new RtspTunnel. pullUrl=%s, pushUrl=%s", uniqueKey, pullUrl, pushUrl)
	return &RtspTunnel{
		pullUrl:       pullUrl,
		pushUrl:       pushUrl,
		pullOverTcp:   pullOverTcp,
		pushOverTcp:   pushOverTcp,
		uniqueKey:     uniqueKey,
		rtpPacketChan: make(chan rtprtcp.RtpPacket, 1024),
		waitChan:      make(chan error, 1),
	}
}

// Start 开启任务，阻塞直到任务开启成功或失败。
//
// @return: 如果为nil，表示任务启动成功，此时数据已经在后台转发
func (r *RtspTunnel) Start() error {
	r.pullSession = rtsp.NewPullSession(r, func(option *rtsp.PullSessionOption) {
		option.PullTimeoutMs = 5000
		option.OverTcp = r.pullOverTcp
	})
	if err := r.pullSession.Pull(r.pullUrl); err != nil {
		nazalog.Errorf("[%s] start pull failed. err=%+v, url=%s", r.uniqueKey, err, r.pullUrl)
		return err
	}
	sdpCtx := r.pullSession.GetSdp()
	nazalog.Debugf("[%s] start pull succ. sdp=%s", r.uniqueKey, string(sdpCtx.RawSdp))

	r.pushSession = rtsp.NewPushSession(func(option *rtsp.PushSessionOption) {
		option.PushTimeoutMs = 5000
		option.OverTcp = r.pushOverTcp
	})
	if err := r.pushSession.Push(r.pushUrl, sdpCtx); err != nil {
		nazalog.Errorf("[%s] start push failed. err=%+v, url=%s", r.uniqueKey, err, r.pushUrl)
		return err
	}
	nazalog.Debugf("[%s] start push succ.", r.uniqueKey)

	go func() {
		t := time.NewTicker(1 * time.Second)
		defer t.Stop()
		for {
			select {
			case pkt := <-r.rtpPacketChan:
				_ = r.pushSession.WriteRtpPacket(pkt)
			case err := <-r.pullSession.WaitChan():
				nazalog.Debugf("[%s] < pullSession.Wait(). err=%+v", r.uniqueKey, err)
				_ = r.dispose(err)
				return
			case err := <-r.pushSession.WaitChan():
				nazalog.Debugf("[%s] < pushSession.Wait(). err=%+v", r.uniqueKey, err)
				_ = r.dispose(err)
				return
			case <-t.C:
				r.pullSession.UpdateStat(1)
				pullStat := r.pullSession.GetStat()
				r.pushSession.UpdateStat(1)
				pushStat := r.pushSession.GetStat()
				nazalog.Debugf("[%s] stat. pull=%+v, push=%+v", r.uniqueKey, pullStat, pushStat)
			}
		}
	}()

	return nil
}

// Dispose 主动关闭tunnel时调用
//
// 注意，只有 Start 成功后的tunnel才能调用，否则行为未定义
//
// 更详细的说明参考 IClientSessionLifecycle interface
func (r *RtspTunnel) Dispose() error {
	return r.dispose(nil)
}

// WaitChan Start 成功后，可使用这个channel来接收tunnel结束的消息
//
// 更详细的说明参考 IClientSessionLifecycle interface
func (r *RtspTunnel) WaitChan() chan error {
	return r.waitChan
}

// ---------------------------------------------------------------------------------------------------------------------

func (r *RtspTunnel) OnRtpPacket(pkt rtprtcp.RtpPacket) {
	r.rtpPacketChan <- pkt
}

func (r *RtspTunnel) OnSdp(sdpCtx sdp.LogicContext) {
	// noop
}

func (r *RtspTunnel) OnAvPacket(pkt base.AvPacket) {
	// noop
}

func (r *RtspTunnel) dispose(err error) error {
	var retErr error
	r.disposeOnce.Do(func() {
		nazalog.Infof("[%s] lifecycle dispose RtspTunnel.", r.uniqueKey)
		e1 := r.pullSession.Dispose()
		e2 := r.pushSession.Dispose()
		retErr = nazaerrors.CombineErrors(e1, e2)
		r.waitChan <- err
	})
	return retErr
}

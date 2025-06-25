package keepalive

import "time"

type KeepAlive struct {
	interval    int64
	timeout     int64
	timer       *time.Timer
	sendFunc    func()
	timeoutFunc func()
	stop        chan struct{}
	pong        chan struct{}
}

func New(interval int64, timeout int64) *KeepAlive {
	return &KeepAlive{
		interval: interval,
		timeout:  timeout,
		timer:    time.NewTimer(0),
		stop:     make(chan struct{}),
		pong:     make(chan struct{}),
	}
}
func (k *KeepAlive) Start() {
	timer := time.NewTimer(time.Duration(k.interval) * time.Second)
	go func() {
		for {
			select {
			case <-k.stop:
				timer.Stop()
				return
			case <-timer.C:
				k.sendPing()
				timer.Reset(time.Duration(k.interval) * time.Second)
			}
		}
	}()
}
func (k *KeepAlive) Stop() {
	k.stop <- struct{}{}
}
func (k *KeepAlive) PingFunc(f func()) {
	k.sendFunc = f
}
func (k *KeepAlive) HandlePong() {
	k.pong <- struct{}{}
}
func (k *KeepAlive) TimeoutFunc(f func()) {
	k.timeoutFunc = f
}
func (k *KeepAlive) sendPing() {
	k.sendFunc()
	k.timer.Stop()
	k.timer.Reset(time.Duration(k.timeout) * time.Second)
	go func() {
		select {
		case <-k.timer.C:
			k.timeoutFunc()
		case <-k.stop:
			return
		case <-k.pong:
			return
		}
	}()
}

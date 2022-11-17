package monitor

import (
	"time"

	"github.com/jucardi/go-titan/logx"
)

var (
	// To validate the interface implementation at compile time.
	_ IMonitor = (*service)(nil)
)

type service struct {
	updateInterval time.Duration
	watchers       []IWatcher
	running        bool
	throttle       *time.Ticker
	stop           chan struct{}
}

func (s *service) AddWatcher(watcher IWatcher) {
	s.watchers = append(s.watchers, watcher)
}

func (s *service) Start() {
	if s.throttle != nil {
		return
	}

	logx.Info("Starting monitor service")
	s.throttle = time.NewTicker(s.updateInterval)
	s.stop = make(chan struct{})
	s.runWatchers()

	for {
		select {
		case <-s.throttle.C:
			s.runWatchers()
		case <-s.stop:
			return
		}
	}
}

func (s *service) StartAsync() {
	go s.Start()
}

func (s *service) Stop() {
	if s.throttle == nil {
		return
	}
	s.throttle.Stop()
	s.stop <- struct{}{}
}

func (s *service) runWatchers() {
	for _, w := range s.watchers {
		w.Run()
	}
}

package reconciliation

import (
	"github.com/rs/zerolog/log"
	"time"
)

type Scheduler interface {
	Run(n int, d time.Duration) error
}

func NewScheduler(r Reconciler) Scheduler {
	return &scheduler{
		r: r,
	}
}

type scheduler struct {
	cancel chan struct{}
	r      Reconciler
}

func (s *scheduler) Run(n int, d time.Duration) error {
	for i := 0; i < n; i++ {
		go s.worker(d)
	}

	return nil
}

func (s *scheduler) worker(d time.Duration) {
	t := time.NewTicker(d)

	for {
		select {
		case <-s.cancel:
			return
		default:
			worked, _ := s.r.Run()
			if !worked {
				log.Trace().Msg("nothing to reconcile, sleeping...")
				t.Reset(d)
				select {
				case <-s.cancel:
					return
				case <-t.C:
				}
			}
		}
	}
}

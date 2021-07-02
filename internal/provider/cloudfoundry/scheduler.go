package cloudfoundry

import "time"

type ReconciliationScheduler interface {
	Run(n int, d time.Duration) error
}

func NewScheduler(r Reconciler) ReconciliationScheduler {
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

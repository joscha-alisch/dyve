package reconciliation

import (
	"sync"
	"testing"
	"time"
)

func TestScheduler(t *testing.T) {
	r := &fakeReconciler{
		mux: &sync.Mutex{},
	}
	s := NewScheduler(r).(*scheduler)

	_ = s.Run(1, 10*time.Millisecond)
	time.Sleep(120 * time.Millisecond)
	close(s.cancel)

	r.mux.Lock()
	if r.times < 10 || r.times > 15 {
		t.Error("expected reconciler to have been called roughly 12 times")
	}
	r.mux.Unlock()
}

type fakeReconciler struct {
	ok    bool
	err   error
	times int
	mux   *sync.Mutex
}

func (f *fakeReconciler) Run() (bool, error) {
	f.mux.Lock()
	f.times++
	f.mux.Unlock()
	return f.ok, f.err
}

func (f *fakeReconciler) Handler(t Type, h ReconcileHandler) {}

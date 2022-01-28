package reconciliation

import (
	"testing"
	"time"
)

func TestScheduler(t *testing.T) {
	r := &fakeReconciler{}
	s := NewScheduler(r).(*scheduler)

	_ = s.Run(1, 10*time.Millisecond)
	time.Sleep(120 * time.Millisecond)
	close(s.cancel)

	if r.times < 10 || r.times > 15 {
		t.Error("expected reconciler to have been called roughly 12 times")
	}
}

type fakeReconciler struct {
	ok    bool
	err   error
	times int
}

func (f *fakeReconciler) Run() (bool, error) {
	f.times++
	return f.ok, f.err
}

func (f *fakeReconciler) Handler(t Type, h ReconcileHandler) {}

package reconciliation

import (
	"github.com/google/go-cmp/cmp"
	"testing"
	"time"
)

func Test_Runner(t *testing.T) {
	p := &fakeJobProvider{}
	r := NewReconciler(p, 2*time.Minute)

	ok, err := r.Run()
	if err != nil {
		t.Error("didn't expect error ", err)
	}
	if ok {
		t.Error("didn't expect to work to have been done")
	}

	var triggered bool
	var recJob Job
	h := ReconcileHandler(func(j Job) error {
		triggered = true
		recJob = j
		return nil
	})
	r.Handler("someType", h)

	ok, err = r.Run()
	if err != nil {
		t.Error("didn't expect error ", err)
	}
	if ok {
		t.Error("didn't expect to work to have been done")
	}
	if triggered {
		t.Error("didn't expect handler to have been triggered")
	}

	p.job = Job{
		Type: "someType",
		Guid: "a",
	}
	p.ok = true

	ok, err = r.Run()
	if err != nil {
		t.Error("didn't expect error ", err)
	}
	if !ok {
		t.Error("expected work to have been done")
	}
	if !triggered {
		t.Error("expected handler to have been triggered")
	}
	if !cmp.Equal(p.job, recJob) {
		t.Errorf("job mismatch: %s\n", cmp.Diff(p.job, recJob))
	}

	p.job = Job{
		Type: "someOtherType",
		Guid: "a",
	}
	p.ok = true

	triggered = false
	ok, err = r.Run()
	if err != nil {
		t.Error("didn't expect error ", err)
	}
	if !ok {
		t.Error("expected work to have been done")
	}
	if triggered {
		t.Error("did not expect handler to have been triggered")
	}
}

type fakeJobProvider struct {
	job      Job
	ok       bool
	recorded time.Duration
}

func (f *fakeJobProvider) AcceptReconcileJob(olderThan time.Duration) (Job, bool) {
	f.recorded = olderThan
	return f.job, f.ok
}

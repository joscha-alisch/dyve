package cloudfoundry

import (
	"errors"
	"fmt"
)

var errNotFound = errors.New("not found")

type errReconcileFailed struct {
	Err error
	Job ReconcileJob
}

func (r *errReconcileFailed) Is(target error) bool {
	if rFailed, ok := target.(*errReconcileFailed); ok {
		return rFailed.Job.Type == r.Job.Type &&
			rFailed.Job.Guid == r.Job.Guid && errors.Is(rFailed.Err, r.Err)
	}
	return false
}

func (r *errReconcileFailed) Unwrap() error {
	return r.Err
}

func (r *errReconcileFailed) Error() string {
	t := ""
	switch r.Job.Type {
	case ReconcileSpaces:
		t = "org"
	case ReconcileApps:
		t = "space"
	}
	return fmt.Sprintf("%s reconcile failed for guid '%s': %s", t, r.Job.Guid, r.Err)
}

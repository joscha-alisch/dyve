package cloudfoundry

import (
	"errors"
	"github.com/joscha-alisch/dyve/internal/reconciliation"
	"testing"
)

func TestErrFailedUnwrap(t *testing.T) {
	wrapped := errors.New("wrapped error")
	err := &errReconcileFailed{Err: wrapped}

	if errors.Unwrap(err) != wrapped {
		t.Fatal("expected wrapped error, but was ", errors.Unwrap(err))
	}
}

func TestErrSpaceFailed_Error(t *testing.T) {
	wrapped := errors.New("wrapped error")
	err := &errReconcileFailed{Err: wrapped, Job: reconciliation.Job{Guid: "a", Type: ReconcileSpaces}}

	if err.Error() != "org reconcile failed for guid 'a': wrapped error" {
		t.Fatal("wrong error string: ", err.Error())
	}
}
func TestErrOrgFailed_Error(t *testing.T) {
	wrapped := errors.New("wrapped error")
	err := &errReconcileFailed{Err: wrapped, Job: reconciliation.Job{Guid: "a", Type: ReconcileApps}}

	if err.Error() != "space reconcile failed for guid 'a': wrapped error" {
		t.Fatal("wrong error string: ", err.Error())
	}
}

func TestErrFailed_Is(t *testing.T) {
	wrapped := errors.New("wrapped error")

	err := &errReconcileFailed{Err: wrapped}

	if !errors.Is(err, &errReconcileFailed{Err: wrapped}) {
		t.Fatal("should be itself, but isnt")
	}

	if errors.Is(&errReconcileFailed{Err: wrapped}, errors.New("some other error")) {
		t.Fatal("should not be with a different wrapped error")
	}

	if !errors.Is(err, &errReconcileFailed{Err: nil}) {
		t.Fatal("should be with a nil wrapped error")
	}
}

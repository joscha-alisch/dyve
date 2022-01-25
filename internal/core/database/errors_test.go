package database

import (
	"errors"
	"testing"
)

func TestMongoFailedUnwrap(t *testing.T) {
	wrapped := errors.New("wrapped error")
	err := mongoFailed(wrapped)

	if errors.Unwrap(err) != wrapped {
		t.Fatal("expected wrapped error, but was ", errors.Unwrap(err))
	}
}

func TestErrMongoQueryFailed_Error(t *testing.T) {
	wrapped := errors.New("wrapped error")
	err := mongoFailed(wrapped)

	if err.Error() != "something went wrong when querying MongoDB: wrapped error" {
		t.Fatal("wrong error string: ", err.Error())
	}
}

func TestErrMongoQueryFailed_Is(t *testing.T) {
	wrapped := errors.New("wrapped error")

	err := mongoFailed(wrapped)

	if !errors.Is(err, mongoFailed(wrapped)) {
		t.Fatal("should be itself, but isnt")
	}

	if errors.Is(mongoFailed(err), errors.New("some other error")) {
		t.Fatal("should not be with a different wrapped error")
	}

	if !errors.Is(err, mongoFailed(nil)) {
		t.Fatal("should be with a nil wrapped error")
	}
}

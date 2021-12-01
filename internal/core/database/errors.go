package database

import (
	"errors"
	"fmt"
)

var ErrNotFound = errors.New("element not found in collection")
var ErrExists = errors.New("element already exists")

type ErrMongoQueryFailed struct {
	Err error
}

func mongoFailed(err error) error {
	return &ErrMongoQueryFailed{Err: err}
}

func (r *ErrMongoQueryFailed) Is(target error) bool {
	if rFailed, ok := target.(*ErrMongoQueryFailed); ok {
		return errors.Is(rFailed.Err, r.Err)
	}
	return false
}

func (r *ErrMongoQueryFailed) Unwrap() error {
	return r.Err
}

func (r *ErrMongoQueryFailed) Error() string {
	return fmt.Sprintf("something went wrong when querying MongoDB: %s", r.Err.Error())
}

package sdk

import "errors"

var ErrNotFound = errors.New("not found")
var ErrPageExceeded = errors.New("specified page exceeds total pages")
var ErrQueryPageMalformed = errors.New("query parameter 'page' is malformed")
var ErrQueryPerPageMalformed = errors.New("query parameter 'perPage' is malformed")
var ErrInternal = errors.New("internal error occurred")

// internal
var errExpectedQueryParamMissing = errors.New("query parameter was expected but is missing")
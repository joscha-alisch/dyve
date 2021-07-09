package provider

import "errors"

var ErrNotFound = errors.New("provider not found")
var ErrExists = errors.New("provider with id already exists")
var ErrNil = errors.New("can't add nil provider")

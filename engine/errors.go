package engine

import "errors"

var (
	ErrInvalidConfig = errors.New("invalid configuration")
	ErrNotRunning    = errors.New("keyboard engine not running")
)

package split

import "errors"

var (
	ErrInsufficientData   = errors.New("insufficient data")
	ErrInvalidHeader      = errors.New("invalid header")
	ErrInvalidChecksum    = errors.New("invalid checksum")
	ErrInvalidMessageType = errors.New("invalid message type")
	ErrTimeout            = errors.New("timeout")
	ErrUARTNotConfigured  = errors.New("UART not configured")
)

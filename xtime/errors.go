package xtime

import "errors"

var (
	ErrNotImplemented      = errors.New("xtime: not implemented")
	ErrInvalidFormatString = errors.New("xtime: invalid format string")
)

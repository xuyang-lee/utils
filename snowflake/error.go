package snowflake

import "errors"

var (
	ErrInvalidWorkerId     = errors.New("worker Id out of range")
	ErrInvalidDatacenterId = errors.New("datacenter Id out of range")
	ErrClockMovedBackwards = errors.New("clock moved backwards")
)

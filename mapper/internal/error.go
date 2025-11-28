package internal

import "errors"

var (
	ErrInvalidMappingDestination = errors.New("mapping destination must be non-nil and addressable")
	ErrInvalidMappingSource      = errors.New("mapping source must be non-nil and addressable")
	ErrMapKeyNotMatch            = errors.New("map's key type doesn't match")
	ErrBadDstAndSrc              = errors.New("bad mapping src and dst")
	ErrNotSupported              = errors.New("not supported")
	ErrNestedTooDeep             = errors.New("nested too deep")
)

package internal

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	ErrInvalidDst               = errors.New("dst should be a struct pointer or slice pointer of struct")
	ErrInvalidSrc               = errors.New("src should be a map, a struct, a slice of structs, or a pointer to one of the aforementioned types")
	ErrUnknownDstSrcCombination = errors.New("unknown combination")
	ErrMapKeyNotString          = errors.New("map key should be string")
	ErrNotMap                   = errors.New("not a map")
	ErrNotStruct                = errors.New("not a struct")
	ErrNotSlice                 = errors.New("not a slice")
)

func CombinationError(dstVal, srcVal reflect.Value) error {
	return fmt.Errorf("%w of dst[%v] and src[%v] types", ErrUnknownDstSrcCombination, dstVal.Kind().String(), srcVal.Kind().String())
}

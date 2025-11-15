package mapper

import (
	"github.com/xuyang-lee/utils/mapper/internal"
	"reflect"
)

func Mapper(dst, src any) error {
	dstVal := reflect.ValueOf(dst)
	srcVal := reflect.ValueOf(src)
	var err error

	// Verify the destination value. now dstVal is struct or slice
	dstVal, err = internal.VerifyDst(dstVal)
	if err != nil {
		return err
	}

	// Verify the source value. now srcVal is struct, map, or slice
	srcVal, err = internal.VerifySrc(srcVal)
	if err != nil {
		return err
	}

	switch {
	case srcVal.Kind() == reflect.Slice && dstVal.Kind() == reflect.Slice:
		return mapperSlice(dstVal, srcVal)
	case srcVal.Kind() == reflect.Map && dstVal.Kind() == reflect.Struct:
		return mapperMap(dstVal, srcVal)
	case srcVal.Kind() == reflect.Struct && dstVal.Kind() == reflect.Struct:
		return mapperStruct(dstVal, srcVal)
	default:
		return internal.CombinationError(dstVal, srcVal)
	}

}

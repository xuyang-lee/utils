package internal

import "reflect"

func VerifyDst(dstVal reflect.Value) (reflect.Value, error) {

	if dstVal.Kind() != reflect.Ptr { //不是指针
		return dstVal, ErrInvalidDst
	}

	for dstVal.Kind() == reflect.Ptr { //循环取指针，直到得到最原本的底层结构
		dstVal = dstVal.Elem() //取指针的值
	}

	// 底层结构不是结构体也不是切片
	// 或 是切片，但切片的元素不是结构体
	if (dstVal.Kind() != reflect.Struct && dstVal.Kind() != reflect.Slice) ||
		(dstVal.Kind() == reflect.Slice && dstVal.Type().Elem().Kind() != reflect.Struct) {
		return dstVal, ErrInvalidDst
	}

	// 返回底层结构
	return dstVal, nil
}

func VerifySrc(srcVal reflect.Value) (reflect.Value, error) {

	if srcVal.Kind() != reflect.Ptr && srcVal.Kind() != reflect.Slice && srcVal.Kind() != reflect.Struct && srcVal.Kind() != reflect.Map {
		return srcVal, ErrInvalidSrc
	}

	for srcVal.Kind() == reflect.Ptr { //循环取指针，直到得到最原本的底层结构
		srcVal = srcVal.Elem() //取指针的值
	}

	// 底层结构不是结构体/切片/map
	// 或 是切片，但切片的元素不是结构体
	if (srcVal.Kind() != reflect.Struct && srcVal.Kind() != reflect.Slice && srcVal.Kind() != reflect.Map) ||
		(srcVal.Kind() == reflect.Slice && srcVal.Type().Elem().Kind() != reflect.Struct) {
		return srcVal, ErrInvalidSrc
	}

	return srcVal, nil
}

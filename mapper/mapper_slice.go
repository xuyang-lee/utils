package mapper

import (
	"reflect"
)

func mapperSlice(dstVal, srcVal reflect.Value) error {

	if dstVal.Type().Elem().AssignableTo(srcVal.Type().Elem()) {
		// 基础元素相同可直接辅助
		for i := 0; i < srcVal.Len(); i++ {
			srcElem := srcVal.Index(i)
			dstElem := reflect.New(dstVal.Type().Elem()).Elem()
			if dstElem.CanSet() {
				dstElem.Set(srcElem)
			}
			dstVal.Set(reflect.Append(dstVal, dstElem))
		}
	} else {
		// 基础元素不同
		for i := 0; i < srcVal.Len(); i++ {
			srcElem := srcVal.Index(i)
			dstElem := reflect.New(dstVal.Type().Elem()).Elem()
			if err := mapperElem(dstElem, srcElem); err != nil {
				return err
			}
			dstVal.Set(reflect.Append(dstVal, dstElem))
		}
	}
	return nil
}

func mapperElem(dstElem, srcElem reflect.Value) error {
	if dstElem.Type().AssignableTo(srcElem.Type()) {
		//能直接赋值就直接赋值
		dstElem.Set(srcElem)
		return nil
	}

	//todo 不同类型的元素，需要进一步映射
	if dstElem.Kind() == reflect.Ptr {
		dstElem = dstElem.Elem()
	}
	// 处理一级指针后，能赋值的赋值
	if dstElem.Type().AssignableTo(srcElem.Type()) {
		//能直接赋值就直接赋值
		dstElem.Set(srcElem)
		return nil
	}

	//处理指针后仍然不是结构体或slice，不做处理
	if dstElem.Kind() != reflect.Struct || dstElem.Kind() != reflect.Slice {
		return nil
	}

	//不同于Mapper，在映射功能的内部，可能出现各种类型的组合，所以这里没有default
	switch {
	case srcElem.Kind() == reflect.Struct && dstElem.Kind() == reflect.Struct:
		return mapperStruct(dstElem, srcElem)
	case srcElem.Kind() == reflect.Map && dstElem.Kind() == reflect.Struct:
		return mapperMap(dstElem, srcElem)
	case srcElem.Kind() == reflect.Slice && dstElem.Kind() == reflect.Slice:
		return mapperSlice(dstElem, srcElem)
	}

	return nil
}

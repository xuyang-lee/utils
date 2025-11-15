package mapper

import (
	"github.com/xuyang-lee/utils/mapper/internal"
	"reflect"
)

func mapperMap(dstVal, srcVal reflect.Value) error {
	if srcVal.Type().Kind() != reflect.Map {
		return internal.ErrNotMap
	}
	if srcVal.Type().Key().Kind() != reflect.String {
		return internal.ErrMapKeyNotString
	}

	for i := 0; i < dstVal.NumField(); i++ {
		fieldInfo := dstVal.Type().Field(i)
		fieldVal := dstVal.Field(i)

		// Skip if field cannot be set (unexported fields are not settable).
		if !fieldVal.CanSet() {
			continue
		}

		tag := fieldInfo.Tag.Get(internal.Tag)
		tag, ignore := internal.ParamsTag(tag)
		if ignore {
			continue
		}

		if tag == "" {
			tag = fieldInfo.Name
		}

		// Get the value associated with the tag or field name from srcVal.
		mapValue := srcVal.MapIndex(reflect.ValueOf(tag))
		if !mapValue.IsValid() {
			continue
		}

		// 在 map 中有这个 key，并且字段类型与 map 中的类型相匹配，就设置它的值
		if mapValue.Type().AssignableTo(fieldVal.Type()) {
			fieldVal.Set(mapValue)
			continue
		}

		//同名情况下，如果类型不匹配，只看看是不是dst是不是一层指针
		if fieldVal.Type().Kind() == reflect.Ptr {
			// 一层指针，可直接赋值：直接赋值
			if mapValue.Type().AssignableTo(fieldVal.Type().Elem()) {
				newPtr := reflect.New(fieldVal.Type().Elem())
				newPtr.Elem().Set(mapValue)
				fieldVal.Set(newPtr)
				continue
			}
			// 一层指针，dst是结构体且，mapValue是map：递归
			if fieldVal.Type().Elem().Kind() == reflect.Struct && mapValue.Kind() == reflect.Map {
				newPtr := reflect.New(fieldVal.Type().Elem())
				if err := mapperMap(newPtr.Elem(), mapValue); err != nil {
					return err
				}
				fieldVal.Set(newPtr)
				continue
			}

			//也不是错了，从string unmarshal过来的map[string]interface{}，如果有null值，会被解析成nil，但跳过就行了
			continue
			//什么，还执行到了到这里？报错吧你
			//return internal.CombinationError(fieldVal, mapValue)
		}

	}
	return nil

}

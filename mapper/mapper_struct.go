package mapper

import (
	"github.com/xuyang-lee/utils/mapper/internal"
	"reflect"
)

func mapperStruct(dstVal, srcVal reflect.Value) error {

	dstType := dstVal.Type()
	dstFieldNameToValue := make(map[string]reflect.Value)
	for i := 0; i < dstVal.Type().NumField(); i++ {
		fieldInfo := dstType.Field(i)
		tag := fieldInfo.Tag.Get(internal.Tag)
		tag, ignore := internal.ParamsTag(tag)
		if ignore {
			continue
		}

		if tag == "" {
			dstFieldNameToValue[fieldInfo.Name] = dstVal.Field(i)
			continue
		} else {
			dstFieldNameToValue[tag] = dstVal.Field(i)
		}

	}

	srcType := srcVal.Type()
	for i := 0; i < srcVal.NumField(); i++ {
		fieldInfo := srcType.Field(i)
		field := srcVal.Field(i)
		tag := fieldInfo.Tag.Get(internal.Tag)
		tag, ignore := internal.ParamsTag(tag)
		if ignore {
			continue
		}

		if tag == "" {
			tag = fieldInfo.Name
		}

		if dstFieldVal, ok := dstFieldNameToValue[tag]; ok && dstFieldVal.CanSet() {
			// 如果字段类型相同，直接赋值
			if fieldInfo.Type.AssignableTo(dstFieldVal.Type()) {
				dstFieldVal.Set(field)
				continue
			}

			// dstFieldVal是指针，且指针指向的类型与srcVal.Field的类型相同
			if dstFieldVal.Type().Kind() == reflect.Ptr && dstFieldVal.Type().Elem().AssignableTo(fieldInfo.Type) {
				newPtr := reflect.New(fieldInfo.Type)
				newPtr.Elem().Set(field)
				dstFieldVal.Set(newPtr)
				continue
			}

			// srcVal.Field是指针，且指针指向的类型与dstFieldVal的类型相同
			if field.Type().Kind() == reflect.Ptr && fieldInfo.Type.Elem().AssignableTo(dstFieldVal.Type()) {
				dstFieldVal.Set(srcVal.Field(i).Elem())
				continue
			}

		}
	}

	return nil
}

package internal

import "reflect"

func GetFieldByTag(s reflect.Value, tag string) (reflect.Value, reflect.StructField, bool) {

	var ok bool

	if s.Kind() != reflect.Struct {
		return reflect.Value{}, reflect.StructField{}, false
	}

	t := s.Type()

	// 外层优先匿名
	for i := 0; i < t.NumField(); i++ {
		sField := t.Field(i)
		fVal := s.Field(i)

		tagName, ignore := ParseTag(sField.Tag.Get(tag))
		if ignore {
			continue
		}

		if tag == tagName {
			return fVal, sField, true
		}

		// 内嵌
		if sField.Anonymous {
			inner := fVal
			innerType, innerDepth := IndirectType(fVal.Type())
			if innerType.Kind() != reflect.Struct {
				continue // 不是结构体就跳过
			}

			inner, ok = AdjustSrcDepth(inner, innerDepth)
			if !ok {
				continue
			}

			// 获取底层结构体重复查找
			if sub, subField, ok := GetFieldByTag(inner, tag); ok {
				return sub, subField, true
			}

		}
	}

	return reflect.Value{}, reflect.StructField{}, false
}

func GetExportFieldByTag(s reflect.Value, tag string) (reflect.Value, reflect.StructField, bool) {

	var ok bool

	if s.Kind() != reflect.Struct {
		return reflect.Value{}, reflect.StructField{}, false
	}

	t := s.Type()

	// 外层优先匿名
	for i := 0; i < t.NumField(); i++ {
		sField := t.Field(i)
		fVal := s.Field(i)

		tagName, ignore := ParseTag(sField.Tag.Get(tag))
		if ignore {
			continue
		}

		if tag == tagName {
			return fVal, sField, true
		}

		// 内嵌
		if sField.Anonymous {
			inner := fVal
			innerType, innerDepth := IndirectType(fVal.Type())
			if innerType.Kind() != reflect.Struct {
				continue // 不是结构体就跳过
			}

			inner, ok = AdjustSrcDepth(inner, innerDepth)
			if !ok {
				continue
			}

			// 获取底层结构体重复查找, 与GetFieldByTag不同，条件增加 字段是导出字段
			if sub, subField, ok := GetExportFieldByTag(inner, tag); ok && subField.IsExported() {
				return sub, subField, true
			}

		}
	}

	return reflect.Value{}, reflect.StructField{}, false
}

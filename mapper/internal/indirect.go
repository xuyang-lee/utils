package internal

import (
	"reflect"
)

func IndirectType(reflectType reflect.Type) (_ reflect.Type, depth int) {
	depth = 0
	for reflectType.Kind() == reflect.Ptr {
		reflectType = reflectType.Elem()
		depth++
	}
	return reflectType, depth
}

func Indirect(reflectValue reflect.Value) (_ reflect.Value, ok bool) {
	// 防止开局崩,譬如开局一个nil
	if !reflectValue.IsValid() {
		return reflectValue, false
	}
	for reflectValue.Kind() == reflect.Ptr {
		// 非空，直接迭代
		if !reflectValue.IsNil() {
			reflectValue = reflectValue.Elem()
			continue
		}
		// 空了
		return reflectValue, false

	}
	return reflectValue, true
}

// AdjustSrcDepth 调整来源深度
func AdjustSrcDepth(reflectValue reflect.Value, depth int) (reflect.Value, bool) {
	if depth < 0 {
		panic("invalid depth")
	}
	for i := 0; i < depth; i++ {
		if !reflectValue.IsValid() {
			return reflectValue, false
		}
		if reflectValue.Kind() == reflect.Ptr {
			reflectValue = reflectValue.Elem()
		}
	}
	return reflectValue, true
}

// AdjustDstDepth 调整目标深度
func AdjustDstDepth(reflectValue reflect.Value, depth int) reflect.Value {
	if depth < 0 {
		panic("invalid depth")
	}

	for i := 0; i < depth; i++ {

		// 非指针，直接返回
		if reflectValue.Type().Kind() != reflect.Ptr {
			return reflectValue
		}

		// 非空，直接迭代
		if !reflectValue.IsNil() {
			reflectValue = reflectValue.Elem()
			continue
		}

		// 空指针

		// 不可设置 直接返回
		if !reflectValue.CanSet() {
			return reflectValue
		}
		// 设置新对象
		newValue := reflect.New(reflectValue.Type().Elem())
		reflectValue.Set(newValue)
		reflectValue = reflectValue.Elem()
	}
	return reflectValue
}

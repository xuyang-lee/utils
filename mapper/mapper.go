package mapper

import (
	"github.com/xuyang-lee/utils/mapper/internal"
	"reflect"
)

func Map(dst any, src any) error {
	return mapping(dst, src, &internal.Option{})
}

func MapWithOptions(dst any, src any, options ...internal.OptionFunc) error {
	opt := &internal.Option{MaxDepth: 100}
	for _, option := range options {
		option(opt)
	}
	return mapping(dst, src, opt)
}

func mapping(dst any, src any, opts *internal.Option) (err error) {
	// 判断来源可用，目的可寻址
	var ok bool
	from := reflect.ValueOf(src)
	to := reflect.ValueOf(dst)

	if to.Kind() != reflect.Ptr {
		return internal.ErrInvalidMappingDestination
	}
	to = to.Elem()

	// 错误结论：为何不用CanSet，对于结构体有未导出字段，但是可以只复制导出字段
	// CanSet对于结构体有未导出字段，也可判断，未导出字段不影响结构体整体的可设置规则
	if !from.IsValid() {
		return internal.ErrInvalidMappingSource
	}
	if !to.CanSet() {
		return internal.ErrInvalidMappingDestination
	}

	// 判断类型 底层类型可匹配
	fromType, fromDepth := internal.IndirectType(from.Type())
	toType, toDepth := internal.IndirectType(to.Type())

	// 底层类型匹配，直接处理
	if fromType.AssignableTo(toType) || fromType.ConvertibleTo(toType) {

		// to和from 都进入最底层
		from, ok = internal.AdjustSrcDepth(from, fromDepth)
		if !ok {
			return nil
		}
		to = internal.AdjustDstDepth(to, toDepth)
		if to.CanSet() {
			to.Set(from.Convert(to.Type()))
		}

		return
	}

	// 底层类型不匹配
	err = mapSwitch(to, from, opts)
	if err != nil {
		return err
	}

	return nil
}

func mapSwitch(to, from reflect.Value, opts *internal.Option) (err error) {
	if opts.IsTooDeep() {
		return internal.ErrNestedTooDeep
	}
	opts.DepthIncr()
	defer opts.DepthDecr()

	fromType, _ := internal.IndirectType(from.Type())
	toType, _ := internal.IndirectType(to.Type())
	// 判断
	switch {
	case toType.Kind() == reflect.Slice && fromType.Kind() == reflect.Slice:
		err = slice2slice(to, from, opts)
	case toType.Kind() == reflect.Struct && fromType.Kind() == reflect.Struct:
		err = struct2struct(to, from, opts)
	case toType.Kind() == reflect.Struct && fromType.Kind() == reflect.Map:
		err = map2struct(to, from, opts)
	case toType.Kind() == reflect.Map && fromType.Kind() == reflect.Map:
		err = map2map(to, from, opts)
	case fromType.Kind() == reflect.Interface:
		err = interface2other(to, from, opts)
		//case toType.Kind() == reflect.Map && fromType.Kind() == reflect.Struct: //
	}
	return
}

func map2map(to, from reflect.Value, opts *internal.Option) error {
	var ok bool
	fromType, fromDepth := internal.IndirectType(from.Type())
	toType, toDepth := internal.IndirectType(to.Type())

	if toType.Kind() != reflect.Map || fromType.Kind() != reflect.Map {
		return internal.ErrBadDstAndSrc
	}

	// to和from 都进入最底层
	from, ok = internal.AdjustSrcDepth(from, fromDepth)
	if !ok {
		return nil
	}
	to = internal.AdjustDstDepth(to, toDepth)

	if !fromType.Key().ConvertibleTo(toType.Key()) {
		return internal.ErrMapKeyNotMatch
	}

	if to.IsNil() {
		to.Set(reflect.MakeMapWithSize(toType, from.Len()))
	}

	for _, k := range from.MapKeys() {
		toKey := reflect.New(toType.Key()).Elem()
		kCopy := k.Convert(toType.Key())
		toKey.Set(kCopy)

		elemType := toType.Elem()
		if elemType.Kind() == reflect.Ptr {
			elemType, _ = internal.IndirectType(elemType)
		}

		toValue := reflect.New(elemType).Elem()

		if err := mapping(toValue.Addr().Interface(), from.MapIndex(k).Interface(), opts); err != nil {
			return err
		}

		for {
			if elemType == toType.Elem() {
				to.SetMapIndex(toKey, toValue)
				break
			}
			elemType = reflect.PointerTo(elemType)
			toValue = toValue.Addr()
		}
	}

	return nil
}

func map2struct(to, from reflect.Value, opts *internal.Option) error {
	var ok bool
	fromType, fromDepth := internal.IndirectType(from.Type())
	toType, toDepth := internal.IndirectType(to.Type())

	if toType.Kind() != reflect.Struct || fromType.Kind() != reflect.Map {
		return internal.ErrBadDstAndSrc
	}

	// key类型不是string，直接返回好了，不用处理
	if fromType.Key().Kind() != reflect.String {
		return nil
	}

	// to和from 都进入最底层
	from, ok = internal.AdjustSrcDepth(from, fromDepth)
	if !ok {
		return nil
	}
	to = internal.AdjustDstDepth(to, toDepth)

	for i := 0; i < to.NumField(); i++ {
		field := to.Field(i)
		if !field.CanSet() {
			continue
		}

		sField := toType.Field(i)
		tag := sField.Tag.Get(internal.Tag)
		fieldName, ignore := internal.ParseTag(tag)
		if ignore {
			continue
		}

		if fieldName == "" {
			fieldName = sField.Name
		}

		// 从 map 取值
		mapVal := from.MapIndex(reflect.ValueOf(fieldName))
		if !mapVal.IsValid() {
			continue
		}

		flag, err := internal.Set(field, mapVal, opts)
		if err != nil {
			return err
		}

		if flag.IsProcessed() {
			continue
		}

		if err = mapSwitch(field, mapVal, opts); err != nil {
			return err
		}

	}

	return nil
}

func struct2struct(to, from reflect.Value, opts *internal.Option) error {
	var ok bool
	fromType, fromDepth := internal.IndirectType(from.Type())
	toType, toDepth := internal.IndirectType(to.Type())

	if toType.Kind() != reflect.Struct || fromType.Kind() != reflect.Struct {
		return internal.ErrBadDstAndSrc
	}

	// to和from 都进入最底层
	from, ok = internal.AdjustSrcDepth(from, fromDepth)
	if !ok {
		return nil
	}
	to = internal.AdjustDstDepth(to, toDepth)

	for i := 0; i < to.NumField(); i++ {
		field := to.Field(i)
		if !field.CanSet() {
			continue
		}

		sField := toType.Field(i)
		tag := sField.Tag.Get(internal.Tag)
		fieldName, ignore := internal.ParseTag(tag)
		if ignore {
			continue
		}

		if fieldName == "" {
			fieldName = sField.Name
		}

		// 当前字段是嵌套字段，而且tag没有内容，就要展开处理
		if sField.Anonymous && tag == "" {
			if err := struct2struct(field, from, opts); err != nil {
				return err
			}
		}

		// 不是嵌套字段了

		// 从 from 取值，tag 优先
		fromValue, _, ok := internal.GetExportFieldByTag(from, fieldName)
		if !ok { // tag 不存在，查找name
			fromValue = from.FieldByName(fieldName)
			//_, ok = from.Type().FieldByName(fieldName)
			if !fromValue.IsValid() {
				continue // from没对应字段，直接返回
			}
		}

		flag, err := internal.Set(field, fromValue, opts)
		if err != nil {
			return err
		}

		if flag.IsProcessed() {
			continue
		}

		if err = mapSwitch(field, fromValue, opts); err != nil {
			return err
		}
	}

	return nil
}

func slice2slice(to, from reflect.Value, opts *internal.Option) error {

	var ok bool
	fromType, fromDepth := internal.IndirectType(from.Type())
	toType, toDepth := internal.IndirectType(to.Type())

	if toType.Kind() != reflect.Slice || fromType.Kind() != reflect.Slice {
		return internal.ErrBadDstAndSrc
	}

	// to和from 都进入最底层
	from, ok = internal.AdjustSrcDepth(from, fromDepth)
	if !ok {
		return nil
	}
	to = internal.AdjustDstDepth(to, toDepth)

	// from 为空
	if from.IsNil() {
		to.Set(reflect.Zero(toType))
		return nil
	}

	newSlice := reflect.MakeSlice(toType, 0, from.Len())
	for i := 0; i < from.Len(); i++ {
		srcElem := from.Index(i)
		dstElem := reflect.New(toType.Elem()).Elem()

		_, srcElemDepth := internal.IndirectType(srcElem.Type())
		srcElem, _ = internal.AdjustSrcDepth(srcElem, srcElemDepth)

		_, dstValDepth := internal.IndirectType(dstElem.Type())
		dstVal := internal.AdjustDstDepth(dstElem, dstValDepth) // 最终类型值，单独保存用于赋值，原始dstElem用于加入切片

		flag, err := internal.Set(dstVal, srcElem, opts)
		if err != nil {
			return err
		}
		if !flag.IsProcessed() {
			if err := mapSwitch(dstVal, srcElem, opts); err != nil {
				return err
			}
		}

		newSlice = reflect.Append(newSlice, dstElem)
	}

	to.Set(newSlice)
	return nil
}

func interface2other(to, from reflect.Value, opts *internal.Option) error {
	var ok bool

	fromType, fromDepth := internal.IndirectType(from.Type())

	if fromType.Kind() != reflect.Interface {
		return internal.ErrBadDstAndSrc
	}

	// from 进入最底层
	from, ok = internal.AdjustSrcDepth(from, fromDepth)
	if !ok {
		return nil
	}
	// from 获取interface的原始值
	from = from.Elem()

	flag, err := internal.Set(to, from, opts)
	if err != nil {
		return err
	}

	if flag.IsProcessed() {
		return nil
	}

	// 其他类型
	err = mapSwitch(to, from, opts)
	return err

}

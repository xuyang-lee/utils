package internal

import (
	"reflect"
)

type SetFlag int

const (
	SetFlagProcessed SetFlag = 1 << iota
	SetFlagSet
	SetFlagDstNotSettable
	SetFlagSrcInvalid
)

func (f *SetFlag) set(flag SetFlag) {
	if f != nil {
		*f |= flag
	}
}

func (f *SetFlag) clear(flag SetFlag) {
	if f != nil {
		*f &= ^flag
	}
}

// Toggle 反转
func (f *SetFlag) Toggle(flag SetFlag) {
	if f != nil {
		*f ^= flag
	}
}

func (f *SetFlag) Has(flag SetFlag) bool {
	if f == nil {
		return false
	}
	return *f&flag != 0
}

func (f *SetFlag) IsSet() bool {
	return f.Has(SetFlagSet)
}

func (f *SetFlag) IsProcessed() bool {
	return f.Has(SetFlagProcessed)
}

func (f *SetFlag) IsValid() bool {
	return f != nil
}

func (f *SetFlag) Processed() {
	f.set(SetFlagProcessed)
}

func (f *SetFlag) BeSet() {
	f.set(SetFlagProcessed | SetFlagSet)
}

func (f *SetFlag) DstNotSettable() {
	f.set(SetFlagProcessed | SetFlagDstNotSettable)
}

func (f *SetFlag) SrcInvalid() {
	f.set(SetFlagProcessed | SetFlagSrcInvalid)
}

func Set(to, from reflect.Value, opts *Option) (flag SetFlag, err error) {
	var ok bool
	// 解引用字段类型
	fromType, fromDepth := IndirectType(from.Type())
	toType, toDepth := IndirectType(to.Type())

	// 类型可转换，直接处理 转换后复制
	if fromType.AssignableTo(toType) || fromType.ConvertibleTo(toType) {
		// to和from 都进入最底层
		from, ok = AdjustSrcDepth(from, fromDepth)
		if !ok {
			flag.SrcInvalid()
			return
		}
		to = AdjustDstDepth(to, toDepth)

		if to.CanSet() {
			to.Set(from.Convert(to.Type()))
			flag.BeSet()
		} else {
			flag.DstNotSettable()
		}
	}

	return
}

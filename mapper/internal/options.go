package internal

import "sync/atomic"

type Option struct {
	DeepCopy bool
	depth    atomic.Int64
	MaxDepth int64
}

type OptionFunc func(*Option)

func WithDeepCopy(deepCopy bool) OptionFunc {
	return func(o *Option) {
		o.DeepCopy = deepCopy
	}
}

func WithMaxDepth(depth int) OptionFunc {
	return func(o *Option) {
		o.MaxDepth = int64(depth)
	}
}

func (o *Option) DepthIncr() {
	if o == nil {
		return
	}
	o.depth.Add(1)
}

func (o *Option) DepthDecr() {
	if o == nil {
		return
	}
	o.depth.Add(-1)
}

func (o *Option) IsTooDeep() bool {
	if o == nil {
		return false
	}
	return o.depth.Load() > o.MaxDepth
}

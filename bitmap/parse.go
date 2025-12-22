package bitmap

func Parse(bytes []byte, order bitOrder, opts ...ParseOption) (*Bitmap, error) {
	bm := &Bitmap{
		size: 0,
		step: defaultStep,
	}
	if len(opts) > 0 {
		for _, opt := range opts {
			opt(bm)
		}
	}

	bm.data = convBasemap(bytes, order)
	if bm.size == 0 {
		bm.size = bm.data.size() * 8
	}

	if bm.step <= 0 {
		return nil, ErrInvalidStep
	}
	if bm.size < 0 || bm.size > bm.data.size()*8 {
		return nil, ErrInvalidSize
	}

	return bm, nil
}

type ParseOption func(b *Bitmap)

func WithOffset(offset int) ParseOption {
	return func(b *Bitmap) {
		b.offset = offset
	}
}

func WithSize(size int) ParseOption {
	return func(b *Bitmap) {
		b.size = size
	}
}

func WithStep(step int) ParseOption {
	return func(b *Bitmap) {
		b.step = step
	}
}

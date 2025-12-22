package bitmap

import (
	"errors"
)

const (
	defaultStep = 1
	badBit      = -1 // bit < 0 means invalid position
	baseBit     = 1
)

var (
	ErrOutOfBounds = errors.New("position out of bounds")
	ErrInvalidStep = errors.New("invalid step, must be greater than 0")
	ErrInvalidSize = errors.New("invalid size, size is too large or negative")
)

// Bitmap is a slice of bytes.
type Bitmap struct {
	//data   []byte
	data   basemap
	size   int // Total number of bits in the bitmap
	step   int // step of two nearly pos
	offset int // first bit index
}

// NewBitmap creates a new Bitmap with the given size in bits.
func NewBitmap(size int, order ...bitOrder) *Bitmap {
	var bo = BitOrderMSB
	if len(order) > 0 {
		bo = order[0]
	}

	return &Bitmap{
		data: newBasemap(size, bo), // Round up to the nearest byte
		size: size,
		step: defaultStep,
	}
}

// NewBitmapWithOffset creates a new Bitmap with the given size in bits and the first pos is offset.
func NewBitmapWithOffset(size int, offset int, order ...bitOrder) *Bitmap {
	var bo = BitOrderMSB
	if len(order) > 0 {
		bo = order[0]
	}
	return &Bitmap{
		data:   newBasemap(size, bo), // Round up to the nearest byte
		size:   size,
		step:   defaultStep,
		offset: offset,
	}
}

// NewBitmapWithRang creates a new Bitmap with step, whose range from begin to end(without end).
func NewBitmapWithRang(begin int, end int, step int, order ...bitOrder) *Bitmap {
	if step <= 0 {
		panic("step must be greater than 0")
	}
	var bo = BitOrderMSB
	if len(order) > 0 {
		bo = order[0]
	}

	return &Bitmap{
		data:   newBasemap((end-begin+step-1)/step, bo), // Round up to the nearest byte
		size:   (end - begin + step - 1) / step,
		offset: begin,
		step:   step,
	}
}

// 将位图根据offset、step进行压缩转化
func (b *Bitmap) encodePos(pos int) int {
	bit := pos - b.offset
	if bit%b.step != 0 {
		panic(ErrOutOfBounds)
	}
	bit = bit / b.step
	if bit < 0 || bit >= b.size {
		panic(ErrOutOfBounds)
	}
	return bit
}

// 将位图根据offset、step进行压缩转化
func (b *Bitmap) tryEncodePos(pos int) (int, error) {
	bit := pos - b.offset
	if bit%b.step != 0 {
		return badBit, ErrOutOfBounds
	}
	bit = bit / b.step
	if bit < 0 || bit >= b.size {
		return badBit, ErrOutOfBounds
	}
	return bit, nil
}

// 将位图根据offset、step进行解压转化
func (b *Bitmap) decodePos(bit int) int {
	return bit*b.step + b.offset
}

// Set sets the bit at pos to 1.
func (b *Bitmap) Set(pos int) {
	bit := b.encodePos(pos)
	b.data.set(bit)
}

// Clear clears the bit at pos to 0.
func (b *Bitmap) Clear(pos int) {
	bit := b.encodePos(pos)
	b.data.clear(bit)
}

// IsSet checks if the bit at pos is set to 1.
func (b *Bitmap) IsSet(pos int) bool {
	bit := b.encodePos(pos)
	return b.data.check(bit)
}

// GetPos returns the positions of all the set bits in the bitmap.
func (b *Bitmap) GetPos() []int {
	poses := make([]int, 0)
	for i := 0; i < b.data.size(); i++ {
		for j := 0; j < 8; j++ {
			cursor := i*8 + j
			if (cursor < b.size) && b.data.check(cursor) {
				poses = append(poses, b.decodePos(cursor))
			}
		}
	}
	return poses
}

func (b *Bitmap) GetNoPos() []int {
	poses := make([]int, 0)
	for i := 0; i < b.data.size(); i++ {
		for j := 0; j < 8; j++ {
			cursor := i*8 + j
			if (cursor < b.size) && !b.data.check(cursor) {
				poses = append(poses, b.decodePos(cursor))
			}
		}
	}
	return poses
}

// TrySet try set the bit at pos to 1. If pos is out of bounds, do nothing.
func (b *Bitmap) TrySet(pos int) error {
	bit, err := b.tryEncodePos(pos)
	if err != nil {
		return err
	}
	b.data.set(bit)
	return err
}

// TryClear clears the bit at pos to 0. If pos is out of bounds, do nothing.
func (b *Bitmap) TryClear(pos int) error {
	bit, err := b.tryEncodePos(pos)
	if err != nil {
		return err
	}
	b.data.clear(bit)
	return nil
}

// IsSetWithErr checks if the bit at pos is set to 1. If pos is out of bounds, return false,err.
func (b *Bitmap) IsSetWithErr(pos int) (bool, error) {
	bit, err := b.tryEncodePos(pos)
	if err != nil {
		return false, err
	}
	return b.data.check(bit), nil
}

func (b *Bitmap) Size() int {
	return b.size
}
func (b *Bitmap) Step() int {
	return b.step
}
func (b *Bitmap) Offset() int {
	return b.offset
}

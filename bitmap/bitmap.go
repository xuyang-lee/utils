package bitmap

import "fmt"

const (
	defaultStep = 1
)

// Bitmap is a slice of bytes.
type Bitmap struct {
	data   []byte
	size   int // Total number of bits in the bitmap
	step   int // step of two nearly pos
	offset int // first bit index
}

// NewBitmap creates a new Bitmap with the given size in bits.
func NewBitmap(size int) *Bitmap {
	return &Bitmap{
		data: make([]byte, (size+7)/8), // Round up to the nearest byte
		size: size,
		step: defaultStep,
	}
}

// NewBitmapWithOffset creates a new Bitmap with the given size in bits and the first pos is offset.
func NewBitmapWithOffset(size int, offset int) *Bitmap {
	return &Bitmap{
		data:   make([]byte, (size+7)/8), // Round up to the nearest byte
		size:   size,
		step:   defaultStep,
		offset: offset,
	}
}

// NewBitmapWithRang creates a new Bitmap from begin to end(without end) with step.
func NewBitmapWithRang(begin int, end int, step int) *Bitmap {
	if step <= 0 {
		panic("step must be greater than 0")
	}
	return &Bitmap{
		data:   make([]byte, (end-begin+7)/8), // Round up to the nearest byte
		size:   (end - begin + step - 1) / step,
		offset: begin,
		step:   step,
	}
}

// 将位图根据offset、step进行压缩转化
func (b *Bitmap) encodePos(pos int) int {
	bit := pos - b.offset
	if bit%b.step != 0 {
		panic("position out of bounds")
	}
	bit = bit / b.step
	if bit < 0 || bit >= b.size {
		panic("position out of bounds")
	}
	return bit
}

// 将位图根据offset、step进行压缩转化
func (b *Bitmap) tryEncodePos(pos int) (int, error) {
	bit := pos - b.offset
	if bit%b.step != 0 {
		return -1, fmt.Errorf("position out of bounds")
	}
	bit = bit / b.step
	if bit < 0 || bit >= b.size {
		return -1, fmt.Errorf("position out of bounds")
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
	byteIndex := bit / 8      // Find the byte index
	bitIndex := uint(bit % 8) // Find the bit index within the byte
	b.data[byteIndex] |= 1 << bitIndex
}

// Clear clears the bit at pos to 0.
func (b *Bitmap) Clear(pos int) {
	bit := b.encodePos(pos)
	byteIndex := bit / 8
	bitIndex := uint(bit % 8)
	// A &^ B 位清空操作符，A的第B位清空
	// A &^ B <==> A & (^B)
	b.data[byteIndex] &^= 1 << bitIndex
}

// IsSet checks if the bit at pos is set to 1.
func (b *Bitmap) IsSet(pos int) bool {
	bit := b.encodePos(pos)
	byteIndex := bit / 8
	bitIndex := uint(bit % 8)
	return (b.data[byteIndex] & (1 << bitIndex)) != 0
}

// GetPos returns the positions of all the set bits in the bitmap.
func (b *Bitmap) GetPos() []int {
	poses := make([]int, 0)
	for i := range b.data {
		for j := 0; j < 8; j++ {
			if (i*8+j < b.size) && b.data[i]&(1<<j) != 0 {
				poses = append(poses, b.decodePos(i*8+j))
			}
		}
	}
	return poses
}

func (b *Bitmap) GetNoPos() []int {
	poses := make([]int, 0)
	for i := range b.data {
		for j := 0; j < 8; j++ {
			if (i*8+j < b.size) && b.data[i]&(1<<j) == 0 {
				poses = append(poses, b.decodePos(i*8+j))
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
	byteIndex := bit / 8      // Find the byte index
	bitIndex := uint(bit % 8) // Find the bit index within the byte
	b.data[byteIndex] |= 1 << bitIndex
	return err
}

// TryClear clears the bit at pos to 0. If pos is out of bounds, do nothing.
func (b *Bitmap) TryClear(pos int) error {
	bit, err := b.tryEncodePos(pos)
	if err != nil {
		return err
	}
	byteIndex := bit / 8
	bitIndex := uint(bit % 8)
	// A &^ B 位清空操作符，A的第B位清空
	// A &^ B <==> A & (^B)
	b.data[byteIndex] &^= 1 << bitIndex
	return nil
}

// IsSetWithErr checks if the bit at pos is set to 1. If pos is out of bounds, return false,err.
func (b *Bitmap) IsSetWithErr(pos int) (bool, error) {
	bit, err := b.tryEncodePos(pos)
	if err != nil {
		return false, err
	}
	byteIndex := bit / 8
	bitIndex := uint(bit % 8)
	return (b.data[byteIndex] & (1 << bitIndex)) != 0, nil
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

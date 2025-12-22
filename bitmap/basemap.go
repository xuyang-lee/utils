package bitmap

const (
	BitOrderMSB bitOrder = iota
	BitOrderLSB
)

type basemap interface {
	set(bit int)
	clear(bit int)
	check(bit int) bool
	size() int
}

type bitOrder int

func newBasemap(size int, bo bitOrder) basemap {
	bytes := make([]byte, (size+7)/8)
	return convBasemap(bytes, bo)
}

func convBasemap(bytes []byte, bo bitOrder) basemap {
	switch bo {
	case BitOrderMSB:
		v := msbMap(bytes)
		return &v
	case BitOrderLSB:
		v := lsbMap(bytes)
		return &v
	default:
		v := msbMap(bytes)
		return &v
	}
}

func byteIndex(bit int) int {
	return bit / 8
}

// lsbMap little-endian bit order
type lsbMap []byte

func (r *lsbMap) set(bit int) {
	(*r)[byteIndex(bit)] |= baseBit << r.bitIndex(bit)
}

func (r *lsbMap) clear(bit int) {
	// A &^ B 位清空操作符，A的第B位清空
	// A &^ B <==> A & (^B)
	(*r)[byteIndex(bit)] &^= baseBit << r.bitIndex(bit)
}

func (r *lsbMap) check(bit int) bool {
	return (*r)[byteIndex(bit)]&(baseBit<<r.bitIndex(bit)) != 0
}

func (r *lsbMap) size() int {
	return len(*r)
}

func (r *lsbMap) bitIndex(bit int) uint {
	return uint(bit % 8)
}

// msbMap big-endian bit order
type msbMap []byte

func (r *msbMap) set(bit int) {
	(*r)[byteIndex(bit)] |= baseBit << r.bitIndex(bit)
}

func (r *msbMap) clear(bit int) {
	// A &^ B 位清空操作符，A的第B位清空
	// A &^ B <==> A & (^B)
	(*r)[byteIndex(bit)] &^= baseBit << r.bitIndex(bit)
}

func (r *msbMap) check(bit int) bool {
	return (*r)[byteIndex(bit)]&(baseBit<<r.bitIndex(bit)) != 0
}

func (r *msbMap) size() int {
	return len(*r)
}

func (r *msbMap) bitIndex(bit int) uint {
	return 7 - uint(bit%8)
}

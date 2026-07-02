package bitset

import "fmt"

type BitSet struct {
	bits []uint64
	size uint
}

// Создание массива битов
func New(size uint) (*BitSet, error) {
	if size < 1 {
		return nil, fmt.Errorf("bitset: размер массива не меньше 1")
	}
	bitarray := make([]uint64, (size+63)/64)
	return &BitSet{bitarray, size}, nil
}

// Установка бита
func (bitset *BitSet) Set(bit uint) {
	if bit >= bitset.size {
		panic("bitset: устанавливаемый bit вне диапазона")
	}
	i := bit / 64
	j := bit % 64
	bitset.bits[i] |= 1 << j
}

// Проверка бита
func (bitset *BitSet) Get(bit uint) bool {
	if bit >= bitset.size {
		panic("bitset: полученный bit вне диапазона")
	}
	i := bit / 64
	j := bit % 64
	return (bitset.bits[i]&(1<<j) != 0)
}

// Размер массива
func (bitset *BitSet) Size() uint {
	return bitset.size
}

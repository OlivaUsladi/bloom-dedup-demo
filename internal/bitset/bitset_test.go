package bitset

import "testing"

func TestNewWithSize(t *testing.T) {
	size := uint(100)
	bs, err := New(size)
	if err != nil {
		t.Fatalf("New вернул неожиданную ошибку: %v", err)
	}
	if bs.Size() != size {
		t.Errorf("ожидали Size()=%d, получили %d", size, bs.Size())
	}
}

func TestNewWithZeroSize(t *testing.T) {
	size := uint(0)
	_, err := New(size)
	if err == nil {
		t.Errorf("ожидали ошибку при size=%d, получили nil", size)
	}
}

func TestGet(t *testing.T) {
	bs, err := New(100)
	if err != nil {
		t.Fatalf("New вернул ошибку: %v", err)
	}
	bs.bits[0] = 1 << 5
	result := bs.Get(5)
	other := bs.Get(6)
	if !result {
		t.Errorf("ожидали true для бита 5, получили %v", result)
	}
	if other {
		t.Errorf("ожидали false для бита 6, получили %v", other)
	}
}

func TestSetAndGet(t *testing.T) {
	bs, err := New(100)
	if err != nil {
		t.Fatalf("New вернул ошибку: %v", err)
	}
	bs.Set(0)
	bs.Set(99)

	firstBit := bs.Get(0)
	lastBit := bs.Get(99)
	untouchedBit := bs.Get(50)

	if !firstBit {
		t.Errorf("ожидали true для бита 0, получили %v", firstBit)
	}
	if !lastBit {
		t.Errorf("ожидали true для бита 99, получили %v", lastBit)
	}
	if untouchedBit {
		t.Errorf("ожидали false для бита 50, получили %v", untouchedBit)
	}
}

func TestSetPanics(t *testing.T) {
	bs, _ := New(10)
	bit := uint(100)
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("ожидали панику при установке бита %d", bit)
		}
	}()
	bs.Set(bit)
}

func TestGetPanics(t *testing.T) {
	bs, _ := New(1)
	bit := uint(2)
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("ожидали панику при установке бита %d", bit)
		}
	}()
	bs.Get(bit)
}

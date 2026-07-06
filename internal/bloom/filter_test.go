package bloom

import "testing"

func TestFilterMap(t *testing.T) {
	a, b, c, err := MapFilter("../../testdata/control/event1.jsonl")
	if err != nil {
		t.Fatal(err)
	}
	if a != 4 || b != 4 || c != 0 {
		t.Errorf("Должно быть 4 значений, 4 уникальных значений и 0 дублей, вывелось %v, %v, %v", a, b, c)
	}
}

func BenchmarkFilterMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MapFilter("../../testdata/control/event1.jsonl")
	}
}

func TestBloomFilter(t *testing.T) {
	total, unique, duplicates, err := BloomFilter(4, 0.01, "../../testdata/control/event1.jsonl")
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}
	if total != 4 {
		t.Errorf("total должен быть 4, получено %v", total)
	}
	if unique != 4 || duplicates != 0 {
		t.Errorf("ожидалось 4 уникальных, 0 дублей, получено unique=%v, duplicates=%v", unique, duplicates)
	}
}

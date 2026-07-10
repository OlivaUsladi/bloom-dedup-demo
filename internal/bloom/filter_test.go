package bloom

import "testing"

func TestFilterMap(t *testing.T) {
	a, b, c, err := MapFilter("../../testdata/control/event.jsonl", true)
	if err != nil {
		t.Fatal(err)
	}
	if a != 4 || b != 4 || c != 0 {
		t.Errorf("Должно быть 4 значений, 4 уникальных значений и 0 дублей, вывелось %v, %v, %v", a, b, c)
	}
}

func TestFilterMap1(t *testing.T) {
	a, b, c, err := MapFilter("../../testdata/control/event1.jsonl", true)
	if err != nil {
		t.Fatal(err)
	}
	if a != 100 || b != 87 || c != 13 {
		t.Errorf("Должно быть 100 значений, 87 уникальных значений и 13 дублей, вывелось %v, %v, %v", a, b, c)
	}
}

func TestFilterMap2(t *testing.T) {
	a, b, c, err := MapFilter("../../testdata/control/event2.jsonl", true)
	if err != nil {
		t.Fatal(err)
	}
	if a != 50 || b != 50 || c != 0 {
		t.Errorf("ожидали total=50, unique=50, duplicates=0, получили %v, %v, %v", a, b, c)
	}
}

func TestFilterMap3(t *testing.T) {
	a, b, c, err := MapFilter("../../testdata/control/event3.jsonl", true)
	if err != nil {
		t.Fatal(err)
	}
	if a != 200 || b != 101 || c != 99 {
		t.Errorf("ожидали total=200, unique=101, duplicates=99, получили %v, %v, %v", a, b, c)
	}
}

func BenchmarkFilterMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		MapFilter("../../testdata/control/event.jsonl", true)
	}
}

func TestBloomFilter(t *testing.T) {
	total, unique, duplicates, err := BloomFilter(0.01, "../../testdata/control/event.jsonl", true)
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

func TestBloomFilter1(t *testing.T) {
	total, unique, duplicates, err := BloomFilter(0.1, "../../testdata/control/event1.jsonl", true)
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}
	if total != 100 {
		t.Errorf("total должен быть 100, получено %v", total)
	}
	if unique != 84 || duplicates != 16 {
		t.Errorf("ожидалось 87 уникальных, 13 дублей, получено unique=%v, duplicates=%v", unique, duplicates)
	}
}

func TestBloomFilter2(t *testing.T) {
	total, unique, duplicates, err := BloomFilter(0.01, "../../testdata/control/event2.jsonl", true)
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}
	if total != 50 {
		t.Errorf("total должен быть 50, получено %v", total)
	}
	if unique != 49 || duplicates != 1 {
		t.Errorf("ожидалось unique=49, duplicates=1, получено unique=%v, duplicates=%v", unique, duplicates)
	}
}

func TestBloomFilter3(t *testing.T) {
	total, unique, duplicates, err := BloomFilter(0.01, "../../testdata/control/event3.jsonl", true)
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}
	if total != 200 {
		t.Errorf("total должен быть 200, получено %v", total)
	}
	if unique != 101 || duplicates != 99 {
		t.Errorf("ожидалось unique=101, duplicates=99, получено unique=%v, duplicates=%v", unique, duplicates)
	}
}

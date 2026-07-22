package bloom

import (
	"bloom-dedup-demo/internal/model"
	"fmt"
	"testing"
)

func TestFilterMap(t *testing.T) {
	events, _, _, err := model.ReadEvents("../../testdata/tests/event.jsonl", true)
	_, b, c, _, _, err := MapFilter(events)
	if err != nil {
		t.Fatal(err)
	}
	if b != 4 || c != 0 {
		t.Errorf("Должно быть 4 уникальных значений и 0 дублей, вывелось %v, %v", b, c)
	}
}

func TestFilterMap1(t *testing.T) {
	events, _, _, err := model.ReadEvents("../../testdata/tests/event1.jsonl", true)
	_, b, c, _, _, err := MapFilter(events)
	if err != nil {
		t.Fatal(err)
	}
	if b != 87 || c != 13 {
		t.Errorf("Должно быть 87 уникальных значений и 13 дублей, вывелось %v, %v", b, c)
	}
}

func TestFilterMap2(t *testing.T) {
	events, _, _, err := model.ReadEvents("../../testdata/tests/event2.jsonl", true)
	_, b, c, _, _, err := MapFilter(events)
	if err != nil {
		t.Fatal(err)
	}
	if b != 50 || c != 0 {
		t.Errorf("ожидали unique=50, duplicates=0, получили %v, %v", b, c)
	}
}

func TestFilterMap3(t *testing.T) {
	events, _, _, err := model.ReadEvents("../../testdata/tests/event3.jsonl", true)
	_, b, c, _, _, err := MapFilter(events)
	if err != nil {
		t.Fatal(err)
	}
	if b != 101 || c != 99 {
		t.Errorf("ожидали unique=101, duplicates=99, получили %v, %v", b, c)
	}
}

func TestMapFilterLargeFile(t *testing.T) {
	events, _, _, err := model.ReadEvents("../../testdata/tests/event_large.jsonl", true)
	if err != nil {
		t.Fatal(err)
	}
	_, _, _, _, memory, err := MapFilter(events)
	if err != nil {
		t.Fatal(err)
	}
	if memory <= 0 {
		t.Errorf("на большом файле ожидали положительное значение памяти")
	}
	fmt.Println(memory)
}

func BenchmarkFilterMap(b *testing.B) {
	events, _, _, _ := model.ReadEvents("../../testdata/tests/event.jsonl", true)
	for i := 0; i < b.N; i++ {
		MapFilter(events)
	}
}

func TestBloomFilter(t *testing.T) {
	events, _, _, err := model.ReadEvents("../../testdata/tests/event.jsonl", true)
	unique, duplicates, _, _, err := BloomFilter(events, 5, "fnv64_double_hashing", 0.01)
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}
	if unique != 4 || duplicates != 0 {
		t.Errorf("ожидалось 4 уникальных, 0 дублей, получено unique=%v, duplicates=%v", unique, duplicates)
	}
}

func TestBloomFilter1(t *testing.T) {
	events, _, _, err := model.ReadEvents("../../testdata/tests/event1.jsonl", true)
	unique, duplicates, _, _, err := BloomFilter(events, 100, "fnv64_double_hashing", 0.1)
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}
	if unique != 84 || duplicates != 16 {
		t.Errorf("ожидалось 87 уникальных, 13 дублей, получено unique=%v, duplicates=%v", unique, duplicates)
	}
}

func TestBloomFilter2(t *testing.T) {
	events, _, _, err := model.ReadEvents("../../testdata/tests/event2.jsonl", true)
	unique, duplicates, _, _, err := BloomFilter(events, 50, "fnv64_double_hashing", 0.01)
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}
	if unique != 49 || duplicates != 1 {
		t.Errorf("ожидалось unique=49, duplicates=1, получено unique=%v, duplicates=%v", unique, duplicates)
	}
}

func TestBloomFilter3(t *testing.T) {
	events, _, _, err := model.ReadEvents("../../testdata/tests/event3.jsonl", true)
	unique, duplicates, _, _, err := BloomFilter(events, 200, "fnv64_double_hashing", 0.01)
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}
	if unique != 101 || duplicates != 99 {
		t.Errorf("ожидалось unique=101, duplicates=99, получено unique=%v, duplicates=%v", unique, duplicates)
	}
}

func TestBloomFilterSHA256(t *testing.T) {
	events, _, _, err := model.ReadEvents("../../testdata/control/demo_events.jsonl", true)
	if err != nil {
		t.Fatal(err)
	}

	cfg, err := model.ReadConfig("../../testdata/tests/conf_sha256.json")
	if err != nil {
		t.Fatal(err)
	}
	unique, duplicates, _, memory, err := BloomFilter(events, cfg.ExpectedItems, cfg.HashFamily, cfg.FalsePositiveRate)
	if err != nil {
		t.Fatalf("ошибка: %v", err)
	}

	if unique+duplicates != len(events) {
		t.Errorf("ожидалось unique + duplicates = %d, получено %d", len(events), unique+duplicates)
	}
	if unique <= 0 {
		t.Errorf("ожидалось положительное число unique, получено %d", unique)
	}
	if duplicates < 0 {
		t.Errorf("ожидалось неотрицательное число duplicates, получено %d", duplicates)
	}
	if memory <= 0 {
		t.Errorf("ожидалось положительное значение памяти, получено %d", memory)
	}
}

func TestBloomFilterBadExpectedItems(t *testing.T) {
	events, _, _, err := model.ReadEvents("../../testdata/tests/event3.jsonl", true)
	if err != nil {
		t.Fatal(err)
	}
	_, _, _, _, err = BloomFilter(events, -1, "fnv64_double_hashing", 0.01)
	if err == nil {
		t.Fatal("ожидали ошибку при expectedItems < 0")
	}
}

func TestCountingBloomFilterBadRate(t *testing.T) {
	events, _, _, err := model.ReadEvents("../../testdata/tests/event3.jsonl", true)
	if err != nil {
		t.Fatal(err)
	}
	_, _, _, _, err = CountingBloomFilter(events, 10, "fnv64_double_hashing", 0)
	if err == nil {
		t.Fatal("ожидали ошибку при p=0")
	}
}

func TestHash(t *testing.T) {
	h11, h12 := Hash("value")
	h21, h22 := Hash("value")

	if h11 != h21 || h12 != h22 {
		t.Errorf("Hash должны быть одинаковыми, получили (%d, %d) и (%d, %d)", h11, h12, h21, h22)
	}
}

func TestGetIndexesRange(t *testing.T) {
	indexes := getIndexes("tratata", 7, 100)

	if len(indexes) != 7 {
		t.Fatalf("ожидали 7 индексов, получили %d", len(indexes))
	}

	for _, idx := range indexes {
		if idx >= 100 {
			t.Errorf("индекс вне диапазона: %d", idx)
		}
	}
}

func TestGetIndexesSHA256Range(t *testing.T) {
	indexes := getIndexesSHA256("tratata", 7, 100)

	if len(indexes) != 7 {
		t.Fatalf("ожидали 7 индексов, получили %d", len(indexes))
	}

	for _, idx := range indexes {
		if idx >= 100 {
			t.Errorf("индекс вне диапазона: %d", idx)
		}
	}
}

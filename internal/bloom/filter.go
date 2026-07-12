package bloom

import (
	"bloom-dedup-demo/internal/bitset"
	"bloom-dedup-demo/internal/model"
	"hash/fnv"
	"runtime"
	"time"
)

// Точная дедупликация событий через map
// Возвращает unique, duplicates, durationMs, memoryBytes, error
func MapFilter(events []model.Event) (int, int, int64, int, error) {
	var m1, m2 runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m1)
	start := time.Now()
	total := len(events)
	eventsMap := make(map[string]model.Event)
	//events, _, total, _, err := model.ReadEvents(path, fs)
	//if err != nil {
	//	return 0, 0, 0, 0, err
	//}
	for i, event := range events {
		eventsMap[event.EventHash] = events[i]
	}
	unique := len(eventsMap)
	duplicates := total - unique
	duration := time.Since(start).Milliseconds()
	runtime.ReadMemStats(&m2)
	memory := int(m2.TotalAlloc - m1.TotalAlloc)
	if memory == 0 {
		memory = estimateMapMemory(eventsMap)
	}
	return unique, duplicates, duration, memory, nil
}

// Рассчёт оценки map в байтах (примерное значение)
func estimateMapMemory(m map[string]model.Event) int {
	const mapBucketOverhead = 48
	total := 0
	for k, v := range m {
		total += len(k)
		total += len(v.EventID) + len(v.Source) + len(v.Timestamp)
		total += 8
		total += mapBucketOverhead
	}
	return total
}

// Вычисляет два независимых хеша строки для дальнейшего использования в схеме двойного хеширования
// Возвращает h1, h2 - базовые хеш-значения
func Hash(str string) (uint64, uint64) {
	h1 := fnv.New64()
	message := []byte(str)
	h1.Write(message)
	h1Value := h1.Sum64()
	h2 := fnv.New64a()
	message2 := []byte(str)
	h2.Write(message2)
	h2Value := h2.Sum64()
	return h1Value, h2Value
}

// Вычисляет k позиций в битовом массиве размера m для ключа key методом двойного хеширования
// Возвращает массив из k индексов битового массива
func getIndexes(key string, k int, m uint64) []uint64 {
	h1, h2 := Hash(key)
	indexes := make([]uint64, k)
	for i := 0; i < k; i++ {
		indexes[i] = (h1 + uint64(i)*h2) % m
	}
	return indexes
}

// Фильтр Блума
// Возвращает unique, duplicates, durationMs, memoryBytes, error
func BloomFilter(events []model.Event, p float64) (int, int, int64, int, error) {
	start := time.Now()
	//events, _, total, _, err2 := model.ReadEvents(path, fs)
	//if err2 != nil {
	//	return 0, 0, err2
	//}
	total := len(events)
	m, k, err1 := Params(total, p)
	if err1 != nil {
		return 0, 0, 0, 0, err1
	}
	bs, err := bitset.New(uint(m))
	if err != nil {
		return 0, 0, 0, 0, err
	}
	duplicates := 0

	for _, event := range events {
		indexes := getIndexes(event.EventHash, k, uint64(m))
		alreadySet := true
		for _, idx := range indexes {
			if !bs.Get(uint(idx)) {
				alreadySet = false
			}
			bs.Set(uint(idx))
		}
		if alreadySet {
			duplicates++
		}
	}
	unique := total - duplicates
	duration := time.Since(start).Milliseconds()
	memory := (m + 7) / 8

	return unique, duplicates, duration, memory, nil
}

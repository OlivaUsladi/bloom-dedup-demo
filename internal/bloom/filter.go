package bloom

import (
	"bloom-dedup-demo/internal/bitset"
	"bloom-dedup-demo/internal/model"
	"hash/fnv"
)

// Точная дедупликация событий через map
// Возвращает total (всего событий), unique (уникальных), duplicates (дубликатов)
func MapFilter(path string) (int, int, int, error) {
	eventsMap := make(map[string]model.Event)
	events, _, total, _, err := model.ReadEvents(path)
	if err != nil {
		return 0, 0, 0, err
	}
	for i, event := range events {
		eventsMap[event.EventHash] = events[i]
	}
	unique := len(eventsMap)
	duplicates := total - unique
	return total, unique, duplicates, nil
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
// Возвращает total (всего событий), unique (уникальных), duplicates (дубликатов)
func BloomFilter(p float64, path string) (int, int, int, error) {
	events, _, total, _, err2 := model.ReadEvents(path)
	if err2 != nil {
		return 0, 0, 0, err2
	}
	m, k, err1 := Params(total, p)
	if err1 != nil {
		return 0, 0, 0, err1
	}
	bs, err := bitset.New(uint(m))
	if err != nil {
		return 0, 0, 0, err
	}
	duplicates := 0

	for _, event := range events {
		indexes := getIndexes(event.EventHash, k, uint64(m))
		alreadySet := true
		//если хоть один из значений - 0, то события ещё не было
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
	return total, unique, duplicates, nil
}

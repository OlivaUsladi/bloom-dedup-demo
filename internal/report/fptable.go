package report

import (
	"bloom-dedup-demo/internal/bloom"
	"bloom-dedup-demo/internal/model"
	"fmt"
)

// Формирование Markdown-таблицы сравнения фильтра Блума для нескольких значений false_positive_rate
func BuildFPTable(events []model.Event, expectedItems int, hash string, rates []float64) error {
	if len(rates) == 0 {
		return fmt.Errorf("список false_positive_rate не может быть пустым")
	}
	if hash != "f64" && hash != "s256" {
		return fmt.Errorf("hash должен быть f64 или s256")
	}
	if hash == "f64" {
		hash = "fnv64_double_hashing"
	} else {
		hash = "sha256_slices"
	}
	_, exactUnique, exactDup, _, _, err := bloom.MapFilter(events)
	if err != nil {
		return err
	}

	fmt.Printf("Сравнение значений false_positive_rate\n\n")
	fmt.Printf("Записей: %d, уникальных (map): %d, дублей (map): %d, expected_items: %d\n", len(events), exactUnique, exactDup, expectedItems)
	fmt.Printf("%-5s %-12s %-5s %-15s %-15s %-18s %-18s\n", "p", "m (бит)", "k", "Память (байт)", "Дубли (Блум)", "Ложные срабатывания", "Реальный FP rate")
	fmt.Printf("|-----|-------|-----|----------|---------|----------|----------|\n")

	for _, p := range rates {
		m, k, err := bloom.Params(expectedItems, p)
		if err != nil {
			return fmt.Errorf("p=%v: %w", p, err)
		}
		_, bloomDup, _, memory, err := bloom.BloomFilter(events, expectedItems, hash, p)
		if err != nil {
			return fmt.Errorf("p=%v: %w", p, err)
		}
		estFP := bloomDup - exactDup
		if estFP < 0 {
			estFP = 0
		}
		fpRate := 0.0
		if exactUnique > 0 {
			fpRate = float64(estFP) / float64(exactUnique)
		}
		fmt.Printf("| %g | %d | %d | %d | %d | %d | %.6f |\n",
			p, m, k, memory, bloomDup, estFP, fpRate)
	}
	return nil
}

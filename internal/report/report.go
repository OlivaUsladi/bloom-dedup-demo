package report

import (
	"bloom-dedup-demo/internal/bloom"
	"bloom-dedup-demo/internal/model"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type SourceStats struct {
	TotalRecords            int `json:"total_records"`
	ExactUnique             int `json:"exact_unique"`
	ExactDuplicates         int `json:"exact_duplicates"`
	BloomMayDuplicate       int `json:"bloom_may_duplicate"`
	EstimatedFalsePositives int `json:"estimated_false_positives"`
}

type Report struct {
	TotalRecords            int                    `json:"total_records"`
	BadLines                int                    `json:"bad_lines"`
	ExactUnique             int                    `json:"exact_unique"`
	ExactDuplicates         int                    `json:"exact_duplicates"`
	BloomNew                int                    `json:"bloom_new"`
	BloomMayDuplicate       int                    `json:"bloom_may_duplicate"`
	EstimatedFalsePositives int                    `json:"estimated_false_positives"`
	RealFalsePositiveRate   float64                `json:"real_false_positive_rate"`
	BloomMemoryBytes        int                    `json:"bloom_memory_bytes"`
	ExactMapMemoryBytes     int                    `json:"exact_map_memory_bytes"`
	MapDurationMs           int64                  `json:"map_duration_ms"`
	BloomDurationMs         int64                  `json:"bloom_duration_ms"`
	BySource                map[string]SourceStats `json:"by_source"`
	InvalidSources          []string               `json:"invalid_sources"`
}

// Создание отчёта
// Входные данные: события, плохие строки, плохие источники, файл с конфигурацией, флаг пропуска
func BuildReport(events []model.Event, badLines []int, badSources []string, pathcfg string) (*Report, error) {
	cfg, err := model.ReadConfig(pathcfg)
	if err != nil {
		return nil, err
	}
	total := len(events)

	_, exactUnique, exactDup, mapDuration, mapMemory, err1 := bloom.MapFilter(events)
	if err1 != nil {
		return nil, err1
	}
	var bloomNew, bloomDup, bloomMemory int
	var bloomDuration int64
	var err3 error
	if cfg.Mode == "bloom" {
		bloomNew, bloomDup, bloomDuration, bloomMemory, err3 = bloom.BloomFilter(events, cfg.ExpectedItems, cfg.HashFamily, cfg.FalsePositiveRate)
	} else {
		bloomNew, bloomDup, bloomDuration, bloomMemory, err3 = bloom.CountingBloomFilter(events, cfg.ExpectedItems, cfg.HashFamily, cfg.FalsePositiveRate)
	}
	if err3 != nil {
		return nil, err3
	}

	estFP := bloomDup - exactDup
	if estFP < 0 {
		estFP = 0
	}
	var fpRate float64
	if exactUnique > 0 {
		fpRate = float64(estFP) / float64(exactUnique)
	}

	invalid := uniqueStrings(badSources)

	bS, err4 := BuildBySource(events, cfg.HashFamily, cfg.Mode, cfg.FalsePositiveRate)
	if err4 != nil {
		return nil, err4
	}
	return &Report{
		TotalRecords:            total,
		BadLines:                len(badLines),
		ExactUnique:             exactUnique,
		ExactDuplicates:         exactDup,
		BloomNew:                bloomNew,
		BloomMayDuplicate:       bloomDup,
		EstimatedFalsePositives: estFP,
		RealFalsePositiveRate:   fpRate,
		BloomMemoryBytes:        bloomMemory,
		ExactMapMemoryBytes:     mapMemory,
		MapDurationMs:           mapDuration,
		BloomDurationMs:         bloomDuration,
		BySource:                bS,
		InvalidSources:          invalid,
	}, nil
}

// Отбор уникальных невалидных источников
func uniqueStrings(in []string) []string {
	seen := make(map[string]bool)
	out := []string{}
	for _, s := range in {
		if !seen[s] {
			seen[s] = true
			out = append(out, s)
		}
	}
	return out
}

// Группировка данных по источнику
func BuildBySource(events []model.Event, hash string, mode string, fPR float64) (map[string]SourceStats, error) {
	bySource := make(map[string]SourceStats)
	grouped := make(map[string][]model.Event)
	for _, e := range events {
		grouped[e.Source] = append(grouped[e.Source], e)
	}
	for src, evs := range grouped {
		_, exactUnique, exactDup, _, _, err1 := bloom.MapFilter(evs)
		if err1 != nil {
			return nil, err1
		}
		var bloomDup int
		if mode == "bloom" {
			_, bloomDup, _, _, err1 = bloom.BloomFilter(evs, len(evs), hash, fPR)
			if err1 != nil {
				return nil, err1
			}
		} else {
			_, bloomDup, _, _, err1 = bloom.CountingBloomFilter(evs, len(evs), hash, fPR)
			if err1 != nil {
				return nil, err1
			}
		}
		estFP := bloomDup - exactDup
		if estFP < 0 {
			estFP = 0
		}

		bySource[src] = SourceStats{
			TotalRecords:            len(evs),
			ExactUnique:             exactUnique,
			ExactDuplicates:         exactDup,
			BloomMayDuplicate:       bloomDup,
			EstimatedFalsePositives: estFP,
		}
	}
	return bySource, nil
}

func SaveJSON(path string, report *Report) error {
	byteValue, err := json.MarshalIndent(report, "", " ")
	if err != nil {
		return fmt.Errorf("не удалось сериализовать отчёт: %w", err)
	}
	err = os.MkdirAll(filepath.Dir(path), 0755)
	if err != nil {
		return fmt.Errorf("не удалось создать директорию для %s: %w", path, err)
	}
	err = os.WriteFile(path, byteValue, 0644)
	if err != nil {
		return fmt.Errorf("не удалось записать файл %s: %w", path, err)
	}
	return nil
}

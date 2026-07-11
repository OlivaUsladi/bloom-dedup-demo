package report

import (
	"bloom-dedup-demo/internal/bloom"
	"bloom-dedup-demo/internal/model"
	"time"
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
	DurationMs              int64                  `json:"duration_ms"`
	BySource                map[string]SourceStats `json:"by_source"`
	InvalidSources          []string               `json:"invalid_sources"`
}

// Создание отчёта
func BuildReport(path string, cfg *model.Config, strict bool) (*Report, error) {
	events, badLines, total, badSources, err := model.ReadEvents(path, strict)
	if err != nil {
		return nil, err
	}
	start := time.Now()
	exactUnique, exactDup, err1 := bloom.MapFilter(path, strict)
	if err1 != nil {
		return nil, err1
	}

	//Если оценка памяти будет внутри фильтров, то убрать этот вызов
	m, _, err2 := bloom.Params(cfg.ExpectedItems, cfg.FalsePositiveRate)

	if err2 != nil {
		return nil, err2
	}
	bloomNew, bloomDub, err3 := bloom.BloomFilter(cfg.FalsePositiveRate, path, strict)
	if err3 != nil {
		return nil, err3
	}
	duration := time.Since(start).Milliseconds()

	estFP := bloomDub - exactDup
	var fpRate float64
	if exactUnique > 0 {
		fpRate = float64(estFP) / float64(exactUnique)
	}

	invalid := uniqueStrings(badSources)

	//память и время перенести внутрь фильтров
	return &Report{
		TotalRecords:            total,
		BadLines:                len(badLines),
		ExactUnique:             exactUnique,
		ExactDuplicates:         exactDup,
		BloomNew:                bloomNew,
		BloomMayDuplicate:       bloomDub,
		EstimatedFalsePositives: estFP,
		RealFalsePositiveRate:   fpRate,
		BloomMemoryBytes:        (m + 7) / 8,      //память
		ExactMapMemoryBytes:     exactUnique * 80, //пмять
		DurationMs:              duration,
		BySource:                BuildBySource(events),
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

// пока нет bloom_may_duplicate/estimated_false_positives, надо делать вызов фильтра Блума
// Группировка данных по источнику
func BuildBySource(events []model.Event) map[string]SourceStats {
	bySource := make(map[string]SourceStats)
	grouped := make(map[string][]model.Event)
	for _, e := range events {
		grouped[e.Source] = append(grouped[e.Source], e)
	}
	for src, evs := range grouped {
		seen := make(map[string]bool)
		dup := 0
		for _, e := range evs {
			if seen[e.EventHash] {
				dup++
			}
			seen[e.EventHash] = true
		}
		bySource[src] = SourceStats{
			TotalRecords:    len(evs),
			ExactUnique:     len(seen),
			ExactDuplicates: dup,
		}
	}
	return bySource
}

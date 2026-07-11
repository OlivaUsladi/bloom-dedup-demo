package report

import (
	"bloom-dedup-demo/internal/model"
	"fmt"
	"testing"
)

func TestBuildReport_Basic(t *testing.T) {
	cfg := &model.Config{
		ExpectedItems:     4,
		FalsePositiveRate: 0.01,
		HashFamily:        "fnv64_double_hashing",
		Mode:              "bloom",
	}

	report, err := BuildReport("../../testdata/control/event.jsonl", cfg, true)
	if err != nil {
		t.Fatalf("BuildReport вернул ошибку: %v", err)
	}

	if report.TotalRecords != 4 {
		t.Errorf("ожидали TotalRecords=4, получили %d", report.TotalRecords)
	}
	if report.ExactUnique != 4 {
		t.Errorf("ожидали ExactUnique=4, получили %d", report.ExactUnique)
	}
	if report.ExactDuplicates != 0 {
		t.Errorf("ожидали ExactDuplicates=0, получили %d", report.ExactDuplicates)
	}
	if report.BloomMemoryBytes <= 0 {
		t.Errorf("ожидали положительный BloomMemoryBytes, получили %d", report.BloomMemoryBytes)
	}
	if report.ExactMapMemoryBytes <= 0 {
		t.Errorf("ожидали положительный ExactMapMemoryBytes, получили %d", report.ExactMapMemoryBytes)
	}
	if report.DurationMs < 0 {
		t.Errorf("DurationMs не может быть отрицательным, получили %d", report.DurationMs)
	}
	if report.BySource == nil {
		t.Errorf("ожидали ненулевую карту BySource")
	}
	if report.InvalidSources == nil {
		t.Errorf("ожидали ненулевой срез InvalidSources")
	}
	fmt.Println(report)
}

package report

import (
	"fmt"
	"testing"
)

func TestBuildReport(t *testing.T) {

	report, err := BuildReport("../../testdata/control/event.jsonl", "../../testdata/control/bloom1.json", true)
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
	if report.MapDurationMs < 0 {
		t.Errorf("MapDurationMs не может быть отрицательным, получили %d", report.MapDurationMs)
	}
	if report.BloomDurationMs < 0 {
		t.Errorf("BloomDurationMs не может быть отрицательным, получили %d", report.BloomDurationMs)
	}
	if report.BySource == nil {
		t.Errorf("ожидали ненулевую карту BySource")
	}
	if report.InvalidSources == nil {
		t.Errorf("ожидали ненулевой срез InvalidSources")
	}
	fmt.Println(report)
}

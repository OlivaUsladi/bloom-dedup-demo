package report

import (
	"bloom-dedup-demo/internal/model"
	"testing"
)

func TestBuildFPTable(t *testing.T) {
	events, _, _, err := model.ReadEvents("../../testdata/tests/event3.jsonl", true)
	if err != nil {
		t.Fatalf("ReadEvents вернул ошибку: %v", err)
	}
	if err := BuildFPTable(events, len(events), "f64", []float64{0.1, 0.05, 0.01, 0.001}); err != nil {
		t.Errorf("BuildFPTable вернул ошибку: %v", err)
	}
}

func TestBuildFPTableError1(t *testing.T) {
	events, _, _, err := model.ReadEvents("../../testdata/tests/event.jsonl", true)
	if err != nil {
		t.Fatalf("ReadEvents вернул ошибку: %v", err)
	}
	if err := BuildFPTable(events, len(events), "f64", nil); err == nil {
		t.Errorf("ожидали ошибку для пустого списка rates")
	}
}

func TestBuildFPTableError2(t *testing.T) {
	events, _, _, err := model.ReadEvents("../../testdata/tests/event.jsonl", true)
	if err != nil {
		t.Fatalf("ReadEvents вернул ошибку: %v", err)
	}
	if err := BuildFPTable(events, len(events), "f64", []float64{1.5}); err == nil {
		t.Errorf("ожидали ошибку для p > 1")
	}
}

func TestBuildFPTableError3(t *testing.T) {
	events, _, _, err := model.ReadEvents("../../testdata/tests/event.jsonl", true)
	if err != nil {
		t.Fatalf("ReadEvents вернул ошибку: %v", err)
	}
	if err := BuildFPTable(events, len(events), "ftratata64", []float64{0.1, 0.01}); err == nil {
		t.Errorf("ожидали ошибку для несуществующего hash")
	}
}

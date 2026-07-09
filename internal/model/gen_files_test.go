package model

import (
	"testing"
)

func TestGenerateEvents1(t *testing.T) {
	path := "../../testdata/control/event1.jsonl"

	n := 100
	_ = GenerateEvents(path, n, 5, 0.1, 42)

	_, badLines, total, sour, err := ReadEvents(path)
	if err != nil {
		t.Fatalf("ReadEvents вернул ошибку: %v", err)
	}
	if len(badLines) != 0 {
		t.Errorf("ожидали 0 битых строк, получили %v", badLines)
	}
	if total != n {
		t.Errorf("ожидали %d событий, получили %d", n, total)
	}
	if len(sour) != 0 {
		t.Errorf("ожидали 0 невалидных источников, получили %v", len(sour))
	}
}

func TestGenerateEvents2(t *testing.T) {
	path := "../../testdata/control/event2.jsonl"

	n := 50
	_ = GenerateEvents(path, n, 1, 0.0, 7)

	events, badLines, total, sour, err := ReadEvents(path)
	if err != nil {
		t.Fatalf("ReadEvents вернул ошибку: %v", err)
	}
	if len(badLines) != 0 {
		t.Errorf("ожидали 0 битых строк, получили %v", badLines)
	}
	if total != n {
		t.Errorf("ожидали %d событий, получили %d", n, total)
	}
	if len(sour) != 0 {
		t.Errorf("ожидали 0 невалидных источников, получили %v", len(sour))
	}
	for _, e := range events {
		if e.Source != "collector_01" {
			t.Errorf("при s=1 source должен быть collector_01, получено %q", e.Source)
		}
	}
}

func TestGenerateEvents3(t *testing.T) {
	path := "../../testdata/control/event3.jsonl"

	n := 200
	_ = GenerateEvents(path, n, 20, 0.5, 123)

	_, badLines, total, sour, err := ReadEvents(path)
	if err != nil {
		t.Fatalf("ReadEvents вернул ошибку: %v", err)
	}
	if len(badLines) != 0 {
		t.Errorf("ожидали 0 битых строк, получили %v", badLines)
	}
	if total != n {
		t.Errorf("ожидали %d событий, получили %d", n, total)
	}
	if len(sour) != 0 {
		t.Errorf("ожидали 0 невалидных источников, получили %v", len(sour))
	}
}

func TestGenerateInvalidPathEvents(t *testing.T) {
	path := ""
	n := 20
	err := GenerateEvents(path, n, 20, 0.5, 123)
	if err == nil {
		t.Errorf("ожидали ошибку аргумента, получили nil")
	}
}

func TestGenerateZeroNEvents(t *testing.T) {
	path := "../../testdata/control/event7.jsonl"
	n := 0
	err := GenerateEvents(path, n, 20, 0.5, 123)
	if err == nil {
		t.Errorf("ожидали ошибку аргумента, получили nil")
	}
}

func TestGenerateInvalidRate(t *testing.T) {
	path := "../../testdata/control/event7.jsonl"
	n := 20
	err := GenerateEvents(path, n, 20, 1.5, 123)
	if err == nil {
		t.Errorf("ожидали ошибку аргумента, получили nil")
	}
}

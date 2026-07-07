package model

import (
	"testing"
)

func TestGenerateEvents1(t *testing.T) {
	path := "../../testdata/control/event1.jsonl"
	//defer os.Remove(path)

	n := 100
	GenerateEvents(path, n, 5, 0.1, 42)

	events, badLines, err := ReadEvents(path)
	if err != nil {
		t.Fatalf("ReadEvents вернул ошибку: %v", err)
	}
	if len(badLines) != 0 {
		t.Errorf("ожидали 0 битых строк, получили %v", badLines)
	}
	if len(events) != n {
		t.Errorf("ожидали %d событий, получили %d", n, len(events))
	}
}

func TestGenerateEvents2(t *testing.T) {
	path := "../../testdata/control/event2.jsonl"

	n := 50
	GenerateEvents(path, n, 1, 0.0, 7)

	events, badLines, err := ReadEvents(path)
	if err != nil {
		t.Fatalf("ReadEvents вернул ошибку: %v", err)
	}
	if len(badLines) != 0 {
		t.Errorf("ожидали 0 битых строк, получили %v", badLines)
	}
	if len(events) != n {
		t.Errorf("ожидали %d событий, получили %d", n, len(events))
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
	GenerateEvents(path, n, 20, 0.5, 123)

	events, badLines, err := ReadEvents(path)
	if err != nil {
		t.Fatalf("ReadEvents вернул ошибку: %v", err)
	}
	if len(badLines) != 0 {
		t.Errorf("ожидали 0 битых строк, получили %v", badLines)
	}
	if len(events) != n {
		t.Errorf("ожидали %d событий, получили %d", n, len(events))
	}
}

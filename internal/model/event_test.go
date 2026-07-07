package model

import "testing"

func TestReadEvents(t *testing.T) {
	events, badLines, err := ReadEvents("../../testdata/control/event.jsonl")
	if err != nil {
		t.Fatalf("ReadEvents вернул ошибку: %v", err)
	}

	if len(events) != 4 {
		t.Errorf("ожидали 4 валидных события, получили %d", len(events))
	}
	if len(badLines) != 1 || badLines[0] != 3 {
		t.Errorf("ожидали одну битую строку с номером 3, получили %v", badLines)
	}

	for _, e := range events {
		if err := e.Validate(); err != nil {
			t.Errorf("событие %s не прошло валидацию: %v", e.EventID, err)
		}
	}
}

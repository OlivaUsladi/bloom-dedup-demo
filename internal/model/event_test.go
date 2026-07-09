package model

import "testing"

func TestReadEvents(t *testing.T) {
	events, badLines, total, sour, err := ReadEvents("../../testdata/control/event.jsonl")
	if err != nil {
		t.Fatalf("ReadEvents вернул ошибку: %v", err)
	}

	if total != 4 {
		t.Errorf("ожидали 4 валидных события, получили %d", total)
	}
	if len(badLines) != 1 || badLines[0] != 3 {
		t.Errorf("ожидали одну битую строку с номером 3, получили %v", badLines)
	}

	if len(sour) != 0 {
		t.Errorf("ожидали 0 невалидных источников, получили %v", len(sour))
	}
	for _, e := range events {
		if err := e.Validate(); err != nil {
			t.Errorf("событие %s не прошло валидацию: %v", e.EventID, err)
		}
	}
}

func TestReadBadEvents(t *testing.T) {
	_, badLines, total, sour, err := ReadEvents("../../testdata/control/bad_lines.jsonl")
	if err != nil {
		t.Fatalf("ReadEvents выдал ошибку: %v", err)
	}
	if len(badLines) != 2 {
		t.Errorf("ожидали 2 битые строки, получили %d", len(badLines))
	}
	if total != 4 {
		t.Errorf("ожидали 4 валидных строк, получили %d", total)
	}
	if len(sour) != 0 {

	}
}

func TestEventValidateSource(t *testing.T) {
	e := Event{Source: "bad"}
	valid, _ := e.ValidateSource()
	if valid {
		t.Errorf("ожидалась ошибка, получено значение true")
	}
}

func TestEventValidateEmptySource(t *testing.T) {
	e := Event{Source: ""}
	valid, _ := e.ValidateSource()
	if !valid {
		t.Errorf("пустой source должен быть валидным, получено invalid")
	}
}

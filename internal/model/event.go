package model

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"
)

type Event struct {
	Seq       int    `json:"seq"`
	EventID   string `json:"event_id"`
	EventHash string `json:"event_hash"`
	Source    string `json:"source"`
	Timestamp string `json:"timestamp"`
}

// Валидация
func (e *Event) Validate() error {
	if e.Seq <= 0 {
		return fmt.Errorf("seq должен быть больше 0, получено %d", e.Seq)
	}
	if e.EventID == "" {
		return fmt.Errorf("event_id обязателен и не может быть пустым")
	}
	if len(e.EventID) < 10 {
		return fmt.Errorf("event_id должен быть не менее 10 символов, получено %q длиной %d", e.EventID, len(e.EventID))
	}
	if e.EventID[:4] != "evt_" {
		return fmt.Errorf("event_id должен начинаться с 'evt_', получено %q", e.EventID)
	}
	if e.EventHash == "" {
		return fmt.Errorf("event_hash обязателен и не может быть пустым")
	}
	if len(e.EventHash) != 16 && len(e.EventHash) != 32 {
		return fmt.Errorf("event_hash должен быть 16 или 32 hex-символа, получено %d символов", len(e.EventHash))
	}
	if len(e.EventID) > 256 {
		return fmt.Errorf("event_id превышает допустимую длину %d символов", 256)
	}
	if len(e.Source) > 256 {
		return fmt.Errorf("source превышает допустимую длину %d символов", 256)
	}
	if e.Timestamp != "" {
		if _, err := time.Parse(time.RFC3339, e.Timestamp); err != nil {
			return fmt.Errorf("timestamp должен быть в формате RFC3339, получено %q", e.Timestamp)
		}
	}
	return nil
}

// Валидация источника
func (e *Event) ValidateSource() (bool, string) {
	if e.Source == "" {
		return true, ""
	}
	if len(e.Source) != 12 {
		return false, e.Source
	}
	if e.Source[:10] != "collector_" {
		return false, e.Source
	}
	i := e.Source[10:]
	iInt, err := strconv.Atoi(i)
	if err != nil || iInt < 1 || iInt > 99 {
		return false, e.Source
	}
	return true, ""
}

// Считывание файла входных данных построчно
// Получает: путь к файлу
// Возвращает: массив распарсенных данных, массив пропущенных битых строк,
// количество элементов, список невалидных источников, ошибку
// Если какого-то элемента выходных данных нет - возвращает nil
func ReadEvents(path string) ([]Event, []int, int, []string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, nil, 0, nil, fmt.Errorf("не удалось открыть файл %s: %w", path, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)
	events := []Event{}
	badLines := []int{}
	badSources := []string{}
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		var event Event
		if err := json.Unmarshal([]byte(line), &event); err != nil {
			badLines = append(badLines, lineNum)
			continue
		}
		err1 := event.Validate()
		flag, str := event.ValidateSource()

		//Вот тут сделать флаг для невалидных источников
		if flag == true {
			if err1 == nil {
				events = append(events, event)
			}
		} else {
			badSources = append(badSources, str)
		}
	}

	n := len(events)
	if err := scanner.Err(); err != nil {
		return nil, nil, 0, nil, fmt.Errorf("ошибка чтения файла %s: %w", path, err)
	}

	return events, badLines, n, badSources, nil
}

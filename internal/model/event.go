package model

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

type Event struct {
	Seq       int    `json:"seq"`
	EventID   string `json:"event_id"`
	EventHash string `json:"event_hash"`
	Source    string `json:"source"`
	Timestamp string `json:"timestamp"`
}

// Считывание файла входных данных построчно
// Получает: путь к файлу
// Возвращает: массив распарсенных данных, массив пропущенных битых строк, ошибку
// Если какого-то элемента выходных данных нет - возвращает nil
func ReadEvents(path string) ([]Event, []int, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, nil, fmt.Errorf("не удалось открыть файл %s: %w", path, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	events := []Event{}
	badLines := []int{}
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		var event Event
		if err := json.Unmarshal([]byte(line), &event); err != nil {
			badLines = append(badLines, lineNum)
			continue
		}
		events = append(events, event)
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, fmt.Errorf("ошибка чтения файла %s: %w", path, err)
	}

	return events, badLines, nil
}

// Валидация
func (e *Event) Validate() error {
	if e.Seq <= 0 {
		return fmt.Errorf("seq должен быть больше 0, получено %d", e.Seq)
	}
	if e.EventID == "" {
		return fmt.Errorf("event_id обязателен и не может быть пустым")
	}
	if e.EventHash == "" {
		return fmt.Errorf("event_hash обязателен и не может быть пустым")
	}
	return nil
}

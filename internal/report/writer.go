package report

import (
	"bloom-dedup-demo/internal/model"
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

// Запись дедуплицированных событий в файл
func WriteEvents(path string, events []model.Event) error {
	if path == "" {
		return fmt.Errorf("неправильный аргумент path")
	}
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("не удалось создать файл %s: %w", path, err)
	}
	defer file.Close()
	writer := bufio.NewWriter(file)

	for _, event := range events {
		data, err := json.Marshal(event)
		if err != nil {
			fmt.Fprintln(os.Stderr, "ошибка преобразования (marshal):", err)
			continue
		}
		if _, err := writer.Write(data); err != nil {
			return fmt.Errorf("ошибка записи события %s: %w", event.EventID, err)
		}
		if err := writer.WriteByte('\n'); err != nil {
			return fmt.Errorf("ошибка записи перевода строки: %w", err)
		}
	}
	if err := writer.Flush(); err != nil {
		return fmt.Errorf("ошибка записи в файл: %w", err)
	}
	return nil
}

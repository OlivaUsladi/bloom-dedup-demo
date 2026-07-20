package report

import (
	"bloom-dedup-demo/internal/model"
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Запись дедуплицированных событий в файл
func WriteEvents(path string, events []model.Event) error {
	if path == "" {
		return fmt.Errorf("неправильный аргумент path")
	}
	err := os.MkdirAll(filepath.Dir(path), 0755)
	if err != nil {
		return fmt.Errorf("не удалось создать директорию для %s: %w", path, err)
	}
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("не удалось создать файл %s: %w", path, err)
	}
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
	err1 := file.Close()
	if err1 != nil {
		return fmt.Errorf("ошибка закрытия файла: %w", err1)
	}
	return nil
}

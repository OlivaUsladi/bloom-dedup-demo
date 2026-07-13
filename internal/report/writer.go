package report

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func SaveJSON(path string, report *Report) error {
	byteValue, err := json.MarshalIndent(report, "", " ")
	if err != nil {
		return fmt.Errorf("не удалось сериализовать отчёт: %w", err)
	}
	err = os.MkdirAll(filepath.Dir(path), 0755)
	if err != nil {
		return fmt.Errorf("не удалось создать директорию для %s: %w", path, err)
	}
	err = os.WriteFile(path, byteValue, 0644)
	if err != nil {
		return fmt.Errorf("не удалось записать файл %s: %w", path, err)
	}
	return nil
}

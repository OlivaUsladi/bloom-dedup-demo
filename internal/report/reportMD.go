package report

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
)

func SaveMarkdown(path string, report *Report) error {
	if path == "" {
		return fmt.Errorf("неправильный аргумент path")
	}
	if report == nil {
		return fmt.Errorf("отчёт не может быть nil")
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

	if _, err := writer.WriteString("# Отчёт о выполненной фильтрации"); err != nil {
		return fmt.Errorf("ошибка записи строки: %w", err)
	}
	if err := writer.WriteByte('\n'); err != nil {
		return fmt.Errorf("ошибка записи перевода строки: %w", err)
	}
	if err := writer.WriteByte('\n'); err != nil {
		return fmt.Errorf("ошибка записи перевода строки: %w", err)
	}

	if report.RealFalsePositiveRate == 0.0 && report.MapDurationMs == 0 && report.ExactMapMemoryBytes == 0 && report.ExactDuplicates == 0 && report.ExactUnique == 0 {
		if _, err := writer.WriteString("# Map фильтр не запускался, все его значения обнулены"); err != nil {
			return fmt.Errorf("ошибка записи строки: %w", err)
		}
		if err := writer.WriteByte('\n'); err != nil {
			return fmt.Errorf("ошибка записи перевода строки: %w", err)
		}
		if err := writer.WriteByte('\n'); err != nil {
			return fmt.Errorf("ошибка записи перевода строки: %w", err)
		}
	}
	if _, err := writer.WriteString("## Общая статистика"); err != nil {
		return fmt.Errorf("ошибка записи строки: %w", err)
	}
	if err := writer.WriteByte('\n'); err != nil {
		return fmt.Errorf("ошибка записи перевода строки: %w", err)
	}
	if err := writer.WriteByte('\n'); err != nil {
		return fmt.Errorf("ошибка записи перевода строки: %w", err)
	}

	if _, err := writer.WriteString("- Общее количество строк: " + strconv.Itoa(report.TotalRecords)); err != nil {
		return fmt.Errorf("ошибка записи строки: %w", err)
	}
	if err := writer.WriteByte('\n'); err != nil {
		return fmt.Errorf("ошибка записи перевода строки: %w", err)
	}
	if _, err := writer.WriteString("- Выявлено битых/невалидных строк: " + strconv.Itoa(report.BadLines)); err != nil {
		return fmt.Errorf("ошибка записи строки: %w", err)
	}
	if err := writer.WriteByte('\n'); err != nil {
		return fmt.Errorf("ошибка записи перевода строки: %w", err)
	}
	if _, err := writer.WriteString("- Точное количество уникальных событий: " + strconv.Itoa(report.ExactUnique)); err != nil {
		return fmt.Errorf("ошибка записи строки: %w", err)
	}
	if err := writer.WriteByte('\n'); err != nil {
		return fmt.Errorf("ошибка записи перевода строки: %w", err)
	}
	if _, err := writer.WriteString("- Точное количество дубликатов: " + strconv.Itoa(report.ExactDuplicates)); err != nil {
		return fmt.Errorf("ошибка записи строки: %w", err)
	}
	if err := writer.WriteByte('\n'); err != nil {
		return fmt.Errorf("ошибка записи перевода строки: %w", err)
	}
	if _, err := writer.WriteString("- Уникальные события по фильтру Блума: " + strconv.Itoa(report.BloomNew)); err != nil {
		return fmt.Errorf("ошибка записи строки: %w", err)
	}
	if err := writer.WriteByte('\n'); err != nil {
		return fmt.Errorf("ошибка записи перевода строки: %w", err)
	}
	if _, err := writer.WriteString("- Дубликаты по фильтру Блума: " + strconv.Itoa(report.BloomMayDuplicate)); err != nil {
		return fmt.Errorf("ошибка записи строки: %w", err)
	}
	if err := writer.WriteByte('\n'); err != nil {
		return fmt.Errorf("ошибка записи перевода строки: %w", err)
	}
	if _, err := writer.WriteString("- Число случаев, когда фильтр Блума ошибочно счёл новый элемент дублем: " + strconv.Itoa(report.EstimatedFalsePositives)); err != nil {
		return fmt.Errorf("ошибка записи строки: %w", err)
	}
	if err := writer.WriteByte('\n'); err != nil {
		return fmt.Errorf("ошибка записи перевода строки: %w", err)
	}
	if _, err := writer.WriteString("- Измеренная доля ложных срабатываний: " + strconv.FormatFloat(report.RealFalsePositiveRate, 'f', 10, 64)); err != nil {
		return fmt.Errorf("ошибка записи строки: %w", err)
	}
	if err := writer.WriteByte('\n'); err != nil {
		return fmt.Errorf("ошибка записи перевода строки: %w", err)
	}
	if _, err := writer.WriteString("- Количество байт для фильтра Блума: " + strconv.Itoa(report.BloomMemoryBytes)); err != nil {
		return fmt.Errorf("ошибка записи строки: %w", err)
	}
	if err := writer.WriteByte('\n'); err != nil {
		return fmt.Errorf("ошибка записи перевода строки: %w", err)
	}
	bloomMiB := float64(report.BloomMemoryBytes) / 1024.0 / 1024.0
	exactMiB := float64(report.ExactMapMemoryBytes) / 1024.0 / 1024.0
	if bloomMiB != 0.0 {
		if _, err := writer.WriteString(fmt.Sprintf("- Размер фильтра Блума: %.2f Mб", bloomMiB)); err != nil {
			return fmt.Errorf("ошибка записи строки: %w", err)
		}
		if err := writer.WriteByte('\n'); err != nil {
			return fmt.Errorf("ошибка записи перевода строки: %w", err)
		}
	}
	if _, err := writer.WriteString("- Количество байт для точного сравнения: " + strconv.Itoa(report.ExactMapMemoryBytes)); err != nil {
		return fmt.Errorf("ошибка записи строки: %w", err)
	}
	if err := writer.WriteByte('\n'); err != nil {
		return fmt.Errorf("ошибка записи перевода строки: %w", err)
	}
	if exactMiB != 0.0 {
		if _, err := writer.WriteString(fmt.Sprintf("- Размер структуры точного сравнения: %.2f Mб", exactMiB)); err != nil {
			return fmt.Errorf("ошибка записи строки: %w", err)
		}
		if err := writer.WriteByte('\n'); err != nil {
			return fmt.Errorf("ошибка записи перевода строки: %w", err)
		}
	}
	if _, err := writer.WriteString("- Время выполнения фильтра Блума (мс): " + strconv.Itoa(int(report.BloomDurationMs))); err != nil {
		return fmt.Errorf("ошибка записи строки: %w", err)
	}
	if err := writer.WriteByte('\n'); err != nil {
		return fmt.Errorf("ошибка записи перевода строки: %w", err)
	}
	if _, err := writer.WriteString("- Время выполнения точного сравнения (мс): " + strconv.Itoa(int(report.MapDurationMs))); err != nil {
		return fmt.Errorf("ошибка записи строки: %w", err)
	}
	if err := writer.WriteByte('\n'); err != nil {
		return fmt.Errorf("ошибка записи перевода строки: %w", err)
	}
	if err := writer.WriteByte('\n'); err != nil {
		return fmt.Errorf("ошибка записи перевода строки: %w", err)
	}

	if _, err := writer.WriteString("## Статистика по источникам"); err != nil {
		return fmt.Errorf("ошибка записи строки: %w", err)
	}
	if err := writer.WriteByte('\n'); err != nil {
		return fmt.Errorf("ошибка записи перевода строки: %w", err)
	}
	if err := writer.WriteByte('\n'); err != nil {
		return fmt.Errorf("ошибка записи перевода строки: %w", err)
	}

	sources := make([]string, 0, len(report.BySource))
	for source := range report.BySource {
		sources = append(sources, source)
	}
	sort.Strings(sources)

	for _, source := range sources {
		bs := report.BySource[source]
		if _, err := writer.WriteString("### " + source); err != nil {
			return fmt.Errorf("ошибка записи строки: %w", err)
		}
		if err := writer.WriteByte('\n'); err != nil {
			return fmt.Errorf("ошибка записи перевода строки: %w", err)
		}
		if err := writer.WriteByte('\n'); err != nil {
			return fmt.Errorf("ошибка записи перевода строки: %w", err)
		}

		if _, err := writer.WriteString("- Общее количество строк: " + strconv.Itoa(bs.TotalRecords)); err != nil {
			return fmt.Errorf("ошибка записи строки: %w", err)
		}
		if err := writer.WriteByte('\n'); err != nil {
			return fmt.Errorf("ошибка записи перевода строки: %w", err)
		}
		if _, err := writer.WriteString("- Точное количество уникальных событий: " + strconv.Itoa(bs.ExactUnique)); err != nil {
			return fmt.Errorf("ошибка записи строки: %w", err)
		}
		if err := writer.WriteByte('\n'); err != nil {
			return fmt.Errorf("ошибка записи перевода строки: %w", err)
		}
		if _, err := writer.WriteString("- Точное количество дубликатов: " + strconv.Itoa(bs.ExactDuplicates)); err != nil {
			return fmt.Errorf("ошибка записи строки: %w", err)
		}
		if err := writer.WriteByte('\n'); err != nil {
			return fmt.Errorf("ошибка записи перевода строки: %w", err)
		}
		if _, err := writer.WriteString("- Дубликаты по фильтру Блума: " + strconv.Itoa(bs.BloomMayDuplicate)); err != nil {
			return fmt.Errorf("ошибка записи строки: %w", err)
		}
		if err := writer.WriteByte('\n'); err != nil {
			return fmt.Errorf("ошибка записи перевода строки: %w", err)
		}
		if _, err := writer.WriteString("- Число случаев, когда фильтр Блума ошибочно счёл новый элемент дублем: " + strconv.Itoa(bs.EstimatedFalsePositives)); err != nil {
			return fmt.Errorf("ошибка записи строки: %w", err)
		}
		if err := writer.WriteByte('\n'); err != nil {
			return fmt.Errorf("ошибка записи перевода строки: %w", err)
		}
		if err := writer.WriteByte('\n'); err != nil {
			return fmt.Errorf("ошибка записи перевода строки: %w", err)
		}
	}

	if _, err := writer.WriteString("## Невалидные источники"); err != nil {
		return fmt.Errorf("ошибка записи строки: %w", err)
	}
	if err := writer.WriteByte('\n'); err != nil {
		return fmt.Errorf("ошибка записи перевода строки: %w", err)
	}
	if err := writer.WriteByte('\n'); err != nil {
		return fmt.Errorf("ошибка записи перевода строки: %w", err)
	}

	if len(report.InvalidSources) == 0 {
		if _, err := writer.WriteString("Все источники были валидными."); err != nil {
			return fmt.Errorf("ошибка записи строки: %w", err)
		}
		if err := writer.WriteByte('\n'); err != nil {
			return fmt.Errorf("ошибка записи перевода строки: %w", err)
		}
	} else {
		for _, s := range report.InvalidSources {
			if _, err := writer.WriteString("- " + s); err != nil {
				return fmt.Errorf("ошибка записи строки: %w", err)
			}
			if err := writer.WriteByte('\n'); err != nil {
				return fmt.Errorf("ошибка записи перевода строки: %w", err)
			}
		}
	}

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("ошибка записи в файл: %w", err)
	}
	err = file.Close()
	if err != nil {
		return fmt.Errorf("ошибка закрытия файла: %w", err)
	}
	return nil
}

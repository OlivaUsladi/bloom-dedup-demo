package model

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	ExpectedItems     int     `json:"expected_items"`
	FalsePositiveRate float64 `json:"false_positive_rate"`
	HashFamily        string  `json:"hash_family"`
	Mode              string  `json:"mode"`
}

// Считывание файла
func ReadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("не удалось открыть файл Config %s: %w", path, err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("не удалось разобрать файл Config %s: %w", path, err)
	}
	if err1 := config.Validate(); err1 != nil {
		return nil, err1
	}

	return &config, nil
}

// Валидация данных
func (c *Config) Validate() error {
	if c.ExpectedItems <= 0 {
		return fmt.Errorf("expected_items должен быть больше 0, получено %d", c.ExpectedItems)
	}
	if c.FalsePositiveRate <= 0 || c.FalsePositiveRate >= 1 {
		return fmt.Errorf("false_positive_rate должен быть между 0 и 1, получено %f", c.FalsePositiveRate)
	}

	if c.HashFamily != "fnv64_double_hashing" && c.HashFamily != "sha256_slices" {
		return fmt.Errorf("hash_family должен быть fnv64_double_hashing или sha256_slices, получено %q", c.HashFamily)
	}

	if c.Mode != "bloom" && c.Mode != "counting_bloom" {
		return fmt.Errorf("mode должен быть bloom или counting_bloom, получено %q", c.Mode)
	}

	return nil
}

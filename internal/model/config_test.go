package model

import "testing"

func TestReadConfig(t *testing.T) {
	config, err := ReadConfig("../../testdata/control/bloom.json")
	if err != nil {
		t.Fatalf("ReadConfig вернул ошибку: %v", err)
	}

	if err := config.Validate(); err != nil {
		t.Errorf("Config не прошёл валидацию: %v", err)
	}

	if config.ExpectedItems != 1000000 {
		t.Errorf("ожидали ExpectedItems=1000000, получили %d", config.ExpectedItems)
	}
}

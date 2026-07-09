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

func TestReadConfigBadHash(t *testing.T) {
	config, err := ReadConfig("../../testdata/control/bad_config_hash.json")
	if err == nil {
		t.Errorf("ожидали ошибку при неправильном hash_family, получили nil и hash = %s", config.HashFamily)
	}
}

func TestReadConfigBadMode(t *testing.T) {
	config, err := ReadConfig("../../testdata/control/bad_config_mode.json")
	if err == nil {
		t.Errorf("ожидали ошибку при неправильном mode, получили nil и mode = %s", config.Mode)
	}
}

func TestReadConfigInvalidRate(t *testing.T) {
	config, err := ReadConfig("../../testdata/control/bad_config_invalid_rate.json")
	if err == nil {
		t.Errorf("ожидали ошибку при неправильном rate, получили nil и rate = %f", config.FalsePositiveRate)
	}
}

func TestReadConfigNegativeItems(t *testing.T) {
	config, err := ReadConfig("../../testdata/control/bad_config_negative_items.json")
	if err == nil {
		t.Errorf("ожидали ошибку при отрицательном expected_items, получили nil и n = %v", config.ExpectedItems)
	}
}

package config

import (
	"encoding/json"
	"os"

	"yadro-impulse/pkg/errors"
)

func LoadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, errors.NewFileNotFoundError(filename, err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	decoder.DisallowUnknownFields()

	var cfg Config
	if err := decoder.Decode(&cfg); err != nil {
		return nil, errors.NewInvalidConfigError("failed to load config", err)
	}

	return &cfg, nil
}

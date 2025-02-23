package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	ClientPath      string `json:"client_path"`
	VersionFileName string `json:"version_file_name"`
	BackendAddr     string `json:"backend_addr"`
}

func Read() (*Config, error) {
	b, err := os.ReadFile("settings/config.json")
	if err != nil {
		return nil, err
	}
	var c Config
	err = json.Unmarshal(b, &c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

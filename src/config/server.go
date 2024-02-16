package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type ServerConfig struct {
	HttpAddr    string `yaml:"http_addr"`
	LogLevel    string `yaml:"log_level"`
	LogFilePath string `yaml:"log_file_path"`
}

func LoadServer(filename string) (*ServerConfig, error) {
	cfg := &ServerConfig{}

	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(content, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

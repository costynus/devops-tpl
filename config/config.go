package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		Agent  `yaml:"agent"`
		Server `yaml:"server"`
		Log    `yaml:"logger"`
	}
	Agent struct {
		Name             string   `yaml:"name"`
		Version          string   `yaml:"version"`
		PollInterval     int64    `yaml:"pollInterval"`
		ReportInterval   int64    `yaml:"reportInterval"`
		ServerURL        string   `yaml:"server_url"`
		MetricFieldNames []string `yaml:"metric_field_names"`
	}
	Server struct {
		Name    string `yaml:"name"`
		Version string `yaml:"version"`
		Port    string `yaml:"port"`
	}
	Log struct {
		Level string `yaml:"log_level"`
	}
)

func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig("./config/config.yml", cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	return cfg, nil
}

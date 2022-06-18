package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		Agent  `yaml:"agent"`
		Server `yaml:"server"`
		Log    `yaml:"logger"`
	}
	Agent struct {
		Name             string        `yaml:"name"`
		Version          string        `yaml:"version"`
		PollInterval     time.Duration `yaml:"pollInterval" env:"POLL_INTERVAL"`
		ReportInterval   time.Duration `yaml:"reportInterval" env:"REPORT_INTERVAL"`
		ServerURL        string        `yaml:"server_url" env:"ADDRESS"`
		MetricFieldNames []string      `yaml:"metric_field_names"`
	}
	Server struct {
		Name    string `yaml:"name"`
		Version string `yaml:"version"`
		Address string `yaml:"address" env:"ADDRESS"`
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
	if err = cleanenv.ReadEnv(cfg); err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	return cfg, nil
}

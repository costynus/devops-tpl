package agent

import (
	"flag"
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		Agent `yaml:"agent"`
		Log   `yaml:"logger"`
	}
	Agent struct {
		Name             string        `yaml:"name"`
		Version          string        `yaml:"version"`
		PollInterval     time.Duration `yaml:"pollInterval" env:"POLL_INTERVAL"`
		ReportInterval   time.Duration `yaml:"reportInterval" env:"REPORT_INTERVAL"`
		ServerURL        string        `yaml:"server_url" env:"ADDRESS"`
		ServerSchema     string        `yaml:"server_schema" env:"SERVER_SCHEMA"`
		MetricFieldNames []string      `yaml:"metric_field_names"`
	}
	Log struct {
		Level string `yaml:"log_level"`
	}
)

func NewConfig() (*Config, error) {
	cfg := &Config{}
	// Flag init -.
	flag.StringVar(&cfg.Agent.ServerURL, "a", cfg.Agent.ServerURL, "server address")
	flag.DurationVar(&cfg.Agent.ReportInterval, "r", cfg.Agent.ReportInterval, "report interval")
	flag.DurationVar(&cfg.Agent.PollInterval, "p", cfg.Agent.PollInterval, "poll interval")

	// YAML Config -.
	err := cleanenv.ReadConfig("./config/config.yml", cfg)
	if err != nil {
		return cfg, fmt.Errorf("config error: %w", err)
	}

	// Flags Config -.
	flag.Parse()

	// Env config -.
	if err = cleanenv.ReadEnv(cfg); err != nil {
		return cfg, fmt.Errorf("config error: %w", err)
	}

	return cfg, nil
}

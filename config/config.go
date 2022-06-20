package config

import (
	"flag"
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
		Name          string        `yaml:"name"`
		Version       string        `yaml:"version"`
		Address       string        `yaml:"address" env:"ADDRESS"`
		StoreInterval time.Duration `yaml:"store_interval" env:"STORE_INTERVAL"`
		StoreFile     string        `yaml:"store_file" env:"STORE_FILE"`
		Restore       bool          `yaml:"restore" env:"RESTORE"`
	}
	Log struct {
		Level string `yaml:"log_level"`
	}
)

func NewConfig() (*Config, error) {
	cfg := &Config{}

	// YAML Config -.
	err := cleanenv.ReadConfig("./config/config.yml", cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	// Flags Config -.
	flag.StringVar(&cfg.Server.Address, "a", cfg.Server.Address, "address to listen on")
	flag.BoolVar(&cfg.Server.Restore, "r", cfg.Server.Restore, "restore data from file")
	flag.DurationVar(&cfg.Server.StoreInterval, "i", cfg.Server.StoreInterval, "store interval")
	flag.StringVar(&cfg.Server.StoreFile, "f", cfg.Server.StoreFile, "store file")

	flag.StringVar(&cfg.Agent.ServerURL, "a", cfg.Agent.ServerURL, "server address")
	flag.DurationVar(&cfg.Agent.ReportInterval, "r", cfg.Agent.ReportInterval, "report interval")
	flag.DurationVar(&cfg.Agent.PollInterval, "p", cfg.Agent.PollInterval, "poll interval")

	// Env config -.
	if err = cleanenv.ReadEnv(cfg); err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}
	return cfg, nil
}

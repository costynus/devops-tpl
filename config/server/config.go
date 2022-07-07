package server

import (
	"flag"
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		Server `yaml:"server"`
		Log    `yaml:"logger"`
	}
	Server struct {
		Name          string        `yaml:"name" env:"NAME"`
		Version       string        `yaml:"version"`
		Address       string        `yaml:"address" env:"ADDRESS"`
		StoreInterval time.Duration `yaml:"store_interval" env:"STORE_INTERVAL"`
		StoreFile     string        `yaml:"store_file" env:"STORE_FILE"`
		Restore       bool          `yaml:"restore" env:"RESTORE"`
		KEY           string        `env:"KEY"`
	}
	Log struct {
		Level string `yaml:"log_level"`
	}
)

func NewConfig() (*Config, error) {
	cfg := &Config{}

	// Flag init -.
	flag.StringVar(&cfg.Server.Address, "a", cfg.Server.Address, "address to listen on")
	flag.BoolVar(&cfg.Server.Restore, "r", cfg.Server.Restore, "restore data from file")
	flag.DurationVar(&cfg.Server.StoreInterval, "i", cfg.Server.StoreInterval, "store interval")
	flag.StringVar(&cfg.Server.StoreFile, "f", cfg.Server.StoreFile, "store file")
	flag.StringVar(&cfg.Server.KEY, "k", cfg.Server.KEY, "crypto key")

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

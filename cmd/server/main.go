package main

import (
	"devops-tpl/config"
	"devops-tpl/internal/app/server"
	"flag"
	"log"
)

func main() {
	cfg := config.NewConfig()
	// Flag init -.
	flag.StringVar(&cfg.Server.Address, "a", cfg.Server.Address, "address to listen on")
	flag.BoolVar(&cfg.Server.Restore, "r", cfg.Server.Restore, "restore data from file")
	flag.DurationVar(&cfg.Server.StoreInterval, "i", cfg.Server.StoreInterval, "store interval")
	flag.StringVar(&cfg.Server.StoreFile, "f", cfg.Server.StoreFile, "store file")
	err := config.Init(cfg)
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}
	server.Run(cfg)
}

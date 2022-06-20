package main

import (
	"devops-tpl/config"
	"devops-tpl/internal/app/agent"
	"flag"

	"log"
)

func main() {
	cfg := config.NewConfig()
	// Flag init -.
	flag.StringVar(&cfg.Agent.ServerURL, "a", cfg.Agent.ServerURL, "server address")
	flag.DurationVar(&cfg.Agent.ReportInterval, "r", cfg.Agent.ReportInterval, "report interval")
	flag.DurationVar(&cfg.Agent.PollInterval, "p", cfg.Agent.PollInterval, "poll interval")
	err := config.Init(cfg)
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}
	agent.Run(cfg)
}

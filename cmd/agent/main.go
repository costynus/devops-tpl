package main

import (
	"devops-tpl/config"
	"devops-tpl/internal/app/agent"

	"log"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}
	agent.Run(cfg)
}

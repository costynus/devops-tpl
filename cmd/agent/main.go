package main

import (
	agent_config "devops-tpl/config/agent"
	"devops-tpl/internal/app/agent"

	"log"
)

func main() {
	cfg, err := agent_config.NewConfig()

	if err != nil {
		log.Fatalf("Config error: %s", err)
	}
	agent.Run(cfg)
}

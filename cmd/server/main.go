package main

import (
	server_config "devops-tpl/config/server"
	"devops-tpl/internal/app/server"
	"log"
)

func main() {
	cfg, err := server_config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}
	server.Run(cfg)
}

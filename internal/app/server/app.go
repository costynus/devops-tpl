package server

import (
	"devops-tpl/config"
	"devops-tpl/internal/infrastructure/repo"
	"devops-tpl/pkg/logger"
	"net"
	"net/http"
)

func Run(cfg *config.Config) {
	// Logger
	l := logger.New(cfg.Log.Level)

	// Handlers
	http.HandleFunc("/healthz", HealthzHandler)
	http.HandleFunc(
		"/update/",
		UpdateMetricViewHandler(
			repo.New(),
		),
	)

	// HTTP Server -.
	l.Fatal(http.ListenAndServe(net.JoinHostPort("", cfg.Server.Port), nil))
}

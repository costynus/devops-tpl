package server

import (
	"devops-tpl/config"
	"devops-tpl/internal/infrastructure/repo"
	"devops-tpl/pkg/logger"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Run(cfg *config.Config) {
	// Logger -.
	l := logger.New(cfg.Log.Level)

	// HTTP Server -.
	handler := chi.NewRouter()
	NewRouter(
		handler,
		repo.New(),
	)
	l.Fatal(http.ListenAndServe(cfg.Address, handler))
}

package server

import (
	"context"
	"devops-tpl/config"
	"devops-tpl/internal/infrastructure/repo"
	"devops-tpl/pkg/logger"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Run(cfg *config.Config) {
	// Logger -.
	l := logger.New(cfg.Log.Level)

	// Repository
	repo := repo.New(cfg.Server.StoreFile, cfg.Server.Restore)

	// Worker -.
	worker := NewWorker(cfg.Server.StoreInterval, repo, l)
	go worker.StoreMetrics(context.Background())

	// HTTP Server -.
	handler := chi.NewRouter()
	NewRouter(
		handler,
		repo,
	)
	l.Fatal(http.ListenAndServe(cfg.Address, handler))
}

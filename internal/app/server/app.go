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
	repo := repo.New(cfg.Server.StoreFile)
	if cfg.Server.Restore {
		err := repo.UploadFromFile()
		if err != nil {
			l.Error("Error with upload from file: %w", err)
		} else {
			l.Info("Upload from file - success")
		}
	}

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

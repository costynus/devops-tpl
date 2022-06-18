package server

import (
	"devops-tpl/config"
	"devops-tpl/internal/infrastructure/repo"
	"devops-tpl/internal/usecase"
	"devops-tpl/pkg/logger"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Run(cfg *config.Config) {
	// Logger -.
	l := logger.New(cfg.Log.Level)

	// UseCase
	uc := usecase.New(
		repo.New(
			cfg.Server.StoreFile,
			cfg.Server.Restore,
		),
	)

	// Worker -.
	// worker := NewWorker(cfg.Server.StoreInterval, repo, l)
	// go worker.StoreMetrics(context.Background())

	// HTTP Server -.
	handler := chi.NewRouter()
	NewRouter(
		handler,
		uc,
		l,
	)

	// Start
	l.Fatal(http.ListenAndServe(cfg.Address, handler))
}

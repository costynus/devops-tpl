package server

import (
	server_config "devops-tpl/config/server"
	"devops-tpl/internal/infrastructure/repo"
	"devops-tpl/internal/usecase"
	"devops-tpl/pkg/logger"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Run(cfg *server_config.Config) {
	// Logger -.
	l := logger.New(cfg.Log.Level)

	// Repo Options
	repoOptions := make([]repo.Option, 0)
	if cfg.Server.StoreFile != " " {
		repoOptions = append(repoOptions, repo.StoreFilePath(cfg.Server.StoreFile))
	}
	if cfg.Server.Restore && cfg.Server.StoreFile != " " {
		repoOptions = append(repoOptions, repo.Restore())
	}

	// UseCase Options
	ucOptions := make([]usecase.Option, 0)
	if cfg.Server.StoreInterval == 0 {
		ucOptions = append(ucOptions, usecase.SynchWriteFile())
	} else {
		ucOptions = append(ucOptions, usecase.WriteFileDuration(cfg.Server.StoreInterval))
	}

	// UseCase
	uc := usecase.New(
		repo.New(
			repoOptions...,
		),
		l,
		ucOptions...,
	)

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

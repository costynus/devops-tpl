package server

import (
	server_config "devops-tpl/config/server"
	"devops-tpl/internal/infrastructure/repo"
	"devops-tpl/internal/usecase"
	"devops-tpl/pkg/logger"
	"devops-tpl/pkg/postgres"
	"fmt"
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
		ucOptions = append(ucOptions, usecase.SyncWriteFile())
	} else {
		ucOptions = append(ucOptions, usecase.WriteFileDuration(cfg.Server.StoreInterval))
	}
	if cfg.Server.KEY != "" {
		ucOptions = append(ucOptions, usecase.CheckSign(cfg.Server.KEY))
	}

	var currRepo usecase.MetricRepo
	switch cfg.PG.URL {
	case "":
		currRepo = repo.New(repoOptions...)
	default:
		err := migration(cfg.PG.URL, cfg.PG.MigDir)
		if err != nil {
			l.Fatal(err)
		}
		pg, err := postgres.New(cfg.PG.URL)
		if err != nil {
			l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
		}
		defer pg.Close()

		currRepo = repo.NewPG(pg)
	}

	// UseCase
	uc := usecase.New(
		currRepo,
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

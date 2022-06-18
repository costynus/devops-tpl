package server

import (
	"context"
	"devops-tpl/internal/usecase"
	"devops-tpl/pkg/logger"
	"time"
)

type Worker struct {
	StoreInterval time.Duration
	repo          usecase.MetricRepo
	l             logger.Interface
}

func NewWorker(StoreInterval time.Duration, Repo usecase.MetricRepo, l logger.Interface) *Worker {
	return &Worker{
		StoreInterval: StoreInterval,
		repo:          Repo,
		l:             l,
	}
}

func (w Worker) StoreMetrics(ctx context.Context) {
	if w.StoreInterval == 0 {
		return
	}
	ticker := time.NewTicker(w.StoreInterval)
	for {
		<-ticker.C
		err := w.repo.StoreToFile(ctx)
		if err != nil {
			w.l.Error("Error while writing to file: %w", err)
		} else {
			w.l.Info("store metrics success")
		}

	}
}

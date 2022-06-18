package server

import (
	"context"
	"devops-tpl/pkg/logger"
	"time"
)

type Worker struct {
	StoreInterval time.Duration
	repo          MetricRepo
	l             logger.Interface
}

func NewWorker(StoreInterval time.Duration, Repo MetricRepo, l logger.Interface) *Worker {
	return &Worker{
		StoreInterval: StoreInterval,
		repo:          Repo,
		l:             l,
	}
}

func (w Worker) StoreMetrics(ctx context.Context) {
	ticker := time.NewTicker(w.StoreInterval)
	for {
		<-ticker.C
		err := w.repo.StoreToFile()
		if err != nil {
			w.l.Error("Error while writing to file: %w", err)
		} else {
			w.l.Info("store metrics success")
		}

	}
}

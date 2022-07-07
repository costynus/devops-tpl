package usecase

import (
	"context"
	"devops-tpl/internal/entity"
	"devops-tpl/internal/infrastructure/repo"
	"devops-tpl/pkg/logger"
	"errors"
	"fmt"
	"time"
)

const (
	Gauge   = "gauge"
	Counter = "counter"
)

type DevOpsUseCase struct {
	repo MetricRepo
	l    logger.Interface

	checkSign               bool
	cryptoKey               string
	writeFileDuration       time.Duration
	writeToFileWithDuration bool
	asynchWriteFile         bool
	synchWriteFile          bool
	C                       chan struct{}
}

func New(repo MetricRepo, l logger.Interface, opts ...Option) *DevOpsUseCase {
	uc := &DevOpsUseCase{
		repo: repo,
		l:    l,
	}

	// Set Options
	for _, opt := range opts {
		opt(uc)
	}

	if uc.writeToFileWithDuration {
		go func() {
			ticker := time.NewTicker(uc.writeFileDuration)
			for {
				<-ticker.C
				uc.C <- struct{}{}
			}
		}()
	}
	if uc.writeToFileWithDuration || uc.asynchWriteFile {
		uc.C = make(chan struct{}, 1)
		go uc.saveStorage()
	}

	return uc
}

func (uc *DevOpsUseCase) saveStorage() {
	for {
		<-uc.C
		// add WG
		err := uc.repo.StoreAll()
		if err != nil {
			uc.l.Error("error while writing to file: %w", err)
		} else {
			uc.l.Info("store metric success")
		}
	}
}

func (uc DevOpsUseCase) GetMetricNames(ctx context.Context) ([]string, error) {
	names := uc.repo.GetMetricNames(ctx)
	return names, nil
}

func (uc *DevOpsUseCase) StoreMetric(ctx context.Context, metric entity.Metric) error {
	if uc.checkSign {
		if !metric.CheckSign(uc.cryptoKey) {
			return ErrSignNotEqual
		}
	}
	switch metric.MType {
	case Gauge:
		if err := uc.repo.StoreMetric(ctx, metric); err != nil {
			if errors.Is(err, repo.ErrNotFound) {
				return ErrNotFound
			}
			return fmt.Errorf("DevOpsUseCase - StoreMetric: %w", err)
		}
	case Counter:
		oldMetric, err := uc.repo.GetMetric(ctx, metric.ID)
		if err != nil {
			if !errors.Is(err, repo.ErrNotFound) {
				return fmt.Errorf("DevOpsUseCase - GetMetric: %w", err)
			}
		} else {
			delta := *metric.Delta + *oldMetric.Delta
			metric.Delta = &delta
		}
		if err := uc.repo.StoreMetric(ctx, metric); err != nil {
			return fmt.Errorf("DevOpsUseCase - StoreMetric: %w", err)
		}
	default:
		return ErrNotImplemented
	}
	if uc.asynchWriteFile {
		uc.C <- struct{}{}
	}
	if uc.synchWriteFile {
		err := uc.repo.StoreAll()
		if err != nil {
			return fmt.Errorf("DevOpsUseCase - StoreMetric - uc.repo.StoreAll: %w", err)
		}
	}
	return nil
}

func (uc *DevOpsUseCase) GetMetric(ctx context.Context, metric entity.Metric) (entity.Metric, error) {
	metric, err := uc.repo.GetMetric(ctx, metric.ID)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return metric, ErrNotFound
		}
		return metric, fmt.Errorf("DevOpsUseCase - GetMetric: %w", err)
	}
	metric.Sign(uc.cryptoKey)
	return metric, nil
}

func (uc *DevOpsUseCase) PingRepo(ctx context.Context) error {
	return uc.repo.Ping(ctx)
}

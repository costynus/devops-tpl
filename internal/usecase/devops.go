package usecase

import (
	"context"
	"devops-tpl/internal/entity"
	"devops-tpl/internal/infrastructure/repo"
	"errors"
	"fmt"
)

const (
	Gauge   = "gauge"
	Counter = "counter"
)

type DevOpsUseCase struct {
	repo MetricRepo
}

func New(repo MetricRepo) *DevOpsUseCase {
	return &DevOpsUseCase{
		repo: repo,
	}
}

func (uc DevOpsUseCase) MetricNames(ctx context.Context) ([]string, error) {
	names := uc.repo.GetMetricNames(ctx)
	return names, nil
}

func (uc *DevOpsUseCase) StoreMetric(ctx context.Context, metric entity.Metric) error {
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
	return nil
}

func (uc *DevOpsUseCase) Metric(ctx context.Context, metric entity.Metric) (entity.Metric, error) {
	metric, err := uc.repo.GetMetric(ctx, metric.ID)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return metric, ErrNotFound
		}
		return metric, fmt.Errorf("DevOpsUseCase - Metric: %w", err)
	}
	return metric, nil
}

func (uc *DevOpsUseCase) SaveStorage(ctx context.Context) error {
	return nil
}

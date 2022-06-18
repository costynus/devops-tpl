package usecase

import (
	"context"
	"devops-tpl/internal/entity"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks_test.go -package=usecase_test
type (
	DevOps interface {
		MetricNames(context.Context) ([]string, error)
		StoreMetric(context.Context, entity.Metric) error
		Metric(context.Context, entity.Metric) (entity.Metric, error)
		SaveStorage(context.Context) error
	}

	MetricRepo interface {
		GetMetricNames(context.Context) []string
		StoreMetric(context.Context, entity.Metric) error
		GetMetric(context.Context, string) (entity.Metric, error)

		StoreToFile(context.Context) error
		UploadFromFile(context.Context) error
	}
)

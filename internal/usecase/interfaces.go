package usecase

import (
	"context"
	"devops-tpl/internal/entity"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks_test.go -package=usecase_test
type (
	DevOps interface {
		GetMetricNames(context.Context) ([]string, error)
		GetMetric(context.Context, entity.Metric) (entity.Metric, error)
		StoreMetric(context.Context, entity.Metric) error
		PingRepo(context.Context) error
	}

	MetricRepo interface {
		GetMetricNames(context.Context) []string
		GetMetric(context.Context, string) (entity.Metric, error)
		StoreMetric(context.Context, entity.Metric) error

		StoreAll() error
		Upload(context.Context) error

		Ping(context.Context) error
	}
)

package repo

import (
	"context"
	"devops-tpl/internal/entity"
	"devops-tpl/pkg/postgres"
)

type MetricPGRepo struct {
	*postgres.Postgres
}

func NewPG(pg *postgres.Postgres) *MetricPGRepo {
	return &MetricPGRepo{pg}
}

func (r *MetricPGRepo) GetMetricNames(ctx context.Context) []string {
	return nil
}

func (r *MetricPGRepo) GetMetric(ctx context.Context, name string) (entity.Metric, error) {
	return entity.Metric{}, nil
}

func (r *MetricPGRepo) StoreMetric(ctx context.Context, metric entity.Metric) error {
	return nil
}

func (r *MetricPGRepo) StoreAll() error {
	return nil
}

func (r *MetricPGRepo) Upload(context.Context) error {
	return nil
}

func (r *MetricPGRepo) Ping(ctx context.Context) error {
	return r.Pool.Ping(ctx)
}

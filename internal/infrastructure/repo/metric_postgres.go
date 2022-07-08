package repo

import (
	"context"
	"devops-tpl/internal/entity"
	"devops-tpl/pkg/postgres"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
)

type MetricPGRepo struct {
	*postgres.Postgres
}

func NewPG(pg *postgres.Postgres) *MetricPGRepo {
	return &MetricPGRepo{pg}
}

func (r *MetricPGRepo) GetMetricNames(ctx context.Context) []string {
	dst := make([]string, 0)

	pgxscan.Select(ctx, r.Pool, &dst, "select name from public.metric;")
	return dst
}

func (r *MetricPGRepo) GetMetric(ctx context.Context, name string) (entity.Metric, error) {
	sql, args, err := r.Builder.
		Select("name", "mtype", "delta", "value", "hash").
		From("public.metric").
		Where(sq.Eq{"name": name}).
		ToSql()

	if err != nil {
		return entity.Metric{}, fmt.Errorf("MetricPGRepo - StoreMetric - r.Builder: %w", err)
	}

	dst := make([]entity.Metric, 0)
	if err = pgxscan.Select(ctx, r.Pool, &dst, sql, args...); err != nil {
		return entity.Metric{}, fmt.Errorf("MetricPGRepo - StoreMetric - pgxscan.Select: %w", err)
	}

	if len(dst) == 0 {
		return entity.Metric{}, ErrNotFound
	}

	return dst[0], nil
}

func (r *MetricPGRepo) StoreMetric(ctx context.Context, metric entity.Metric) error {
	updateSQL, updateArgs, err := r.Builder.
		Update("public.metric").
		Set("delta", metric.Delta).
		Set("value", metric.Value).
		Where(sq.Eq{"name": metric.ID}).
		ToSql()

	if err != nil {
		return fmt.Errorf("MetricPGRepo - StoreMetric - r.Builder: %w", err)
	}

	insertSQL, insertArgs, err := r.Builder.
		Insert("public.metric").
		Columns("name", "mtype", "delta", "value", "hash").
		Values(metric.ID, metric.MType, metric.Delta, metric.Value, metric.Hash).
		ToSql()
	if err != nil {
		return fmt.Errorf("MetricPGRepo - StoreMetric - r.Builder: %w", err)
	}

	_, err = r.Pool.Exec(ctx, insertSQL, insertArgs...)
	if err != nil {
		_, err = r.Pool.Exec(ctx, updateSQL, updateArgs...)
		if err != nil {
			return fmt.Errorf("MetricPGRepo - StoreMetric - r.Pool.Exec: %w", err)
		}
	}

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

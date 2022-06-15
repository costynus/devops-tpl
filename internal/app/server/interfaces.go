package server

import (
	"devops-tpl/internal/entity"
)

type (
	MetricRepo interface {
		StoreGauge(string, entity.Gauge) error
		AddCounter(string, entity.Counter) error
		GetMetric(string) (interface{}, error)
		GetMetricNames() []string
	}
)

package server

import (
	"devops-tpl/internal/entity"
)

type (
	MetricRepo interface {
		StoreGauge(string, entity.Gauge) error
		StoreCounter(string, entity.Counter) error
	}
)

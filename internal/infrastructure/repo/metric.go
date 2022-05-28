package repo

import (
	"devops-tpl/internal/entity"
)

type MetricRepo struct {
	data map[string]interface{}
}

func New() *MetricRepo {
	metricRepo := MetricRepo{}
	metricRepo.data = make(map[string]interface{})
	return &metricRepo
}

func (r *MetricRepo) StoreGauge(name string, value entity.Gauge) error {
	r.data[name] = value
	return nil
}

func (r *MetricRepo) StoreCounter(name string, value entity.Counter) error {
	oldValue, ok := r.data[name]
	if ok {
		r.data[name] = value + oldValue.(entity.Counter)
	} else {
		r.data[name] = value
	}
	return nil
}

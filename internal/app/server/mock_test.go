package server

import "devops-tpl/internal/entity"

type MockMetricRepo struct {
	GetMetrics interface{}
	Err        error
}

func (m *MockMetricRepo) StoreGauge(string, entity.Gauge) error   { return m.Err }
func (m *MockMetricRepo) AddCounter(string, entity.Counter) error { return m.Err }
func (m *MockMetricRepo) GetMetric(string) (interface{}, error)   { return m.GetMetrics, m.Err }
func (m *MockMetricRepo) GetMetricNames() []string                { return nil }

package server

import "devops-tpl/internal/entity"

type MockMetricRepo struct {
	GetMetrics entity.Metrics
	Err        error
}

func (m *MockMetricRepo) StoreGauge(string, entity.Gauge) error    { return m.Err }
func (m *MockMetricRepo) AddCounter(string, entity.Counter) error  { return m.Err }
func (m *MockMetricRepo) GetMetric(string) (entity.Metrics, error) { return m.GetMetrics, m.Err }
func (m *MockMetricRepo) GetMetricNames() []string                 { return nil }
func (m *MockMetricRepo) StoreToFile() error                       { return nil }
func (m *MockMetricRepo) UploadFromFile() error                    { return nil }

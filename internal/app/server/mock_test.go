package server

import "devops-tpl/internal/entity"

type MockMetricRepo struct{}

func (m *MockMetricRepo) StoreGauge(string, entity.Gauge) error     { return nil }
func (m *MockMetricRepo) StoreCounter(string, entity.Counter) error { return nil }

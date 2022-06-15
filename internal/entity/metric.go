package entity

import (
	"fmt"
	"strconv"
)

type (
	Gauge   float64
	Counter int64

	Metric struct {
		Name  string
		Value interface{}
	}
	Metrics struct {
		ID    string   `json:"id"`              // имя метрики
		MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
		Delta *Counter `json:"delta,omitempty"` // значение метрики в случае передачи counter
		Value *Gauge   `json:"value,omitempty"` // значение метрики в случае передачи gauge
	}
)

func ParseGauge(value string) (Gauge, error) {
	s, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, fmt.Errorf("Gauge - ParseGauge - strconv.ParseFloat: %w", err)
	}

	return Gauge(s), nil
}

func ParseCounter(value string) (Counter, error) {
	s, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("Gauge - ParseGauge - strconv.Atoi: %w", err)
	}
	return Counter(s), nil
}

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

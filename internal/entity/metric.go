package entity

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"strconv"
)

type (
	Gauge   float64
	Counter int64
	Metric  struct {
		ID    string   `json:"id" db:"name"`    // имя метрики
		MType string   `json:"type" db:"mtype"` // параметр, принимающий значение gauge или counter
		Delta *Counter `json:"delta,omitempty"` // значение метрики в случае передачи counter
		Value *Gauge   `json:"value,omitempty"` // значение метрики в случае передачи gauge
		Hash  string   `json:"hash,omitempty"`  // значение хеш-функции
	}
)

func (m Metric) hash(key string) string {
	var msg string
	switch m.MType {
	case "counter":
		msg = fmt.Sprintf("%s:%s:%d", m.ID, m.MType, *m.Delta)
	case "gauge":
		msg = fmt.Sprintf("%s:%s:%f", m.ID, m.MType, *m.Value)
	}
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(msg))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (m *Metric) Sign(key string) {
	if key == "" {
		return
	}
	m.Hash = m.hash(key)
}

func (m Metric) CheckSign(key string) bool {
	return m.hash(key) == m.Hash
}

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

func (g Gauge) TypeString() string {
	return "gauge"
}

func (c Counter) TypeString() string {
	return "counter"
}

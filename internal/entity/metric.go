package entity

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
)

type (
	Gauge   float64
	Counter int64
	Metric  struct {
		ID    string   `json:"id"`              // имя метрики
		MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
		Delta *Counter `json:"delta,omitempty"` // значение метрики в случае передачи counter
		Value *Gauge   `json:"value,omitempty"` // значение метрики в случае передачи gauge
		Hash  string   `json:"hash,omitempty"`  // значение хеш-функции
	}
)

func (m Metric) hash(key string) string {
	var src string
	switch m.MType {
	case "gauge":
		if *m.Value == Gauge(int(*m.Value)) {
			src = fmt.Sprintf("%s:gauge:%d", m.ID, int(*m.Value))
		} else {
			src = fmt.Sprintf("%s:gauge:%f", m.ID, *m.Value)
		}
	case "counter":
		src = fmt.Sprintf("%s:counter:%d", m.ID, *m.Delta)
	}
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(src))
	dst := h.Sum(nil)
	return hex.EncodeToString(dst)
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

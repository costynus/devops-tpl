package agent

import (
	"math/rand"
	"runtime"
	"sync"
)

type (
	gauge   float64
	counter int64

	Metrics struct {
		PollCount   counter
		RandomValue gauge
		Mutex       *sync.Mutex
		*collector
	}

	collector struct {
		Alloc         gauge
		BuckHashSys   gauge
		Frees         gauge
		GCCPUFraction gauge
		GCSys         gauge
		HeapAlloc     gauge
		HeapIdle      gauge
		HeapInuse     gauge
		HeapObjects   gauge
		HeapReleased  gauge
		HeapSys       gauge
		LastGC        gauge
		Lookups       gauge
		MCacheInuse   gauge
		MCacheSys     gauge
		MSpanInuse    gauge
		MSpanSys      gauge
		Mallocs       gauge
		NextGC        gauge
		NumForcedGC   gauge
		NumGC         gauge
		OtherSys      gauge
		PauseTotalNs  gauge
		StackInuse    gauge
		StackSys      gauge
		Sys           gauge
		TotalAlloc    gauge
	}
)

func NewMetrics() *Metrics {
	return &Metrics{
		Mutex:     &sync.Mutex{},
		collector: &collector{},
	}
}

func (m *Metrics) UpdateMetrics() {
	memStats := &runtime.MemStats{}
	runtime.ReadMemStats(memStats)

	m.Mutex.Lock()
	m.collectMetrics(memStats)
	m.Mutex.Unlock()

	m.PollCount += 1
	m.RandomValue = gauge(rand.Float64())
}

func (m *Metrics) collectMetrics(memStats *runtime.MemStats) {
	m.collector.Alloc = gauge(memStats.Alloc)
	m.collector.BuckHashSys = gauge(memStats.BuckHashSys)
	m.collector.Frees = gauge(memStats.Frees)
	m.collector.GCCPUFraction = gauge(memStats.GCCPUFraction)
	m.collector.GCSys = gauge(memStats.GCSys)
	m.collector.HeapAlloc = gauge(memStats.HeapAlloc)
	m.collector.HeapIdle = gauge(memStats.HeapIdle)
	m.collector.HeapInuse = gauge(memStats.HeapInuse)
	m.collector.HeapObjects = gauge(memStats.HeapObjects)
	m.collector.HeapReleased = gauge(memStats.HeapReleased)
	m.collector.HeapSys = gauge(memStats.HeapSys)
	m.collector.LastGC = gauge(memStats.LastGC)
	m.collector.Lookups = gauge(memStats.Lookups)
	m.collector.MCacheInuse = gauge(memStats.MCacheInuse)
	m.collector.MCacheSys = gauge(memStats.MCacheSys)
	m.collector.MSpanInuse = gauge(memStats.MSpanInuse)
	m.collector.MSpanSys = gauge(memStats.MSpanSys)
	m.collector.Mallocs = gauge(memStats.Mallocs)
	m.collector.NextGC = gauge(memStats.NextGC)
	m.collector.NumForcedGC = gauge(memStats.NumForcedGC)
	m.collector.NumGC = gauge(memStats.NumGC)
	m.collector.OtherSys = gauge(memStats.OtherSys)
	m.collector.PauseTotalNs = gauge(memStats.PauseTotalNs)
	m.collector.StackInuse = gauge(memStats.StackInuse)
	m.collector.StackSys = gauge(memStats.StackSys)
	m.collector.Sys = gauge(memStats.Sys)
	m.collector.TotalAlloc = gauge(memStats.TotalAlloc)
}

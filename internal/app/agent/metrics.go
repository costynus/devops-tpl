package agent

import (
	"devops-tpl/internal/entity"
	"math/rand"
	"runtime"
	"sync"
)

type (
	Metrics struct {
		PollCount   entity.Counter
		RandomValue entity.Gauge
		Mutex       *sync.Mutex
		*collector
	}

	collector struct {
		Alloc         entity.Gauge
		BuckHashSys   entity.Gauge
		Frees         entity.Gauge
		GCCPUFraction entity.Gauge
		GCSys         entity.Gauge
		HeapAlloc     entity.Gauge
		HeapIdle      entity.Gauge
		HeapInuse     entity.Gauge
		HeapObjects   entity.Gauge
		HeapReleased  entity.Gauge
		HeapSys       entity.Gauge
		LastGC        entity.Gauge
		Lookups       entity.Gauge
		MCacheInuse   entity.Gauge
		MCacheSys     entity.Gauge
		MSpanInuse    entity.Gauge
		MSpanSys      entity.Gauge
		Mallocs       entity.Gauge
		NextGC        entity.Gauge
		NumForcedGC   entity.Gauge
		NumGC         entity.Gauge
		OtherSys      entity.Gauge
		PauseTotalNs  entity.Gauge
		StackInuse    entity.Gauge
		StackSys      entity.Gauge
		Sys           entity.Gauge
		TotalAlloc    entity.Gauge
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
	m.RandomValue = entity.Gauge(rand.Float64())
}

func (m *Metrics) collectMetrics(memStats *runtime.MemStats) {
	m.collector.Alloc = entity.Gauge(memStats.Alloc)
	m.collector.BuckHashSys = entity.Gauge(memStats.BuckHashSys)
	m.collector.Frees = entity.Gauge(memStats.Frees)
	m.collector.GCCPUFraction = entity.Gauge(memStats.GCCPUFraction)
	m.collector.GCSys = entity.Gauge(memStats.GCSys)
	m.collector.HeapAlloc = entity.Gauge(memStats.HeapAlloc)
	m.collector.HeapIdle = entity.Gauge(memStats.HeapIdle)
	m.collector.HeapInuse = entity.Gauge(memStats.HeapInuse)
	m.collector.HeapObjects = entity.Gauge(memStats.HeapObjects)
	m.collector.HeapReleased = entity.Gauge(memStats.HeapReleased)
	m.collector.HeapSys = entity.Gauge(memStats.HeapSys)
	m.collector.LastGC = entity.Gauge(memStats.LastGC)
	m.collector.Lookups = entity.Gauge(memStats.Lookups)
	m.collector.MCacheInuse = entity.Gauge(memStats.MCacheInuse)
	m.collector.MCacheSys = entity.Gauge(memStats.MCacheSys)
	m.collector.MSpanInuse = entity.Gauge(memStats.MSpanInuse)
	m.collector.MSpanSys = entity.Gauge(memStats.MSpanSys)
	m.collector.Mallocs = entity.Gauge(memStats.Mallocs)
	m.collector.NextGC = entity.Gauge(memStats.NextGC)
	m.collector.NumForcedGC = entity.Gauge(memStats.NumForcedGC)
	m.collector.NumGC = entity.Gauge(memStats.NumGC)
	m.collector.OtherSys = entity.Gauge(memStats.OtherSys)
	m.collector.PauseTotalNs = entity.Gauge(memStats.PauseTotalNs)
	m.collector.StackInuse = entity.Gauge(memStats.StackInuse)
	m.collector.StackSys = entity.Gauge(memStats.StackSys)
	m.collector.Sys = entity.Gauge(memStats.Sys)
	m.collector.TotalAlloc = entity.Gauge(memStats.TotalAlloc)
}

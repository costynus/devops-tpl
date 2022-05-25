package agent

import (
	"context"
	"devops-tpl/pkg/logger"
	"time"
)

type Worker struct {
	l       logger.Interface
	metrics *Metrics
	webAPI  *WebAPI
}

func NewWorker(l logger.Interface, metrics *Metrics, webAPI *WebAPI) *Worker {
	return &Worker{
		l:       l,
		metrics: metrics,
		webAPI:  webAPI,
	}
}

func (w *Worker) UpdateMetrics(ctx context.Context, ticker *time.Ticker) {
	for {
		<-ticker.C
		w.metrics.UpdateMetrics()
		w.l.Info("metrics updated")
	}
}

func (w *Worker) SendMetrics(ctx context.Context, ticker *time.Ticker) {
	for {
		<-ticker.C

		err := w.webAPI.SendGaugeMetric("Alloc", w.metrics.collector.Alloc)
		if err != nil {
			w.l.Error(err)
		}

		err = w.webAPI.SendGaugeMetric("BuckHashSys", w.metrics.collector.BuckHashSys)
		if err != nil {
			w.l.Error(err)
		}

		err = w.webAPI.SendGaugeMetric("Frees", w.metrics.collector.Frees)
		if err != nil {
			w.l.Error(err)
		}

		err = w.webAPI.SendGaugeMetric("GCCPUFraction", w.metrics.collector.GCCPUFraction)
		if err != nil {
			w.l.Error(err)
		}

		err = w.webAPI.SendGaugeMetric("GCSys", w.metrics.collector.GCSys)
		if err != nil {
			w.l.Error(err)
		}

		err = w.webAPI.SendGaugeMetric("HeapAlloc", w.metrics.collector.HeapAlloc)
		if err != nil {
			w.l.Error(err)
		}

		err = w.webAPI.SendGaugeMetric("HeapIdle", w.metrics.collector.HeapIdle)
		if err != nil {
			w.l.Error(err)
		}

		err = w.webAPI.SendGaugeMetric("HeapInuse", w.metrics.collector.HeapInuse)
		if err != nil {
			w.l.Error(err)
		}

		err = w.webAPI.SendGaugeMetric("HeapObjects", w.metrics.collector.HeapObjects)
		if err != nil {
			w.l.Error(err)
		}

		err = w.webAPI.SendGaugeMetric("HeapReleased", w.metrics.collector.HeapReleased)
		if err != nil {
			w.l.Error(err)
		}

		err = w.webAPI.SendGaugeMetric("HeapSys", w.metrics.collector.HeapSys)
		if err != nil {
			w.l.Error(err)
		}

		err = w.webAPI.SendGaugeMetric("LastGC", w.metrics.collector.LastGC)
		if err != nil {
			w.l.Error(err)
		}

		err = w.webAPI.SendGaugeMetric("Lookups", w.metrics.collector.Lookups)
		if err != nil {
			w.l.Error(err)
		}

		err = w.webAPI.SendGaugeMetric("MCacheInuse", w.metrics.collector.MCacheInuse)
		if err != nil {
			w.l.Error(err)
		}

		err = w.webAPI.SendGaugeMetric("MCacheSys", w.metrics.collector.MCacheSys)
		if err != nil {
			w.l.Error(err)
		}

		err = w.webAPI.SendGaugeMetric("MSpanInuse", w.metrics.collector.MSpanInuse)
		if err != nil {
			w.l.Error(err)
		}

		err = w.webAPI.SendGaugeMetric("MSpanSys", w.metrics.collector.MSpanSys)
		if err != nil {
			w.l.Error(err)
		}

		err = w.webAPI.SendGaugeMetric("Mallocs", w.metrics.collector.Mallocs)
		if err != nil {
			w.l.Error(err)
		}

		err = w.webAPI.SendGaugeMetric("NextGC", w.metrics.collector.NextGC)
		if err != nil {
			w.l.Error(err)
		}

		err = w.webAPI.SendGaugeMetric("NumForcedGC", w.metrics.collector.NumForcedGC)
		if err != nil {
			w.l.Error(err)
		}

		err = w.webAPI.SendGaugeMetric("NumGC", w.metrics.collector.NumGC)
		if err != nil {
			w.l.Error(err)
		}

		err = w.webAPI.SendGaugeMetric("OtherSys", w.metrics.collector.OtherSys)
		if err != nil {
			w.l.Error(err)
		}

		err = w.webAPI.SendGaugeMetric("PauseTotalNs", w.metrics.collector.PauseTotalNs)
		if err != nil {
			w.l.Error(err)
		}

		err = w.webAPI.SendGaugeMetric("StackInuse", w.metrics.collector.StackInuse)
		if err != nil {
			w.l.Error(err)
		}

		err = w.webAPI.SendGaugeMetric("StackSys", w.metrics.collector.StackSys)
		if err != nil {
			w.l.Error(err)
		}

		err = w.webAPI.SendGaugeMetric("Sys", w.metrics.collector.Sys)
		if err != nil {
			w.l.Error(err)
		}

		err = w.webAPI.SendGaugeMetric("TotalAlloc", w.metrics.collector.TotalAlloc)
		if err != nil {
			w.l.Error(err)
		}

		err = w.webAPI.SendGaugeMetric("RandomValue", w.metrics.RandomValue)
		if err != nil {
			w.l.Error(err)
		}

		err = w.webAPI.SendCounterMetric("PollCount", w.metrics.PollCount)
		if err != nil {
			w.l.Error(err)
		}
	}
}

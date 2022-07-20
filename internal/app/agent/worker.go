package agent

import (
	"context"
	"devops-tpl/internal/entity"
	"devops-tpl/pkg/logger"
	"reflect"
	"strings"
	"time"
)

type Worker struct {
	webAPI           AgentWebAPI
	metrics          *Metrics
	metricFieldNames []string
	l                logger.Interface
}

func NewWorker(webAPI AgentWebAPI, metrics *Metrics, metricFieldNames []string, l logger.Interface) *Worker {
	return &Worker{
		webAPI:           webAPI,
		metrics:          metrics,
		metricFieldNames: metricFieldNames,
		l:                l,
	}
}

func (w *Worker) UpdateMetrics(ctx context.Context, ticker *time.Ticker) {
	for {
		<-ticker.C
		go w.metrics.CollectMetrics()
		go w.metrics.CollectAdditionalMetrics()
		w.l.Info("metrics updated")
	}
}

func (w *Worker) SendMetrics(ctx context.Context, ticker *time.Ticker) {
	for {
		<-ticker.C

		for _, metricFieldName := range w.metricFieldNames {
			r := reflect.ValueOf(w.metrics)
			f := reflect.Indirect(r).FieldByName(metricFieldName)
			if !f.IsValid() {
				w.l.Error("field `%s` is not valid", metricFieldName)
				continue
			}

			fType := strings.ToLower(f.Type().Name())
			var valueGauge *entity.Gauge
			var valueCounter *entity.Counter

			switch fType {
			case "gauge":
				value := entity.Gauge(f.Float())
				valueGauge = &value
			case "counter":
				value := entity.Counter(f.Int())
				valueCounter = &value
			default:
				w.l.Error("field `%s` has not valid type")
				continue
			}
			go func(metricName, metricType string, Value *entity.Gauge, Delta *entity.Counter) {
				err := w.webAPI.SendMetric(metricName, metricType, Value, Delta)
				if err != nil {
					w.l.Error(err)
				}
			}(metricFieldName, fType, valueGauge, valueCounter)
		}
	}
}

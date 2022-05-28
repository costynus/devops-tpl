package agent

import (
	"context"
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
		w.metrics.UpdateMetrics()
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
			go func(metricName, metricType string, metricValue interface{}) {
				err := w.webAPI.SendMetric(metricName, metricType, metricValue)
				if err != nil {
					w.l.Error(err)
				}
			}(metricFieldName, strings.ToLower(f.Type().Name()), f)
		}
	}
}

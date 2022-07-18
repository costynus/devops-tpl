package server

import (
	"context"
	"devops-tpl/internal/entity"
	"devops-tpl/internal/usecase"
	"devops-tpl/pkg/logger"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

const (
	Gauge   = "gauge"
	Counter = "counter"
)

func NewRouter(handler *chi.Mux, uc usecase.DevOps, l logger.Interface) {
	// Options
	handler.Use(middleware.RequestID)
	handler.Use(middleware.RealIP)
	handler.Use(middleware.Logger)
	handler.Use(middleware.Recoverer)
	handler.Use(gzipReadHandle)
	handler.Use(gzipWriteHandle)

	// checker
	handler.Get("/healthz", healthzHandler())

	handler.Get("/", getMetricNamesHandler(uc, l))

	// updater
	handler.Route("/update", func(r chi.Router) {
		r.Post("/", updateMetric(uc, l))
		r.Post("/gauge/{metricName}/{metricValue}", updateGaugeMetric(uc, l))
		r.Post("/counter/{metricName}/{metricValue}", updateCounterMetric(uc, l))
		r.Post("/{metricType}/{metricName}/{metricValue}", notImplemented())
	})

	// value
	handler.Route("/value", func(r chi.Router) {
		r.Post("/", getValue(uc, l))
		r.Get("/gauge/{metricName}", getGaugeValue(uc, l))
		r.Get("/counter/{metricName}", getCounterValue(uc, l))
		r.Get("/{metricType}/{metricName}/{metricValue}", notImplemented())
	})
}

func healthzHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) }
}

func notImplemented() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		http.Error(w, "not implemented", http.StatusNotImplemented)
	}
}

func getMetricNamesHandler(uc usecase.DevOps, l logger.Interface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		names, err := uc.GetMetricNames(context.TODO())
		if err != nil {
			errorHandler(w, err)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(strings.Join(names, "\n")))
	}
}

func updateMetric(uc usecase.DevOps, l logger.Interface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var metric entity.Metric

		if err := json.NewDecoder(r.Body).Decode(&metric); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		err := uc.StoreMetric(context.TODO(), metric)
		if err != nil {
			l.Error(err)
			errorHandler(w, err)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
	}
}

func updateGaugeMetric(uc usecase.DevOps, l logger.Interface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		value, err := entity.ParseGauge(chi.URLParam(r, "metricValue"))
		if err != nil {
			l.Error(err)
			http.Error(w, "parsing error", http.StatusBadRequest)
			return
		}

		metric := entity.Metric{
			ID:    chi.URLParam(r, "metricName"),
			MType: value.TypeString(),
			Value: &value,
		}

		err = uc.StoreMetric(context.TODO(), metric)
		if err != nil {
			l.Error(err)
			errorHandler(w, err)
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
	}
}

func updateCounterMetric(uc usecase.DevOps, l logger.Interface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		value, err := entity.ParseCounter(chi.URLParam(r, "metricValue"))
		if err != nil {
			l.Error(err)
			http.Error(w, "parsing error", http.StatusBadRequest)
			return
		}

		metric := entity.Metric{
			ID:    chi.URLParam(r, "metricName"),
			MType: value.TypeString(),
			Delta: &value,
		}

		err = uc.StoreMetric(context.TODO(), metric)
		if err != nil {
			l.Error(err)
			errorHandler(w, err)
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
	}
}

func getValue(uc usecase.DevOps, l logger.Interface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var metric entity.Metric

		if err := json.NewDecoder(r.Body).Decode(&metric); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		metric, err := uc.GetMetric(context.TODO(), metric)
		if err != nil {
			l.Error(err)
			errorHandler(w, err)
			return
		}

		jsonResp, err := json.Marshal(metric)
		if err != nil {
			l.Error(err)
			errorHandler(w, err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResp)
	}
}

func getGaugeValue(uc usecase.DevOps, l logger.Interface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metric := entity.Metric{
			ID:    chi.URLParam(r, "metricName"),
			MType: Gauge,
		}

		metric, err := uc.GetMetric(context.TODO(), metric)
		if err != nil {
			l.Error(err)
			errorHandler(w, err)
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(fmt.Sprintf("%g", *metric.Value)))
	}
}

func getCounterValue(uc usecase.DevOps, l logger.Interface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metric := entity.Metric{
			ID:    chi.URLParam(r, "metricName"),
			MType: Counter,
		}

		metric, err := uc.GetMetric(context.TODO(), metric)
		if err != nil {
			l.Error(err)
			errorHandler(w, err)
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(fmt.Sprintf("%d", *metric.Delta)))
	}
}

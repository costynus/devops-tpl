package server

import (
	"context"
	"devops-tpl/internal/entity"
	"devops-tpl/internal/usecase"
	"devops-tpl/pkg/logger"
	"encoding/json"
	"fmt"
	"io"
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

	// checker
	handler.Get("/healthz", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })

	handler.Get("/", func(w http.ResponseWriter, r *http.Request) {
		names, err := uc.MetricNames(context.Background())
		if err != nil {
			http.Error(w, "internal server problem", http.StatusInternalServerError)
			return
		}
		_, err = io.WriteString(w, strings.Join(names, "\n"))
		if err != nil {
			http.Error(w, "internal server problem", http.StatusInternalServerError)
			return
		}
	})

	// updater
	handler.Route("/update", func(r chi.Router) {
		r.Post(
			"/",
			func(w http.ResponseWriter, r *http.Request) {
				var metric entity.Metric

				if err := json.NewDecoder(r.Body).Decode(&metric); err != nil {
					http.Error(w, "bad request", http.StatusBadRequest)
					return
				}

				err := uc.StoreMetric(context.Background(), metric)
				if err != nil {
					l.Error(err)
					http.Error(w, "storage problem", http.StatusInternalServerError)
					return
				}
				w.WriteHeader(http.StatusOK)
			},
		)
		r.Post(
			"/gauge/{metricName}/{metricValue}",
			func(w http.ResponseWriter, r *http.Request) {
				value, err := entity.ParseGauge(chi.URLParam(r, "metricValue"))
				if err != nil {
					l.Error(err)
					http.Error(w, "parsing error", http.StatusBadRequest)
					return
				}

				metric := entity.Metric{
					ID:    chi.URLParam(r, "metricName"),
					MType: value.String(),
					Value: &value,
				}

				err = uc.StoreMetric(context.Background(), metric)
				if err != nil {
					l.Error(err)
					http.Error(w, "storage problem", http.StatusInternalServerError)
					return
				}

				w.WriteHeader(http.StatusOK)
			},
		)
		r.Post(
			"/counter/{metricName}/{metricValue}",
			func(w http.ResponseWriter, r *http.Request) {
				value, err := entity.ParseCounter(chi.URLParam(r, "metricValue"))
				if err != nil {
					l.Error(err)
					http.Error(w, "parsing error", http.StatusBadRequest)
					return
				}

				metric := entity.Metric{
					ID:    chi.URLParam(r, "metricName"),
					MType: value.String(),
					Delta: &value,
				}

				err = uc.StoreMetric(context.Background(), metric)
				if err != nil {
					l.Error(err)
					http.Error(w, "storage problem", http.StatusInternalServerError)
					return
				}

				w.WriteHeader(http.StatusOK)
			},
		)
	})

	// value
	handler.Route("/value", func(r chi.Router) {
		r.Post(
			"/",
			func(w http.ResponseWriter, r *http.Request) {
				var metric entity.Metric

				if err := json.NewDecoder(r.Body).Decode(&metric); err != nil {
					http.Error(w, "bad request", http.StatusBadRequest)
					return
				}

				metric, err := uc.Metric(context.Background(), metric)
				if err != nil {
					l.Error(err)
					http.Error(w, "metric not found", http.StatusNotFound)
					return
				}

				jsonResp, err := json.Marshal(metric)
				w.Header().Set("Content-Type", "application/json")
				w.Write(jsonResp)
			},
		)
		r.Get(
			"/gauge/{metricName}",
			func(w http.ResponseWriter, r *http.Request) {
				metric := entity.Metric{
					ID:    chi.URLParam(r, "metricName"),
					MType: Gauge,
				}

				metric, err := uc.Metric(context.Background(), metric)
				if err != nil {
					l.Error(err)
					http.Error(w, "metric not found", http.StatusNotFound)
					return
				}

				w.Write([]byte(fmt.Sprintf("%g", *metric.Value)))
			},
		)
		r.Get(
			"/counter/{metricName}",
			func(w http.ResponseWriter, r *http.Request) {
				metric := entity.Metric{
					ID:    chi.URLParam(r, "metricName"),
					MType: Counter,
				}

				metric, err := uc.Metric(context.Background(), metric)
				if err != nil {
					l.Error(err)
					http.Error(w, "metric not found", http.StatusNotFound)
					return
				}

				w.Write([]byte(fmt.Sprintf("%v", *metric.Delta)))
			},
		)
	})
}

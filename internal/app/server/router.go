package server

import (
	"devops-tpl/internal/entity"
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

func NewRouter(handler *chi.Mux, repo MetricRepo) {
	// Options
	handler.Use(middleware.RequestID)
	handler.Use(middleware.RealIP)
	handler.Use(middleware.Logger)
	handler.Use(middleware.Recoverer)

	// checker
	handler.Get("/healthz", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })

	handler.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := io.WriteString(w, strings.Join(repo.GetMetricNames(), "\n"))
		if err != nil {
			panic(err)
		}
	})

	// updater
	handler.Route("/update", func(r chi.Router) {
		r.Post(
			"/",
			func(w http.ResponseWriter, r *http.Request) {
				var metrics entity.Metrics

				if err := json.NewDecoder(r.Body).Decode(&metrics); err != nil {
					http.Error(w, "bad request", http.StatusBadRequest)
					return
				}

				switch metrics.MType {
				case Gauge:
					if err := repo.StoreGauge(metrics.ID, *metrics.Value); err != nil {
						http.Error(w, "storage problem", http.StatusInternalServerError)
					}
				case Counter:
					if err := repo.AddCounter(metrics.ID, *metrics.Delta); err != nil {
						http.Error(w, "storage problem", http.StatusInternalServerError)
					}
				}
				w.WriteHeader(http.StatusOK)
			},
		)
		r.Post(
			"/{metricType}/{metricName}/{metricValue}",
			func(w http.ResponseWriter, r *http.Request) {
				metricType := chi.URLParam(r, "metricType")
				metricName := chi.URLParam(r, "metricName")
				metricValue := chi.URLParam(r, "metricValue")

				switch metricType {
				case Gauge:
					value, err := entity.ParseGauge(metricValue)
					if err != nil {
						http.Error(w, "bad value type", http.StatusBadRequest)
						return
					}
					err = repo.StoreGauge(metricName, value)
					if err != nil {
						http.Error(w, "storage problem", http.StatusInternalServerError)
					}
				case Counter:
					value, err := entity.ParseCounter(metricValue)
					if err != nil {
						http.Error(w, "bad value type", http.StatusBadRequest)
						return
					}
					err = repo.AddCounter(metricName, value)
					if err != nil {
						http.Error(w, "storage problem", http.StatusInternalServerError)
					}
				default:
					http.Error(w, "metric type is not found", http.StatusNotImplemented)
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
				var metrics entity.Metrics

				if err := json.NewDecoder(r.Body).Decode(&metrics); err != nil {
					http.Error(w, "bad request", http.StatusBadRequest)
					return
				}

				repoMetric, err := repo.GetMetric(metrics.ID)
				if err != nil {
					http.Error(w, "metric not found", http.StatusNotFound)
					return
				}
				switch metrics.MType {
				case Gauge:
					metrics.Value = repoMetric.Value
				case Counter:
					metrics.Delta = repoMetric.Delta
				}

				jsonResp, err := json.Marshal(metrics)
				if err != nil {
					http.Error(w, "server error", http.StatusInternalServerError)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				w.Write(jsonResp)
			},
		)
		r.Get("/{metricType}/{metricName}", func(w http.ResponseWriter, r *http.Request) {
			metricType := chi.URLParam(r, "metricType")
			metricName := chi.URLParam(r, "metricName")

			metric, err := repo.GetMetric(metricName)
			if err != nil {
				http.Error(w, "metric not found", http.StatusNotFound)
				return
			}
			switch metricType {
			case Gauge:
				w.Write([]byte(fmt.Sprintf("%g", *metric.Value)))
			case Counter:
				w.Write([]byte(fmt.Sprintf("%d", *metric.Delta)))
			default:
				http.Error(w, "metric type is not found", http.StatusNotImplemented)
			}
		})
	})
}

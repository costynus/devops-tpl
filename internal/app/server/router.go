package server

import (
	"devops-tpl/internal/entity"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func NewRouter(handler *chi.Mux, repo MetricRepo) {
	// Options
	handler.Use(middleware.RequestID)
	handler.Use(middleware.RealIP)
	handler.Use(middleware.Logger)
	handler.Use(middleware.Recoverer)

	// checker
	handler.Get("/healthz", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })

	// updater
	handler.Route("/update", func(r chi.Router) {
		r.Post(
			"/{metricType}/{metricName}/{metricValue}",
			func(w http.ResponseWriter, r *http.Request) {
				metricType := chi.URLParam(r, "metricType")
				metricName := chi.URLParam(r, "metricName")
				metricValue := chi.URLParam(r, "metricValue")

				switch metricType {
				case "gauge":
					value, err := entity.ParseGauge(metricValue)
					if err != nil {
						http.Error(w, "bad value type", http.StatusBadRequest)
						return
					}
					err = repo.StoreGauge(metricName, value)
					if err != nil {
						http.Error(w, "storage problem", http.StatusInternalServerError)
					}
				case "counter":
					value, err := entity.ParseCounter(metricValue)
					if err != nil {
						http.Error(w, "bad value type", http.StatusBadRequest)
						return
					}
					err = repo.StoreCounter(metricName, value)
					if err != nil {
						http.Error(w, "storage problem", http.StatusInternalServerError)
					}
				default:
					http.Error(w, "Metric type is not found", http.StatusNotImplemented)
				}
				w.WriteHeader(http.StatusOK)
			},
		)
	})

	// value
	handler.Route("/value", func(r chi.Router) {
		r.Get("/{metricType}/{metricName}", func(w http.ResponseWriter, r *http.Request) {
			metricType := chi.URLParam(r, "metricType")
			metricName := chi.URLParam(r, "metricName")

			switch metricType {
			case "gauge":
				value, err := repo.GetMetric(metricName)
				if err != nil {
					http.Error(w, "Metric not found", http.StatusNotFound)
					return
				}
				w.Write([]byte(fmt.Sprintf("%f", value)))
			case "counter":
				value, err := repo.GetMetric(metricName)
				if err != nil {
					http.Error(w, "Metric not found", http.StatusNotFound)
					return
				}
				w.Write([]byte(fmt.Sprintf("%d", value)))
			default:
				http.Error(w, "Metric type is not found", http.StatusNotImplemented)
			}
		})
	})
}

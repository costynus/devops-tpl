package server

import (
	"devops-tpl/internal/entity"
	"devops-tpl/pkg/logger"
	"net/http"
	"strings"
)

func HealthzHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func UpdateMetricViewHandler(repo MetricRepo, l logger.Interface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
			return
		} else if r.Header.Get("Content-Type") != "text/plain" {
			http.Error(w, "Only `text/plain` content-type is allowed!", http.StatusUnsupportedMediaType)
			return
		}

		params := strings.Split(
			strings.TrimPrefix(r.URL.Path, "/update/"),
			"/",
		)
		if len(params) != 3 {
			http.Error(w, "Pattern is type/name/value", http.StatusBadRequest)
			return
		}

		switch params[0] {
		case "gauge":
			value, err := entity.ParseGauge(params[2])
			if err != nil {
				l.Error(err, "http - /update - gauge")
				http.Error(w, "value parser problem", http.StatusInternalServerError)
				return
			}
			err = repo.StoreGauge(params[1], value)
			if err != nil {
				l.Error(err, "http - /update - gauge")
				http.Error(w, "storage problem", http.StatusInternalServerError)
			}
		case "counter":
			value, err := entity.ParseCounter(params[2])
			if err != nil {
				l.Error(err, "http - /update - counter")
				http.Error(w, "value parser problem", http.StatusInternalServerError)
				return
			}
			err = repo.StoreCounter(params[1], value)
			if err != nil {
				l.Error(err, "http - /update - counter")
				http.Error(w, "storage problem", http.StatusInternalServerError)
			}
		default:
			http.Error(w, "Metric type is not found", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

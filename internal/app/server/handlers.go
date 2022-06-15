package server

import (
	"devops-tpl/internal/entity"
	"net/http"
	"strings"
)

func HealthzHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func UpdateMetricViewHandler(repo MetricRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
			return
		}

		params := strings.Split(
			strings.TrimPrefix(r.URL.Path, "/update/"),
			"/",
		)
		if len(params) != 3 {
			http.Error(w, "Pattern is type/name/value", http.StatusNotFound)
			return
		}

		switch params[0] {
		case "gauge":
			value, err := entity.ParseGauge(params[2])
			if err != nil {
				http.Error(w, "bad value type", http.StatusBadRequest)
				return
			}
			err = repo.StoreGauge(params[1], value)
			if err != nil {
				http.Error(w, "storage problem", http.StatusInternalServerError)
			}
		case "counter":
			value, err := entity.ParseCounter(params[2])
			if err != nil {
				http.Error(w, "bad value type", http.StatusBadRequest)
				return
			}
			err = repo.StoreCounter(params[1], value)
			if err != nil {
				http.Error(w, "storage problem", http.StatusInternalServerError)
			}
		default:
			http.Error(w, "Metric type is not found", http.StatusNotImplemented)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

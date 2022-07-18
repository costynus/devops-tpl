package agent

import (
	"devops-tpl/internal/entity"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
)

type WebAPI struct {
	client *resty.Client
	Key    string
}

func NewWebAPI(client *resty.Client, key string) *WebAPI {
	return &WebAPI{
		client: client,
		Key:    key,
	}
}

func (webAPI *WebAPI) SendMetric(metricName, metricType string, Value *entity.Gauge, Delta *entity.Counter) error {
	metrics := entity.Metric{
		ID:    metricName,
		MType: metricType,
		Value: Value,
		Delta: Delta,
	}
	metrics.Sign(webAPI.Key)
	resp, err := webAPI.client.
		R().
		SetHeader("Content-Type", "application/json").
		SetBody(metrics).
		Post(
			"/update/",
		)
	if err != nil {
		return fmt.Errorf("WebAPI - SendMetric - webAPI.client.R().Post: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("WebAPI - SendMetric - webAPI.client.R().Post: cant't send metric. Status code <> 200")
	}
	return nil
}

func (webAPI *WebAPI) SendMetrics(metrics []entity.Metric) error {
	resp, err := webAPI.client.
		R().
		SetHeader("Content-Type", "application/json").
		SetBody(metrics).
		Post("/updates/")
	if err != nil {
		return fmt.Errorf("WebAPI - SendMetrics - webAPI.client.R().Post: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("WebAPI - SendMetrics - webAPI.client.R().Post: cant't send metric. Status code <> 200")
	}
	return nil
}

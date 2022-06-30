package agent

import (
	"devops-tpl/internal/entity"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
)

type WebAPI struct {
	client *resty.Client
}

func NewWebAPI(client *resty.Client) *WebAPI {
	return &WebAPI{
		client: client,
	}
}

func (webAPI *WebAPI) SendMetric(metricName, metricType string, Value *entity.Gauge, Delta *entity.Counter) error {
	metrics := entity.Metric{
		ID:    metricName,
		MType: metricType,
		Value: Value,
		Delta: Delta,
	}
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

package agent

import (
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

func (webAPI *WebAPI) SendMetric(metricName, metricType string, metricValue interface{}) error {
	resp, err := webAPI.client.
		R().
		SetHeader("Content-Type", "text/plain").
		Post(
			fmt.Sprintf("/update/%s/%s/%v", metricType, metricName, metricValue),
		)
	if err != nil {
		return fmt.Errorf("WebAPI - SendMetric - webAPI.client.R().Post: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("WebAPI - SendMetric - webAPI.client.R().Post: cant't send metric. Status code <> 200")
	}
	return nil
}

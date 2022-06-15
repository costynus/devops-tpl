package agent

import (
	"devops-tpl/internal/entity"
	"devops-tpl/pkg/logger"
	"fmt"
	"net/http"
)

type WebAPI struct {
	l      logger.Interface
	host   string
	client *http.Client
}

func NewWebAPI(l logger.Interface, host string) *WebAPI {
	return &WebAPI{
		l:    l,
		host: host,
		client: &http.Client{
			Transport: &http.Transport{},
		},
	}
}

func (webAPI *WebAPI) SendGaugeMetric(metricName string, metricValue entity.Gauge) error {
	request, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%supdate/gauge/%s/%f", webAPI.host, metricName, metricValue),
		nil,
	)
	if err != nil {
		return fmt.Errorf("WebAPI - SendGaugeMetric - http.NewRequest: %w", err)
	}

	request.Header.Add("Content-Type", "text/plain")
	response, err := webAPI.client.Do(request)
	if err != nil {
		return fmt.Errorf("WebAPI - SendGaugeMetric - http.NewRequest: %w", err)
	}

	defer response.Body.Close()
	webAPI.l.Info("send with status: %s", response.Status)
	return nil
}

func (webAPI *WebAPI) SendCounterMetric(metricName string, metricValue entity.Counter) error {
	request, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%supdate/counter/%s/%d", webAPI.host, metricName, metricValue),
		nil,
	)
	if err != nil {
		return fmt.Errorf("WebAPI - SendCounterMetric - http.NewRequest: %w", err)
	}

	request.Header.Add("Content-Type", "text/plain")
	response, err := webAPI.client.Do(request)
	if err != nil {
		return fmt.Errorf("WebAPI - SendCounterMetric - http.NewRequest: %w", err)
	}

	defer response.Body.Close()
	webAPI.l.Info("send with status: %s", response.Status)
	return nil
}

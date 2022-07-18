package agent

import "devops-tpl/internal/entity"

type (
	AgentWebAPI interface {
		SendMetric(string, string, *entity.Gauge, *entity.Counter) error
		SendMetrics([]entity.Metric) error
	}
)

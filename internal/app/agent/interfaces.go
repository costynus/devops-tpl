package agent

type (
	AgentWebAPI interface {
		SendMetric(string, string, interface{}) error
	}
)

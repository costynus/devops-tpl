package integrationtest_test

import (
	"crypto/hmac"
	"crypto/sha256"
	"devops-tpl/internal/entity"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/suite"
)

type CryptoSuite struct {
	suite.Suite

	serverAddress string
	key           []byte
}

func (suite *CryptoSuite) SetupSuite() {
	suite.serverAddress = serverAddress
	suite.key = []byte(key)
}

func (suite *CryptoSuite) SetHBody(r *resty.Request, m *entity.Metric) *resty.Request {
	hash := suite.Hash(m)
	m.Hash = hash
	return r.SetBody(m)
}

func (suite *CryptoSuite) Hash(m *entity.Metric) string {
	var data string
	switch m.MType {
	case "counter":
		data = fmt.Sprintf("%s:%s:%d", m.ID, m.MType, *m.Delta)
	case "gauge":
		data = fmt.Sprintf("%s:%s:%f", m.ID, m.MType, *m.Value)
	}
	h := hmac.New(sha256.New, suite.key)
	h.Write([]byte(data))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (suite *CryptoSuite) TestCollectAgentMetrics() {
	tests := []struct {
		name   string
		method string
		value  float64
		delta  int64
		update int
		ok     bool
		static bool
	}{
		{method: "counter", name: "PollCount"},
		{method: "gauge", name: "RandomValue"},
		{method: "gauge", name: "Alloc"},
		{method: "gauge", name: "BuckHashSys", static: true},
		{method: "gauge", name: "Frees"},
		{method: "gauge", name: "GCCPUFraction", static: true},
		{method: "gauge", name: "GCSys", static: true},
		{method: "gauge", name: "HeapAlloc"},
		{method: "gauge", name: "HeapIdle"},
		{method: "gauge", name: "HeapInuse"},
		{method: "gauge", name: "HeapObjects"},
		{method: "gauge", name: "HeapReleased", static: true},
		{method: "gauge", name: "HeapSys", static: true},
		{method: "gauge", name: "LastGC", static: true},
		{method: "gauge", name: "Lookups", static: true},
		{method: "gauge", name: "MCacheInuse", static: true},
		{method: "gauge", name: "MCacheSys", static: true},
		{method: "gauge", name: "MSpanInuse", static: true},
		{method: "gauge", name: "MSpanSys", static: true},
		{method: "gauge", name: "Mallocs"},
		{method: "gauge", name: "NextGC", static: true},
		{method: "gauge", name: "NumForcedGC", static: true},
		{method: "gauge", name: "NumGC", static: true},
		{method: "gauge", name: "OtherSys", static: true},
		{method: "gauge", name: "PauseTotalNs", static: true},
		{method: "gauge", name: "StackInuse", static: true},
		{method: "gauge", name: "StackSys", static: true},
		{method: "gauge", name: "Sys", static: true},
		{method: "gauge", name: "TotalAlloc"},
	}

	httpc := resty.New().SetBaseURL(suite.serverAddress)
	req := httpc.R().SetHeader("Content-Type", "application/json")
	for _, tt := range tests {
		var result entity.Metric
		resp, err := req.SetBody(&entity.Metric{
			ID:    tt.name,
			MType: tt.method,
		}).
			SetResult(&result).
			Post(valuePath)

		suite.Assert().NoError(err, "Error with /value")
		suite.Assert().Equal(resp.StatusCode(), http.StatusOK)

		suite.Equal(
			suite.Hash(&result),
			result.Hash,
			fmt.Sprintf(
				"Хеш-сумма не соответствует расчетной (name: %s, method: %s)",
				tt.name,
				tt.method,
			))
	}
}

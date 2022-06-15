package server

import (
	"bytes"
	"devops-tpl/internal/entity"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

func TestRouter(t *testing.T) {
	type (
		args struct {
			repo MetricRepo
		}
		want struct {
			code int
			body string
		}
	)
	tests := []struct {
		name    string
		args    args
		request string
		method  string
		want    want
	}{
		{
			name:    "simple",
			args:    args{&MockMetricRepo{}},
			request: "/",
			method:  http.MethodGet,
			want: want{
				code: 200,
				body: "",
			},
		},
		{
			name:    "simple update gauge value",
			args:    args{&MockMetricRepo{}},
			request: "/update/gauge/Alloc/1.1",
			method:  http.MethodPost,
			want: want{
				code: 200,
				body: "",
			},
		},
		{
			name:    "simple update counter value",
			args:    args{&MockMetricRepo{}},
			request: "/update/counter/Alloc/1",
			method:  http.MethodPost,
			want: want{
				code: 200,
				body: "",
			},
		},
		{
			name:    "simple get gauge value",
			args:    args{&MockMetricRepo{GetMetrics: entity.Gauge(777)}},
			request: "/value/gauge/Alloc",
			method:  http.MethodGet,
			want: want{
				code: 200,
				body: "777",
			},
		},
		{
			name:    "simple get counter value",
			args:    args{&MockMetricRepo{GetMetrics: entity.Counter(777)}},
			request: "/value/counter/Alloc",
			method:  http.MethodGet,
			want: want{
				code: 200,
				body: "777",
			},
		},
		{
			name:    "simple not impl type post",
			args:    args{&MockMetricRepo{}},
			request: "/update/blabla/Alloc/1.1",
			method:  http.MethodPost,
			want: want{
				code: 501,
				body: "metric type is not found\n",
			},
		},
		{
			name:    "simple update error",
			args:    args{&MockMetricRepo{Err: errors.New("error")}},
			request: "/update/gauge/Alloc/1.1",
			method:  http.MethodPost,
			want: want{
				code: 500,
				body: "storage problem\n",
			},
		},
		{
			name:    "simple get error",
			args:    args{&MockMetricRepo{Err: errors.New("error")}},
			request: "/value/gauge/Alloc",
			method:  http.MethodGet,
			want: want{
				code: 404,
				body: "metric not found\n",
			},
		},
	}
	for _, tt := range tests {
		r := chi.NewRouter()
		NewRouter(r, tt.args.repo)
		ts := httptest.NewServer(r)
		defer ts.Close()

		req, err := http.NewRequest(tt.method, ts.URL+tt.request, nil)
		require.NoError(t, err)

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)

		respBody, err := ioutil.ReadAll(resp.Body)
		require.NoError(t, err)

		defer resp.Body.Close()
		body := string(respBody)

		require.Equal(t, tt.want.code, resp.StatusCode)
		require.Equal(t, tt.want.body, body)
	}
}

func TestRouterJSON(t *testing.T) {
	gaugeValue := entity.Gauge(1.1)
	counterValue := entity.Counter(1)
	type (
		args struct {
			repo    MetricRepo
			metrics entity.Metrics
		}
		want struct {
			code  int
			Value *entity.Gauge
			Delta *entity.Counter
		}
	)
	tests := []struct {
		name    string
		args    args
		method  string
		request string
		want    want
	}{
		{
			name: "simple post gauge update",
			args: args{
				&MockMetricRepo{},
				entity.Metrics{ID: "Alloc", MType: "gauge", Delta: nil, Value: &gaugeValue},
			},
			method:  http.MethodPost,
			request: "/update/",
			want: want{
				code:  200,
				Value: nil,
				Delta: nil,
			},
		},
		{
			name: "simple post counter update",
			args: args{
				&MockMetricRepo{},
				entity.Metrics{ID: "Count", MType: "counter", Delta: &counterValue, Value: nil},
			},
			method:  http.MethodPost,
			request: "/update/",
			want: want{
				code:  200,
				Value: nil,
				Delta: nil,
			},
		},
		{
			name: "simple post gauge value",
			args: args{
				&MockMetricRepo{GetMetrics: gaugeValue},
				entity.Metrics{ID: "Alloc", MType: "gauge", Delta: nil, Value: nil},
			},
			method:  http.MethodPost,
			request: "/value/",
			want: want{
				code:  200,
				Value: &gaugeValue,
				Delta: nil,
			},
		},
		{
			name: "simple post gauge value",
			args: args{
				&MockMetricRepo{Err: errors.New("error")},
				entity.Metrics{ID: "Alloc", MType: "gauge", Delta: nil, Value: nil},
			},
			method:  http.MethodPost,
			request: "/value/",
			want: want{
				code:  404,
				Value: nil,
				Delta: nil,
			},
		},
	}

	for _, tt := range tests {
		r := chi.NewRouter()
		NewRouter(r, tt.args.repo)
		ts := httptest.NewServer(r)
		defer ts.Close()

		reqJson, err := json.Marshal(tt.args.metrics)
		require.NoError(t, err)

		req, err := http.NewRequest(tt.method, ts.URL+tt.request, bytes.NewBuffer(reqJson))
		require.NoError(t, err)

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)

		defer resp.Body.Close()

		require.Equal(t, tt.want.code, resp.StatusCode)

		if tt.want.Value != nil || tt.want.Delta != nil {
			var respJson entity.Metrics

			err = json.NewDecoder(resp.Body).Decode(&respJson)
			require.NoError(t, err)

			require.Equal(t, respJson.Value, tt.want.Value)
			require.Equal(t, respJson.Delta, tt.want.Delta)
		}
	}
}

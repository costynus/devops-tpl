package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHealthzHandler(t *testing.T) {
	type want struct {
		StatusCode int
	}
	tests := []struct {
		name    string
		request string
		want    want
	}{
		{
			name:    "simple healthz",
			request: "/healthz",
			want:    want{200},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, tt.request, nil)

			w := httptest.NewRecorder()
			h := http.HandlerFunc(HealthzHandler)
			h.ServeHTTP(w, request)
			res := w.Result()
			defer res.Body.Close()

			require.Equal(t, res.StatusCode, tt.want.StatusCode)
		})
	}
}

func TestUpdateMetricViewHandler(t *testing.T) {
	type (
		want struct {
			StatusCode int
		}
		args struct {
			repo MetricRepo
		}
		test struct {
			name    string
			args    args
			want    want
			request string
			method  string
		}
	)

	tests := []test{
		{
			name: "simple success",
			args: args{
				repo: &MockMetricRepo{},
			},
			want:    want{200},
			request: "/update/gauge/Alloc/1.1",
			method:  http.MethodPost,
		},
		{
			name: "simple method error",
			args: args{
				repo: &MockMetricRepo{},
			},
			want:    want{405},
			request: "/update/gauge/Alloc/1.1",
			method:  http.MethodGet,
		},
		{
			name: "simple not found",
			args: args{
				repo: &MockMetricRepo{},
			},
			want:    want{404},
			request: "/update/gauge/",
			method:  http.MethodPost,
		},
		{
			name: "simple bad request 1",
			args: args{
				repo: &MockMetricRepo{},
			},
			want:    want{400},
			request: "/update/counter/testCounter/none",
			method:  http.MethodPost,
		},
		{
			name: "simple bad request 2",
			args: args{
				repo: &MockMetricRepo{},
			},
			want:    want{400},
			request: "/update/gauge/testCounter/none",
			method:  http.MethodPost,
		},
		{
			name: "bad type",
			args: args{
				repo: &MockMetricRepo{},
			},
			want:    want{404},
			request: "/update/blabla/testCounter/none",
			method:  http.MethodPost,
		},
		{
			name: "only type and name",
			args: args{
				repo: &MockMetricRepo{},
			},
			want:    want{404},
			request: "/update/gauge/testCounter",
			method:  http.MethodPost,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, tt.request, nil)

			w := httptest.NewRecorder()
			h := http.HandlerFunc(UpdateMetricViewHandler(tt.args.repo))
			h.ServeHTTP(w, request)
			res := w.Result()
			defer res.Body.Close()

			require.Equal(t, tt.want.StatusCode, res.StatusCode)
		})
	}
}

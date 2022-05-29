package server

import (
	"devops-tpl/internal/entity"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	respBody, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)

	defer resp.Body.Close()

	return resp, string(respBody)
}

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

		resp, respBody := testRequest(t, ts, tt.method, tt.request)
		require.Equal(t, tt.want.code, resp.StatusCode)
		require.Equal(t, tt.want.body, respBody)
	}
}

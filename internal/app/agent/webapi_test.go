package agent

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/require"
)

func TestNewWebAPI(t *testing.T) {
	type args struct {
		client *resty.Client
	}
	tests := []struct {
		name string
		args args
		want *WebAPI
	}{
		{
			name: "simple",
			args: args{&resty.Client{}},
			want: &WebAPI{client: &resty.Client{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewWebAPI(tt.args.client)
			require.Equal(t, got, tt.want)
		})
	}
}

func TestWebAPI_SendMetric(t *testing.T) {
	type (
		args struct {
			metricName  string
			metricType  string
			metricValue interface{}
		}
	)
	tests := []struct {
		name           string
		args           args
		response       int
		withTestServer bool
		wantErr        bool
	}{
		{
			name:           "simple",
			args:           args{"Alloc", "gauge", 1.1},
			response:       http.StatusOK,
			withTestServer: true,
			wantErr:        false,
		},
		{
			name:           "simple error code",
			args:           args{"Alloc", "gauge", 1.1},
			response:       http.StatusBadRequest,
			withTestServer: true,
			wantErr:        true,
		},
		{
			name:           "simple without server",
			args:           args{"Alloc", "gauge", 1.1},
			withTestServer: false,
			wantErr:        true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.response)
			}))
			defer testServer.Close()

			var serverURL string
			if tt.withTestServer {
				serverURL = testServer.URL
			}

			webAPI := &WebAPI{
				client: resty.New().SetBaseURL(serverURL),
			}

			err := webAPI.SendMetric(tt.args.metricName, tt.args.metricType, tt.args.metricValue)
			if !tt.wantErr {
				return
			}
			require.Error(t, err)
		})
	}
}

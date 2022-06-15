package agent

import (
	"devops-tpl/internal/entity"
	"runtime"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMetrics_UpdateMetrics(t *testing.T) {
	type fields struct {
		PollCount   entity.Counter
		RandomValue entity.Gauge
		Mutex       *sync.Mutex
		collector   *collector
	}
	tests := []struct {
		name   string
		fields fields
		want   entity.Counter
	}{
		{
			name: "simple UpdateMetrics",
			fields: fields{
				100,
				100,
				&sync.Mutex{},
				&collector{},
			},
			want: 101,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Metrics{
				PollCount:   tt.fields.PollCount,
				RandomValue: tt.fields.RandomValue,
				Mutex:       tt.fields.Mutex,
				collector:   tt.fields.collector,
			}
			m.UpdateMetrics()
			require.Equal(t, m.PollCount, tt.want)
		})
	}
}

func TestMetrics_collectMetrics(t *testing.T) {
	type fields struct {
		PollCount   entity.Counter
		RandomValue entity.Gauge
		Mutex       *sync.Mutex
		collector   *collector
	}
	type args struct {
		memStats *runtime.MemStats
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   entity.Gauge
	}{
		{
			name: "simple collectMetrics",
			fields: fields{
				100,
				100,
				&sync.Mutex{},
				&collector{},
			},
			args: args{
				&runtime.MemStats{Alloc: 1000},
			},
			want: entity.Gauge(1000),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Metrics{
				PollCount:   tt.fields.PollCount,
				RandomValue: tt.fields.RandomValue,
				Mutex:       tt.fields.Mutex,
				collector:   tt.fields.collector,
			}
			m.collectMetrics(tt.args.memStats)
			require.Equal(t, m.Alloc, tt.want)
		})
	}
}

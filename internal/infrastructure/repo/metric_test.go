package repo

import (
	"devops-tpl/internal/entity"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMetricRepo_StoreGauge(t *testing.T) {
	type fields struct {
		data map[string]interface{}
	}
	type args struct {
		name  string
		value entity.Gauge
	}
	type want struct {
		value entity.Gauge
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name: "simple",
			fields: fields{
				data: map[string]interface{}{},
			},
			args: args{"Alloc", 1.1},
			want: want{entity.Gauge(1.1)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &MetricRepo{
				data:  tt.fields.data,
				Mutex: &sync.Mutex{},
			}
			r.StoreGauge(tt.args.name, tt.args.value)
			got, ok := r.data[tt.args.name]
			require.True(t, ok)
			require.Equal(t, got, tt.want.value)
		})
	}
}

func TestMetricRepo_StoreCounter(t *testing.T) {
	type fields struct {
		data map[string]interface{}
	}
	type args struct {
		name  string
		value entity.Counter
	}
	type want struct {
		value entity.Counter
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   want
	}{
		{
			name:   "simple",
			fields: fields{data: map[string]interface{}{}},
			args:   args{"Total", 1},
			want:   want{entity.Counter(1)},
		},
		{
			name:   "simple add",
			fields: fields{data: map[string]interface{}{"Total": entity.Counter(1)}},
			args:   args{"Total", 1},
			want:   want{entity.Counter(2)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &MetricRepo{
				data:  tt.fields.data,
				Mutex: &sync.Mutex{},
			}
			r.StoreCounter(tt.args.name, tt.args.value)
			got, ok := r.data[tt.args.name]
			require.True(t, ok)
			require.Equal(t, got, tt.want.value)
		})
	}
}

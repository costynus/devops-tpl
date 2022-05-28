package entity

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseGauge(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name    string
		args    args
		want    Gauge
		wantErr bool
	}{
		{
			name:    "simple",
			args:    args{"1.1"},
			want:    Gauge(1.1),
			wantErr: false,
		},
		{
			name:    "simple error",
			args:    args{"none"},
			want:    Gauge(0),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseGauge(tt.args.value)
			if !tt.wantErr {
				require.Equal(t, got, tt.want)
				return
			}
			require.Error(t, err)
		})
	}
}

func TestParseCounter(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name    string
		args    args
		want    Counter
		wantErr bool
	}{
		{
			name:    "simple",
			args:    args{"1"},
			want:    Counter(1),
			wantErr: false,
		},
		{
			name:    "simple error",
			args:    args{"none"},
			want:    Counter(0),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseCounter(tt.args.value)
			if !tt.wantErr {
				require.Equal(t, got, tt.want)
				return
			}
			require.Error(t, err)
		})
	}
}

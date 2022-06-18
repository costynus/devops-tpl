package usecase_test

import (
	"context"
	"devops-tpl/internal/entity"
	"devops-tpl/internal/usecase"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

var errFromRepo = errors.New("some error")

type test struct {
	name string
	mock func()
	arg  entity.Metric
	res  interface{}
	err  error
}

func devops(t *testing.T) (*usecase.DevOpsUseCase, *MockMetricRepo) {
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	repo := NewMockMetricRepo(mockCtl)

	devops := usecase.New(repo)

	return devops, repo
}

func TestMetricNames(t *testing.T) {
	t.Parallel()

	devops, repo := devops(t)

	tests := []test{
		{
			name: "simple empty result",
			mock: func() {
				repo.EXPECT().GetMetricNames(context.Background()).Return(nil)
			},
			res: []string(nil),
			err: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tt.mock()

			res, err := devops.MetricNames(context.Background())

			require.ErrorIs(t, err, tt.err)
			require.Equal(t, res, tt.res)
		})
	}
}

func TestStoreMetric(t *testing.T) {
	t.Parallel()

	delta1 := entity.Counter(1)
	delta2 := entity.Counter(2)

	devops, repo := devops(t)

	tests := []test{
		{
			name: "simple empty result",
			mock: func() {
				repo.EXPECT().StoreMetric(context.Background(), entity.Metric{}).Return(nil)
			},
			arg: entity.Metric{},
			err: nil,
		},
		{
			name: "simple repo error",
			mock: func() {
				repo.EXPECT().StoreMetric(context.Background(), entity.Metric{}).Return(errFromRepo)
			},
			arg: entity.Metric{},
			err: errFromRepo,
		},
		{
			name: "gauge metric",
			mock: func() {
				repo.EXPECT().StoreMetric(context.Background(), entity.Metric{MType: usecase.Gauge}).Return(nil)
			},
			arg: entity.Metric{MType: usecase.Gauge},
			err: nil,
		},
		{
			name: "gauge repo error",
			mock: func() {
				repo.EXPECT().StoreMetric(context.Background(), entity.Metric{MType: usecase.Gauge}).Return(errFromRepo)
			},
			arg: entity.Metric{MType: usecase.Gauge},
			err: errFromRepo,
		},
		{
			name: "counter new metric",
			mock: func() {
				repo.EXPECT().GetMetric(context.Background(), "").Return(entity.Metric{}, errFromRepo)
				repo.EXPECT().StoreMetric(context.Background(), entity.Metric{MType: usecase.Counter}).Return(nil)
			},
			arg: entity.Metric{MType: usecase.Counter},
			err: nil,
		},
		{
			name: "counter old metric",
			mock: func() {
				repo.EXPECT().GetMetric(context.Background(), "").Return(entity.Metric{MType: usecase.Counter, Delta: &delta1}, nil)
				repo.EXPECT().StoreMetric(context.Background(), entity.Metric{MType: usecase.Counter, Delta: &delta2}).Return(nil)
			},
			arg: entity.Metric{MType: usecase.Counter, Delta: &delta1},
			err: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tt.mock()

			err := devops.StoreMetric(context.Background(), tt.arg)

			require.ErrorIs(t, err, tt.err)
		})
	}
}

func TestMetric(t *testing.T) {
	t.Parallel()

	devops, repo := devops(t)

	tests := []test{
		{
			name: "simple empty result",
			mock: func() {
				repo.EXPECT().GetMetric(context.Background(), "").Return(entity.Metric{}, nil)
			},
			res: entity.Metric{},
			err: nil,
		},
		{
			name: "repo error",
			mock: func() {
				repo.EXPECT().GetMetric(context.Background(), "").Return(entity.Metric{}, errFromRepo)
			},
			res: entity.Metric{},
			err: errFromRepo,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tt.mock()

			res, err := devops.Metric(context.Background(), entity.Metric{})

			require.ErrorIs(t, err, tt.err)
			require.EqualValues(t, res, tt.res)
		})
	}
}

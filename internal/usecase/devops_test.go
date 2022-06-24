package usecase_test

import (
	"context"
	"devops-tpl/internal/entity"
	"devops-tpl/internal/infrastructure/repo"
	"devops-tpl/internal/usecase"
	"devops-tpl/pkg/logger"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

type test struct {
	name string
	mock func()
	arg  entity.Metric
	res  interface{}
	err  error
}

func devops(t *testing.T) (*usecase.DevOpsUseCase, *MockMetricRepo) {
	l := logger.New("debug")
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	repoMock := NewMockMetricRepo(mockCtl)

	devops := usecase.New(repoMock, l)

	return devops, repoMock
}

func TestMetricNames(t *testing.T) {
	t.Parallel()

	devops, repoMock := devops(t)

	tests := []test{
		{
			name: "simple empty result",
			mock: func() {
				repoMock.EXPECT().GetMetricNames(context.Background()).Return(nil)
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

			res, err := devops.GetMetricNames(context.Background())

			require.ErrorIs(t, err, tt.err)
			require.Equal(t, res, tt.res)
		})
	}
}

func TestStoreMetric(t *testing.T) {
	t.Parallel()

	delta1 := entity.Counter(1)
	delta2 := entity.Counter(2)

	devops, repoMock := devops(t)

	tests := []test{
		{
			name: "simple empty error",
			mock: func() {
				repoMock.EXPECT().StoreMetric(context.Background(), entity.Metric{}).Return(nil)
			},
			arg: entity.Metric{},
			err: usecase.ErrNotImplemented,
		},
		{
			name: "simple repo error",
			mock: func() {
				repoMock.EXPECT().StoreMetric(context.Background(), entity.Metric{MType: usecase.Gauge}).Return(repo.ErrNotFound)
			},
			arg: entity.Metric{MType: usecase.Gauge},
			err: usecase.ErrNotFound,
		},
		{
			name: "gauge metric",
			mock: func() {
				repoMock.EXPECT().StoreMetric(context.Background(), entity.Metric{MType: usecase.Gauge}).Return(nil)
			},
			arg: entity.Metric{MType: usecase.Gauge},
			err: nil,
		},
		{
			name: "gauge repo error",
			mock: func() {
				repoMock.EXPECT().StoreMetric(context.Background(), entity.Metric{MType: usecase.Gauge}).Return(repo.ErrNotFound)
			},
			arg: entity.Metric{MType: usecase.Gauge},
			err: usecase.ErrNotFound,
		},
		{
			name: "counter new metric",
			mock: func() {
				repoMock.EXPECT().GetMetric(context.Background(), "").Return(entity.Metric{}, repo.ErrNotFound)
				repoMock.EXPECT().StoreMetric(context.Background(), entity.Metric{MType: usecase.Counter}).Return(nil)
			},
			arg: entity.Metric{MType: usecase.Counter},
			err: nil,
		},
		{
			name: "counter old metric",
			mock: func() {
				repoMock.EXPECT().GetMetric(context.Background(), "").Return(entity.Metric{MType: usecase.Counter, Delta: &delta1}, nil)
				repoMock.EXPECT().StoreMetric(context.Background(), entity.Metric{MType: usecase.Counter, Delta: &delta2}).Return(nil)
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

	devops, repoMock := devops(t)

	tests := []test{
		{
			name: "simple empty result",
			mock: func() {
				repoMock.EXPECT().GetMetric(context.Background(), "").Return(entity.Metric{}, nil)
			},
			res: entity.Metric{},
			err: nil,
		},
		{
			name: "repo error",
			mock: func() {
				repoMock.EXPECT().GetMetric(context.Background(), "").Return(entity.Metric{}, repo.ErrNotFound)
			},
			res: entity.Metric{},
			err: usecase.ErrNotFound,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tt.mock()

			res, err := devops.GetMetric(context.Background(), entity.Metric{})

			require.ErrorIs(t, err, tt.err)
			require.EqualValues(t, res, tt.res)
		})
	}
}

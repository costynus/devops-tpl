package repo

import (
	"bufio"
	"context"
	"devops-tpl/internal/entity"
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type MetricRepo struct {
	data          map[string]entity.Metric
	StoreFilePath string
	Restore       bool
	Mutex         *sync.Mutex
}

func New(opts ...Option) *MetricRepo {
	metricRepo := &MetricRepo{
		Mutex: &sync.Mutex{},
	}
	metricRepo.data = make(map[string]entity.Metric)

	// Set Options
	for _, opt := range opts {
		opt(metricRepo)
	}

	if metricRepo.Restore {
		metricRepo.UploadFromFile(context.Background())
	}

	return metricRepo
}

func (r MetricRepo) StoreToFile() error {
	file, err := os.OpenFile(r.StoreFilePath, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return fmt.Errorf("MetricRepo.StoreToFile - os.OpenFile: %w", err)
	}
	defer file.Close()
	writer := bufio.NewWriter(file)

	r.Mutex.Lock()
	data, err := json.Marshal(r.data)
	r.Mutex.Unlock()
	if err != nil {
		return fmt.Errorf("MetricRepo.StoreToFile - json.Marshal: %w", err)
	}

	_, err = writer.Write(data)
	if err != nil {
		return fmt.Errorf("MetricRepo.StoreToFile - writer.Write: %w", err)
	}

	if err := writer.WriteByte('\n'); err != nil {
		return fmt.Errorf("MetricRepo.StoreToFile - writer.WriteByte: %w", err)
	}
	writer.Flush()
	return nil
}

func (r *MetricRepo) UploadFromFile(ctx context.Context) error {
	file, err := os.OpenFile(r.StoreFilePath, os.O_RDONLY, 0777)
	if err != nil {
		return fmt.Errorf("MetricRepo.UploadFromFile - os.OpenFile: %w", err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	data, err := reader.ReadBytes('\n')
	if err != nil {
		return fmt.Errorf("MetricRepo.UploadFromFile - reader.ReadBytes: %w", err)
	}
	err = json.Unmarshal(data, &r.data)
	if err != nil {
		return fmt.Errorf("MetricRepo.UploadFromFile - json.Unmarshal: %w", err)
	}
	return nil
}

func (r *MetricRepo) GetMetricNames(ctx context.Context) []string {
	var list []string
	for name := range r.data {
		list = append(list, name)
	}
	return list
}

func (r *MetricRepo) StoreMetric(ctx context.Context, metric entity.Metric) error {
	r.Mutex.Lock()
	r.data[metric.ID] = metric
	r.Mutex.Unlock()
	return nil
}

func (r *MetricRepo) GetMetric(ctx context.Context, name string) (entity.Metric, error) {
	r.Mutex.Lock()
	metric, ok := r.data[name]
	r.Mutex.Unlock()
	if !ok {
		return entity.Metric{}, ErrNotFound
	}
	return metric, nil
}

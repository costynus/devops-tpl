package repo

import (
	"bufio"
	"devops-tpl/internal/entity"
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type MetricRepo struct {
	data      map[string]entity.Metrics
	StoreFile string
	Mutex     *sync.Mutex
	storeMode bool
}

func New(StoreFile string, Restore bool) *MetricRepo {
	metricRepo := MetricRepo{
		StoreFile: StoreFile,
		Mutex:     &sync.Mutex{},
	}
	metricRepo.data = make(map[string]entity.Metrics)
	metricRepo.storeMode = (StoreFile == " ")
	if Restore {
		metricRepo.UploadFromFile()
	}
	return &metricRepo
}

func (r MetricRepo) StoreToFile() error {
	if !r.storeMode {
		return nil
	}
	file, err := os.OpenFile(r.StoreFile, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return fmt.Errorf("MetricRepo.StoreToFile - os.OpenFile: %w", err)
	}
	defer file.Close()
	writer := bufio.NewWriter(file)

	data, err := json.Marshal(r.data)
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
	return nil
}

func (r *MetricRepo) UploadFromFile() error {
	if !r.storeMode {
		return nil
	}
	file, err := os.OpenFile(r.StoreFile, os.O_RDONLY, 0777)
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
	fmt.Println(r.data)
	return nil
}

func (r *MetricRepo) GetMetricNames() []string {
	var list []string
	for name := range r.data {
		list = append(list, name)
	}
	return list
}

func (r *MetricRepo) StoreGauge(name string, value entity.Gauge) error {
	r.Mutex.Lock()
	r.data[name] = entity.Metrics{ID: name, MType: value.String(), Value: &value}
	r.Mutex.Unlock()
	return nil
}

func (r *MetricRepo) AddCounter(name string, value entity.Counter) error {
	r.Mutex.Lock()
	oldMetric, ok := r.data[name]
	if ok {
		delta := value + *r.data[name].Delta
		oldMetric.Delta = &delta
		r.data[name] = oldMetric
	} else {
		r.data[name] = entity.Metrics{ID: name, MType: value.String(), Delta: &value}
	}
	r.Mutex.Unlock()
	return nil
}

func (r *MetricRepo) GetMetric(name string) (entity.Metrics, error) {
	value, ok := r.data[name]
	if !ok {
		return entity.Metrics{}, fmt.Errorf("not Found (%s)", name)
	}
	return value, nil
}

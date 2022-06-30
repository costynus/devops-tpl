package repo

type Option func(*MetricRepo)

func Restore() Option {
	return func(repo *MetricRepo) {
		repo.Restore = true
	}
}

func StoreFilePath(StoreFilePath string) Option {
	return func(repo *MetricRepo) {
		repo.StoreFilePath = StoreFilePath
	}
}

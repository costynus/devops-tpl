package usecase

import "time"

type Option func(*DevOpsUseCase)

func WriteFileDuration(duration time.Duration) Option {
	return func(uc *DevOpsUseCase) {
		uc.writeToFileWithDuration = true
		uc.writeFileDuration = duration
	}
}

func AsynchWriteFile() Option {
	return func(uc *DevOpsUseCase) {
		uc.asynchWriteFile = true
	}
}

func SyncWriteFile() Option {
	return func(uc *DevOpsUseCase) {
		uc.synchWriteFile = true
	}
}

func CheckSign(key string) Option {
	return func(uc *DevOpsUseCase) {
		uc.checkSign = true
		uc.cryptoKey = key
	}
}

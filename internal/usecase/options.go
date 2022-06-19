package usecase

import "time"

type Option func(*DevOpsUseCase)

func WriteFileDuration(duration time.Duration) Option {
	return func(uc *DevOpsUseCase) {
		uc.writeToFileWithDuration = true
		uc.writeFileDuration = duration
	}
}

func SynchWriteFile() Option {
	return func(uc *DevOpsUseCase) {
		uc.synchWriteFile = true
	}
}

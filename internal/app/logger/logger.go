package logger

import (
	"sync"

	"go.uber.org/zap"
)

var once sync.Once
var instance *zap.SugaredLogger
var loggerErr error

func getInstance() (*zap.SugaredLogger, error) {
	once.Do(func() {
		logger, loggerErr := zap.NewProduction()
		if loggerErr != nil {
			return
		}

		instance = logger.Sugar()
	})

	return instance, loggerErr
}

func NewLogger() (*zap.SugaredLogger, error) {
	return getInstance()
}

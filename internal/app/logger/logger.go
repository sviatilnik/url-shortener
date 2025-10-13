package logger

import (
	"sync"

	"go.uber.org/zap"
)

var once sync.Once
var instance *zap.SugaredLogger

func getInstance() *zap.SugaredLogger {
	once.Do(func() {
		logger, err := zap.NewProduction()
		if err != nil {
			panic(err)
		}

		instance = logger.Sugar()
	})

	return instance
}

func NewLogger() *zap.SugaredLogger {
	return getInstance()
}

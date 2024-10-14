// Модуль логера для приложения.
package logger

import (
	"go.uber.org/zap"
)

// Конструктор логгера для приложения. Используется zap.
func NewLogger() (*zap.SugaredLogger, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	defer logger.Sync()

	return logger.Sugar(), nil
}

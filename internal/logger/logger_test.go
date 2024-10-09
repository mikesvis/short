// Модуль логера для приложения.
package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestNewLogger(t *testing.T) {
	prodLogger, err := zap.NewProduction()
	defer prodLogger.Sync()
	tests := []struct {
		name string
		want *zap.SugaredLogger
	}{
		{
			name: "Init new logger no error",
			want: prodLogger.Sugar(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.NoError(t, err)
			logger, err := NewLogger()
			require.NoError(t, err)
			assert.ObjectsAreEqual(prodLogger, logger)
		})
	}
}

// Пакет приложения сокращателя ссылок.
package app

import (
	"testing"

	"github.com/mikesvis/short/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	config := &config.Config{
		ServerAddress:   "127.0.0.1",
		BaseURL:         "http://short.go",
		FileStoragePath: "",
		DatabaseDSN:     "",
	}
	app := New(config)
	tests := []struct {
		name string
		want *App
	}{
		{
			name: "App initiates",
			want: app,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := New(config)
			assert.ObjectsAreEqual(tt.want, result)
		})
	}
}

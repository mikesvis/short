// Пакет приложения сокращателя ссылок.
package app

import (
	"testing"

	"github.com/mikesvis/short/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestApp_Run(t *testing.T) {
	type fields struct {
		config *config.Config
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "App fails to start",
			fields: fields{&config.Config{
				ServerAddress:   "1000.0.0.100",
				BaseURL:         "http://short.go",
				FileStoragePath: "",
				DatabaseDSN:     "",
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := New(tt.fields.config)
			err := a.Run()
			require.Error(t, err)
		})
	}
}

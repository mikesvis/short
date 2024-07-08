package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want *Config
	}{
		{
			name: "Default config with empty FILE_STORAGE_PATH env variable",
			want: &Config{
				ServerAddress:   "localhost:8080",
				BaseURL:         "http://localhost:8080",
				FileStoragePath: "",
			},
		},
	}
	t.Setenv("FILE_STORAGE_PATH", "")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := NewConfig()
			assert.EqualValues(t, tt.want, config)
		})
	}
}

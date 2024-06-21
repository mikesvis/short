package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetServerAddress(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "Default server address no envs or flags",
			want: "localhost:8080",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, GetServerAddress())
		})
	}
}

func TestGetBaseURL(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "Default base url no envs or flags",
			want: "http://localhost:8080",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, GetBaseURL())
		})
	}
}

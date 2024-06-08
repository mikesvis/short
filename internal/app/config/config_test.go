package config

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitConfig(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Args = append(os.Args, "--addr=http://b.com:566/something.avsc")
		})
	}
}

func TestGetServerHostAddr(t *testing.T) {
	type arg struct {
		key   string
		value string
	}

	tests := []struct {
		name      string
		arg       arg
		want      string
		wantError bool
	}{
		{
			name: "Server address from a flag",
			arg: arg{
				key:   "a",
				value: "example.com:8888",
			},
			want: "example.com:8888",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Args = append(os.Args, fmt.Sprintf("-%s=%s", tt.arg.key, tt.arg.value))
			InitConfig()
			assert.Equal(t, tt.want, GetServerHostAddr())
		})
	}
}

func TestGetShortLinkAddr(t *testing.T) {
	type arg struct {
		key   string
		value string
	}

	tests := []struct {
		name      string
		arg       arg
		want      string
		wantError bool
	}{
		{
			name: "Server address from a flag",
			arg: arg{
				key:   "b",
				value: "https://example.com:8888",
			},
			want: "https://example.com:8888",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Args = append(os.Args, fmt.Sprintf("-%s=%s", tt.arg.key, tt.arg.value))
			InitConfig()
			assert.Equal(t, tt.want, GetShortLinkAddr())
		})
	}
}

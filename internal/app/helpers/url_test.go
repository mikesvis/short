package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFormattedURL(t *testing.T) {
	myURLOptions = URLOptions{
		scheme: "http",
		host:   "localhost",
		port:   "8080",
	}

	SetURLOptions(myURLOptions.scheme, myURLOptions.host, myURLOptions.port)
	type args struct {
		shortKey string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Formatted Url",
			args: args{
				shortKey: "IddQd",
			},
			want: "http://localhost:8080/IddQd",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, GetFormattedURL(tt.args.shortKey))
		})
	}
}

func TestValidateURL(t *testing.T) {
	tests := []struct {
		name    string
		URL     string
		wantErr bool
	}{
		{
			name:    "Valid url",
			URL:     "http://ya.ru",
			wantErr: false,
		}, {
			name:    "Invalid not empty url",
			URL:     "://ya.ru",
			wantErr: true,
		}, {
			name:    "Empty string provided as url",
			URL:     "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateURL(tt.URL)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestSanitizeURL(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  string
	}{
		{
			name:  "Valid string no spaces",
			value: "abcde",
			want:  "abcde",
		}, {
			name:  "Valid string left spaces",
			value: "    abcde",
			want:  "abcde",
		}, {
			name:  "Valid string right spaces",
			value: "abcde     ",
			want:  "abcde",
		}, {
			name:  "Valid string spaces both sides",
			value: "    abcde    ",
			want:  "abcde",
		}, {
			name:  "Empty string",
			value: "",
			want:  "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, SanitizeURL(tt.value))
		})
	}
}

func TestSetURLOptions(t *testing.T) {
	type args struct {
		scheme string
		host   string
		port   string
	}
	tests := []struct {
		name string
		args args
		want URLOptions
	}{
		{
			name: "Simple set url options",
			args: args{
				scheme: "http",
				host:   "localhost",
				port:   "8080",
			},
			want: URLOptions{
				scheme: "http",
				host:   "localhost",
				port:   "8080",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			settedOptions := SetURLOptions(tt.args.scheme, tt.args.host, tt.want.port)
			assert.NotEmpty(t, settedOptions)
			assert.Equal(t, tt.want, settedOptions)
		})
	}
}

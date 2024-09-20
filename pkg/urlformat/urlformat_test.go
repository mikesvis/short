package urlformat

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFormattedURL(t *testing.T) {
	type args struct {
		linkServerAddress string
		shortKey          string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Formatted Url",
			args: args{
				linkServerAddress: "http://example.com",
				shortKey:          "IddQd",
			},
			want: "http://example.com/IddQd",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, FormatURL(tt.args.linkServerAddress, tt.args.shortKey))
		})
	}
}

func BenchmarkGetFormattedURL(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FormatURL("http://example.com", "IddQd")
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

func BenchmarkValidateURL(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ValidateURL("http://ya.ru")
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

func BenchmarkSanitizeURL(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SanitizeURL("abcde")
	}
}

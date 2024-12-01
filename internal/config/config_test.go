package config

import (
	"os"
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
				ServerAddress:     "localhost:8080",
				BaseURL:           "http://localhost:8080",
				FileStoragePath:   "",
				DatabaseDSN:       "",
				EnableHTTPS:       false,
				ServerKeyPath:     "",
				ServerCertPath:    "",
				GRPCServerAddress: "localhost:8082",
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

func Test_parseFile(t *testing.T) {
	tmpFile, _ := os.CreateTemp(os.TempDir(), "dbtest*.json")
	tmpFile.Write([]byte(`{"ssssooooomemmemm,stif`))
	tmpFile.Close()

	type args struct {
		c  *Config
		fp string
	}
	tests := []struct {
		name       string
		args       args
		wantOutput string
		wantErr    bool
	}{
		{
			name: "Cant open file",
			args: args{
				c:  &Config{},
				fp: "!",
			},
			wantOutput: "Unable open config file",
			wantErr:    true,
		},
		{
			name: "Cant decode file",
			args: args{
				c:  &Config{},
				fp: tmpFile.Name(),
			},
			wantOutput: "Unable parse config file",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		if tt.wantErr {
			assert.Panics(t, func() { parseFile(tt.args.c, tt.args.fp) })
		}
	}
}

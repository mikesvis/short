package config

import (
	"github.com/caarlos0/env"
	flag "github.com/spf13/pflag"
)

type Config struct {
	ServerAddress   string `env:"SERVER_ADDRESS"`
	BaseURL         string `env:"BASE_URL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"  envDefault:"/tmp/short-url-db.json"`
	DatabaseDSN     string `env:"DATABASE_DSN"`
}

func NewConfig() *Config {
	config := Config{
		ServerAddress:   "localhost:8080",
		BaseURL:         "http://localhost:8080",
		FileStoragePath: "",
		DatabaseDSN:     "",
	}

	parseFlags(&config)
	env.Parse(&config)

	return &config
}

func parseFlags(c *Config) {
	flag.StringVarP(&c.ServerAddress, "address", "a", "", "address of shortener service server")
	flag.StringVarP(&c.BaseURL, "basepath", "b", "", "address of short link basepath")
	flag.StringVarP(&c.FileStoragePath, "file_storage_path", "f", "", "path to file storage of URLs")
	flag.StringVarP(&c.DatabaseDSN, "database_dsn", "d", "", "db connection string")
	flag.Parse()
}

package config

import (
	"fmt"

	"github.com/caarlos0/env"
	flag "github.com/spf13/pflag"
)

type Address string

type Config struct {
	ServerAddress   Address `env:"SERVER_ADDRESS"`
	BaseURL         Address `env:"BASE_URL"`
	FileStoragePath string  `env:"FILE_STORAGE_PATH"  envDefault:"/tmp/short-url-db.json"`
	DatabaseDSN     string  `env:"DATABASE_DSN"`
}

func (a *Address) Set(flagValue string) error {
	*a = Address(flagValue)
	return nil
}

func (a *Address) String() string {
	return string(*a)
}

func (a *Address) Type() string {
	return "string"
}

func (a *Address) UnmarshalText(envValue []byte) error {
	if len(envValue) == 0 {
		return fmt.Errorf("cannot be empty")
	}
	*a = Address(string(envValue))
	return nil
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
	flag.VarP(&c.ServerAddress, "address", "a", "address of shortener service server")
	flag.VarP(&c.BaseURL, "basepath", "b", "address of short link basepath")
	c.FileStoragePath = *flag.StringP("file_storage_path", "f", "", "path to file storage of URLs")
	c.DatabaseDSN = *flag.StringP("database_dsn", "d", "", "db connection string")
	flag.Parse()
}

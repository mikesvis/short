package config

import (
	"fmt"

	"github.com/caarlos0/env"
	flag "github.com/spf13/pflag"
)

type Address string

type FilePath string

type Config struct {
	ServerAddress   Address  `env:"SERVER_ADDRESS"`
	BaseURL         Address  `env:"BASE_URL"`
	FileStoragePath FilePath `env:"FILE_STORAGE_PATH"  envDefault:"/tmp/short-url-db.json"`
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

func (s *FilePath) Set(flagValue string) error {
	*s = FilePath(string(flagValue))
	return nil
}

func (s *FilePath) String() string {
	return string(*s)
}

func (s *FilePath) Type() string {
	return "string"
}

func New() *Config {
	config := Config{
		ServerAddress:   "localhost:8080",
		BaseURL:         "http://localhost:8080",
		FileStoragePath: "",
	}

	parseFlags(&config)
	env.Parse(&config)

	return &config
}

func parseFlags(c *Config) {
	flag.VarP(&c.ServerAddress, "address", "a", "address of shortener service server")
	flag.VarP(&c.BaseURL, "basepath", "b", "address of short link basepath")
	flag.VarP(&c.FileStoragePath, "file_storage_path", "f", "path to file storage of URLs")
	flag.Parse()
}

package config

import (
	"fmt"

	"github.com/caarlos0/env"
	"github.com/mikesvis/short/internal/logger"
	flag "github.com/spf13/pflag"
)

type Address string

type Config struct {
	ServerAddress Address `env:"SERVER_ADDRESS"`
	BaseURL       Address `env:"BASE_URL"`
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

var config Config = Config{
	ServerAddress: "localhost:8080",
	BaseURL:       "http://localhost:8080",
}

func InitConfig() {
	parseFlags(&config)
	env.Parse(&config)
	logger.Log.Infow("Config initialized", "config", config)
}

func parseFlags(config *Config) {
	flag.VarP(&config.ServerAddress, "address", "a", "address of shortener service server")
	flag.VarP(&config.BaseURL, "basepath", "b", "address of short link basepath")
	flag.Parse()
}

func GetServerAddress() string {
	return string(config.ServerAddress)
}

func GetBaseURL() string {
	return string(config.BaseURL)
}

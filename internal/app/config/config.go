package config

import (
	"fmt"
	"log"

	"github.com/caarlos0/env"
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

func (a Address) String() string {
	return string(a)
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
	log.Printf("initialized config %+v", config)
}

func parseFlags(config *Config) {
	flag.VarP(&config.ServerAddress, "address", "a", "address of shortener service server")
	flag.VarP(&config.BaseURL, "basepath", "b", "address of short link basepath")
	flag.Parse()
}

// func parseEnvs(c Config) {
// 	// if _, ok := os.LookupEnv("SERVER_ADDRESS"); ok {
// 	// 	env.Parse(&c.ServerAddress)
// 	// }
// 	// if _, ok := os.LookupEnv("BASE_URL"); ok {
// 	env.Parse(c)(c)
// 	// }
// }

func GetServerAddress() string {
	return string(config.ServerAddress)
}

func GetBaseURL() string {
	return string(config.BaseURL)
}

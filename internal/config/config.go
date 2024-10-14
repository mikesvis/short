// Модуль кофигурации приложения.
package config

import (
	"log"

	"github.com/caarlos0/env"
	flag "github.com/spf13/pflag"
)

// Config конфиг приложения
type Config struct {
	// ServerAddress - адрес сервера приложения. По-умолчанию localhost:8080.
	ServerAddress string `env:"SERVER_ADDRESS"`

	// BaseURL - адрес сервера для коротких URL. По-умолчанию http://localhost:8080.
	BaseURL string `env:"BASE_URL"`

	// FileStoragePath - путь для файла storage, нужен при выборе движка хранения коротких ссылок в файле.
	FileStoragePath string `env:"FILE_STORAGE_PATH" envDefault:"/tmp/short-url-db.json"`

	// DatabaseDSN - адрес подключения к базе postgres, нужен при выборе движка хранения коротких ссылок в базе.
	DatabaseDSN string `env:"DATABASE_DSN"`

	// EnableHTTPS - использовать HTTPS на сервере
	EnableHTTPS bool `env:"ENABLE_HTTPS"`

	// ServerKeyPath - ключ для сертификата
	ServerKeyPath string `env:"SERVER_KEY_PATH"`

	// ServerCertPath - сертификат
	ServerCertPath string `env:"SERVER_CERT_PATH"`
}

// Конструктор конфигурации приложения.
func NewConfig() *Config {
	var config Config

	parseFlags(&config)
	env.Parse(&config)

	if !config.EnableHTTPS {
		return &config
	}

	if config.ServerKeyPath == "" || config.ServerCertPath == "" {
		log.Fatalf("Unable to run HTTPS mode without both server key and sertificate files set")
	}

	return &config
}

func parseFlags(c *Config) {
	flag.StringVarP(&c.ServerAddress, "address", "a", "localhost:8080", "address of shortener service server")
	flag.StringVarP(&c.BaseURL, "basepath", "b", "http://localhost:8080", "address of short link basepath")
	flag.StringVarP(&c.FileStoragePath, "file_storage_path", "f", "", "path to file storage of URLs")
	flag.StringVarP(&c.DatabaseDSN, "database_dsn", "d", "", "db connection string")
	flag.BoolVarP(&c.EnableHTTPS, "enable_https", "s", false, "use HTTPS connection")
	flag.StringVarP(&c.ServerKeyPath, "server_key_path", "k", "", "path to server key file")
	flag.StringVarP(&c.ServerCertPath, "server_cert_path", "c", "", "path to server certificate file")
	flag.Parse()
}

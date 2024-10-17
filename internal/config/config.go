// Модуль кофигурации приложения.
package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/caarlos0/env"
	flag "github.com/spf13/pflag"
)

// Config конфиг приложения
type Config struct {
	// ServerAddress - адрес сервера приложения. По-умолчанию localhost:8080.
	ServerAddress string `env:"SERVER_ADDRESS" json:"server_address"`

	// BaseURL - адрес сервера для коротких URL. По-умолчанию http://localhost:8080.
	BaseURL string `env:"BASE_URL" json:"base_url"`

	// FileStoragePath - путь для файла storage, нужен при выборе движка хранения коротких ссылок в файле.
	FileStoragePath string `env:"FILE_STORAGE_PATH" envDefault:"/tmp/short-url-db.json" json:"file_storage_path"`

	// DatabaseDSN - адрес подключения к базе postgres, нужен при выборе движка хранения коротких ссылок в базе.
	DatabaseDSN string `env:"DATABASE_DSN" json:"database_dsn"`

	// EnableHTTPS - использовать HTTPS на сервере
	EnableHTTPS bool `env:"ENABLE_HTTPS" json:"enable_https"`

	// ServerKeyPath - ключ для сертификата
	ServerKeyPath string `env:"SERVER_KEY_PATH" json:"server_key_path"`

	// ServerCertPath - сертификат
	ServerCertPath string `env:"SERVER_CERT_PATH" json:"server_cert_path"`

	// ConfigFilePath - путь к файлу конфига в формате json
	ConfigFilePath string `env:"CONFIG"`
}

// Конструктор конфигурации приложения.
func NewConfig() *Config {
	var config Config
	var configFile Config

	parseFlags(&config)
	env.Parse(&config)

	if len(config.ConfigFilePath) > 0 {
		parseFile(&configFile, config.ConfigFilePath)
	}

	if config.ServerAddress == "" && len(configFile.ServerAddress) > 0 {
		config.ServerAddress = configFile.ServerAddress
	}

	// setting default value if still empty
	if config.ServerAddress == "" {
		config.ServerAddress = "localhost:8080"
	}

	if config.BaseURL == "" && len(configFile.BaseURL) > 0 {
		config.BaseURL = configFile.BaseURL
	}

	// setting default value if still empty
	if config.BaseURL == "" {
		config.BaseURL = "http://localhost:8080"
	}

	if config.FileStoragePath == "" && len(configFile.FileStoragePath) > 0 {
		config.FileStoragePath = configFile.FileStoragePath
	}

	if config.DatabaseDSN == "" && len(configFile.DatabaseDSN) > 0 {
		config.DatabaseDSN = configFile.DatabaseDSN
	}

	if !config.EnableHTTPS && configFile.EnableHTTPS {
		config.EnableHTTPS = true
	}

	if config.ServerKeyPath == "" && len(configFile.ServerKeyPath) > 0 {
		config.ServerKeyPath = configFile.ServerKeyPath
	}

	if config.ServerCertPath == "" && len(configFile.ServerCertPath) > 0 {
		config.ServerCertPath = configFile.ServerCertPath
	}

	if !config.EnableHTTPS {
		return &config
	}

	if config.ServerKeyPath == "" || config.ServerCertPath == "" {
		log.Fatalf("Unable to run HTTPS mode without both server key and sertificate files set")
	}

	return &config
}

func parseFlags(c *Config) {
	flag.StringVarP(&c.ServerAddress, "address", "a", "", "address of shortener service server (default: localhost:8080)")
	flag.StringVarP(&c.BaseURL, "basepath", "b", "", "address of short link basepath (default: http://localhost:8080)")
	flag.StringVarP(&c.FileStoragePath, "file_storage_path", "f", "", "path to file storage of URLs")
	flag.StringVarP(&c.DatabaseDSN, "database_dsn", "d", "", "db connection string")
	flag.BoolVarP(&c.EnableHTTPS, "enable_https", "s", false, "use HTTPS connection")
	flag.StringVarP(&c.ServerKeyPath, "server_key_path", "k", "", "path to server key file")
	flag.StringVarP(&c.ServerCertPath, "server_cert_path", "e", "", "path to server certificate file")
	flag.StringVarP(&c.ConfigFilePath, "config", "c", "", "path to config file in json format")
	flag.Parse()
}

func parseFile(c *Config, fp string) {
	fmt.Println(c)
	file, err := os.Open(fp)
	if err != nil {
		log.Fatalf("Unable open config file %v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(c)
	if err != nil {
		log.Fatalf("Unable parse config file %v", err)
	}
}

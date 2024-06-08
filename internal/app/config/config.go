package config

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/pflag"
)

type ServerConfig struct {
	host string
	port int
}

type LinkConfig struct {
	scheme,
	host string
	port int
}

type Config struct {
	ServerConfig ServerConfig
	LinkConfig   LinkConfig
}

func (conf *ServerConfig) String() string {
	return fmt.Sprintf("%s:%d", conf.host, conf.port)
}

func (conf *ServerConfig) Type() string {
	return "string"
}

func (conf *ServerConfig) Set(flagValue string) error {
	s := strings.Split(flagValue, ":")
	if len(s) != 2 {
		return errors.New("server address shoud be <host:port>")
	}
	port, err := strconv.Atoi(s[1])
	if err != nil {
		return err
	}
	conf.host = s[0]
	conf.port = port
	return nil
}

func (conf *LinkConfig) String() string {
	return fmt.Sprintf("%s://%s:%d", conf.scheme, conf.host, conf.port)
}

func (conf *LinkConfig) Type() string {
	return "string"
}

func (conf *LinkConfig) Set(flagValue string) error {
	s := strings.Split(flagValue, ":")
	if len(s) != 3 {
		return errors.New("link address shoud be <scheme://host:port>")
	}
	if s[0] != "http" && s[0] != "https" {
		return fmt.Errorf("unknown scheme in link config %s", s[1])
	}
	port, err := strconv.Atoi(s[2])
	if err != nil {
		return err
	}
	conf.scheme = s[0]
	conf.host = strings.TrimLeft(s[1], "/")
	conf.port = port
	return nil
}

var config Config = Config{
	ServerConfig: ServerConfig{
		host: "localhost",
		port: 8080,
	},
	LinkConfig: LinkConfig{
		scheme: "http",
		host:   "localhost",
		port:   8080,
	},
}

func init() {
	pflag.VarP(&config.ServerConfig, "address", "a", "address of shortener service server")
	pflag.VarP(&config.LinkConfig, "basepath", "b", "address of short link basepath")
}

func InitConfig() {
	pflag.Parse()
}

func GetServerHostAddr() string {
	return config.ServerConfig.String()
}

func GetShortLinkAddr() string {
	return config.LinkConfig.String()
}

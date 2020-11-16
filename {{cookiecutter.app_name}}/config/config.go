package config

import (
	"github.com/kelseyhightower/envconfig"
)

var defaultConfig Config

type Config struct {
	// App configuration
	GRPCPort int    `envconfig:"GRPC_PORT" default:"9090"`
	HTTPPort int    `envconfig:"HTTP_PORT" default:"9091"`
	AppName  string `envconfig:"APP_NAME" default:"{{cookiecutter.app_name}}"`
	LogLevel string `envconfig:"LOG_LEVEL" default:"debug"`
	JSONLogs bool   `envconfig:"JSON_LOGS" default:"true"`
	Prefix   string `envconfig:"PREFIX" default:"got"`
}

func init() {
	envconfig.Process("", &defaultConfig)
	// fail on error
}

func Get() Config {
	return defaultConfig
}

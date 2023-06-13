package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

var Config Configuration

type Configuration struct {
	Environment EnvironmentConfig
	Server      ServerConfig
}

type EnvironmentConfig struct {
	Environment string
	LogLevel    string
}

type ServerConfig struct {
	Host string
	Port int
}

type variablesKeys struct {
	envPath  string
	logLevel string
	host     string
	port     string
}

func init() {
	keys := setVariablesKeys()
	env := getEnvironment()
	vr := viper.New()

	vr.SetConfigFile(fmt.Sprintf(keys.envPath, env))
	if err := vr.ReadInConfig(); err != nil {
		panic("environment file not found")
	}

	vr.SetDefault(keys.logLevel, "error")
	vr.SetDefault(keys.host, "localhost")
	vr.SetDefault(keys.port, "8080")

	Config = Configuration{
		Environment: EnvironmentConfig{
			Environment: env,
			LogLevel:    vr.GetString(keys.logLevel),
		},
		Server: ServerConfig{
			Host: vr.GetString(keys.host),
			Port: vr.GetInt(keys.port),
		},
	}
}

func setVariablesKeys() variablesKeys {
	return variablesKeys{
		envPath:  "./environment/%s.env",
		logLevel: "LOG_LEVEL",
		host:     "HOST",
		port:     "PORT",
	}
}

func getEnvironment() string {
	if value := os.Getenv("ENV"); value != "" {
		return value
	}
	return "local"
}

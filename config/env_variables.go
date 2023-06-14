package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"time"
)

var Config Configuration

type Configuration struct {
	Environment EnvironmentConfig
	Server      ServerConfig
	Token       TokenConfig
	Redis       RedisConfig
}

type EnvironmentConfig struct {
	Environment string
	LogLevel    string
}

type ServerConfig struct {
	Host string
	Port int
}

type TokenConfig struct {
	TokenDuration time.Duration
	SecretKeyPath string
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
}

type variablesKeys struct {
	envPath       string
	logLevel      string
	host          string
	port          string
	duration      string
	secretKeyPath string
	redisHost     string
	redisPort     string
	redisPassword string
}

func init() {
	keys := setVariablesKeys()
	env := getEnvironment()
	vr := viper.New()

	vr.SetConfigFile(fmt.Sprintf(keys.envPath, env))
	_ = vr.ReadInConfig()

	vr.SetDefault(keys.logLevel, "error")
	vr.SetDefault(keys.host, "localhost")
	vr.SetDefault(keys.port, 8080)
	vr.SetDefault(keys.duration, "24h")
	vr.SetDefault(keys.secretKeyPath, "./test_resources/mock_RS256.key")
	vr.SetDefault(keys.redisHost, "localhost")
	vr.SetDefault(keys.redisPort, 6379)
	vr.SetDefault(keys.redisPassword, "")

	tokenDuration, err := time.ParseDuration(vr.GetString(keys.duration))
	if err != nil {
		tokenDuration, _ = time.ParseDuration("24h")
	}

	Config = Configuration{
		Environment: EnvironmentConfig{
			Environment: env,
			LogLevel:    vr.GetString(keys.logLevel),
		},
		Server: ServerConfig{
			Host: vr.GetString(keys.host),
			Port: vr.GetInt(keys.port),
		},
		Token: TokenConfig{
			TokenDuration: tokenDuration,
			SecretKeyPath: vr.GetString(keys.secretKeyPath),
		},
		Redis: RedisConfig{
			Host:     vr.GetString(keys.redisHost),
			Port:     vr.GetInt(keys.redisPort),
			Password: vr.GetString(keys.redisPassword),
		},
	}
}

func setVariablesKeys() variablesKeys {
	return variablesKeys{
		envPath:       "./environment/%s.env",
		logLevel:      "LOG_LEVEL",
		host:          "HOST",
		port:          "PORT",
		duration:      "DURATION",
		secretKeyPath: "SECRET_KEY_PATH",
		redisHost:     "REDIS_HOST",
		redisPort:     "REDIS_PORT",
		redisPassword: "REDIS_PASSWORD",
	}
}

func getEnvironment() string {
	if value := os.Getenv("ENV"); value != "" {
		return value
	}
	return "local"
}

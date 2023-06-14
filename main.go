package main

import (
	"os"
	"rkaneko/endava-coding-exercise/api"
	"rkaneko/endava-coding-exercise/config"
	"rkaneko/endava-coding-exercise/db"
)

var service *api.Service

func init() {
	service = &api.Service{
		Config:        config.Config.Server,
		Log:           config.InitLogrus(config.Config.Environment.LogLevel),
		TokenDuration: config.Config.Token.TokenDuration,
		SecretKeyPath: config.Config.Token.SecretKeyPath,
	}

	rds, err := db.RunRedis(
		config.Config.Redis.Password,
		config.Config.Redis.Host,
		config.Config.Redis.Port,
	)
	if err != nil {
		service.Log.Errorf("there was an error running redis: %s", err.Error())
		os.Exit(1)
	}
	service.RedisService = &db.RedisService{
		Client: rds,
	}
}

func main() {
	service.Log.Infof("Running redis in %s:%d", service.Config.Host, 6379)
	service.Log.Infof("Running server in %s:%d", service.Config.Host, service.Config.Port)
	service.Init()
}

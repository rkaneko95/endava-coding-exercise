package main

import (
	"rkaneko/endava-coding-exercise/api"
	"rkaneko/endava-coding-exercise/config"
)

var service *api.Service

func init() {
	service = &api.Service{
		Config: config.Config.Server,
		Log:    config.InitLogrus(config.Config.Environment.LogLevel),
	}
}

func main() {
	service.Log.Infof("Running server %s:%d", service.Config.Host, service.Config.Port)
	service.Init()
}

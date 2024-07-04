package main

import (
	"context"
	"encoding/json"

	myconfig "rummy-session/config"
	"rummy-session/controller"
	"rummy-session/repository"
	"rummy-session/repository/redismiddleware"
	"rummy-session/service"

	"bitbucket.org/junglee_games/getsetgo/logger"
	"bitbucket.org/junglee_games/getsetgo/monitoring/monitoringfactory"
)

// type AppConfig struct {
// 	App string
// }

//	func (c *AppConfig) GetBuild() string {
//		return "local"
//	}
func main() {
	ctx := context.Background()

	// var appconfig AppConfig
	// if err := config.LoadConfig("application", "./", &appconfig); err != nil {
	// logger.Panic(ctx, "failed to load a config, %v", err.Error())
	// }
	conf := myconfig.GetConfig(ctx, "application", "./")
	logger.Config{AppName: conf.Name, Build: conf.Build, Level: logger.LogLevel(conf.LogLevel)}.InitiateLogger()
	bytes, err := json.MarshalIndent(conf, "", "    ")
	if err != nil {
		logger.Panic(ctx, "failed to marshal  config: %v", err.Error())
	}
	logger.Debug(ctx, "config : %s", bytes)
	monitoringAgent, err := monitoringfactory.GetMonitoringAgent(conf.Monitoring)
	if err != nil {
		logger.Panic(ctx, "creating monitoring agent : %v", err.Error())
	}
	repository := repository.NewRepository(redismiddleware.Config{
		Addr:     conf.Repository.Addr,
		Password: conf.Repository.Password,
	}, monitoringAgent)
	service := service.NewService(repository, conf.Monitoring)
	controller := controller.NewController(conf.Controller, service)
	controller.StartListening()
}

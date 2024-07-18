package config

import (
	"context"
	"errors"
	ac "rummy-session/business/config"
	"rummy-session/consul"
	"rummy-session/controller"
	"rummy-session/repository/redismiddleware"
	"sync"

	"bitbucket.org/junglee_games/getsetgo/config"
	"bitbucket.org/junglee_games/getsetgo/configs"
	"bitbucket.org/junglee_games/getsetgo/logger"
)

var (
	ErrInvalidConfig = errors.New("invalid config")
)

type Config struct {
	Application *ac.AppConf
	Repository  redismiddleware.Config
	Controller  controller.Config
	LogLevel    string
	ConsulURI   string
	ConsulToken string
	KVPath      string
	Name        string
	Build       string
	Monitoring  *configs.DefaultMonitoringConfig
}

func (c *Config) GetBuild() string {
	return c.Build
}

func GetConfig(ctx context.Context, env, path string) *Config {
	var conf Config
	var once sync.Once
	once.Do(func() {
		err := config.LoadConfig(env, path, &conf)
		if err != nil {
			logger.Panic(ctx, "reading config : %v", err.Error())
		}
	})
	ag, err := consul.New(conf.KVPath, consul.Config{
		Address: conf.ConsulURI,
		Name:    conf.Name,
		Token:   conf.ConsulToken,
	})
	if err != nil {
		logger.Panic(ctx, "error in creating consul agent", err.Error())

	}
	var appConf = ac.New()
	err = ag.InitAndGetConfig(appConf)
	if err != nil {
		// logger.Panic(ctx, "while getting config from consul", err.Error())
		// bypassing consul configuration
		appConf = &ac.AppConf{
			SessionTokenLength: 5,
		}
	}
	conf.Application = appConf
	return &conf
}

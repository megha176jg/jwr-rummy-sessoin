package controller

import (
	"rummy-session/controller/http"
	"rummy-session/service"
)

type Controller interface {
	StartListening() error
}

type Config struct {
	HTTP http.HttpConfig
}

func NewController(c Config, s service.Service) Controller {
	return http.NewHttpController(c.HTTP, s)
}

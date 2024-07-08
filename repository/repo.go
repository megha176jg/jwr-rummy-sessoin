package repository

import (
	"rummy-session/repository/models"
	"rummy-session/repository/redismiddleware"

	"bitbucket.org/junglee_games/getsetgo/monitoring"
)

var (
	ERRORNOTFOUND = "invalid name"
)

type Repository interface {
	GetAuthToken(name string) models.AuthToken
	DeleteAuthToken(name string) error
	CreateAuthToken(name, value string) error
}

func NewRepository(c redismiddleware.Config, ma monitoring.Agent) Repository {
	return redismiddleware.NewRedisRepository(c, ma)
}

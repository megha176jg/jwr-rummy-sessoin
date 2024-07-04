package repository

import (
	"rummy-session/repository/redismiddleware"

	"bitbucket.org/junglee_games/getsetgo/monitoring"
)

var (
	ERRORNOTFOUND = "invalid name"
)

type Repository interface {
	GetTitle(name string) (string, error)
}

func NewRepository(c redismiddleware.Config, ma monitoring.Agent) Repository {
	return redismiddleware.NewRedisRepository(c, ma)
}

package service

import (
	"rummy-session/repository" // Add the import statement for the "repository" package

	"bitbucket.org/junglee_games/getsetgo/configs"
	"github.com/gin-gonic/gin"
)

type Service interface {
	Login(ctx *gin.Context)
}

type service struct {
	repo             repository.Repository
	monitoringConfig *configs.DefaultMonitoringConfig
}

type result struct {
	Name string
	Age  string
}

func NewService(repo repository.Repository, m *configs.DefaultMonitoringConfig) *service {
	return &service{
		repo:             repo,
		monitoringConfig: m,
	}
}

// @Summary		Login with authtoken
// @Description	get authtoken by
// @Accept			*/*
// @Produce		json
// @Param			userId		query		string	true	"userid"
// @Param			authToken	query		string	true	"authToken"
// @Success		200			{string}	string	"ok"
// @Router			/greet/ [get]
func (s *service) Login(ctx *gin.Context) {

	userid, ok := ctx.GetQuery("userId")
	cmauthtoken, ok := ctx.GetQuery("authToken")
	if ok != true {
		ctx.JSON(400, gin.H{"error": "name is required"})
		return
	}
	authtoken, err := s.repo.GetTitle(userid)
	if err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	// ctx.JSON(200, gin.H{"message": "Hello, " + userid + " " + title})
	if cmauthtoken == authtoken {
		ctx.JSON(200, gin.H{"message": "Hello, " + userid + " " + authtoken})
	} else {
		ctx.JSON(400, gin.H{"error": "Invalid Token"})
	}
}

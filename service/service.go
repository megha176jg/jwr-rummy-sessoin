package service

import (
	"net/http"
	"rummy-session/repository" // Add the import statement for the "repository" package

	"rummy-session/service/model"

	"bitbucket.org/junglee_games/getsetgo/configs"
	"github.com/gin-gonic/gin"
)

type Service interface {
	Validate(ctx *gin.Context)
	Invalidate(ctx *gin.Context)
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
// @Router			/api/v1/session/user/{user_id}/validate [get]
func (s *service) Validate(ctx *gin.Context) {

	userid, ok := ctx.GetQuery("userId")
	if ok != true {
		ctx.JSON(400, gin.H{"error": "userid is required"})
		return
	}
	cmauthtoken, ok := ctx.GetQuery("authToken")
	if ok != true {
		ctx.JSON(400, gin.H{"error": "authtoken is required"})
		return
	}
	result := s.repo.GetAuthToken(userid)
	if result.Error != nil {
		ErrResponse := model.Validate{
			Err: result.Error,
		}
		ctx.JSON((http.StatusNotFound), ErrResponse)
		return
	}
	if cmauthtoken == result.AuthToken {
		Response := model.Validate{
			IsTokenValid: true,
			Err:          nil,
		}
		ctx.JSON(http.StatusOK, Response)
		return
	} else {
		Response := model.Validate{
			IsTokenValid: false,
			Err:          nil,
		}
		ctx.JSON(http.StatusOK, Response)
		return
	}
}

// @Summary		Logout with authtoken
// @Description	get authtoken by
// @Accept			*/*
// @Produce		json
// @Param			userId		query		string	true	"userid"
// @Param			authToken	query		string	true	"authToken"
// @Success		200			{string}	string	"ok"
// @Router			/api/v1/session/user/{user_id}/invalidate [delete]
func (s *service) Invalidate(ctx *gin.Context) {
	userid, ok := ctx.GetQuery("userId")
	if ok != true {
		ctx.JSON(400, gin.H{"error": "userid is required"})
		return
	}
	cmauthtoken, ok := ctx.GetQuery("authToken")
	if ok != true {
		ctx.JSON(400, gin.H{"error": "authtoken is required"})
		return
	}
	result := s.repo.GetAuthToken(userid)
	if result.Error != nil {
		ErrResponse := model.Validate{
			Err: result.Error,
		}
		ctx.JSON((http.StatusNotFound), ErrResponse)
		return
	}
	if cmauthtoken == result.AuthToken {
		res := s.repo.DeleteAuthToken(userid)
		ctx.JSON(http.StatusOK, res)
		return
	} else {
		Response := model.Validate{
			IsTokenValid: false,
			Err:          nil,
		}
		ctx.JSON(http.StatusOK, Response)
		return
	}
}

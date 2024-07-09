package service

import (
	"math/rand"
	"net/http"
	"rummy-session/repository" // Add the import statement for the "repository" package

	"rummy-session/service/model"

	"bitbucket.org/junglee_games/getsetgo/configs"
	"github.com/gin-gonic/gin"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type Service interface {
	Create(ctx *gin.Context)
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

// @Summary		Creation of authtoken
// @Description	get authtoken by
// @Accept			*/*
// @Produce		json
// @Param			userId		path		string	true	"userid"
// @Success		200			{string}	string	"ok"
// @Router			/api/v1/session/user/{user_id} [get]
func (s *service) Create(ctx *gin.Context) {
	userid := ctx.Param("user_id")
	err := s.repo.CreateAuthToken(userid, "Bearer "+RandStringRunes(20))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	result := s.repo.GetAuthToken(userid)
	if result.Error != nil {
		ErrResponse := model.Validate{
			Err: result.Error,
		}
		ctx.JSON((http.StatusNotFound), ErrResponse)
		return
	} else {
		ctx.JSON(http.StatusOK, result)
		return
	}
}

// @Summary		Login with authtoken
// @Description	get authtoken by
// @Accept			*/*
// @Produce		json
// @Param			userId		path		string	true	"userid"
// @Param			authToken	path		string	true	"authToken"
// @Success		200			{string}	string	"ok"
// @Router			/api/v1/session/user/validate [get]
func (s *service) Validate(ctx *gin.Context) {

	userid := ctx.Param("user_id")
	cmauthtoken := ctx.GetHeader("Authorization")

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
// @Accept		*/*
// @Produce		json
// @Param		userId		query		string	true	"userid"
// @Param		authToken	query		string	true	"authToken"
// @Success		200			{string}	string	"ok"
// @Router		/api/v1/session/user/invalidate [delete]
func (s *service) Invalidate(ctx *gin.Context) {
	userid := ctx.Param("user_id")
	cmauthtoken := ctx.GetHeader("Authorization")

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
		deleteResult := model.Invalidate{
			Err: res,
		}
		ctx.JSON(http.StatusOK, deleteResult)
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

func RandStringRunes(n int) string {
	b := make([]byte, n)
	for i := range b {
		k := rand.Intn(len(letterBytes))
		b[i] = letterBytes[k]
	}
	return string(b)
}

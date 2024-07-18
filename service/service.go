package service

import (
	"math/rand"
	"net/http"
	"rummy-session/repository" // Add the import statement for the "repository" package
	"rummy-session/service/model"
	"strings"

	ac "rummy-session/business/config"

	"bitbucket.org/junglee_games/getsetgo/configs"
	"bitbucket.org/junglee_games/getsetgo/logger"
	"bitbucket.org/junglee_games/getsetgo/monitoring/monitoringfactory"
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
	appConf          ac.AppConf
}

type result struct {
	Name string
	Age  string
}

func NewService(repo repository.Repository, m *configs.DefaultMonitoringConfig, a ac.AppConf) *service {
	return &service{
		repo:             repo,
		monitoringConfig: m,
		appConf:          a,
	}
}

func (s *service) Create(ctx *gin.Context) {
	agent, err := monitoringfactory.GetMonitoringAgent(s.monitoringConfig)
	if err != nil {
		panic(err)
	}
	defer agent.StartTransaction("Create").End()
	userid := ctx.Param("user_id")
	if s.repo.GetAuthToken(userid).Error != nil {
		err := s.repo.CreateAuthToken(userid, RandStringRunes((s.appConf.SessionTokenLength)))
		if err != nil {
			logger.Error(ctx, ErrCreatingToken, err.Error())
			ctx.JSON(http.StatusInternalServerError, model.TokenCreate{
				model.Err{
					Message: err.Error(),
				},
			})
			return
		}
		result := s.repo.GetAuthToken(userid)
		if result.Error != nil {
			logger.Error(ctx, ErrGetAuthToken, result.Error.Error())
			ErrResponse := model.TokenGet{
				Error: model.Err{
					Message: result.Error.Error(),
				},
			}
			ctx.JSON((http.StatusNotFound), ErrResponse)
			return
		} else {
			ctx.JSON(http.StatusOK, model.TokenGet{
				AuthToken: result.AuthToken,
			})
			return
		}
	} else {
		logger.Error(ctx, ErrTokenAlreadyExists, "Token already exists")
		ctx.JSON(http.StatusOK, model.TokenGet{
			AuthToken: s.repo.GetAuthToken(userid).AuthToken,
		})
		return

	}
}

func (s *service) Validate(ctx *gin.Context) {
	agent, err := monitoringfactory.GetMonitoringAgent(s.monitoringConfig)
	if err != nil {
		panic(err)
	}
	defer agent.StartTransaction("Validate").End()
	userid := ctx.Param("user_id")
	cmauthorization := ctx.GetHeader("Authorization")
	authorizationparts := strings.Split(cmauthorization, " ")
	if len(authorizationparts) != 2 || authorizationparts[0] != "Bearer" {
		logger.Error(ctx, ErrInvalidAuthHeader, "Invalid Authorization Header")
		ErrResponse := model.Err{
			Message: "Invalid Authorization Header",
		}
		ctx.JSON((http.StatusNotFound), ErrResponse)
		return
	}
	result := s.repo.GetAuthToken(userid)
	if result.Error != nil {
		logger.Error(ctx, ErrGetAuthToken, result.Error.Error())
		ErrResponse := model.TokenGet{
			Error: model.Err{
				Message: result.Error.Error(),
			},
		}
		ctx.JSON((http.StatusNotFound), ErrResponse)
		return
	}
	if authorizationparts[1] == result.AuthToken {
		Response := model.TokenValidation{
			IsTokenValid: true,
		}
		ctx.JSON(http.StatusOK, Response)
		return
	} else {
		logger.Error(ctx, ErrIncorrectAuthtoken, "Incorrect Authtoken")
		Response := model.TokenValidation{
			IsTokenValid: false,
			Error: model.Err{
				Message: "Incorrect Auth Token",
			},
		}
		ctx.JSON(http.StatusOK, Response)
		return
	}
}

func (s *service) Invalidate(ctx *gin.Context) {
	agent, err := monitoringfactory.GetMonitoringAgent(s.monitoringConfig)
	if err != nil {
		panic(err)
	}
	defer agent.StartTransaction("Invalidate").End()
	userid := ctx.Param("user_id")
	if s.repo.GetAuthToken(userid).Error != nil {
		logger.Error(ctx, ErrInvalidUserId, "userid does not exist")
		ErrResponse := model.Err{
			Message: "Invalid User Id",
		}
		ctx.JSON((http.StatusNotFound), ErrResponse)
	} else {
		cmauthorization := ctx.GetHeader("Authorization")
		authorizationparts := strings.Split(cmauthorization, " ")
		if len(authorizationparts) != 2 || authorizationparts[0] != "Bearer" {
			ErrResponse := model.Err{
				Message: "Invalid Authorization Header",
			}
			ctx.JSON((http.StatusNotFound), ErrResponse)
			return
		}
		result := s.repo.GetAuthToken(userid)
		if result.Error != nil {
			ErrResponse := model.TokenGet{
				Error: model.Err{
					Message: result.Error.Error(),
				},
			}
			ctx.JSON((http.StatusNotFound), ErrResponse)
			return
		}

		if authorizationparts[1] == result.AuthToken {
			res := s.repo.DeleteAuthToken(userid)
			if res != nil {
				deleteResult := model.TokenDeletion{
					Error: model.Err{
						Message: res.Error(),
					},
				}
				ctx.JSON(http.StatusOK, deleteResult)
				return
			} else {
				deleteResult := model.TokenDeletion{}
				ctx.JSON(http.StatusOK, deleteResult)
				return

			}
		} else {
			Response := model.TokenValidation{
				IsTokenValid: false,
				Error: model.Err{
					Message: "Incorrect Auth Token",
				},
			}
			ctx.JSON(http.StatusOK, Response)
			return
		}
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

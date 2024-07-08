package http

import (
	"net/http"
	_ "rummy-session/docs"
	"rummy-session/service"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type HttpConfig struct {
	Port                string
	TimeOutInSeconds    int
	MaxIdleConns        int
	MaxIdleConnsPerHost int
}

type HttpController struct {
	config  HttpConfig
	service service.Service
}

func NewHttpController(config HttpConfig, s service.Service) *HttpController {
	return &HttpController{
		config:  config,
		service: s,
	}
}

func (c *HttpController) StartListening() error {
	router := gin.Default()
	router.GET("/api/v1/session/user/validate", func(ctx *gin.Context) {
		c.service.Validate(ctx)
	})
	router.DELETE("/api/v1/session/user/invalidate", func(ctx *gin.Context) {
		c.service.Invalidate(ctx)
	})
	router.GET("/rummysession/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	err := http.ListenAndServe(":"+c.config.Port, router)
	if err != nil {
		return err
	}
	return nil
}

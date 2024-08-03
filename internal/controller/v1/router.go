package v1

import (
	"log/slog"
	"net/http"

	"github.com/romanchechyotkin/avito_test_task/internal/controller/v1/middleware"
	"github.com/romanchechyotkin/avito_test_task/internal/service"

	"github.com/gin-gonic/gin"
)

func NewRouter(log *slog.Logger, router *gin.Engine, services *service.Services) {
	router.Use(middleware.CORS())
	router.Use(middleware.Log(log))

	router.GET("/status", func(c *gin.Context) {
		c.String(http.StatusOK, "ok\n")
	})

	authGroup := router.Group("/auth")
	{
		newAuthRoutes(log, authGroup, services.Auth)
	}

	houseGroup := router.Group("/house")
	{
		newHouseRoutes(log, houseGroup, services.House, services.Auth)
	}

}

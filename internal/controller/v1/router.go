package v1

import (
	"log/slog"
	"net/http"

	_ "github.com/romanchechyotkin/avito_test_task/docs"
	"github.com/romanchechyotkin/avito_test_task/internal/controller/v1/middleware"
	"github.com/romanchechyotkin/avito_test_task/internal/service"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func NewRouter(log *slog.Logger, router *gin.Engine, services *service.Services) {
	router.Use(middleware.Log(log))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	router.GET("/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	authMiddleware := middleware.NewAuthMiddleware(services.Auth)

	authGroup := router.Group("/auth")
	{
		newAuthRoutes(log, authGroup, services.Auth)
	}

	v1 := router.Group("/v1")
	{
		newHouseRoutes(log, v1.Group("/house"), services.House, authMiddleware)
		newFlatRoutes(log, v1.Group("/flat"), services.Flat, authMiddleware)
	}

}

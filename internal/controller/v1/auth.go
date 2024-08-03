package v1

import (
	"log/slog"
	"net/http"

	"github.com/romanchechyotkin/avito_test_task/internal/controller/v1/request"
	"github.com/romanchechyotkin/avito_test_task/internal/service"
	"github.com/romanchechyotkin/avito_test_task/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type authRoutes struct {
	log *slog.Logger

	authService service.Auth
}

func newAuthRoutes(log *slog.Logger, g *gin.RouterGroup, authService service.Auth) {
	r := &authRoutes{
		log:         log,
		authService: authService,
	}

	g.POST("/registration", r.Registration)
	g.POST("/login", r.Login)
}

func (r *authRoutes) Registration(c *gin.Context) {
	var req request.Registration

	if err := c.ShouldBindJSON(&req); err != nil {
		r.log.Error("failed to read request data", logger.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	// todo custom validator with russian responses
	if err := validator.New().Struct(req); err != nil {
		r.log.Error("failed to validate request data", logger.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusCreated, req)
}

func (r *authRoutes) Login(c *gin.Context) {

}

package v1

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/romanchechyotkin/avito_test_task/internal/controller/v1/request"
	"github.com/romanchechyotkin/avito_test_task/internal/controller/v1/response"
	"github.com/romanchechyotkin/avito_test_task/internal/repo/repoerrors"
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

	g.POST("/register", r.Registration)
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

	userID, err := r.authService.CreateUser(c, &service.AuthCreateUserInput{
		Email:    req.Email,
		Password: req.Password,
		UserType: req.UserType,
	})
	if err != nil {
		if errors.Is(err, repoerrors.ErrUserExists) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		r.log.Error("failed to create user", logger.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusCreated, response.Registration{
		UserID: userID,
	})
}

func (r *authRoutes) Login(c *gin.Context) {
	var req request.Login

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

	token, err := r.authService.GenerateToken(c, &service.AuthGenerateTokenInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		if errors.Is(err, service.ErrWrongPassword) || errors.Is(err, repoerrors.ErrUserNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		r.log.Error("failed to generate user token", logger.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, response.Login{
		Token: token,
	})
}

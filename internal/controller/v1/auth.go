package v1

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/romanchechyotkin/avito_test_task/internal/controller/v1/request"
	"github.com/romanchechyotkin/avito_test_task/internal/controller/v1/response"
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
	log = log.With(slog.String("component", "auth routes"))

	r := &authRoutes{
		log:         log,
		authService: authService,
	}

	g.POST("/register", r.registration)
	g.POST("/login", r.login)
}

// @Summary Registration
// @Description Registration
// @Tags auth
// @Accept json
// @Produce json
// @Param input body request.Registration true "input"
// @Success 201 {object} response.Registration
// @Router /auth/register [post]
func (r *authRoutes) registration(c *gin.Context) {
	var req request.Registration

	if err := c.ShouldBindJSON(&req); err != nil {
		r.log.Error("failed to read request data", logger.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

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
		if errors.Is(err, service.ErrUserExists) {
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

// @Summary Login
// @Description Login
// @Tags auth
// @Accept json
// @Produce json
// @Param input body request.Login true "input"
// @Success 201 {object} response.Login
// @Router /auth/login [post]
func (r *authRoutes) login(c *gin.Context) {
	var req request.Login

	if err := c.ShouldBindJSON(&req); err != nil {
		r.log.Error("failed to read request data", logger.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

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
		if errors.Is(err, service.ErrWrongPassword) || errors.Is(err, service.ErrUserNotFound) {
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

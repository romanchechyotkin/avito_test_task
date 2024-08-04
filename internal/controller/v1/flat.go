package v1

import (
	"log/slog"
	"net/http"

	"github.com/romanchechyotkin/avito_test_task/internal/controller/v1/middleware"
	"github.com/romanchechyotkin/avito_test_task/internal/controller/v1/request"
	"github.com/romanchechyotkin/avito_test_task/internal/service"
	"github.com/romanchechyotkin/avito_test_task/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type flatRoutes struct {
	log *slog.Logger

	flatService service.Flat
}

func newFlatRoutes(log *slog.Logger, g *gin.RouterGroup, flatService service.Flat, authMiddleware *middleware.AuthMiddleware) {
	r := &flatRoutes{
		log:         log,
		flatService: flatService,
	}

	g.POST("/create", authMiddleware.AuthOnly(), r.createFlat)
	g.PATCH("/update", authMiddleware.ModeratorsOnly(), r.updateFlat)
}

func (r *flatRoutes) createFlat(c *gin.Context) {
	var req request.CreateFlat

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

	flat, err := r.flatService.CreateFlat(c, &service.FlatCreateInput{
		Number:      req.Number,
		HouseID:     req.HouseID,
		Price:       req.Price,
		RoomsAmount: req.RoomsAmount,
	})
	if err != nil {
		r.log.Error("failed to create flat", logger.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusCreated, flat)
}

func (r *flatRoutes) updateFlat(c *gin.Context) {
	var req request.UpdateFlat

	userID, ok := c.Get("userID")
	if !ok {
		r.log.Error("failed to get key from context", slog.String("key", "userType"))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get key from context",
		})

		return
	}

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

	flat, err := r.flatService.UpdateFlat(c, &service.FlatUpdateInput{
		ID:          req.ID,
		Status:      req.Status,
		ModeratorID: userID.(string),
	})
	if err != nil {
		r.log.Error("failed to update flat status", logger.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, flat)
}

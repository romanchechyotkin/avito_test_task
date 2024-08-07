package v1

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/romanchechyotkin/avito_test_task/internal/controller/v1/middleware"
	"github.com/romanchechyotkin/avito_test_task/internal/controller/v1/request"
	"github.com/romanchechyotkin/avito_test_task/internal/controller/v1/response"
	"github.com/romanchechyotkin/avito_test_task/internal/service"
	"github.com/romanchechyotkin/avito_test_task/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type houseRoutes struct {
	log *slog.Logger

	houseService service.House
}

func newHouseRoutes(log *slog.Logger, g *gin.RouterGroup, houseService service.House, authMiddleware *middleware.AuthMiddleware) {
	r := &houseRoutes{
		log:          log,
		houseService: houseService,
	}

	g.POST("/create", authMiddleware.ModeratorsOnly(), r.createHouse)
	g.GET("/:id", authMiddleware.AuthOnly(), r.getHouseFlats)
	g.POST("/:id/subscribe", authMiddleware.ClientsOnly(), r.subscribe)
}

// @Summary Create House
// @Description Create House
// @Tags house
// @Accept json
// @Produce json
// @Param input body request.CreateHouse true "input"
// @Success 201 {object} response.House
// @Security JWT
// @Router /v1/house/create [post]
func (r *houseRoutes) createHouse(c *gin.Context) {
	var req request.CreateHouse

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

	house, err := r.houseService.CreateHouse(c, &service.HouseCreateInput{
		Address:   req.Address,
		Year:      req.Year,
		Developer: req.Developer,
	})
	if err != nil {
		if errors.Is(err, service.ErrHouseExists) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		r.log.Error("failed to create house", logger.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusCreated, response.BuildHouse(house))
}

// @Summary Get House Flats
// @Description Get House Flats
// @Tags house
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} response.HouseFlats
// @Security JWT
// @Router /v1/house/{id} [get]
func (r *houseRoutes) getHouseFlats(c *gin.Context) {
	houseID := c.Param("id")

	userType, ok := c.Get("userType")
	if !ok {
		r.log.Error("failed to get key from context", slog.String("key", "userType"))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get key from context",
		})

		return
	}

	flats, err := r.houseService.GetHouseFlats(c, &service.GetHouseFlatsInput{
		UserType: userType.(string),
		HouseID:  houseID,
	})
	if err != nil {
		if errors.Is(err, service.ErrHouseFlatsNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "no flats found",
			})

			return
		}

		r.log.Error("failed to get house flats", logger.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, response.HouseFlats{Flats: flats})
}

// @Summary Subscribe for house updates
// @Description Subscribe for house updates
// @Tags house
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Success 204
// @Security JWT
// @Router /v1/house/{id}/subscribe [post]
func (r *houseRoutes) subscribe(c *gin.Context) {
	houseID := c.Param("id")

	userID, ok := c.Get("userID")
	if !ok {
		r.log.Error("failed to get key from context", slog.String("key", "userType"))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to get key from context",
		})

		return
	}

	if err := r.houseService.CreateSubscription(c, &service.CreateSubscriptionInput{
		HouseID: houseID,
		UserID:  userID.(string),
	}); err != nil {
		if errors.Is(err, service.ErrHouseNotFound) || errors.Is(err, service.ErrHouseSubscriptionExists) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		r.log.Error("failed to create subscription", slog.String("house id", houseID), slog.String("userid", userID.(string)))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusNoContent, nil)
}

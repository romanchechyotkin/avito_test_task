package v1

import (
	"github.com/romanchechyotkin/avito_test_task/internal/controller/v1/request"
	"github.com/romanchechyotkin/avito_test_task/internal/controller/v1/response"
	"github.com/romanchechyotkin/avito_test_task/internal/service"
	"github.com/romanchechyotkin/avito_test_task/pkg/logger"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type houseRoutes struct {
	log *slog.Logger

	houseService service.House
}

func newHouseRoutes(log *slog.Logger, g *gin.RouterGroup, houseService service.House) {
	r := &houseRoutes{
		log:          log,
		houseService: houseService,
	}

	g.POST("/create", r.createHouse)
}

func (r *houseRoutes) createHouse(c *gin.Context) {
	var req request.CreateHouse

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

	house, err := r.houseService.CreateHouse(c, &service.HouseCreateInput{
		Address:   req.Address,
		Year:      req.Year,
		Developer: req.Developer,
	})
	if err != nil {
		r.log.Error("failed to create house", logger.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusCreated, response.BuildHouse(house))
}

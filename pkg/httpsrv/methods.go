package httpsrv

import (
	"github.com/romanchechyotkin/avito_test_task/pkg/httpsrv/request"
	"net/http"

	"github.com/romanchechyotkin/avito_test_task/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func (srv *Server) Registration(c *gin.Context) {
	var req request.Registration

	if err := c.ShouldBindJSON(&req); err != nil {
		srv.log.Error("failed to read request data", logger.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	// todo custom validator with russian responses
	if err := validator.New().Struct(req); err != nil {
		srv.log.Error("failed to validate request data", logger.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusCreated, req)
}

func (srv *Server) Login(c *gin.Context) {

}

package v1

import (
	"io"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestMain(t *testing.M) {
	gin.SetMode(gin.TestMode)
	gin.DefaultWriter = io.Discard
}

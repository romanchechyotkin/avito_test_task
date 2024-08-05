package utils

import (
	"io"
	"sync"

	"github.com/gin-gonic/gin"
)

func NewTestRouter() *gin.Engine {
	var mu sync.Mutex

	mu.Lock()
	defer mu.Unlock()

	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.Discard

	return gin.New()
}

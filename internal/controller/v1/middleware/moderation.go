package middleware

import (
	"net/http"
	"strings"

	"github.com/romanchechyotkin/avito_test_task/internal/service"

	"github.com/gin-gonic/gin"
)

type authMiddleware struct {
	authService service.Auth
}

func NewAuthMiddleware(authService service.Auth) *authMiddleware {
	return &authMiddleware{authService: authService}
}

func (m *authMiddleware) ModeratorsOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		parts := strings.Split(header, " ")

		if parts[0] != "Bearer" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		claims, err := m.authService.ParseToken(parts[1])
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if claims.UserType != "moderator" {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Next()
	}
}

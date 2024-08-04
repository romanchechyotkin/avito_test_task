package middleware

import (
	"net/http"
	"strings"

	"github.com/romanchechyotkin/avito_test_task/internal/service"

	"github.com/gin-gonic/gin"
)

// todo refactoring

type AuthMiddleware struct {
	authService service.Auth
}

func NewAuthMiddleware(authService service.Auth) *AuthMiddleware {
	return &AuthMiddleware{authService: authService}
}

func (m *AuthMiddleware) ModeratorsOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		parts := strings.Split(header, " ")

		if parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "no authorization",
			})
			return
		}

		claims, err := m.authService.ParseToken(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "no authorization",
			})
			return
		}

		if claims.UserType != "moderator" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "no roots",
			})
			return
		}

		c.Set("userType", claims.UserType)
		c.Set("userID", claims.UserID)

		c.Next()
	}
}

func (m *AuthMiddleware) ClientsOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		parts := strings.Split(header, " ")

		if parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "no authorization",
			})
			return
		}

		claims, err := m.authService.ParseToken(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "no authorization",
			})
			return
		}

		if claims.UserType != "client" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "no roots",
			})
			return
		}

		c.Set("userType", claims.UserType)
		c.Set("userID", claims.UserID)

		c.Next()
	}
}

func (m *AuthMiddleware) AuthOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		parts := strings.Split(header, " ")

		if parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "no authorization",
			})
			return
		}

		claims, err := m.authService.ParseToken(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "no authorization",
			})
			return
		}

		if claims.UserType == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "no authorization",
			})
			return
		}

		c.Set("userType", claims.UserType)
		c.Set("userID", claims.UserID)

		c.Next()
	}
}

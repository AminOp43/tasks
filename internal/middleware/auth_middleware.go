package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
	"strings"
)

func AuthMiddleware(c *gin.Context) {
	header := c.GetHeader("Authorization")

	parts := strings.SplitN(header, " ", 2)

	if len(parts) != 2 || parts[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid authorization header",
		})
		c.Abort()
		return
	}

	tokenString := parts[1]
	claims := jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (any, error) {
			secret := os.Getenv("JWT_SECRET")
			if secret == "" {
				return nil, errors.New("JWT_SECRET is not configured")
			}

			return []byte(secret), nil
		},
	)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid or expired token",
		})
		c.Abort()
		return
	}

	if token == nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid token",
		})
		c.Abort()
		return
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid user id",
		})
		c.Abort()
		return
	}

	c.Set("user_id", int64(userID))

	c.Next()
}

package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("supersecretkey") // Should match auth_controller

// AuthMiddleware validates JWT token and adds email to context
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("auth_token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(cookie, func(t *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			// Add email and user_id to context
			c.Set("email", claims["email"].(string))
			c.Set("user_id", claims["user_id"].(string))
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}
	}
}

// GetEmailFromContext retrieves email from gin context
func GetEmailFromContext(c *gin.Context) (string, bool) {
	email, exists := c.Get("email")
	if !exists {
		return "", false
	}
	return email.(string), true
}

// GetUserIDFromContext retrieves user_id from gin context
func GetUserIDFromContext(c *gin.Context) (string, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return "", false
	}
	return userID.(string), true
}



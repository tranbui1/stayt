package auth

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func VerifyJWTToken(token string) (int64, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	tokenString, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return 0, err
	}

	claims, ok := tokenString.Claims.(jwt.MapClaims)
	if !ok || !tokenString.Valid {
		return 0, fmt.Errorf("invalid token")
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return 0, fmt.Errorf("invalid user_id in token")
	}

	return int64(userIDFloat), nil
}

func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Takes in provided jwt token from client & verifies jwt
		// + extracts embedded user_id
		// verifies jwt token

		token, err := c.Cookie("token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			c.Abort()
			return
		}

		userID, err := VerifyJWTToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}

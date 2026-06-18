package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}

	return func(c *gin.Context) {
		var tokenString string

		// 1. Пробуем взять из заголовка Authorization: Bearer <token>
		authHeader := c.GetHeader("Authorization")
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			tokenString = parts[1]
		} else {
			// 2. Fallback для WebSocket: query-параметр ?token=...
			tokenString = c.Query("token")
		}

		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization token is required"})
			c.Abort()
			return
		}

		// Проверяем алгоритм подписи — защита от подмены alg
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			c.Abort()
			return
		}

		// Безопасное извлечение user_id с учётом того, что MapClaims парсит числа как float64
		userID, err := extractUserID(claims)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}

func extractUserID(claims jwt.MapClaims) (string, error) {
	raw, exists := claims["user_id"]
	if !exists {
		return "", fmt.Errorf("user_id claim is missing")
	}

	switch v := raw.(type) {
	case string:
		if v == "" {
			return "", fmt.Errorf("user_id claim is empty")
		}
		return v, nil
	case float64:
		// MapClaims парсит JSON-числа как float64
		return fmt.Sprintf("%.0f", v), nil
	default:
		return "", fmt.Errorf("user_id claim has unexpected type: %T", raw)
	}
}

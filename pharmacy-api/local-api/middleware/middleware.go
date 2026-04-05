package middleware

import (
	"errors"
	"log"
	"net/http"
	"os"
	local_models "pharmacy/local-api/models"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

var (
	jwtSecret = []byte(os.Getenv("ACCESSTOKEN_SECRET"))
)

func AuthMiddleware(db *gorm.DB, requiredPermission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			return
		}
		tokenString := parts[1]

		claims := &local_models.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			if err == jwt.ErrTokenExpired {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token_expired"})
				return
			} else {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
				return
			}
		}
		if requiredPermission != "" {
			var role local_models.Role
			if err := db.Preload("Permissions").
				Where("roles.name = ?", claims.Role).Find(&role).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					log.Println("Middleware error [" + claims.Username + "]" + err.Error())
					c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Role not found"})
					return
				}
				log.Println("Middleware error [" + claims.Username + "]" + err.Error())
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			for _, permission := range role.Permissions {
				if permission.Action == requiredPermission {
					c.Set("user", claims)
					c.Next()
					return
				}
			}

			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
		} else {
			c.Set("user", claims)
			c.Next()
			return
		}
	}
}

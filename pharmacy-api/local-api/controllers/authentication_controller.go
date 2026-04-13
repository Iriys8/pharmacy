package controllers

import (
	"log"
	"net/http"
	"os"
	"time"

	models "pharmacy-api/shared/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	jwtSecret     = []byte(os.Getenv("ACCESSTOKEN_SECRET"))
	refreshSecret = []byte(os.Getenv("REFRESHTOKEN_SECRET"))
)

func Login(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginData struct {
			Login    string `json:"login"`
			Password string `json:"password"`
		}

		if err := c.ShouldBindJSON(&loginData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		var user models.User
		if err := db.Preload("Role.Permissions").Where("login = ?", loginData.Login).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		var userPermissions []string
		for _, permission := range user.Role.Permissions {
			userPermissions = append(userPermissions, permission.Action)
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(loginData.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		accessToken, err := generateAccessToken(&user)
		if err != nil {
			log.Println("Loggin error [" + user.UserName + "]" + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
			return
		}

		refreshToken, err := generateRefreshToken(&user)
		if err != nil {
			log.Println("Loggin error [" + user.UserName + "]" + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate refresh token"})
			return
		}

		c.SetCookie("pharmacy_refresh_token", refreshToken, 8*60*60, "/", "localhost", true, true)

		log.Println("Loggin [" + user.UserName + "]")
		c.JSON(http.StatusOK, gin.H{
			"access_token": accessToken,
			"token_type":   "Bearer",
			"expires_in":   15 * 60,
			"user": gin.H{
				"username":    user.UserName,
				"permissions": userPermissions,
			},
		})
	}
}

func Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.SetCookie("pharmacy_refresh_token", "", -1, "/", "localhost", true, true)
		c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
	}
}

func RefreshToken(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		refreshToken, err := c.Cookie("pharmacy_refresh_token")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token required"})
			return
		}

		token, err := jwt.ParseWithClaims(refreshToken, &models.RefreshClaims{}, func(token *jwt.Token) (any, error) {
			return refreshSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
			return
		}

		claims, ok := token.Claims.(*models.RefreshClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}

		var user models.User
		if err := db.Preload("Role.Permissions").First(&user, claims.UserID).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			return
		}

		var userPermissions []string
		for _, permission := range user.Role.Permissions {
			userPermissions = append(userPermissions, permission.Action)
		}

		newAccessToken, err := generateAccessToken(&user)
		if err != nil {
			log.Println("Refresh error: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
			return
		}

		newRefreshToken, err := generateRefreshToken(&user)
		if err != nil {
			log.Println("Refresh error: " + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate refresh token"})
			return
		}

		c.SetCookie("pharmacy_refresh_token", newRefreshToken, 2*60*60, "/", "localhost", true, true)

		c.JSON(http.StatusOK, gin.H{
			"access_token": newAccessToken,
			"token_type":   "Bearer",
			"expires_in":   15 * 60,
			"user": gin.H{
				"username":    user.UserName,
				"permissions": userPermissions,
			},
		})
	}
}

func generateAccessToken(user *models.User) (string, error) {
	expirationTime := time.Now().Add(15 * time.Minute)

	claims := &models.Claims{
		UserID:   user.ID,
		Username: user.Login,
		Role:     user.Role.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func generateRefreshToken(user *models.User) (string, error) {
	expirationTime := time.Now().Add(8 * time.Hour)

	claims := &models.RefreshClaims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(refreshSecret)
}

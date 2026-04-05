package controllers

import (
	"log"
	"net/http"
	local_models "pharmacy/local-api/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func GetUsers(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := c.Query("q")
		pageStr := c.Query("page")
		limitStr := c.Query("limit")

		user, _ := c.Get("user")
		claims := user.(*local_models.Claims)

		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit < 1 {
			limit = 10
		} else if limit > 40 {
			limit = 40
		}

		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			page = 1
		}
		offset := (page - 1) * limit

		var users []local_models.User
		var totalCount int64

		if query != "" {
			if err := db.Preload("Role").Where("user_name LIKE ?", "%"+query+"%").Find(&users).Error; err != nil {
				log.Println("Users GET error [" + claims.Username + "]" + err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			totalCount = int64(len(users))
			start := offset
			end := offset + limit
			if start > len(users) {
				start = len(users)
			}
			if end > len(users) {
				end = len(users)
			}
			users = users[start:end]
		} else {
			db.Model(&local_models.User{}).Count(&totalCount)
			if err := db.Preload("Role").Order("id DESC").Limit(limit).Offset(offset).Find(&users).Error; err != nil {
				log.Println("Users GET error [" + claims.Username + "]" + err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}

		totalPages := (totalCount + int64(limit) - 1) / int64(limit)
		log.Println("Users GET [" + claims.Username + "]")
		c.JSON(http.StatusOK, gin.H{
			"Items":       users,
			"TotalPages":  totalPages,
			"CurrentPage": page,
		})
	}
}

func GetUserByID(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var response local_models.User
		id := c.Param("id")

		user, _ := c.Get("user")
		claims := user.(*local_models.Claims)

		if id != "0" {
			if err := db.Preload("Role").First(&response, id).Error; err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
				return
			}

			response.PasswordHash = []byte("")
		}
		var roles []local_models.Role
		var requstedUser local_models.User
		if err := db.Preload("Role.Permissions").First(&requstedUser, claims.UserID).Error; err != nil {
			log.Println("User GET error [" + claims.Username + "]" + err.Error())
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Requested user not found"})
			return
		}
		for _, permission := range requstedUser.Role.Permissions {
			if permission.Action == "Change_Roles" {
				if err := db.Find(&roles).Error; err != nil {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "Can't found roles"})
					return
				}
				break
			}
		}

		log.Println("User GET [" + claims.Username + "]")
		c.JSON(http.StatusOK, gin.H{
			"User":  response,
			"Roles": roles,
		})
	}
}

func CreateUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var newUserRequest local_models.UserUpdateRequest
		if err := c.ShouldBindJSON(&newUserRequest); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		user, _ := c.Get("user")
		claims := user.(*local_models.Claims)

		pass, err := bcrypt.GenerateFromPassword([]byte(newUserRequest.Password), bcrypt.DefaultCost)

		if err != nil {
			log.Println("User POST error [" + claims.Username + "]" + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Can't generate password"})
			return
		}
		newUser := local_models.User{
			Login:        newUserRequest.Login,
			UserName:     newUserRequest.UserName,
			RoleID:       newUserRequest.RoleID,
			PasswordHash: pass,
		}

		if err := db.Create(&newUser).Error; err != nil {
			log.Println("User POST error [" + claims.Username + "]" + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		log.Println("User POST [" + claims.Username + "]")
		c.JSON(http.StatusOK, gin.H{"message": "User created successfully"})
	}
}

func DeleteUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		user, _ := c.Get("user")
		claims := user.(*local_models.Claims)

		if err := db.Delete(&local_models.User{}, id).Error; err != nil {
			log.Println("User DELETE error [" + claims.Username + "]" + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		log.Println("User DELETE [" + claims.Username + "]")
		c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
	}
}

func UpdateUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var userUpdateRequest local_models.UserUpdateRequest

		user, _ := c.Get("user")
		claims := user.(*local_models.Claims)

		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var updateData local_models.User

		if userUpdateRequest.Password != "" {
			pass, err := bcrypt.GenerateFromPassword([]byte(userUpdateRequest.Password), bcrypt.DefaultCost)

			if err != nil {
				log.Println("User PATCH error [" + claims.Username + "]" + err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
				return
			}
			updateData = local_models.User{
				Login:        userUpdateRequest.Login,
				UserName:     userUpdateRequest.UserName,
				RoleID:       userUpdateRequest.RoleID,
				PasswordHash: pass,
			}
		} else {
			updateData = local_models.User{
				Login:    userUpdateRequest.Login,
				UserName: userUpdateRequest.UserName,
				RoleID:   userUpdateRequest.RoleID,
			}
		}

		if err := db.Model(&local_models.User{}).Where("id = ?", id).Updates(updateData).Error; err != nil {
			log.Println("User PATCH error [" + claims.Username + "]" + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order: " + err.Error()})
			return
		}

		log.Println("User PATCH [" + claims.Username + "]")
		c.JSON(http.StatusOK, gin.H{"message": "Order and items updated"})
	}
}

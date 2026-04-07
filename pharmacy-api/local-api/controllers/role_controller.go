package controllers

import (
	"log"
	"net/http"
	local_models "pharmacy-api/local-api/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetRoles(db *gorm.DB) gin.HandlerFunc {
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

		var roles []local_models.Role
		var totalCount int64

		if query != "" {
			if err := db.Where("name LIKE ?", "%"+query+"%").Find(&roles).Error; err != nil {
				log.Println("Roles GET error [" + claims.Username + "]" + err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			totalCount = int64(len(roles))
			start := offset
			end := offset + limit
			if start > len(roles) {
				start = len(roles)
			}
			if end > len(roles) {
				end = len(roles)
			}
			roles = roles[start:end]
		} else {
			db.Model(&local_models.Role{}).Count(&totalCount)
			if err := db.Preload("Permissions").Order("id DESC").Limit(limit).Offset(offset).Find(&roles).Error; err != nil {
				log.Println("Roles GET error [" + claims.Username + "]" + err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}

		totalPages := (totalCount + int64(limit) - 1) / int64(limit)
		log.Println("Roles GET [" + claims.Username + "]")
		c.JSON(http.StatusOK, gin.H{
			"Items":       roles,
			"TotalPages":  totalPages,
			"CurrentPage": page,
		})
	}
}

func GetRoleByID(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var role local_models.Role

		user, _ := c.Get("user")
		claims := user.(*local_models.Claims)

		if err := db.Preload("Permissions").First(&role, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
			return
		}

		log.Println("Role GET [" + claims.Username + "]")
		c.JSON(http.StatusOK, role)
	}
}

func CreateRole(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var role local_models.Role

		user, _ := c.Get("user")
		claims := user.(*local_models.Claims)

		if err := c.ShouldBindJSON(&role); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		if err := db.Create(&role).Error; err != nil {
			log.Println("Role POST error [" + claims.Username + "]" + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if err := db.Model(&role).Association("Permissions").Replace(role.Permissions); err != nil {
			log.Println("Role POST error [" + claims.Username + "]" + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to bind permissions: " + err.Error()})
			return
		}

		log.Println("Role POST [" + claims.Username + "]")
		c.JSON(http.StatusOK, gin.H{"message": "Role created successfully"})
	}
}

func DeleteRole(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		user, _ := c.Get("user")
		claims := user.(*local_models.Claims)

		var role local_models.Role
		if err := db.Preload("Permissions").First(&role, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
			return
		}

		if err := db.Model(&role).Association("Permissions").Clear(); err != nil {
			log.Println("Role DELETE error [" + claims.Username + "]" + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear permissions: " + err.Error()})
			return
		}

		if err := db.Delete(&role).Error; err != nil {
			log.Println("Role DELETE error [" + claims.Username + "]" + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		log.Println("Role DELETE [" + claims.Username + "]")
		c.JSON(http.StatusOK, gin.H{"message": "Role deleted"})
	}
}

func UpdateRole(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var role local_models.Role

		user, _ := c.Get("user")
		claims := user.(*local_models.Claims)

		if err := c.ShouldBindJSON(&role); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var existingRole local_models.Role
		if err := db.Preload("Permissions").First(&existingRole, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
			return
		}

		if err := db.Model(&existingRole).Updates(local_models.Role{Name: role.Name}).Error; err != nil {
			log.Println("Role PATCH error [" + claims.Username + "]" + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update role: " + err.Error()})
			return
		}

		if err := db.Model(&existingRole).Association("Permissions").Replace(role.Permissions); err != nil {
			log.Println("Role PATCH error [" + claims.Username + "]" + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update permissions: " + err.Error()})
			return
		}

		log.Println("Role PATCH [" + claims.Username + "]")
		c.JSON(http.StatusOK, gin.H{"message": "Role and permissions updated"})
	}
}

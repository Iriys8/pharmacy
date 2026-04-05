package controllers

import (
	"fmt"
	"log"
	"net/http"
	local_models "pharmacy/local-api/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetPermissions(db *gorm.DB) gin.HandlerFunc {
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

		var permissions []local_models.Permission
		var totalCount int64

		if query != "" {
			if err := db.Where("action LIKE ?", "%"+query+"%").Find(&permissions).Error; err != nil {
				log.Println("Permissions GET error [" + claims.Username + "]" + err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			totalCount = int64(len(permissions))
			start := offset
			end := offset + limit
			if start > len(permissions) {
				start = len(permissions)
			}
			if end > len(permissions) {
				end = len(permissions)
			}
			fmt.Println(permissions)
			permissions = permissions[start:end]
		} else {
			db.Model(&local_models.Permission{}).Count(&totalCount)
			if err := db.Order("id DESC").Limit(limit).Offset(offset).Find(&permissions).Error; err != nil {
				log.Println("Permissions GET error [" + claims.Username + "]" + err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}

		totalPages := (totalCount + int64(limit) - 1) / int64(limit)
		log.Println("Permissions GET [" + claims.Username + "]")
		c.JSON(http.StatusOK, gin.H{
			"Items":       permissions,
			"TotalPages":  totalPages,
			"CurrentPage": page,
		})
	}
}

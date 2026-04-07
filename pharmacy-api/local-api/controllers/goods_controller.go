package controllers

import (
	"log"
	"net/http"
	local_models "pharmacy-api/local-api/models"
	shared_models "pharmacy-api/shared/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// расширить до полного crud
func UpdateGoods(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var updateData local_models.GoodsUpdateRequest

		user, _ := c.Get("user")
		claims := user.(*local_models.Claims)

		if err := c.ShouldBindJSON(&updateData); err != nil {
			log.Println("Good PATH error [" + claims.Username + "]" + err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var existingGood shared_models.Goods
		if err := db.First(&existingGood, id).Error; err != nil {
			log.Println("Good PATH error [" + claims.Username + "]" + err.Error())
			c.JSON(http.StatusNotFound, gin.H{"error": "Good not found"})
			return
		}

		existingGood.Name = updateData.Name
		existingGood.Instruction = updateData.Instruction
		existingGood.Description = updateData.Description
		existingGood.IsPrescriptionNeeded = updateData.IsPrescriptionNeeded
		existingGood.IsInStock = updateData.IsInStock
		existingGood.Price = updateData.Price

		if err := db.Save(&existingGood).Error; err != nil {
			log.Println("Good PATH error [" + claims.Username + "]" + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		log.Println("Good PATH [" + claims.Username + "]")
		c.JSON(http.StatusOK, gin.H{"message": "good updated", "data": existingGood})
	}
}

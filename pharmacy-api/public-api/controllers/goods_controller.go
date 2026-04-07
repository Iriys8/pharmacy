package controllers

import (
	"math/rand"
	"net/http"
	public_models "pharmacy-api/public-api/models"
	shared_models "pharmacy-api/shared/models"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func GetPromoItems(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var goods []shared_models.Goods
		if err := db.Where("is_in_stock LIKE ?", "1").Find(&goods).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		totalGoods := len(goods)
		if totalGoods == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		rand.NewSource(time.Now().UnixNano())
		indices := rand.Perm(totalGoods)[:min(5, totalGoods)]
		var promoItems []public_models.PromoItem
		for _, index := range indices {
			promoItems = append(promoItems, public_models.PromoItem{
				ID:          goods[index].ID,
				Name:        goods[index].Name,
				Description: goods[index].Description,
				Price:       goods[index].Price,
				Image:       goods[index].Image,
			})
		}

		c.JSON(http.StatusOK, promoItems)
	}
}

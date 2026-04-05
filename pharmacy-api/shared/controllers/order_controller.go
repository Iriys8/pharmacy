package controllers

import (
	"net/http"

	shared_models "pharmacy/shared/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateOrder(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestData shared_models.OrderResponse

		if err := c.ShouldBindJSON(&requestData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}
		order := shared_models.Order{
			ClientFIO:   requestData.Name,
			ClientEmail: requestData.Email,
			ClientPhone: requestData.Phone,
		}
		db.Create(&order)
		var lastOrder shared_models.Order
		db.Where("client_fio = ? AND client_email = ? AND client_phone = ?", requestData.Name, requestData.Email, requestData.Phone).
			Order("id DESC").First(&lastOrder)
		for _, item := range requestData.Items {
			db.Create(&shared_models.GoodsOrders{
				OrderID:  lastOrder.ID,
				GoodsID:  item.ID,
				Quantity: item.Quantity,
			})
		}
		c.JSON(http.StatusOK, gin.H{"message": "Order created successfully", "order_id": lastOrder.ID})
	}
}

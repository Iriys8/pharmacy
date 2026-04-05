package controllers

import (
	"log"
	"net/http"
	local_models "pharmacy/local-api/models"
	shared_models "pharmacy/shared/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetOrders(db *gorm.DB) gin.HandlerFunc {
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

		var orders []shared_models.Order
		var totalCount int64
		var response []shared_models.OrderResponse

		if query != "" {
			if err := db.Where("client_fio LIKE ?", "%"+query+"%").Find(&orders).Error; err != nil {
				log.Println("Orders GET error [" + claims.Username + "]" + err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			totalCount = int64(len(orders))
			start := offset
			end := offset + limit
			if start > len(orders) {
				start = len(orders)
			}
			if end > len(orders) {
				end = len(orders)
			}
			orders = orders[start:end]
		} else {
			db.Model(&shared_models.Order{}).Count(&totalCount)
			if err := db.Order("id DESC").Limit(limit).Offset(offset).Find(&orders).Error; err != nil {
				log.Println("Orders GET error [" + claims.Username + "]" + err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}

		for _, order := range orders {
			frontendOrder := shared_models.OrderResponse{
				ID:    order.ID,
				Name:  order.ClientFIO,
				Phone: order.ClientPhone,
				Email: order.ClientEmail,
			}
			response = append(response, frontendOrder)
		}

		totalPages := (totalCount + int64(limit) - 1) / int64(limit)
		log.Println("Orders GET [" + claims.Username + "]")
		c.JSON(http.StatusOK, gin.H{
			"Items":       response,
			"TotalPages":  totalPages,
			"CurrentPage": page,
		})
	}
}

func GetOrderByID(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var order shared_models.Order
		id := c.Param("id")

		user, _ := c.Get("user")
		claims := user.(*local_models.Claims)

		if err := db.Preload("Goods.Goods").First(&order, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
			return
		}

		response := shared_models.OrderResponse{
			ID:    order.ID,
			Name:  order.ClientFIO,
			Email: order.ClientEmail,
			Phone: order.ClientPhone,
			Items: make([]shared_models.OrderedItem, 0),
		}

		for _, goItem := range order.Goods {
			item := shared_models.OrderedItem{
				ID:          goItem.Goods.ID,
				Name:        goItem.Goods.Name,
				Image:       goItem.Goods.Image,
				Description: goItem.Goods.Description,
				Price:       goItem.Goods.Price,
				Quantity:    goItem.Quantity,
			}
			response.Items = append(response.Items, item)
		}

		log.Println("Order GET [" + claims.Username + "]")
		c.JSON(http.StatusOK, response)
	}
}

func DeleteOrder(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		user, _ := c.Get("user")
		claims := user.(*local_models.Claims)

		if err := db.Delete(&shared_models.Order{}, id).Error; err != nil {
			log.Println("Order DELETE error [" + claims.Username + "]" + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		log.Println("Order DELETE [" + claims.Username + "]")
		c.JSON(http.StatusOK, gin.H{"message": "Order deleted"})
	}
}

func UpdateOrder(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var input shared_models.OrderResponse

		user, _ := c.Get("user")
		claims := user.(*local_models.Claims)

		if err := c.ShouldBindJSON(&input); err != nil {
			log.Println("Order PATCH error [" + claims.Username + "]" + err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		updateData := map[string]interface{}{
			"ClientFIO":   input.Name,
			"ClientEmail": input.Email,
			"ClientPhone": input.Phone,
		}
		if err := db.Model(&shared_models.Order{}).Where("id = ?", id).Updates(updateData).Error; err != nil {
			log.Println("Order PATCH error [" + claims.Username + "]" + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order: " + err.Error()})
			return
		}

		if err := db.Where("order_id = ?", id).Delete(&shared_models.GoodsOrders{}).Error; err != nil {
			log.Println("Order PATCH error [" + claims.Username + "]" + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete old goods: " + err.Error()})
			return
		}

		for _, item := range input.Items {
			newGoodsOrder := shared_models.GoodsOrders{
				OrderID:  input.ID,
				GoodsID:  item.ID,
				Quantity: item.Quantity,
			}
			if err := db.Create(&newGoodsOrder).Error; err != nil {
				log.Println("Order PATCH error [" + claims.Username + "]" + err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create goods order: " + err.Error()})
				return
			}
		}

		log.Println("Order PATCH [" + claims.Username + "]")
		c.JSON(http.StatusOK, gin.H{"message": "Order and items updated"})
	}
}

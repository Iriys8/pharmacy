package controller

import (
	"log"
	"net/http"
	models "pharmacy-api/shared/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetOrders(db *gorm.DB, query string, pageStr string, limitStr string, claims models.Claims) (result map[string]any, err error) {

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

	var orders []models.Order
	var totalCount int64
	var response []models.OrderResponse

	if query != "" {
		if err = db.Where("client_fio LIKE ?", "%"+query+"%").Find(&orders).Error; err != nil {
			log.Println("Orders GET error [" + claims.Username + "]" + err.Error())
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
		db.Model(&models.Order{}).Count(&totalCount)
		if err = db.Order("id DESC").Limit(limit).Offset(offset).Find(&orders).Error; err != nil {
			log.Println("Orders GET error [" + claims.Username + "]" + err.Error())
			return
		}
	}

	for _, order := range orders {
		frontendOrder := models.OrderResponse{
			ID:    order.ID,
			Name:  order.ClientFIO,
			Phone: order.ClientPhone,
			Email: order.ClientEmail,
		}
		response = append(response, frontendOrder)
	}

	totalPages := (totalCount + int64(limit) - 1) / int64(limit)
	log.Println("Orders GET [" + claims.Username + "]")

	result = make(map[string]any)
	result["Items"] = response
	result["TotalPages"] = totalPages
	result["CurrentPage"] = page
	return
}

func GetOrderByID(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var order models.Order
		id := c.Param("id")

		user, _ := c.Get("user")
		claims := user.(*models.Claims)

		if err := db.Preload("Goods.Goods").First(&order, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
			return
		}

		response := models.OrderResponse{
			ID:    order.ID,
			Name:  order.ClientFIO,
			Email: order.ClientEmail,
			Phone: order.ClientPhone,
			Items: make([]models.OrderedItem, 0),
		}

		for _, goItem := range order.Goods {
			item := models.OrderedItem{
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

func CreateOrder(db *gorm.DB, requestData models.OrderResponse) (result map[string]any, err error) {

	order := models.Order{
		ClientFIO:   requestData.Name,
		ClientEmail: requestData.Email,
		ClientPhone: requestData.Phone,
	}
	db.Create(&order)
	var lastOrder models.Order
	db.Where("client_fio = ? AND client_email = ? AND client_phone = ?", requestData.Name, requestData.Email, requestData.Phone).
		Order("id DESC").First(&lastOrder)
	for _, item := range requestData.Items {
		db.Create(&models.GoodsOrders{
			OrderID:  lastOrder.ID,
			GoodsID:  item.ID,
			Quantity: item.Quantity,
		})
	}

	result = make(map[string]any)
	result["Message"] = "Order created successfully"
	result["Order_id"] = lastOrder.ID
	return
}

func DeleteOrder(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		user, _ := c.Get("user")
		claims := user.(*models.Claims)

		if err := db.Delete(&models.Order{}, id).Error; err != nil {
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
		var input models.OrderResponse

		user, _ := c.Get("user")
		claims := user.(*models.Claims)

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
		if err := db.Model(&models.Order{}).Where("id = ?", id).Updates(updateData).Error; err != nil {
			log.Println("Order PATCH error [" + claims.Username + "]" + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order: " + err.Error()})
			return
		}

		if err := db.Where("order_id = ?", id).Delete(&models.GoodsOrders{}).Error; err != nil {
			log.Println("Order PATCH error [" + claims.Username + "]" + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete old goods: " + err.Error()})
			return
		}

		for _, item := range input.Items {
			newGoodsOrder := models.GoodsOrders{
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

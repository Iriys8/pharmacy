package controller

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func Pickup(redisDB *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {

		key := c.Query("id")

		if key == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"Error": "key is required",
			})
			return
		}

		val, err := redisDB.HGetAll(c.Request.Context(), "public_task:"+key).Result()

		if len(val) == 0 {
			c.JSON(http.StatusNotFound, gin.H{
				"Status": "not_found",
			})
			return

		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"Error": "Internal error",
			})
			return

		}

		switch val["status"] {
		case "completed":
			err := redisDB.Del(c.Request.Context(), "public_task:"+key).Err()
			if err != nil {
				log.Printf("Error: %v", err)
				c.JSON(http.StatusOK, gin.H{
					"Status": val["status"],
					"Value":  val["result"],
				})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"Status": val["status"],
				"Value":  val["result"],
			})
			return
		case "pending":
			c.JSON(http.StatusOK, gin.H{
				"Status": val["status"],
			})
		case "error":
			c.JSON(http.StatusOK, gin.H{
				"Status": val["status"],
			})
		}

	}
}

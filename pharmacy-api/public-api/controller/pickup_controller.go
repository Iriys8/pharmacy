package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func Pickup(redisDB *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.Param("key")

		if key == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "key is required",
			})
			return
		}

		val, err := redisDB.Get(c.Request.Context(), key).Result()

		if err == redis.Nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "key not found",
				"key":   key,
			})
			return

		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "redis error: " + err.Error(),
				"key":   key,
			})
			return

		}
		delErr := redisDB.Del(c.Request.Context(), key).Err()
		if delErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to delete key: " + delErr.Error(),
				"key":   key,
				"value": val,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"key":   key,
			"value": val,
		})
	}
}

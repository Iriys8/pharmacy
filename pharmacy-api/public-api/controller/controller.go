package controller

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"

	shared_models "pharmacy-api/shared/models"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	amqp "github.com/rabbitmq/amqp091-go"
)

func generateTaskID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func MakeTask(route string, task string, redisDB *redis.Client, broker *shared_models.Broker) gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx := c.Request.Context()

		var taskContext shared_models.RequestContext

		query := make(map[string]any)
		switch task {
		case "get":
			query["id"] = c.Query("id")
			query["q"] = c.Query("q")
			query["page"] = c.Query("page")
			query["limit"] = c.Query("limit")

			queryJSON, _ := json.Marshal(query)
			taskContext.Query = queryJSON
		case "schedule_dated":
			if c.Query("start") == "" || c.Query("end") == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Start and end dates are required"})
				return
			}
			_, err := time.Parse("2006-01-02", c.Query("start"))
			_, err = time.Parse("2006-01-02", c.Query("end"))
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
				return
			}
			query["start"] = c.Query("start")
			query["end"] = c.Query("end")

			queryJSON, _ := json.Marshal(query)
			taskContext.Query = queryJSON
		case "post":
			if err := c.ShouldBindJSON(&taskContext.Context); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error":   "Invalid request body",
					"details": err.Error(),
				})
				return
			}
		}
		taskContextJSON, _ := json.Marshal(taskContext)

		taskID, err := generateTaskID()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal error",
			})
			return
		}

		key := "public_task:" + taskID
		err = redisDB.HSet(ctx, key, map[string]interface{}{
			"id":      taskID,
			"status":  "pending",
			"task":    task,
			"route":   route,
			"context": taskContextJSON,
			"api":     "public",
		}).Err()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal error: " + err.Error(),
			})
			return
		}

		err = redisDB.Expire(ctx, key, 10*time.Minute).Err()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal error: " + err.Error(),
			})
			return
		}

		err = broker.Channel.Publish(
			"public_exchange",
			route,
			false,
			false,
			amqp.Publishing{
				ContentType: "application/json",
				Body:        []byte(taskID),
				Timestamp:   time.Now(),
			},
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal error: " + err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"task_id": taskID,
		})
	}
}

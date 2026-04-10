package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	controllers "pharmacy-api/shared/controllers"
	shared_models "pharmacy-api/shared/models"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	amqp "github.com/rabbitmq/amqp091-go"
)

func MakeTask(route string, task string, redisDB *redis.Client, broker *shared_models.Broker) gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx := c.Request.Context()

		var taskContext shared_models.RequestContext

		query := make(map[string]any)
		switch task {
		case "get":
			var err error
			query["id"], err = strconv.Atoi(c.Query("id"))
			query["q"] = c.Query("q")
			query["page"] = c.Query("page")
			query["limit"] = c.Query("limit")

			if err != nil && c.Query("id") != "" {
				c.JSON(http.StatusBadRequest, gin.H{
					"error":   "Invalid request",
					"details": err.Error(),
				})
				return
			}

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

		taskID, err := controllers.RandomGen()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal error",
			})
			log.Print(err.Error())
			return
		}

		key := "public_task:" + taskID

		if err = redisDB.HSet(ctx, key, map[string]interface{}{
			"status":  "pending",
			"task":    task,
			"route":   route,
			"result":  "",
			"context": taskContextJSON,
		}).Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal error",
			})
			log.Print(err.Error())
			return
		}

		if err = redisDB.Expire(ctx, key, 5*time.Minute).Err(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal error",
			})
			log.Print(err.Error())
			return
		}

		err = broker.Channel.Publish(
			"public_exchange",
			route,
			false,
			false,
			amqp.Publishing{
				ContentType: "application/json",
				Body:        []byte(key),
				Timestamp:   time.Now(),
			},
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal error",
			})
			log.Print(err.Error())
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"task_id": key,
		})
	}
}

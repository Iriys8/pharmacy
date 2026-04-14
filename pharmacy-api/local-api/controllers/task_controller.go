package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	controllers "pharmacy-api/shared/controllers"
	models "pharmacy-api/shared/models"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	amqp "github.com/rabbitmq/amqp091-go"
)

func MakeTask(route string, task string, redisDB *redis.Client, broker *models.Broker) gin.HandlerFunc {
	return func(c *gin.Context) {

		ctx := c.Request.Context()

		var taskContext models.RequestContext

		var err error

		query := make(map[string]any)

		user, _ := c.Get("user")
		claims := user.(*models.Claims)

		if taskContext.Claims, err = json.Marshal(claims); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal error",
			})
			log.Print(err.Error())
		}

		switch task {
		case "get":
			var idInt int
			idInt, _ = strconv.Atoi(c.Query("id"))
			if route != "users" {
				query["id"] = idInt
			} else {
				query["id"] = c.Query("id")
			}
			query["q"] = c.Query("q")
			query["page"] = c.Query("page")
			query["limit"] = c.Query("limit")

			if err != nil && c.Query("id") != "" {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Invalid request body",
				})
				log.Print(err.Error())
				return
			}

			if taskContext.Query, err = json.Marshal(query); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Internal error",
				})
				log.Print(err.Error())
			}
		case "post":
			if err := c.ShouldBindJSON(&taskContext.Context); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Invalid request body",
				})
				log.Print(err.Error())
				return
			}
		case "patch":
			query["id"], err = strconv.Atoi(c.Query("id"))
			if err != nil && c.Query("id") != "" {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Invalid request body",
				})
				log.Print(err.Error())
				return
			}
			if taskContext.Query, err = json.Marshal(query); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Internal error",
				})
				log.Print(err.Error())
			}

			if err := c.ShouldBindJSON(&taskContext.Context); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Invalid request body",
				})
				log.Print(err.Error())
				return
			}
		case "delete":
			query["id"], err = strconv.Atoi(c.Query("id"))
			if err != nil && c.Query("id") != "" {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Invalid request body",
				})
				log.Print(err.Error())
				return
			}
			if taskContext.Query, err = json.Marshal(query); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Internal error",
				})
				log.Print(err.Error())
			}
		}
		var taskContextJSON []byte
		if taskContextJSON, err = json.Marshal(taskContext); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal error",
			})
			log.Print(err.Error())
		}
		taskID, err := controllers.RandomGen()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal error",
			})
			log.Print(err.Error())
			return
		}

		key := "local_task:" + taskID

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
			"local_exchange",
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
			"TaskID": taskID,
		})
	}
}

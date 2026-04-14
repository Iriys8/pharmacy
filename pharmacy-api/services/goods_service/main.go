package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	service_controller "pharmacy-api/services/goods_service/controller"
	shared_controller "pharmacy-api/shared/controllers"
	models "pharmacy-api/shared/models"
	setup "pharmacy-api/shared/setup"

	"github.com/go-redis/redis/v8"
	amqp "github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
)

func main() {
	uname_public, _ := shared_controller.RandomGen()
	uname_local, _ := shared_controller.RandomGen()

	logFile := setup.SetupLogs(uname_public + "_" + uname_local)
	defer logFile.Close()
	log.Printf("Service public name: %v, local name: %v", uname_public, uname_local)
	fmt.Printf("Service public name: %v, local name: %v", uname_public, uname_local)

	db := setup.ConnectDB()
	redisDB := setup.ConnectRedis()
	defer redisDB.Close()
	broker := setup.ConnectBroker()
	defer broker.Close()

	go consumeMessages(broker.Channel, "public_goods_queue", redisDB, db, uname_public)
	go consumeMessages(broker.Channel, "local_goods_queue", redisDB, db, uname_local)

	log.Println("Service is running.")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}

func consumeMessages(ch *amqp.Channel, queueName string, redisDB *redis.Client, db *gorm.DB, consumerName string) {
	msgs, err := ch.Consume(
		queueName,
		consumerName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to register consumer for %s: %v", queueName, err)
	}

	for msg := range msgs {
		func() {
			defer msg.Ack(false)

			ctx := context.Background()

			taskData, err := redisDB.HGetAll(ctx, string(msg.Body)).Result()
			if err == redis.Nil {
				log.Printf("Task %s not found in Redis", string(msg.Body))
				return
			} else if err != nil {
				log.Printf("Error getting task from Redis: %v", err)
				return
			}

			if taskData["status"] != "pending" {
				log.Printf("Task %s is not pending (status: %s), skipping", string(msg.Body), taskData["status"])
				return
			}

			var result map[string]any
			var execErr error

			switch taskData["task"] {
			case "get":
				var taskContext struct {
					Query struct {
						Id    int    `json:"id"`
						Q     string `json:"q"`
						Page  string `json:"page"`
						Limit string `json:"limit"`
					} `json:"Query"`
				}
				if err := json.Unmarshal([]byte(taskData["context"]), &taskContext); err != nil {
					log.Printf("Error parsing context for get: %v", err)
					return
				}
				if taskContext.Query.Id == 0 {
					result, execErr = service_controller.GetGoods(db, taskContext.Query.Q, taskContext.Query.Page, taskContext.Query.Limit)
				} else {
					result, execErr = service_controller.GetGoodsByID(db, taskContext.Query.Id)
				}

			case "advert":
				result, execErr = service_controller.GetPromoItems(db)

			case "patch":
				var taskContext struct {
					Query struct {
						Id int `json:"id"`
					} `json:"query"`
					Context models.GoodsUpdateRequest `json:"context"`
					Claims  models.Claims             `json:"claims"`
				}
				if err := json.Unmarshal([]byte(taskData["context"]), &taskContext); err != nil {
					log.Printf("Error parsing context for patch: %v", err)
					return
				}
				result, execErr = service_controller.UpdateGoods(db, taskContext.Query.Id, taskContext.Context, taskContext.Claims)

			default:
				log.Printf("Unknown task type: %s", taskData["task"])
				execErr = errors.New("Unknown task")
				return
			}

			var JSONResult []byte
			if result["Response"] != nil {
				JSONResult, execErr = json.Marshal(result["Response"])
			} else {
				JSONResult, execErr = json.Marshal(result)
			}

			if execErr != nil {
				taskData["status"] = "error"

			} else {
				taskData["status"] = "completed"
			}

			taskData["result"] = string(JSONResult)

			if err := redisDB.HSet(ctx, string(msg.Body), taskData).Err(); err != nil {
				log.Printf("Error updating task result in Redis: %v", err)
				return
			}

			if err = redisDB.Expire(ctx, string(msg.Body), 20*time.Second).Err(); err != nil {
				log.Printf("Error updating task result in Redis: %v", err)
				return
			}

			log.Printf("Task %s completed successfully", string(msg.Body))
		}()
	}
}

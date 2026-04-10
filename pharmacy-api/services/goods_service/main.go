package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	goods_controller "pharmacy-api/services/goods_service/controller"
	shared_controller "pharmacy-api/shared/controllers"
	shared_models "pharmacy-api/shared/models"
	setup "pharmacy-api/shared/setup"

	"github.com/go-redis/redis/v8"
	amqp "github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
)

func main() {
	uname, _ := shared_controller.Random_gen()

	setup.SetupLogs(uname)
	log.Printf("Service name: %v", uname)

	db := setup.ConnectDB()
	redisDB := setup.ConnectRedis()
	defer redisDB.Close()
	broker := setup.ConnectBroker()
	defer broker.Close()

	go consumeMessages(broker.Channel, "public_goods_queue", redisDB, db, uname)
	//go consumeMessages(broker.Channel, "private_goods_queue", redisDB, db, uname)

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
					Context []byte `json:"Context"`
				}
				if err := json.Unmarshal([]byte(taskData["context"]), &taskContext); err != nil {
					log.Printf("Error parsing context for get: %v", err)
					return
				}
				if taskContext.Query.Id == 0 {
					result, execErr = goods_controller.GetGoods(taskContext.Query.Q, taskContext.Query.Page, taskContext.Query.Limit, db)
				} else {
					result, execErr = goods_controller.GetGoodsByID(taskContext.Query.Id, db)
				}

			case "get_promo_items":
				result, execErr = goods_controller.GetPromoItems(db)

			// пиздец
			case "update_goods":
				var context struct {
					ID         int                              `json:"id"`
					UpdateData shared_models.GoodsUpdateRequest `json:"update_data"`
					Claims     shared_models.Claims             `json:"claims"`
				}
				if err := json.Unmarshal([]byte(taskData["context"]), &context); err != nil {
					log.Printf("Error parsing context for update_goods: %v", err)
					return
				}
				result, execErr = goods_controller.UpdateGoods(context.ID, context.UpdateData, context.Claims, db)

			default:
				log.Printf("Unknown task type: %s", taskData["task"])
				return
			}

			var JSONResult []byte
			JSONResult, execErr = json.Marshal(result)

			if execErr != nil {
				taskData["status"] = "error"
				taskData["result"] = execErr.Error()

			} else {
				taskData["status"] = "completed"
				taskData["result"] = string(JSONResult)
			}

			if err := redisDB.HSet(ctx, string(msg.Body), taskData).Err(); err != nil {
				log.Printf("Error updating task result in Redis: %v", err)
				return
			}

			if err = redisDB.Expire(ctx, string(msg.Body), 2*time.Minute).Err(); err != nil {
				log.Printf("Error updating task result in Redis: %v", err)
				return
			}

			log.Printf("Task %s completed successfully", string(msg.Body))
		}()
	}
}

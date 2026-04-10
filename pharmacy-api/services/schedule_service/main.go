package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	service_controller "pharmacy-api/services/schedule_service/controller"
	shared_controller "pharmacy-api/shared/controllers"
	setup "pharmacy-api/shared/setup"

	"github.com/go-redis/redis/v8"
	amqp "github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
)

func main() {
	uname, _ := shared_controller.RandomGen()

	setup.SetupLogs(uname)
	log.Printf("Service name: %v", uname)
	fmt.Printf("Service name: %v", uname)

	db := setup.ConnectDB()
	redisDB := setup.ConnectRedis()
	defer redisDB.Close()
	broker := setup.ConnectBroker()
	defer broker.Close()

	go consumeMessages(broker.Channel, "public_schedule_queue", redisDB, db, uname)
	//go consumeMessages(broker.Channel, "private_schedule_queue", redisDB, db, uname)

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
			case "schedule_dated":
				var taskContext struct {
					Query struct {
						Start string `json:"start"`
						End   string `json:"end"`
					} `json:"Query"`
				}
				if err := json.Unmarshal([]byte(taskData["context"]), &taskContext); err != nil {
					log.Printf("Error parsing context for get: %v", err)
					return
				}
				result, execErr = service_controller.GetScheduleDated(db, taskContext.Query.Start, taskContext.Query.End)

			default:
				log.Printf("Unknown task type: %s", taskData["task"])
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

			if err = redisDB.Expire(ctx, string(msg.Body), 2*time.Minute).Err(); err != nil {
				log.Printf("Error updating task result in Redis: %v", err)
				return
			}

			log.Printf("Task %s completed successfully", string(msg.Body))
		}()
	}
}

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

	service_controller "pharmacy-api/services/roles_service/controller"
	shared_controller "pharmacy-api/shared/controllers"
	"pharmacy-api/shared/models"
	setup "pharmacy-api/shared/setup"

	"github.com/go-redis/redis/v8"
	amqp "github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
)

func main() {
	uname_local, _ := shared_controller.RandomGen()

	logFile := setup.SetupLogs(uname_local)
	defer logFile.Close()
	log.Printf("Service local name: %v", uname_local)
	fmt.Printf("Service local name: %v", uname_local)

	db := setup.ConnectDB()
	redisDB := setup.ConnectRedis()
	defer redisDB.Close()
	broker := setup.ConnectBroker()
	defer broker.Close()

	go consumeMessages(broker.Channel, "local_roles_queue", redisDB, db, uname_local)

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
					} `json:"query"`
					Claims models.Claims `json:"claims"`
				}
				if err := json.Unmarshal([]byte(taskData["context"]), &taskContext); err != nil {
					log.Printf("Error parsing context for get: %v", err)
					return
				}
				if taskContext.Query.Id == 0 {
					result, execErr = service_controller.GetRoles(db, taskContext.Query.Q, taskContext.Query.Page, taskContext.Query.Limit, taskContext.Claims)
				} else {
					result, execErr = service_controller.GetRoleByID(db, taskContext.Query.Id, taskContext.Claims)
				}
			case "post":
				var taskContext struct {
					Context models.Role   `json:"context"`
					Claims  models.Claims `json:"claims"`
				}
				if err := json.Unmarshal([]byte(taskData["context"]), &taskContext); err != nil {
					log.Printf("Error parsing context for post: %v", err)
					return
				}
				result, execErr = service_controller.CreateRole(db, taskContext.Context, taskContext.Claims)
			case "patch":
				var taskContext struct {
					Query struct {
						Id int `json:"id"`
					} `json:"query"`
					Context models.Role   `json:"context"`
					Claims  models.Claims `json:"claims"`
				}
				if err := json.Unmarshal([]byte(taskData["context"]), &taskContext); err != nil {
					log.Printf("Error parsing context for patch: %v", err)
					return
				}
				result, execErr = service_controller.UpdateRole(db, taskContext.Query.Id, taskContext.Context, taskContext.Claims)
			case "delete":
				var taskContext struct {
					Query struct {
						Id int `json:"id"`
					} `json:"query"`
					Claims models.Claims `json:"claims"`
				}
				if err := json.Unmarshal([]byte(taskData["context"]), &taskContext); err != nil {
					log.Printf("Error parsing context for delete: %v", err)
					return
				}
				result, execErr = service_controller.DeleteRole(db, taskContext.Query.Id, taskContext.Claims)
			case "permissions":
				var taskContext struct {
					Query struct {
						Q     string `json:"q"`
						Page  string `json:"page"`
						Limit string `json:"limit"`
					} `json:"query"`
					Claims models.Claims `json:"claims"`
				}
				if err := json.Unmarshal([]byte(taskData["context"]), &taskContext); err != nil {
					log.Printf("Error parsing context for get: %v", err)
					return
				}
				result, execErr = service_controller.GetPermissions(db, taskContext.Query.Q, taskContext.Query.Page, taskContext.Query.Limit, taskContext.Claims)
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

			if err = redisDB.Expire(ctx, string(msg.Body), 20*time.Second).Err(); err != nil {
				log.Printf("Error updating task result in Redis: %v", err)
				return
			}

			log.Printf("Task %s completed successfully", string(msg.Body))
		}()
	}
}

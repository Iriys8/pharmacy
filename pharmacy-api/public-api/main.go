package main

import (
	"log"
	"pharmacy-api/public-api/routes"

	"github.com/gin-gonic/gin"

	setup "pharmacy-api/shared/setup"
)

func main() {
	db := setup.ConnectDB()

	setup.SetupDB(db)

	redisDB := setup.ConnectRedis()
	defer redisDB.Close()

	broker := setup.ConnectBroker()
	defer broker.Close()

	err := setupBroker(broker)
	if err != nil {
		log.Fatalf("failed to setup broker: %v", err)
	}

	router := gin.Default()

	routes.SetupRoutes(router, db, redisDB, broker)

	router.Run(":8080")
}

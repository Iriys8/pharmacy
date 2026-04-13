package main

import (
	"log"

	"pharmacy-api/local-api/routes"
	shared_controlers "pharmacy-api/shared/controllers"
	setup "pharmacy-api/shared/setup"

	"github.com/gin-gonic/gin"
)

func main() {
	name, _ := shared_controlers.RandomGen()
	logFile := setup.SetupLogs("local-api" + name)
	defer logFile.Close()
	log.Println("Started")

	db := setup.ConnectDB()
	log.Println("Database connected")
	setup.SetupDB(db)

	redisDB := setup.ConnectRedis()
	defer redisDB.Close()
	log.Println("Database initialized")

	broker := setup.ConnectBroker()
	defer broker.Close()

	err := setupBroker(broker)
	if err != nil {
		log.Fatalf("failed to setup broker: %v", err)
	}

	router := gin.Default()
	routes.SetupRoutes(router, db, redisDB, broker)
	log.Println("Routers initialized")

	setup.SetupAdmin(db)
	log.Println("Admin created")

	log.Println("Initialized")
	router.Run(":8080")
}

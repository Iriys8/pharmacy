package main

import (
	"pharmacy-api/public-api/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	db := setupDatabase()

	redisDB := setupRedis()

	broker := setupBroker()
	defer broker.Close()

	router := gin.Default()

	routes.SetupRoutes(router, db, redisDB, broker)

	router.Run(":8080")
}

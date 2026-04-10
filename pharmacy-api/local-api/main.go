package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"pharmacy-api/local-api/routes"
	shared_models "pharmacy-api/shared/models"
)

func main() {
	logFile := setupLogs()
	defer logFile.Close()
	log.Println("Started")

	db := setupDatabase()
	log.Println("Database initialized")

	router := gin.Default()
	routes.SetupRoutes(router, db)
	log.Println("Routers initialized")

	setupAdmin(db)
	log.Println("Admin created")

	var itemsAmount int64
	db.Model(&shared_models.Goods{}).Count(&itemsAmount)
	if itemsAmount == 0 {
		test_data(db)
		log.Println("TEST DATA USED!")
	}

	log.Println("Initialized")
	router.Run(":8080")
}

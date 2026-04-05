package main

import (
	"log"

	"pharmacy/local-api/routes"

	shared_models "pharmacy/shared/models"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}

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
		log.Fatalln("TEST DATA USED!")
	}

	log.Println("Initialized")
	router.Run(":8080")
}

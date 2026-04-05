package main

import (
	"log"

	"pharmacy/public-api/routes"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	db := setupDatabase()

	router := gin.Default()

	routes.SetupRoutes(router, db)

	router.Run(":8080")
}

package main

import (
	"pharmacy/public-api/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	db := setupDatabase()

	router := gin.Default()

	routes.SetupRoutes(router, db)

	router.Run(":8080")
}

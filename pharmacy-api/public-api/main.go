package main

import (
	"log"
	"pharmacy-api/public-api/routes"

	"github.com/gin-gonic/gin"

	shared_models "pharmacy-api/shared/models"
	setup "pharmacy-api/shared/setup"
)

func main() {
	db := setup.ConnectDB()

	var itemsAmount int64
	db.Model(&shared_models.Goods{}).Count(&itemsAmount)
	if itemsAmount == 0 {
		test_data(db)
		log.Println("TEST DATA USED!")
	}

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

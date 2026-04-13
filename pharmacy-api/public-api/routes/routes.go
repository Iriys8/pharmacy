package routes

import (
	controllers "pharmacy-api/public-api/controllers"
	shared_controllers "pharmacy-api/shared/controllers"
	models "pharmacy-api/shared/models"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func setupGoodsRouter(ApiGroup *gin.RouterGroup, redisDB *redis.Client, broker *models.Broker) {
	ApiGroup.GET("/goods", controllers.MakeTask("goods", "get", redisDB, broker))
	ApiGroup.GET("/goods/advert", controllers.MakeTask("goods", "advert", redisDB, broker))
}

func setupSheduleRouter(ApiGroup *gin.RouterGroup, redisDB *redis.Client, broker *models.Broker) {
	ApiGroup.GET("/schedule", controllers.MakeTask("schedule", "schedule_dated", redisDB, broker))
}

func setupOrderRouter(ApiGroup *gin.RouterGroup, redisDB *redis.Client, broker *models.Broker) {
	ApiGroup.POST("/order", controllers.MakeTask("orders", "post", redisDB, broker))
}

func setupAnnounceRouter(ApiGroup *gin.RouterGroup, redisDB *redis.Client, broker *models.Broker) {
	ApiGroup.GET("/announces", controllers.MakeTask("announces", "get", redisDB, broker))
}

func setupPickupRouter(ApiGroup *gin.RouterGroup, redisDB *redis.Client) {
	ApiGroup.GET("/pickup", controllers.Pickup(redisDB))
}

func SetupRoutes(router *gin.Engine, redisDB *redis.Client, broker *models.Broker) {
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5000", "http://127.0.0.1:5000", "http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           8 * time.Hour,
	}))

	ApiGroup := router.Group("/api/")

	ApiGroup.GET("/image", shared_controllers.GetImage)

	setupGoodsRouter(ApiGroup, redisDB, broker)

	setupSheduleRouter(ApiGroup, redisDB, broker)

	setupOrderRouter(ApiGroup, redisDB, broker)

	setupAnnounceRouter(ApiGroup, redisDB, broker)

	setupPickupRouter(ApiGroup, redisDB)
}

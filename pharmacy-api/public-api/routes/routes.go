package routes

import (
	controller "pharmacy-api/public-api/controller"
	shared_controllers "pharmacy-api/shared/controllers"
	shared_models "pharmacy-api/shared/models"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

func setupGoodsRouter(ApiGroup *gin.RouterGroup, redisDB *redis.Client, broker *shared_models.Broker) {
	ApiGroup.GET("/goods", controller.MakeTask("goods", "get", redisDB, broker))
	ApiGroup.GET("/goods/advert", controller.MakeTask("goods", "advert", redisDB, broker))
}

func setupWorkTimeRouter(ApiGroup *gin.RouterGroup, redisDB *redis.Client, broker *shared_models.Broker) {
	ApiGroup.GET("/schedule", controller.MakeTask("schedule", "schedule_dated", redisDB, broker))
}

func setupOrderRouter(ApiGroup *gin.RouterGroup, redisDB *redis.Client, broker *shared_models.Broker) {
	ApiGroup.POST("/order", controller.MakeTask("orders", "post", redisDB, broker))
}

func setupAnnounceRouter(ApiGroup *gin.RouterGroup, redisDB *redis.Client, broker *shared_models.Broker) {
	ApiGroup.GET("/announces", controller.MakeTask("announces", "get", redisDB, broker))
}

func SetupRoutes(router *gin.Engine, db *gorm.DB, redisDB *redis.Client, broker *shared_models.Broker) {
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5000", "http://127.0.0.1:5000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           8 * time.Hour,
	}))

	ApiGroup := router.Group("/api/")

	ApiGroup.GET("/image", shared_controllers.GetImage)

	setupGoodsRouter(ApiGroup, redisDB, broker)

	setupWorkTimeRouter(ApiGroup, redisDB, broker)

	setupOrderRouter(ApiGroup, redisDB, broker)

	setupAnnounceRouter(ApiGroup, redisDB, broker)
}

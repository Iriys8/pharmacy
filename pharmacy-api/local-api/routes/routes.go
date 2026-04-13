package routes

import (
	local_controllers "pharmacy-api/local-api/controllers"
	shared_controllers "pharmacy-api/shared/controllers"
	models "pharmacy-api/shared/models"

	"pharmacy-api/local-api/middleware"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

func setupGoodsRouter(ApiGroup *gin.RouterGroup, db *gorm.DB, redisDB *redis.Client, broker *models.Broker) {
	ApiGroup.GET("/goods", middleware.AuthMiddleware(db, ""), local_controllers.MakeTask("goods", "get", redisDB, broker))
	ApiGroup.PATCH("/goods", middleware.AuthMiddleware(db, "Update_Goods"), local_controllers.MakeTask("goods", "patch", redisDB, broker))
}

func setupSheduleRouter(ApiGroup *gin.RouterGroup, db *gorm.DB, redisDB *redis.Client, broker *models.Broker) {
	ApiGroup.GET("/schedule", middleware.AuthMiddleware(db, ""), local_controllers.MakeTask("schedule", "get", redisDB, broker))
	ApiGroup.POST("/schedule", middleware.AuthMiddleware(db, "Create_Shedule"), local_controllers.MakeTask("schecule", "post", redisDB, broker))
	ApiGroup.PATCH("/schedule", middleware.AuthMiddleware(db, "Update_Shedule"), local_controllers.MakeTask("schedule", "patch", redisDB, broker))
	ApiGroup.DELETE("/schedule", middleware.AuthMiddleware(db, "Delete_Shedule"), local_controllers.MakeTask("schedule", "delete", redisDB, broker))
}

func setupAnnounceRouter(ApiGroup *gin.RouterGroup, db *gorm.DB, redisDB *redis.Client, broker *models.Broker) {
	ApiGroup.GET("/announce", middleware.AuthMiddleware(db, ""), local_controllers.MakeTask("announces", "get", redisDB, broker))
	ApiGroup.POST("/announce", middleware.AuthMiddleware(db, "Create_Announces"), local_controllers.MakeTask("announces", "post", redisDB, broker))
	ApiGroup.PATCH("/announce", middleware.AuthMiddleware(db, "Update_Announces"), local_controllers.MakeTask("announces", "patch", redisDB, broker))
	ApiGroup.DELETE("/announce", middleware.AuthMiddleware(db, "Delete_Announces"), local_controllers.MakeTask("announces", "delete", redisDB, broker))
}

func setupOrderRouter(ApiGroup *gin.RouterGroup, db *gorm.DB, redisDB *redis.Client, broker *models.Broker) {
	ApiGroup.GET("/order", middleware.AuthMiddleware(db, "Read_Orders"), local_controllers.MakeTask("orders", "get", redisDB, broker))
	ApiGroup.POST("/order", middleware.AuthMiddleware(db, "Create_Orders"), local_controllers.MakeTask("orders", "post", redisDB, broker))
	ApiGroup.PATCH("/order", middleware.AuthMiddleware(db, "Update_Orders"), local_controllers.MakeTask("orders", "patch", redisDB, broker))
	ApiGroup.DELETE("/order", middleware.AuthMiddleware(db, "Delete_Orders"), local_controllers.MakeTask("orders", "delete", redisDB, broker))
}

func setupUserRouter(ApiGroup *gin.RouterGroup, db *gorm.DB, redisDB *redis.Client, broker *models.Broker) {
	ApiGroup.GET("/user", middleware.AuthMiddleware(db, "Change_Users"), local_controllers.MakeTask("users", "get", redisDB, broker))
	ApiGroup.POST("/user", middleware.AuthMiddleware(db, "Change_Users"), local_controllers.MakeTask("users", "post", redisDB, broker))
	ApiGroup.PATCH("/user", middleware.AuthMiddleware(db, "Change_Users"), local_controllers.MakeTask("users", "patch", redisDB, broker))
	ApiGroup.DELETE("/user", middleware.AuthMiddleware(db, "Change_Users"), local_controllers.MakeTask("users", "delete", redisDB, broker))
}

func setupRoleRouter(ApiGroup *gin.RouterGroup, db *gorm.DB, redisDB *redis.Client, broker *models.Broker) {
	ApiGroup.GET("/role", middleware.AuthMiddleware(db, "Change_Roles"), local_controllers.MakeTask("roles", "get", redisDB, broker))
	ApiGroup.POST("/role", middleware.AuthMiddleware(db, "Change_Roles"), local_controllers.MakeTask("roles", "post", redisDB, broker))
	ApiGroup.PATCH("/role", middleware.AuthMiddleware(db, "Change_Roles"), local_controllers.MakeTask("roles", "patch", redisDB, broker))
	ApiGroup.DELETE("/role", middleware.AuthMiddleware(db, "Change_Roles"), local_controllers.MakeTask("roles", "delete", redisDB, broker))
	ApiGroup.GET("/permission", middleware.AuthMiddleware(db, "Change_Roles"), local_controllers.MakeTask("permissions", "get", redisDB, broker))
}

func setupOtherRouter(ApiGroup *gin.RouterGroup, db *gorm.DB) {
	ApiGroup.GET("/logs", middleware.AuthMiddleware(db, "Download_Logs"), local_controllers.GetLogs)
	ApiGroup.GET("/log", middleware.AuthMiddleware(db, "Download_Logs"), local_controllers.GetLog)
}

func setupPickupRouter(ApiGroup *gin.RouterGroup, redisDB *redis.Client) {
	ApiGroup.GET("/pickup", local_controllers.Pickup(redisDB))
}

func SetupRoutes(router *gin.Engine, db *gorm.DB, redisDB *redis.Client, broker *models.Broker) {
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5001", "http://127.0.0.1:5001", "http://localhost:5174", "http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           8 * time.Hour,
	}))

	ApiGroup := router.Group("/api/")

	ApiGroup.POST("/login", local_controllers.Login(db))
	ApiGroup.POST("/refresh", local_controllers.RefreshToken(db))
	ApiGroup.POST("/logout", local_controllers.Logout())
	ApiGroup.GET("/image", shared_controllers.GetImage)

	setupGoodsRouter(ApiGroup, db, redisDB, broker)

	setupSheduleRouter(ApiGroup, db, redisDB, broker)

	setupAnnounceRouter(ApiGroup, db, redisDB, broker)

	setupOrderRouter(ApiGroup, db, redisDB, broker)

	setupUserRouter(ApiGroup, db, redisDB, broker)

	setupRoleRouter(ApiGroup, db, redisDB, broker)

	setupOtherRouter(ApiGroup, db)

	setupPickupRouter(ApiGroup, redisDB)
}

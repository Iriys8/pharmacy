package routes

import (
	public_controllers "pharmacy/public-api/controllers"
	shared_controllers "pharmacy/shared/controllers"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func setupGoodsRouter(GoodsGroup *gin.RouterGroup, db *gorm.DB) {
	GoodsGroup.GET("/goods", shared_controllers.GetGoods(db))
	GoodsGroup.GET("/goods/:id", shared_controllers.GetGoodsByID(db))
	GoodsGroup.GET("/goods/advert", public_controllers.GetPromoItems(db))
}

func setupWorkTimeRouter(WorkTimeGroup *gin.RouterGroup, db *gorm.DB) {
	WorkTimeGroup.GET("/schedule", public_controllers.GetWorkTimesDated(db))
}

func setupOrderRouter(OrderGroup *gin.RouterGroup, db *gorm.DB) {
	OrderGroup.POST("/order/:id", shared_controllers.CreateOrder(db))
}

func setupAnnounceRouter(AnnounceGroup *gin.RouterGroup, db *gorm.DB) {
	AnnounceGroup.GET("/announce", shared_controllers.GetAnnounces(db))
}

func SetupRoutes(router *gin.Engine, db *gorm.DB) {
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:5174"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           8 * time.Hour,
	}))

	ApiGroup := router.Group("/api/")

	ApiGroup.GET("/image", shared_controllers.GetImage)

	setupGoodsRouter(ApiGroup, db)

	setupWorkTimeRouter(ApiGroup, db)

	setupOrderRouter(ApiGroup, db)

	setupAnnounceRouter(ApiGroup, db)
}

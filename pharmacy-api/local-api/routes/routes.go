package routes

import (
	local_controllers "pharmacy-api/local-api/controllers"
	shared_controllers "pharmacy-api/shared/controllers"

	"pharmacy-api/local-api/middleware"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func setupGoodsRouter(GoodsGroup *gin.RouterGroup, db *gorm.DB) {
	GoodsGroup.GET("/goods", middleware.AuthMiddleware(db, ""), shared_controllers.GetGoods(db))
	GoodsGroup.GET("/goods/:id", middleware.AuthMiddleware(db, ""), shared_controllers.GetGoodsByID(db))
	GoodsGroup.PATCH("/goods/:id", middleware.AuthMiddleware(db, "Update_Goods"), local_controllers.UpdateGoods(db))
}

func setupWorkTimeRouter(SheduleGroup *gin.RouterGroup, db *gorm.DB) {
	SheduleGroup.GET("/schedule", middleware.AuthMiddleware(db, ""), local_controllers.GetWorkTimes(db))
	SheduleGroup.GET("/schedule/:id", middleware.AuthMiddleware(db, ""), local_controllers.GetWorkTimeByID(db))
	SheduleGroup.POST("/schedule/:id", middleware.AuthMiddleware(db, "Create_WorkTime"), local_controllers.CreateWorkTime(db))
	SheduleGroup.PATCH("/schedule/:id", middleware.AuthMiddleware(db, "Update_WorkTime"), local_controllers.UpdateWorkTime(db))
	SheduleGroup.DELETE("/schedule/:id", middleware.AuthMiddleware(db, "Delete_WorkTime"), local_controllers.DeleteWorkTime(db))
}

func setupAnnounceRouter(AnnounceGroup *gin.RouterGroup, db *gorm.DB) {
	AnnounceGroup.GET("/announce", middleware.AuthMiddleware(db, ""), shared_controllers.GetAnnounces(db))
	AnnounceGroup.GET("/announce/:id", middleware.AuthMiddleware(db, ""), local_controllers.GetAnnounceByID(db))
	AnnounceGroup.POST("/announce/:id", middleware.AuthMiddleware(db, "Create_Announces"), local_controllers.CreateAnnounce(db))
	AnnounceGroup.PATCH("/announce/:id", middleware.AuthMiddleware(db, "Update_Announces"), local_controllers.UpdateAnnounce(db))
	AnnounceGroup.DELETE("/announce/:id", middleware.AuthMiddleware(db, "Delete_Announces"), local_controllers.DeleteAnnounce(db))
}

func setupOrderRouter(OrderGroup *gin.RouterGroup, db *gorm.DB) {
	OrderGroup.GET("/order", middleware.AuthMiddleware(db, "Read_Orders"), local_controllers.GetOrders(db))
	OrderGroup.GET("/order/:id", middleware.AuthMiddleware(db, "Read_Orders"), local_controllers.GetOrderByID(db))
	OrderGroup.POST("/order/:id", middleware.AuthMiddleware(db, "Create_Orders"), shared_controllers.CreateOrder(db))
	OrderGroup.PATCH("/order/:id", middleware.AuthMiddleware(db, "Update_Orders"), local_controllers.UpdateOrder(db))
	OrderGroup.DELETE("/order/:id", middleware.AuthMiddleware(db, "Delete_Orders"), local_controllers.DeleteOrder(db))
}

func setupUserRouter(UserGroup *gin.RouterGroup, db *gorm.DB) {
	UserGroup.GET("/user", middleware.AuthMiddleware(db, "Change_Users"), local_controllers.GetUsers(db))
	UserGroup.GET("/user/:id", middleware.AuthMiddleware(db, "Change_Users"), local_controllers.GetUserByID(db))
	UserGroup.POST("/user/:id", middleware.AuthMiddleware(db, "Change_Users"), local_controllers.CreateUser(db))
	UserGroup.PATCH("/user/:id", middleware.AuthMiddleware(db, "Change_Users"), local_controllers.UpdateUser(db))
	UserGroup.DELETE("/user/:id", middleware.AuthMiddleware(db, "Change_Users"), local_controllers.DeleteUser(db))
}

func setupRoleRouter(RoleGroup *gin.RouterGroup, db *gorm.DB) {
	RoleGroup.GET("/role", middleware.AuthMiddleware(db, "Change_Roles"), local_controllers.GetRoles(db))
	RoleGroup.GET("/role/:id", middleware.AuthMiddleware(db, "Change_Roles"), local_controllers.GetRoleByID(db))
	RoleGroup.POST("/role/:id", middleware.AuthMiddleware(db, "Change_Roles"), local_controllers.CreateRole(db))
	RoleGroup.PATCH("/role/:id", middleware.AuthMiddleware(db, "Change_Roles"), local_controllers.UpdateRole(db))
	RoleGroup.DELETE("/role/:id", middleware.AuthMiddleware(db, "Change_Roles"), local_controllers.DeleteRole(db))
}

func setupOtherRouter(OtherGroup *gin.RouterGroup, db *gorm.DB) {
	OtherGroup.GET("/permission", middleware.AuthMiddleware(db, "Change_Roles"), local_controllers.GetPermissions(db))
	OtherGroup.GET("/logs", middleware.AuthMiddleware(db, "Download_Logs"), local_controllers.GetLogs)
	OtherGroup.GET("/log", middleware.AuthMiddleware(db, "Download_Logs"), local_controllers.GetLog)
}

func SetupRoutes(router *gin.Engine, db *gorm.DB) {
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5001", "http://127.0.0.1:5001"},
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

	setupGoodsRouter(ApiGroup, db)

	setupWorkTimeRouter(ApiGroup, db)

	setupAnnounceRouter(ApiGroup, db)

	setupOrderRouter(ApiGroup, db)

	setupUserRouter(ApiGroup, db)

	setupRoleRouter(ApiGroup, db)

	setupOtherRouter(ApiGroup, db)
}

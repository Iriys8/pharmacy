package setup

import (
	"log"
	"pharmacy-api/shared/models"

	"gorm.io/gorm"
)

func SetupDB(db *gorm.DB) {
	if err := db.AutoMigrate(
		&models.Schedule{},
		&models.Producer{},
		&models.Tag{},
		&models.GoodsOrders{},
		&models.Order{},
		&models.Goods{},
		&models.Announcement{},
		&models.User{},
		&models.Role{},
		&models.Permission{},
	); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	// MYSQL or GORM moment idk, on Postgres all alright
	db.Exec("ALTER TABLE schedules MODIFY COLUMN time_start TIME")
	db.Exec("ALTER TABLE schedules MODIFY COLUMN time_end TIME")
	db.Exec("ALTER TABLE schedules MODIFY COLUMN date DATE")

	var persmissionsAmount int64
	db.Model(&models.Permission{}).Count(&persmissionsAmount)
	if persmissionsAmount == 0 {
		permissions := []models.Permission{
			{Action: "Update_Goods"},
			{Action: "Read_Orders"},
			{Action: "Update_Orders"},
			{Action: "Create_Orders"},
			{Action: "Delete_Orders"},
			{Action: "Update_Schedule"},
			{Action: "Create_Schedule"},
			{Action: "Delete_Schedule"},
			{Action: "Update_Announces"},
			{Action: "Create_Announces"},
			{Action: "Delete_Announces"},
			{Action: "Change_Users"},
			{Action: "Change_Roles"},
			{Action: "Download_Logs"},
		}
		db.Create(&permissions)
	}

	var itemsAmount int64
	db.Model(&models.Goods{}).Count(&itemsAmount)
	if itemsAmount == 0 {
		test_data(db)
		log.Println("TEST DATA USED!")
	}
}

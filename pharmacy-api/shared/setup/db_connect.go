package setup

import (
	"log"
	"os"

	shared_models "pharmacy-api/shared/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectDB() (db *gorm.DB) {
	dsn := os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASSWORD") + "@tcp(" + os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT") + ")/" + os.Getenv("DB_NAME") + "?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(
		&shared_models.WorkTime{},
		&shared_models.Producer{},
		&shared_models.Tag{},
		&shared_models.GoodsOrders{},
		&shared_models.Order{},
		&shared_models.Goods{},
		&shared_models.Announcement{},
		&shared_models.User{},
		&shared_models.Role{},
		&shared_models.Permission{},
	); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	// MYSQL or GORM moment idk, on Postgres all alright
	db.Exec("ALTER TABLE work_times MODIFY COLUMN time_start TIME")
	db.Exec("ALTER TABLE work_times MODIFY COLUMN time_end TIME")
	db.Exec("ALTER TABLE work_times MODIFY COLUMN date DATE")

	var persmissionsAmount int64
	db.Model(&shared_models.Permission{}).Count(&persmissionsAmount)
	if persmissionsAmount == 0 {
		permissions := []shared_models.Permission{
			{Action: "Update_Goods"},
			{Action: "Read_Orders"},
			{Action: "Update_Orders"},
			{Action: "Create_Orders"},
			{Action: "Delete_Orders"},
			{Action: "Update_WorkTime"},
			{Action: "Create_WorkTime"},
			{Action: "Delete_WorkTime"},
			{Action: "Update_Announces"},
			{Action: "Create_Announces"},
			{Action: "Delete_Announces"},
			{Action: "Change_Users"},
			{Action: "Change_Roles"},
			{Action: "Download_Logs"},
		}
		db.Create(&permissions)

	}

	return
}

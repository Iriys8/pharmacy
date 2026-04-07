package main

import (
	"log"
	"os"

	local_models "pharmacy-api/local-api/models"
	shared_models "pharmacy-api/shared/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func setupDatabase() (db *gorm.DB) {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	dsn := dbUser + ":" + dbPassword + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	log.Println("Database connected")

	if err := db.AutoMigrate(
		&shared_models.WorkTime{},
		&shared_models.Producer{},
		&shared_models.Tag{},
		&shared_models.GoodsOrders{},
		&shared_models.Order{},
		&shared_models.Goods{},
		&shared_models.Announcement{},
		&local_models.User{},
		&local_models.Role{},
		&local_models.Permission{},
	); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	// MYSQL or GORM moment idk, on Postgres all alright
	db.Exec("ALTER TABLE work_times MODIFY COLUMN time_start TIME")
	db.Exec("ALTER TABLE work_times MODIFY COLUMN time_end TIME")
	db.Exec("ALTER TABLE work_times MODIFY COLUMN date DATE")

	var persmissionsAmount int64
	db.Model(&local_models.Permission{}).Count(&persmissionsAmount)
	if persmissionsAmount == 0 {
		permissions := []local_models.Permission{
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

func setupAdmin(db *gorm.DB) {
	pass, err := bcrypt.GenerateFromPassword([]byte(os.Getenv("CONTROL_PANEL_ADMIN_PASSWORD")), bcrypt.DefaultCost)

	if err != nil {
		log.Fatalf("failed to generate admin password: %v", err)
	}

	var permissions []local_models.Permission
	db.Find(&permissions)

	adminRole := local_models.Role{
		Name:        "Admin",
		Permissions: permissions,
	}

	user := local_models.User{
		Login:        os.Getenv("CONTROL_PANEL_ADMIN_LOGIN"),
		Role:         adminRole,
		UserName:     "Admin",
		PasswordHash: pass,
	}
	db.Create(&user)
}

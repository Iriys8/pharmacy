package controllers

import (
	"log"
	models "pharmacy-api/shared/models"
	"strconv"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func GetUsers(db *gorm.DB, query string, pageStr string, limitStr string, claims models.Claims) (result map[string]any, err error) {

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	} else if limit > 40 {
		limit = 40
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	offset := (page - 1) * limit

	var users []models.User
	var totalCount int64

	if query != "" {
		if err = db.Preload("Role").Where("user_name LIKE ?", "%"+query+"%").Find(&users).Error; err != nil {
			log.Println("Users GET error [" + claims.Username + "]" + err.Error())
			return
		}
		totalCount = int64(len(users))
		start := offset
		end := offset + limit
		if start > len(users) {
			start = len(users)
		}
		if end > len(users) {
			end = len(users)
		}
		users = users[start:end]
	} else {
		db.Model(&models.User{}).Count(&totalCount)
		if err = db.Preload("Role").Order("id DESC").Limit(limit).Offset(offset).Find(&users).Error; err != nil {
			log.Println("Users GET error [" + claims.Username + "]" + err.Error())
			return
		}
	}

	totalPages := (totalCount + int64(limit) - 1) / int64(limit)
	log.Println("Users GET [" + claims.Username + "]")

	result = make(map[string]any)
	result["Items"] = users
	result["TotalPages"] = totalPages
	result["CurrentPage"] = page
	return
}

func GetUserByID(db *gorm.DB, id int, claims models.Claims) (result map[string]any, err error) {
	var response models.User

	if id != 0 {
		if err = db.Preload("Role").First(&response, id).Error; err != nil {
			return
		}

		response.PasswordHash = []byte("")
	}
	var roles []models.Role
	var requstedUser models.User
	if err = db.Preload("Role.Permissions").First(&requstedUser, claims.UserID).Error; err != nil {
		log.Println("User GET error [" + claims.Username + "]" + err.Error())
		return
	}
	for _, permission := range requstedUser.Role.Permissions {
		if permission.Action == "Change_Roles" {
			if err = db.Find(&roles).Error; err != nil {
				return
			}
			break
		}
	}

	log.Println("User GET [" + claims.Username + "]")

	resultData := make(map[string]any)
	resultData["User"] = response
	resultData["Roles"] = roles

	result = make(map[string]any)
	result["Response"] = resultData
	return
}

func CreateUser(db *gorm.DB, newUserRequest models.UserUpdateRequest, claims models.Claims) (result map[string]any, err error) {

	pass, err := bcrypt.GenerateFromPassword([]byte(newUserRequest.Password), bcrypt.DefaultCost)

	if err != nil {
		log.Println("User POST error [" + claims.Username + "]" + err.Error())
		return
	}
	newUser := models.User{
		Login:        newUserRequest.Login,
		UserName:     newUserRequest.UserName,
		RoleID:       newUserRequest.RoleID,
		PasswordHash: pass,
	}

	if err = db.Create(&newUser).Error; err != nil {
		log.Println("User POST error [" + claims.Username + "]" + err.Error())
		return
	}

	log.Println("User POST [" + claims.Username + "]")

	result = make(map[string]any)
	result["Response"] = "schedule created"
	return
}

func DeleteUser(db *gorm.DB, id int, claims models.Claims) (result map[string]any, err error) {
	if err = db.Delete(&models.User{}, id).Error; err != nil {
		log.Println("User DELETE error [" + claims.Username + "]" + err.Error())
		return
	}

	log.Println("User DELETE [" + claims.Username + "]")

	result = make(map[string]any)
	result["Response"] = "schedule deleted"
	return
}

func UpdateUser(db *gorm.DB, id int, userUpdateRequest models.UserUpdateRequest, claims models.Claims) (result map[string]any, err error) {
	var updateData models.User

	if userUpdateRequest.Password != "" {
		var pass []byte
		pass, err = bcrypt.GenerateFromPassword([]byte(userUpdateRequest.Password), bcrypt.DefaultCost)

		if err != nil {
			log.Println("User PATCH error [" + claims.Username + "]" + err.Error())
			return
		}
		updateData = models.User{
			Login:        userUpdateRequest.Login,
			UserName:     userUpdateRequest.UserName,
			RoleID:       userUpdateRequest.RoleID,
			PasswordHash: pass,
		}
	} else {
		updateData = models.User{
			Login:    userUpdateRequest.Login,
			UserName: userUpdateRequest.UserName,
			RoleID:   userUpdateRequest.RoleID,
		}
	}

	if err = db.Model(&models.User{}).Where("id = ?", id).Updates(updateData).Error; err != nil {
		log.Println("User PATCH error [" + claims.Username + "]" + err.Error())
		return
	}

	log.Println("User PATCH [" + claims.Username + "]")

	result = make(map[string]any)
	result["Response"] = "schedule deleted"
	return
}

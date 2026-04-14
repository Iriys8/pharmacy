package controllers

import (
	"log"
	models "pharmacy-api/shared/models"
	"strconv"

	"gorm.io/gorm"
)

func GetRoles(db *gorm.DB, query string, pageStr string, limitStr string, claims models.Claims) (result map[string]any, err error) {
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

	var roles []models.Role
	var totalCount int64

	if query != "" {
		if err = db.Where("name LIKE ?", "%"+query+"%").Find(&roles).Error; err != nil {
			log.Println("Roles GET error [" + claims.Username + "]" + err.Error())
			return
		}
		totalCount = int64(len(roles))
		start := offset
		end := offset + limit
		if start > len(roles) {
			start = len(roles)
		}
		if end > len(roles) {
			end = len(roles)
		}
		roles = roles[start:end]
	} else {
		db.Model(&models.Role{}).Count(&totalCount)
		if err = db.Preload("Permissions").Order("id DESC").Limit(limit).Offset(offset).Find(&roles).Error; err != nil {
			log.Println("Roles GET error [" + claims.Username + "]" + err.Error())
			return
		}
	}

	totalPages := (totalCount + int64(limit) - 1) / int64(limit)
	log.Println("Roles GET [" + claims.Username + "]")

	result = make(map[string]any)
	result["Items"] = roles
	result["TotalPages"] = totalPages
	result["CurrentPage"] = page
	return
}

func GetRoleByID(db *gorm.DB, id int, claims models.Claims) (result map[string]any, err error) {
	var role models.Role

	if err = db.Preload("Permissions").First(&role, id).Error; err != nil {
		return
	}

	log.Println("Role GET [" + claims.Username + "]")

	result = make(map[string]any)
	result["Response"] = role
	return
}

func CreateRole(db *gorm.DB, role models.Role, claims models.Claims) (result map[string]any, err error) {
	if err = db.Create(&role).Error; err != nil {
		log.Println("Role POST error [" + claims.Username + "]" + err.Error())
		return
	}

	if err = db.Model(&role).Association("Permissions").Replace(role.Permissions); err != nil {
		log.Println("Role POST error [" + claims.Username + "]" + err.Error())
		return
	}

	log.Println("Role POST [" + claims.Username + "]")

	result = make(map[string]any)
	result["Response"] = "role created"
	return
}

func DeleteRole(db *gorm.DB, id int, claims models.Claims) (result map[string]any, err error) {
	var role models.Role
	if err = db.Preload("Permissions").First(&role, id).Error; err != nil {
		return
	}

	if err = db.Model(&role).Association("Permissions").Clear(); err != nil {
		log.Println("Role DELETE error [" + claims.Username + "]" + err.Error())
		return
	}

	if err = db.Delete(&role).Error; err != nil {
		log.Println("Role DELETE error [" + claims.Username + "]" + err.Error())
		return
	}

	log.Println("Role DELETE [" + claims.Username + "]")

	result = make(map[string]any)
	result["Response"] = "role deleted"
	return
}

func UpdateRole(db *gorm.DB, id int, role models.Role, claims models.Claims) (result map[string]any, err error) {
	var existingRole models.Role
	if err = db.Preload("Permissions").First(&existingRole, id).Error; err != nil {
		return
	}

	if err = db.Model(&existingRole).Updates(models.Role{Name: role.Name}).Error; err != nil {
		log.Println("Role PATCH error [" + claims.Username + "]" + err.Error())
		return
	}

	if err = db.Model(&existingRole).Association("Permissions").Replace(role.Permissions); err != nil {
		log.Println("Role PATCH error [" + claims.Username + "]" + err.Error())
		return
	}

	log.Println("Role PATCH [" + claims.Username + "]")

	result = make(map[string]any)
	result["Response"] = "role updated"
	return
}

func GetPermissions(db *gorm.DB, query string, pageStr string, limitStr string, claims models.Claims) (result map[string]any, err error) {
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

	var permissions []models.Permission
	var totalCount int64

	if query != "" {
		if err = db.Where("action LIKE ?", "%"+query+"%").Find(&permissions).Error; err != nil {
			log.Println("Permissions GET error [" + claims.Username + "]" + err.Error())
			return
		}
		totalCount = int64(len(permissions))
		start := offset
		end := offset + limit
		if start > len(permissions) {
			start = len(permissions)
		}
		if end > len(permissions) {
			end = len(permissions)
		}
		permissions = permissions[start:end]
	} else {
		db.Model(&models.Permission{}).Count(&totalCount)
		if err = db.Order("id DESC").Limit(limit).Offset(offset).Find(&permissions).Error; err != nil {
			log.Println("Permissions GET error [" + claims.Username + "]" + err.Error())
			return
		}
	}

	totalPages := (totalCount + int64(limit) - 1) / int64(limit)
	log.Println("Permissions GET [" + claims.Username + "]")

	result = make(map[string]any)
	result["Items"] = permissions
	result["TotalPages"] = totalPages
	result["CurrentPage"] = page
	return
}

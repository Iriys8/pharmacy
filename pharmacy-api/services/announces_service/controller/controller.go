package controller

import (
	"log"
	models "pharmacy-api/shared/models"
	"strconv"
	"time"

	"gorm.io/gorm"
)

func GetAnnounces(db *gorm.DB, query string, pageStr string, limitStr string) (result map[string]any, err error) {

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

	var Announce []models.Announcement
	var totalCount int64
	var response []models.AnnouncementResponse

	if query != "" {
		if err = db.Where("date_time LIKE ?", "%"+query+"%").Find(&Announce).Error; err != nil {
			return
		}
		totalCount = int64(len(Announce))
		start := offset
		end := offset + limit
		if start > len(Announce) {
			start = len(Announce)
		}
		if end > len(Announce) {
			end = len(Announce)
		}
		Announce = Announce[start:end]
	} else {
		db.Model(&models.Announcement{}).Count(&totalCount)
		if err = db.Order("date_time ASC").Limit(limit).Offset(offset).Find(&Announce).Error; err != nil {
			return
		}
	}

	for _, announce := range Announce {
		AnnounceResponse := models.AnnouncementResponse{
			ID:       announce.ID,
			DateTime: announce.DateTime.String()[:16],
			From:     announce.From,
			Announce: announce.Announce,
		}
		response = append(response, AnnounceResponse)
	}

	totalPages := (totalCount + int64(limit) - 1) / int64(limit)

	result = make(map[string]any)
	result["Items"] = response
	result["TotalPages"] = totalPages
	result["CurrentPage"] = page
	return
}

func GetAnnounceByID(db *gorm.DB, id int, claims models.Claims) (result map[string]any, err error) {

	var Announce models.Announcement
	if err = db.Find(&Announce, id).Error; err != nil {
		return
	}
	AnnounceResponse := models.AnnouncementResponse{
		ID:       Announce.ID,
		DateTime: Announce.DateTime.String()[:16],
		From:     Announce.From,
		Announce: Announce.Announce,
	}
	log.Println("Announce GET [" + claims.Username + "]")

	result = make(map[string]any)
	result["Response"] = AnnounceResponse
	return
}

func CreateAnnounce(db *gorm.DB, announcementResponse models.AnnouncementResponse, claims models.Claims) (result map[string]any, err error) {
	now := time.Now().Local()
	dateTime := time.Date(
		now.Year(), now.Month(), now.Day(),
		now.Hour(), now.Minute(), 0, 0, time.Local,
	)

	announce := models.Announcement{
		DateTime: dateTime,
		From:     claims.Username,
		Announce: announcementResponse.Announce,
	}

	if err = db.Create(&announce).Error; err != nil {
		log.Println("Announce POST error [" + claims.Username + "]" + err.Error())
		return
	}
	log.Println("Announce POST [" + claims.Username + "]")

	result = make(map[string]any)
	result["Response"] = "announce created"
	return
}

func UpdateAnnounce(db *gorm.DB, id int, announcementResponse models.AnnouncementResponse, claims models.Claims) (result map[string]any, err error) {
	announce := models.Announcement{
		ID:       announcementResponse.ID,
		Announce: announcementResponse.Announce,
	}

	if err = db.Model(&models.Announcement{}).Where("id = ?", id).Updates(announce).Error; err != nil {
		log.Println("Announce PATCH error [" + claims.Username + "]" + err.Error())
		return
	}
	log.Println("Announce PATCH [" + claims.Username + "]")

	result = make(map[string]any)
	result["Response"] = "announce updated"
	return
}

func DeleteAnnounce(db *gorm.DB, id int, claims models.Claims) (result map[string]any, err error) {
	if err = db.Delete(&models.Announcement{}, id).Error; err != nil {
		log.Println("Announce DELETE error [" + claims.Username + "]" + err.Error())
		return
	}
	log.Println("Announce DELETE [" + claims.Username + "]")

	result = make(map[string]any)
	result["Response"] = "announce deleted"
	return
}

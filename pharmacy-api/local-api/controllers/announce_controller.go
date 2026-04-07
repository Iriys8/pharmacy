package controllers

import (
	"log"
	"net/http"
	local_models "pharmacy-api/local-api/models"
	shared_models "pharmacy-api/shared/models"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetAnnounceByID(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		user, _ := c.Get("user")
		claims := user.(*local_models.Claims)

		var Announce shared_models.Announcement
		if err := db.Find(&Announce, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		AnnounceResponse := shared_models.AnnouncementResponse{
			ID:       Announce.ID,
			DateTime: Announce.DateTime.String()[:16],
			From:     Announce.From,
			Announce: Announce.Announce,
		}
		log.Println("Announce GET [" + claims.Username + "]")
		c.JSON(http.StatusOK, AnnounceResponse)
	}
}

func CreateAnnounce(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var announcementResponse shared_models.AnnouncementResponse
		user, _ := c.Get("user")
		claims := user.(*local_models.Claims)

		if err := c.ShouldBindJSON(&announcementResponse); err != nil {
			log.Println("Announce POST error [" + claims.Username + "]" + err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		now := time.Now().Local()
		dateTime := time.Date(
			now.Year(), now.Month(), now.Day(),
			now.Hour(), now.Minute(), 0, 0, time.Local,
		)

		announce := shared_models.Announcement{
			DateTime: dateTime,
			From:     claims.Username,
			Announce: announcementResponse.Announce,
		}

		if err := db.Create(&announce).Error; err != nil {
			log.Println("Announce POST error [" + claims.Username + "]" + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		log.Println("Announce POST [" + claims.Username + "]")
		c.JSON(http.StatusCreated, gin.H{"message": "Announce created", "data": announce})
	}
}

func UpdateAnnounce(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var announcementResponse shared_models.AnnouncementResponse

		user, _ := c.Get("user")
		claims := user.(*local_models.Claims)

		if err := c.ShouldBindJSON(&announcementResponse); err != nil {
			log.Println("Announce PATCH error [" + claims.Username + "]" + err.Error())
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		announce := shared_models.Announcement{
			ID:       announcementResponse.ID,
			Announce: announcementResponse.Announce,
		}

		if err := db.Model(&shared_models.Announcement{}).Where("id = ?", id).Updates(announce).Error; err != nil {
			log.Println("Announce PATCH error [" + claims.Username + "]" + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		log.Println("Announce PATCH [" + claims.Username + "]")
		c.JSON(http.StatusOK, gin.H{"message": "Announce updated", "data": announce})
	}
}

func DeleteAnnounce(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		user, _ := c.Get("user")
		claims := user.(*local_models.Claims)

		if err := db.Delete(&shared_models.Announcement{}, id).Error; err != nil {
			log.Println("Announce DELETE error [" + claims.Username + "]" + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		log.Println("Announce DELETE [" + claims.Username + "]")
		c.JSON(http.StatusOK, gin.H{"message": "Announce deleted"})
	}
}

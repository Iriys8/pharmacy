package controllers

import (
	"net/http"
	shared_models "pharmacy/shared/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetAnnounces(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := c.Query("q")
		pageStr := c.Query("page")
		limitStr := c.Query("limit")

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

		var Announce []shared_models.Announcement
		var totalCount int64
		var response []shared_models.AnnouncementResponse

		if query != "" {
			if err := db.Where("date_time LIKE ?", "%"+query+"%").Find(&Announce).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
			db.Model(&shared_models.Announcement{}).Count(&totalCount)
			if err := db.Order("date_time ASC").Limit(limit).Offset(offset).Find(&Announce).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}

		for _, announce := range Announce {
			AnnounceResponse := shared_models.AnnouncementResponse{
				ID:       announce.ID,
				DateTime: announce.DateTime.String()[:16],
				From:     announce.From,
				Announce: announce.Announce,
			}
			response = append(response, AnnounceResponse)
		}

		totalPages := (totalCount + int64(limit) - 1) / int64(limit)
		c.JSON(http.StatusOK, gin.H{
			"Items":       response,
			"TotalPages":  totalPages,
			"CurrentPage": page,
		})
	}
}

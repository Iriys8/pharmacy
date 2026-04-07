package controllers

import (
	"database/sql"
	"log"
	"net/http"
	local_models "pharmacy-api/local-api/models"
	shared_models "pharmacy-api/shared/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetWorkTimes(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := c.Query("q")
		pageStr := c.Query("page")
		limitStr := c.Query("limit")

		user, _ := c.Get("user")
		claims := user.(*local_models.Claims)

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

		var workTimes []shared_models.WorkTime
		var totalCount int64
		var response []shared_models.WorkTimeResponse

		if query != "" {
			if err := db.Where("date LIKE ?", "%"+query+"%").Find(&workTimes).Error; err != nil {
				log.Println("WorkTimes GET error [" + claims.Username + "]" + err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			totalCount = int64(len(workTimes))
			start := offset
			end := offset + limit
			if start > len(workTimes) {
				start = len(workTimes)
			}
			if end > len(workTimes) {
				end = len(workTimes)
			}
			workTimes = workTimes[start:end]
		} else {
			db.Model(&shared_models.WorkTime{}).Count(&totalCount)
			if err := db.Order("date ASC").Limit(limit).Offset(offset).Find(&workTimes).Error; err != nil {
				log.Println("WorkTimes GET error [" + claims.Username + "]" + err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
		for _, worktime := range workTimes {
			workTimeResponse := shared_models.WorkTimeResponse{
				ID:        worktime.ID,
				Date:      worktime.Date.String[:10],
				TimeStart: worktime.TimeStart.String[:5],
				TimeEnd:   worktime.TimeEnd.String[:5],
				IsOpened:  worktime.IsOpened,
			}
			response = append(response, workTimeResponse)
		}

		totalPages := (totalCount + int64(limit) - 1) / int64(limit)
		log.Println("WorkTimes GET [" + claims.Username + "]")
		c.JSON(http.StatusOK, gin.H{
			"Items":       response,
			"TotalPages":  totalPages,
			"CurrentPage": page,
		})
	}
}

func GetWorkTimeByID(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var workTime shared_models.WorkTime

		user, _ := c.Get("user")
		claims := user.(*local_models.Claims)

		if err := db.Find(&workTime, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
			return
		}
		workTimeResponse := shared_models.WorkTimeResponse{
			ID:        workTime.ID,
			Date:      workTime.Date.String[:10],
			TimeStart: workTime.TimeStart.String[:5],
			TimeEnd:   workTime.TimeEnd.String[:5],
			IsOpened:  workTime.IsOpened,
		}
		log.Println("WorkTime GET [" + claims.Username + "]")
		c.JSON(http.StatusOK, workTimeResponse)
	}
}

func CreateWorkTime(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var workTimeString shared_models.WorkTimeResponse

		user, _ := c.Get("user")
		claims := user.(*local_models.Claims)

		if err := c.ShouldBindJSON(&workTimeString); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		workTime := shared_models.WorkTime{
			IsOpened:  workTimeString.IsOpened,
			Date:      sql.NullString{String: workTimeString.Date, Valid: true},
			TimeStart: sql.NullString{String: workTimeString.TimeStart, Valid: true},
			TimeEnd:   sql.NullString{String: workTimeString.TimeEnd, Valid: true},
		}

		if err := db.Create(&workTime).Error; err != nil {
			log.Println("WorkTime POST error [" + claims.Username + "]" + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		log.Println("WorkTime POST [" + claims.Username + "]")
		c.JSON(http.StatusCreated, gin.H{"message": "workTime created", "data": workTime})
	}
}

func UpdateWorkTime(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var workTimeString shared_models.WorkTimeResponse

		user, _ := c.Get("user")
		claims := user.(*local_models.Claims)

		if err := c.ShouldBindJSON(&workTimeString); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		workTime := shared_models.WorkTime{
			IsOpened:  workTimeString.IsOpened,
			Date:      sql.NullString{String: workTimeString.Date, Valid: true},
			TimeStart: sql.NullString{String: workTimeString.TimeStart, Valid: true},
			TimeEnd:   sql.NullString{String: workTimeString.TimeEnd, Valid: true},
		}
		if err := db.Model(&shared_models.WorkTime{}).Where("id = ?", id).Updates(workTime).Error; err != nil {
			log.Println("WorkTime PATCH error [" + claims.Username + "]" + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		log.Println("WorkTime PATCH [" + claims.Username + "]")
		c.JSON(http.StatusOK, gin.H{"message": "workTime updated", "data": workTime})
	}
}

func DeleteWorkTime(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		user, _ := c.Get("user")
		claims := user.(*local_models.Claims)

		if err := db.Delete(&shared_models.WorkTime{}, id).Error; err != nil {
			log.Println("WorkTime DELETE error [" + claims.Username + "]" + err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		log.Println("WorkTime DELETE [" + claims.Username + "]")
		c.JSON(http.StatusOK, gin.H{"message": "workTime deleted"})
	}
}

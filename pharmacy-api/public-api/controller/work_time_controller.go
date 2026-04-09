package controller

import (
	"net/http"
	shared_models "pharmacy-api/shared/models"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetWorkTimesDated(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var workTimes []shared_models.WorkTime
		var response []shared_models.WorkTimeResponse

		startDate := c.Query("start")
		endDate := c.Query("end")

		if startDate == "" || endDate == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Start and end dates are required"})
			return
		}

		startParsed, err := time.Parse("2006-01-02", startDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format"})
			return
		}
		endParsed, err := time.Parse("2006-01-02", endDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format"})
			return
		}

		if endParsed.Before(startParsed) || endParsed.Compare(startParsed.Add(time.Hour*24*31)) == 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format"})
			return
		}

		if err := db.Where("date BETWEEN ? AND ?", startDate, endDate).Find(&workTimes).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		for current := startParsed; !current.After(endParsed); current = current.AddDate(0, 0, 1) {
			dateStr := current.Format("2006-01-02")
			var found bool

			for _, wt := range workTimes {
				wtDate, err := time.Parse("2006-01-02", wt.Date.String[:10])
				if err != nil {
					continue
				}
				if wtDate.Equal(current) {
					response = append(response, shared_models.WorkTimeResponse{
						ID:        wt.ID,
						Date:      dateStr,
						TimeStart: wt.TimeStart.String[:5],
						TimeEnd:   wt.TimeEnd.String[:5],
						IsOpened:  wt.IsOpened,
					})
					found = true
					break

				}
			}
			if !found {
				response = append(response, shared_models.WorkTimeResponse{
					ID:        0,
					Date:      dateStr,
					TimeStart: "8:00",
					TimeEnd:   "22:00",
					IsOpened:  true,
				})
			}
		}

		c.JSON(http.StatusOK, response)
	}
}

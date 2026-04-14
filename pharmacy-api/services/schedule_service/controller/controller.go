package controller

import (
	"database/sql"
	"log"
	"strconv"
	"time"

	models "pharmacy-api/shared/models"

	"gorm.io/gorm"
)

func GetScheduleDated(db *gorm.DB, startDate string, endDate string) (result map[string]any, err error) {
	var schedule []models.Schedule
	var response []models.ScheduleResponse

	startParsed, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return
	}
	endParsed, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return
	}

	if endParsed.Before(startParsed) || endParsed.Compare(startParsed.Add(time.Hour*24*31)) == 1 {
		return
	}

	if err = db.Where("date BETWEEN ? AND ?", startDate, endDate).Find(&schedule).Error; err != nil {
		return
	}

	for current := startParsed; !current.After(endParsed); current = current.AddDate(0, 0, 1) {
		dateStr := current.Format("2006-01-02")
		var found bool

		for _, wt := range schedule {
			wtDate, err := time.Parse("2006-01-02", wt.Date.String[:10])
			if err != nil {
				continue
			}
			if wtDate.Equal(current) {
				response = append(response, models.ScheduleResponse{
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
			response = append(response, models.ScheduleResponse{
				ID:        0,
				Date:      dateStr,
				TimeStart: "8:00",
				TimeEnd:   "22:00",
				IsOpened:  true,
			})
		}
	}

	result = make(map[string]any)
	result["Response"] = response
	return
}

func GetSchedule(db *gorm.DB, query string, pageStr string, limitStr string, claims models.Claims) (result map[string]any, err error) {

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

	var schedule []models.Schedule
	var totalCount int64
	var response []models.ScheduleResponse

	if query != "" {
		if err = db.Where("date LIKE ?", "%"+query+"%").Find(&schedule).Error; err != nil {
			log.Println("Schedule GET error [" + claims.Username + "]" + err.Error())
			return
		}
		totalCount = int64(len(schedule))
		start := offset
		end := offset + limit
		if start > len(schedule) {
			start = len(schedule)
		}
		if end > len(schedule) {
			end = len(schedule)
		}
		schedule = schedule[start:end]
	} else {
		db.Model(&models.Schedule{}).Count(&totalCount)
		if err = db.Order("date ASC").Limit(limit).Offset(offset).Find(&schedule).Error; err != nil {
			log.Println("Schedule GET error [" + claims.Username + "]" + err.Error())
			return
		}
	}
	for _, schedule_iterable := range schedule {
		scheduleResponse := models.ScheduleResponse{
			ID:        schedule_iterable.ID,
			Date:      schedule_iterable.Date.String[:10],
			TimeStart: schedule_iterable.TimeStart.String[:5],
			TimeEnd:   schedule_iterable.TimeEnd.String[:5],
			IsOpened:  schedule_iterable.IsOpened,
		}
		response = append(response, scheduleResponse)
	}

	totalPages := (totalCount + int64(limit) - 1) / int64(limit)
	log.Println("Schedule GET [" + claims.Username + "]")

	result = make(map[string]any)
	result["Items"] = response
	result["TotalPages"] = totalPages
	result["CurrentPage"] = page
	return
}

func GetScheduleByID(db *gorm.DB, id int, claims models.Claims) (result map[string]any, err error) {
	var schedule models.Schedule

	if err = db.Find(&schedule, id).Error; err != nil {
		return
	}
	scheduleResponse := models.ScheduleResponse{
		ID:        schedule.ID,
		Date:      schedule.Date.String[:10],
		TimeStart: schedule.TimeStart.String[:5],
		TimeEnd:   schedule.TimeEnd.String[:5],
		IsOpened:  schedule.IsOpened,
	}
	log.Println("Schedule GET [" + claims.Username + "]")

	result = make(map[string]any)
	result["Response"] = scheduleResponse
	return
}

func CreateSchedule(db *gorm.DB, scheduleString models.ScheduleResponse, claims models.Claims) (result map[string]any, err error) {
	schedule := models.Schedule{
		IsOpened:  scheduleString.IsOpened,
		Date:      sql.NullString{String: scheduleString.Date, Valid: true},
		TimeStart: sql.NullString{String: scheduleString.TimeStart, Valid: true},
		TimeEnd:   sql.NullString{String: scheduleString.TimeEnd, Valid: true},
	}

	if err = db.Create(&schedule).Error; err != nil {
		log.Println("Schedule POST error [" + claims.Username + "]" + err.Error())
		return
	}
	log.Println("Schedule POST [" + claims.Username + "]")

	result = make(map[string]any)
	result["Response"] = "schedule created"
	return
}

func UpdateSchedule(db *gorm.DB, id int, scheduleString models.ScheduleResponse, claims models.Claims) (result map[string]any, err error) {
	schedule := models.Schedule{
		IsOpened:  scheduleString.IsOpened,
		Date:      sql.NullString{String: scheduleString.Date, Valid: true},
		TimeStart: sql.NullString{String: scheduleString.TimeStart, Valid: true},
		TimeEnd:   sql.NullString{String: scheduleString.TimeEnd, Valid: true},
	}
	if err = db.Model(&models.Schedule{}).Where("id = ?", id).Updates(schedule).Error; err != nil {
		log.Println("Schedule PATCH error [" + claims.Username + "]" + err.Error())
		return
	}
	log.Println("Schedule PATCH [" + claims.Username + "]")

	result = make(map[string]any)
	result["Response"] = "schedule updated"
	return
}

func DeleteSchedule(db *gorm.DB, id int, claims models.Claims) (result map[string]any, err error) {

	if err = db.Delete(&models.Schedule{}, id).Error; err != nil {
		log.Println("Schedule DELETE error [" + claims.Username + "]" + err.Error())
		return
	}

	log.Println("Schedule DELETE [" + claims.Username + "]")

	result = make(map[string]any)
	result["Response"] = "schedule deleted"
	return
}

package controllers

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	models "pharmacy-api/shared/models"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

const LOGS_FOLDER = "./shared/pharmacy-content/logs"

func GetLogs(c *gin.Context) {
	query := c.Query("q")
	pageStr := c.Query("page")
	limitStr := c.Query("limit")

	user, _ := c.Get("user")
	claims := user.(*models.Claims)

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

	files, err := os.ReadDir(LOGS_FOLDER)
	if err != nil {
		log.Println("Logs GET error [" + claims.Username + "] " + err.Error())
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	var allLogs []models.LogsResponse
	for _, file := range files {
		if !file.IsDir() {
			name := file.Name()
			if query == "" || strings.Contains(strings.ToLower(name), strings.ToLower(query)) {
				allLogs = append(allLogs, models.LogsResponse{Name: name})
			}
		}
	}

	totalCount := len(allLogs)

	start := offset
	end := offset + limit
	if start > totalCount {
		start = totalCount
	}
	if end > totalCount {
		end = totalCount
	}
	pagedLogs := allLogs[start:end]

	totalPages := (totalCount + limit - 1) / limit

	log.Println("Logs GET [" + claims.Username + "]")
	c.JSON(http.StatusOK, gin.H{
		"Items":       pagedLogs,
		"TotalPages":  totalPages,
		"CurrentPage": page,
	})
}

func GetLog(c *gin.Context) {
	user, _ := c.Get("user")
	claims := user.(*models.Claims)

	logFile := c.Query("name")
	if logFile == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": ""})
		return
	}

	fullPath := filepath.Join(LOGS_FOLDER, logFile)
	if _, err := os.Stat(fullPath); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	}

	log.Println("Log GET [" + claims.Username + "]")
	c.FileAttachment(fullPath, logFile)
}

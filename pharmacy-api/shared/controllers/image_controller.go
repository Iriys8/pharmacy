package controllers

import (
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

const IMAGE_FOLDER = "./shared/pharmacy-content/images"

func GetImage(c *gin.Context) {
	imageName := c.Query("name")
	if imageName == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": ""})
		return
	}

	filePath := filepath.Join(IMAGE_FOLDER, imageName)
	c.File(filePath)
}

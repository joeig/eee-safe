package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type statusBody struct {
	MaxBackupBytes uint `json:"maxBackupBytes"`
	RetentionDays  uint `json:"retentionDays"`
}

// GetConfigHandler Gin route
func GetConfigHandler(c *gin.Context) {
	c.JSON(http.StatusOK, statusBody{
		MaxBackupBytes: config.Server.Backups.MaxBackupBytes,
		RetentionDays:  config.Server.Backups.RetentionDays,
	})
}

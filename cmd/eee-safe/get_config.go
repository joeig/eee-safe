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
func (a *AppCtx) GetConfigHandler(c *gin.Context) {
	c.JSON(http.StatusOK, statusBody{
		MaxBackupBytes: a.Config.Server.Backups.MaxBackupBytes,
		RetentionDays:  a.Config.Server.Backups.RetentionDays,
	})
}

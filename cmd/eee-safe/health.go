package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// HealthStatus contains information regarding the healthiness of the application
type HealthStatus struct {
	ApplicationRunning bool `json:"applicationRunning"`
}

// GetHealthHandler Gin route
func GetHealthHandler(c *gin.Context) {
	hs := &HealthStatus{
		ApplicationRunning: true,
	}
	c.JSON(http.StatusOK, hs)
}

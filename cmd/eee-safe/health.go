package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthStatus contains information regarding the healthiness of the application
type HealthStatus struct {
	ApplicationRunning bool `json:"applicationRunning"`
}

// GetHealthHandler Gin route
func (a *AppContext) GetHealthHandler(c *gin.Context) {
	hs := &HealthStatus{
		ApplicationRunning: true,
	}
	c.JSON(http.StatusOK, hs)
}

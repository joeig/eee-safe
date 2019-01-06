package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joeig/eee-safe/threema"
	"log"
	"net/http"
	"strings"
)

// Adds an unique request ID to every single Gin request
func requestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		uuid4, err := uuid.NewRandom()
		if err != nil {
			log.Fatal("Unable to generate request ID")
			return
		}
		requestID := uuid4.String()
		c.Set("RequestID", requestID)
		c.Header("X-Request-ID", requestID)
		log.SetPrefix(fmt.Sprintf("[%s] ", requestID))
		log.Printf("Set request ID to \"%s\"", requestID)
		c.Next()
	}
}

// Validates the Threema user agent
func validateUserAgentMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userAgent := c.GetHeader("User-Agent")
		if !strings.Contains(userAgent, threema.ValidUserAgentSubstring) {
			log.Printf("Invalid user agent header \"%s\"", userAgent)
			c.String(http.StatusBadRequest, "")
			c.Abort()
			return
		}
		log.Printf("User agent header is valid: \"%s\"", userAgent)
		c.Next()
	}
}

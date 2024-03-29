package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RequestIDGenerator defines a function which generates a request ID.
type RequestIDGenerator func() string

// NewRandomRequestIDGenerator creates a function which returns a random request ID.
func NewRandomRequestIDGenerator() RequestIDGenerator {
	return func() string {
		randomBytes := make([]byte, 12)

		if _, err := rand.Read(randomBytes); err != nil {
			panic(err)
		}

		return base64.RawURLEncoding.EncodeToString(randomBytes)
	}
}

const (
	// GinRequestIDSymbol contains a symbol which stores the request ID for this request.
	GinRequestIDSymbol = "RequestID"

	// RequestIDHeader contains a header filed name which stores the request ID for this request.
	RequestIDHeader = "X-Request-ID"
)

// requestIDMiddleware adds an unique request ID to every single Gin request.
func requestIDMiddleware(generateRequestID RequestIDGenerator) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		requestID := generateRequestID()

		ginCtx.Header(RequestIDHeader, requestID)
		ginCtx.Set(GinRequestIDSymbol, requestID)

		log.SetPrefix(fmt.Sprintf("[%s] ", requestID))
		log.Printf("Set request ID to %q", requestID)
	}
}

// ThreemaHeaderParams contains the header parameters for validateUserAgentMiddleware.
type ThreemaHeaderParams struct {
	UserAgent string `header:"User-Agent" binding:"required,eq=Threema"`
}

// validateUserAgentMiddleware validates the expected user agent.
func validateUserAgentMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var headerParams ThreemaHeaderParams
		if err := c.ShouldBindHeader(&headerParams); err != nil {
			log.Println(err)

			c.String(http.StatusBadRequest, "")
			c.Abort()

			return
		}

		log.Printf("User agent header is valid: %q", headerParams.UserAgent)

		c.Next()
	}
}

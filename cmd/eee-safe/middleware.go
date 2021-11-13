package main

import (
	cryptoRand "crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joeig/eee-safe/pkg/threema"
)

// initialiseSeed globally initialises math/rand with a cryptographically strong seed.
func initialiseSeed() {
	seed := make([]byte, 8)

	_, err := cryptoRand.Read(seed)
	if err != nil {
		panic(err)
	}

	rand.Seed(int64(binary.LittleEndian.Uint64(seed)))
}

// RequestIDGenerator defines a function which generates a request ID.
type RequestIDGenerator func() string

// NewRandomRequestIDGenerator creates a function which returns a random request ID.
func NewRandomRequestIDGenerator() RequestIDGenerator {
	return func() string {
		randomBytes := make([]byte, 12)

		// Ignore gosec for the following line, because math/rand is supposed to be used after seed initialization.
		// #nosec
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
		log.Printf("Set request ID to \"%s\"", requestID)
	}
}

// validateUserAgentMiddleware validates the Threema user agent.
func validateUserAgentMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userAgent := c.GetHeader("User-Agent")
		if err := threema.ValidateUserAgent(userAgent); err != nil {
			log.Println(err)

			c.String(http.StatusBadRequest, "")
			c.Abort()

			return
		}

		log.Printf("User agent header is valid: \"%s\"", userAgent)

		c.Next()
	}
}

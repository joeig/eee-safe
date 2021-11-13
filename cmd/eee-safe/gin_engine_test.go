package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// NewMockRequestIDGenerator creates a RequestIDGenerator which always returns a fixed request ID.
func NewMockRequestIDGenerator(fixture string) RequestIDGenerator {
	return func() string {
		return fixture
	}
}

func TestRequestIDMiddleware(t *testing.T) {
	const fixture = "fixed-request-ID"

	ginCtx, _ := gin.CreateTestContext(httptest.NewRecorder())
	ginCtx.Request = &http.Request{Header: map[string][]string{}}
	requestIDMiddleware(NewMockRequestIDGenerator(fixture))(ginCtx)

	if value := ginCtx.Writer.Header().Get(RequestIDHeader); value != fixture {
		t.Errorf("wrong request ID header: %q != %q", value, fixture)
	}

	if value, _ := ginCtx.Get(GinRequestIDSymbol); value != fixture {
		t.Errorf("wrong request ID: %q != %q", value, fixture)
	}
}

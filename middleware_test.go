package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequestID(t *testing.T) {
	configFile := "config.dist.yml"
	parseConfig(&config, &configFile)
	router := getGinEngine()
	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	req.SetBasicAuth("jonathan", "byers")
	req.Header.Set("User-Agent", "Threema")
	router.ServeHTTP(res, req)
	if res.Header().Get("X-Request-ID") == "" {
		t.Errorf("X-Request-ID is missing")
	}
}

func TestUnauthorizedAccess(t *testing.T) {
	configFile := "config.dist.yml"
	parseConfig(&config, &configFile)
	router := getGinEngine()
	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	req.SetBasicAuth("foo", "bar")
	req.Header.Set("User-Agent", "Threema")
	router.ServeHTTP(res, req)
	if res.Code != http.StatusUnauthorized {
		t.Errorf("HTTP request returns %d instead of %d", res.Code, http.StatusUnauthorized)
	}
}

func TestInvalidUserAgent(t *testing.T) {
	configFile := "config.dist.yml"
	parseConfig(&config, &configFile)
	router := getGinEngine()
	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	req.SetBasicAuth("jonathan", "byers")
	req.Header.Set("User-Agent", "foo")
	router.ServeHTTP(res, req)
	if res.Code != http.StatusBadRequest {
		t.Errorf("HTTP request returns %d instead of %d", res.Code, http.StatusBadRequest)
	}
}

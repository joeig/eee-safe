package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetHealthHandler(t *testing.T) {
	configFile := "../../configs/config.dist.yml"
	parseConfig(&config, &configFile)
	router := getGinEngine()
	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	req.SetBasicAuth("jonathan", "byers")
	req.Header.Set("User-Agent", "Threema")
	router.ServeHTTP(res, req)
	if res.Code != http.StatusOK {
		t.Errorf("HTTP request does not return %v", http.StatusOK)
	}
}

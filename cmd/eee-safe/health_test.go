package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/spf13/viper"
)

func TestGetHealthHandler(t *testing.T) {
	appCtx := &AppCtx{
		Config: &Config{},
	}
	_ = appCtx.Config.Read(viper.New(), "../../configs/config.dist.yml")

	router := getGinEngine(appCtx)
	res := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/health", nil)
	req.SetBasicAuth("jonathan", "byers")
	req.Header.Set("User-Agent", "Threema")

	router.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Errorf("HTTP request does not return %v", http.StatusOK)
	}
}

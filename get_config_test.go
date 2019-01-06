package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetConfigHandler(t *testing.T) {
	configFile := "config.dist.yml"
	parseConfig(&config, &configFile)
	router := getGinEngine()
	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/config", nil)
	req.SetBasicAuth("jonathan", "byers")
	req.Header.Set("User-Agent", "Threema")
	router.ServeHTTP(res, req)
	if res.Code != http.StatusOK {
		t.Errorf("HTTP request does not return %v", http.StatusOK)
	}
	var resStatusBody statusBody
	err := json.Unmarshal(res.Body.Bytes(), &resStatusBody)
	if err != nil {
		t.Errorf("Cannot unmarshal status response body: %v", err)
	}
	if resStatusBody.RetentionDays != config.Server.Backups.RetentionDays || resStatusBody.MaxBackupBytes != config.Server.Backups.MaxBackupBytes {
		t.Errorf("Config endpoint returns invalid values")
	}
}

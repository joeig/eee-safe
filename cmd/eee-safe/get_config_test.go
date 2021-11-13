package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/spf13/viper"
)

func TestGetConfigHandler(t *testing.T) {
	appCtx := &AppCtx{
		Config: &Config{},
	}
	_ = appCtx.Config.Read(viper.New(), "../../configs/config.dist.yml")

	router := getGinEngine(appCtx)
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

	if resStatusBody.RetentionDays != appCtx.Config.Server.Backups.RetentionDays || resStatusBody.MaxBackupBytes != appCtx.Config.Server.Backups.MaxBackupBytes {
		t.Errorf("Config endpoint returns invalid values")
	}
}

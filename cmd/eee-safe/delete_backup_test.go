package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/joeig/eee-safe/pkg/threema"
	"github.com/spf13/viper"
)

func assertDeleteBackupHandlerComponent(t *testing.T, router *gin.Engine, backupID string, assertedCode int) { // nolint:interfacer
	url := fmt.Sprintf("/backups/%s", backupID)

	req, _ := http.NewRequest(http.MethodDelete, url, nil)
	req.SetBasicAuth("jonathan", "byers")
	req.Header.Set("User-Agent", "Threema")

	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Code != assertedCode {
		t.Errorf("HTTP request to \"%s\" returned %d instead of %d", url, res.Code, assertedCode)
	}
}

func TestDeleteBackupHandler(t *testing.T) {
	appCtx := &AppContext{
		Config:             &Config{},
		RequestIDGenerator: NewMockRequestIDGenerator("foo"),
	}
	_ = appCtx.Config.Read(viper.New(), "../../configs/config.dist.yml")
	_ = appCtx.InitializeStorageBackend()

	router := getGinEngine(appCtx)

	// OK
	t.Run("TestDeleteValidBackup", func(t *testing.T) {
		backupID, _ := threema.ConvertToBackupID("87df5aaa32e3de72426e04e845d1251d87df5aaa32e3de72426e04e845d1251d")
		backupInput := &threema.BackupInput{BackupID: backupID, EncryptedBackup: threema.EncryptedBackup("Foo")}
		_ = appCtx.StorageBackend.PutBackup(context.Background(), backupInput)
		assertDeleteBackupHandlerComponent(t, router, "87df5aaa32e3de72426e04e845d1251d87df5aaa32e3de72426e04e845d1251d", http.StatusOK)
	})

	// BadRequest
	t.Run("TestTooShortBackupID", func(t *testing.T) {
		assertDeleteBackupHandlerComponent(t, router, "87df5aaa32e3de72426e04e8", http.StatusBadRequest)
	})
	t.Run("TestTooLongBackupID", func(t *testing.T) {
		assertDeleteBackupHandlerComponent(t, router, "87df5aaa32e3de72426e04e845d1251d87df5aaa32e3de72426e04e845d1251daaaa", http.StatusBadRequest)
	})

	// NotFound
	t.Run("TestDeleteNonExistingBackupID", func(t *testing.T) {
		assertDeleteBackupHandlerComponent(t, router, "87df5aaa32e3de72426e04e845d1251d87df5aaa32e3de72426e04e845d1251d", http.StatusNotFound)
	})
}

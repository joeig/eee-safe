package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/joeig/eee-safe/pkg/threema"
	"github.com/spf13/viper"
)

func assertPutBackupHandlerComponent(t *testing.T, router *gin.Engine, backupID string, encryptedBackup []byte, contentType string, assertedCode int) { // nolint:interfacer
	url := fmt.Sprintf("/backups/%s", backupID)
	body := strings.NewReader(string(encryptedBackup))

	req, _ := http.NewRequest(http.MethodPut, url, body)
	req.SetBasicAuth("jonathan", "byers")
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("User-Agent", "Threema")

	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Code != assertedCode {
		t.Errorf("HTTP request to \"%s\" returned %d instead of %d", url, res.Code, assertedCode)
	}
}

func TestPutBackupHandler(t *testing.T) {
	appCtx := &AppContext{
		Config:             &Config{},
		RequestIDGenerator: NewMockRequestIDGenerator("foo"),
	}
	_ = appCtx.Config.Read(viper.New(), "../../configs/config.dist.yml")
	_ = appCtx.InitializeStorageBackend()

	router := getGinEngine(appCtx)

	// OK
	t.Run("TestPutValidBackup", func(t *testing.T) {
		assertPutBackupHandlerComponent(t, router, "87df5aaa32e3de72426e04e845d1251d87df5aaa32e3de72426e04e845d1251d", []byte("Foo"), "application/octet-stream", http.StatusOK)
	})
	t.Run("TestOverwriteValidBackup", func(t *testing.T) {
		assertPutBackupHandlerComponent(t, router, "87df5aaa32e3de72426e04e845d1251d87df5aaa32e3de72426e04e845d1251d", []byte("Bar"), "application/octet-stream", http.StatusOK)
	})

	// BadRequest
	t.Run("TestWrongContentType", func(t *testing.T) {
		assertPutBackupHandlerComponent(t, router, "87df5aaa32e3de72426e04e845d1251d87df5aaa32e3de72426e04e845d1251d", []byte("Foo"), "text/plain", http.StatusBadRequest)
	})
	t.Run("TestTooShortBackupID", func(t *testing.T) {
		assertPutBackupHandlerComponent(t, router, "87df5aaa32e3de72426e04e8", []byte("Foo"), "application/octet-stream", http.StatusBadRequest)
	})
	t.Run("TestTooLongBackupID", func(t *testing.T) {
		assertPutBackupHandlerComponent(t, router, "87df5aaa32e3de72426e04e845d1251d87df5aaa32e3de72426e04e845d1251daaaa", []byte("Foo"), "application/octet-stream", http.StatusBadRequest)
	})
	t.Run("TestMissingEncryptedBackup", func(t *testing.T) {
		assertPutBackupHandlerComponent(t, router, "bbbb5aaa32e3de72426e04e845d1251d87df5aaa32e3de72426e04e845d1251d", []byte{}, "application/octet-stream", http.StatusBadRequest)
	})

	// Clean up
	backupID, _ := threema.ConvertToBackupID("87df5aaa32e3de72426e04e845d1251d87df5aaa32e3de72426e04e845d1251d")
	_ = appCtx.StorageBackend.DeleteBackup(context.Background(), backupID)
}

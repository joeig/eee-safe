package main

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joeig/eee-safe/pkg/threema"
	"net/http"
	"net/http/httptest"
	"testing"
)

func assertGetBackupHandlerComponent(t *testing.T, router *gin.Engine, backupID string, assertedEncryptedBackup []byte, assertedCode int) *httptest.ResponseRecorder {
	url := fmt.Sprintf("/backups/%s", backupID)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.SetBasicAuth("jonathan", "byers")
	req.Header.Set("Accept", "application/octet-stream")
	req.Header.Set("User-Agent", "Threema")
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)
	if res.Code != assertedCode {
		t.Errorf("HTTP request to \"%s\" returned %d instead of %d", url, res.Code, assertedCode)
	}
	if !bytes.Equal(res.Body.Bytes(), assertedEncryptedBackup) {
		t.Errorf("HTTP response payload does not match")
	}
	return res
}

func TestGetBackupHandler(t *testing.T) {
	configFile := "../../configs/config.dist.yml"
	parseConfig(&config, &configFile)
	setStorageBackend(&storageBackend)
	router := getGinEngine()

	// Initialization
	backupIDString := "c8df5aaa32e3de72426e04e845d1251d87df5aaa32e3de72426e04e845d1251d"
	backupID, _ := threema.ConvertToBackupID(backupIDString)
	backupInput := &threema.BackupInput{
		BackupID:        backupID,
		EncryptedBackup: threema.EncryptedBackup("Foo"),
		RetentionDays:   config.Server.Backups.RetentionDays,
	}
	_ = storageBackend.PutBackup(backupInput)

	// OK
	t.Run("TestGetValidBackup", func(t *testing.T) {
		assertGetBackupHandlerComponent(t, router, backupIDString, []byte("Foo"), http.StatusOK)
	})

	// BadRequest
	t.Run("TestTooShortBackupID", func(t *testing.T) {
		assertGetBackupHandlerComponent(t, router, "87df5aaa32e3de72426e04e8", []byte{}, http.StatusBadRequest)
	})
	t.Run("TestTooLongBackupID", func(t *testing.T) {
		assertGetBackupHandlerComponent(t, router, "87df5aaa32e3de72426e04e845d1251d87df5aaa32e3de72426e04e845d1251daaaa", []byte{}, http.StatusBadRequest)
	})

	// NotFound
	t.Run("TestGetNonExistingBackupID", func(t *testing.T) {
		assertGetBackupHandlerComponent(t, router, "bbbb5aaa32e3de72426e04e845d1251d87df5aaa32e3de72426e04e845d1251d", []byte{}, http.StatusNotFound)
	})

	// Clean up
	_ = storageBackend.DeleteBackup(backupID)
}

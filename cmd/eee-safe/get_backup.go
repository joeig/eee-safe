package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joeig/eee-safe/pkg/storage"
	"github.com/joeig/eee-safe/pkg/threema"
)

// GetBackupHandler Gin route
func GetBackupHandler(c *gin.Context) {
	threemaSafeBackupID, err := threema.ConvertToBackupID(c.Param("threemaSafeBackupID"))
	if err != nil {
		log.Println(err)
		c.Data(http.StatusBadRequest, "", []byte{})

		return
	}

	backup, err := storageBackend.GetBackup(threemaSafeBackupID)
	if err != nil {
		log.Println(err)

		switch err.(type) {
		case *storage.BackupIDNotFoundError:
			c.Data(http.StatusNotFound, "", []byte{})
		default:
			c.Data(http.StatusInternalServerError, "", []byte{})
		}

		return
	}

	if !backup.CreationTime.IsZero() {
		c.Header("Date", backup.CreationTime.Format(http.TimeFormat))
	}

	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Data(http.StatusOK, "application/octet-stream", backup.EncryptedBackup)
}

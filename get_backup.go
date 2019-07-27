package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joeig/eee-safe/storage"
	"github.com/joeig/eee-safe/threema"
	"log"
	"net/http"
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
			break
		default:
			c.Data(http.StatusInternalServerError, "", []byte{})
		}
		return
	}
	if !backup.CreationTime.IsZero() {
		c.Header("Date", backup.CreationTime.Format(http.TimeFormat))
	}
	if !backup.ExpirationTime.IsZero() {
		c.Header("Expires", backup.ExpirationTime.Format(http.TimeFormat))
	}
	c.Data(http.StatusOK, "application/octet-stream", backup.EncryptedBackup)
	return
}

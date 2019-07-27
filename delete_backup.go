package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joeig/eee-safe/storage"
	"github.com/joeig/eee-safe/threema"
	"log"
	"net/http"
)

// DeleteBackupHandler Gin route
func DeleteBackupHandler(c *gin.Context) {
	threemaSafeBackupID, err := threema.ConvertToBackupID(c.Param("threemaSafeBackupID"))
	if err != nil {
		log.Println(err)
		c.Data(http.StatusBadRequest, "", []byte{})
		return
	}
	if err := storageBackend.DeleteBackup(threemaSafeBackupID); err != nil {
		log.Println(err)
		switch err.(type) {
		case *storage.BackupIDNotFoundError:
			c.Data(http.StatusNotFound, "", []byte{})
		default:
			c.Data(http.StatusInternalServerError, "", []byte{})
		}
		return
	}
	c.Data(http.StatusOK, "", []byte{})
}

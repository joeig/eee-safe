package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joeig/eee-safe/pkg/storage"
	"github.com/joeig/eee-safe/pkg/threema"
)

// DeleteBackupHandler Gin route
func (a *AppCtx) DeleteBackupHandler(c *gin.Context) {
	threemaSafeBackupID, err := threema.ConvertToBackupID(c.Param("threemaSafeBackupID"))
	if err != nil {
		log.Println(err)
		c.Data(http.StatusBadRequest, "", []byte{})

		return
	}

	if err := a.StorageBackend.DeleteBackup(c.Request.Context(), threemaSafeBackupID); err != nil {
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

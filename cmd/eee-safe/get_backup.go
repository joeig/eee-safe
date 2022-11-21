package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joeig/eee-safe/pkg/storage"
	"github.com/joeig/eee-safe/pkg/threema"
)

// GetBackupPathParams contains path parameters for GetBackupHandler.
type GetBackupPathParams struct {
	ThreemaSafeBackupID string `uri:"threemaSafeBackupID" binding:"required,hexadecimal,len=64"`
}

// GetBackupHandler Gin route
func (a *AppContext) GetBackupHandler(c *gin.Context) {
	var pathParams GetBackupPathParams
	if err := c.ShouldBindUri(&pathParams); err != nil {
		log.Println(err)
		c.Data(http.StatusBadRequest, "", []byte{})

		return
	}

	threemaSafeBackupID, err := threema.ConvertToBackupID(pathParams.ThreemaSafeBackupID)
	if err != nil {
		log.Println(err)
		c.Data(http.StatusBadRequest, "", []byte{})

		return
	}

	backup, err := a.StorageBackend.GetBackup(c.Request.Context(), threemaSafeBackupID)
	if err != nil {
		log.Println(err)

		switch err.(type) {
		case *storage.ErrBackupIDNotFound:
			c.Data(http.StatusNotFound, "", []byte{})
		default:
			c.Data(http.StatusInternalServerError, "", []byte{})
		}

		return
	}

	if !backup.CreationTime.IsZero() {
		c.Header("Last-Modified", backup.CreationTime.Format(http.TimeFormat))
	}

	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Data(http.StatusOK, "application/octet-stream", backup.EncryptedBackup)
}

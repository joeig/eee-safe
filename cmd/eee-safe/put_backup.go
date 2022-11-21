package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joeig/eee-safe/pkg/threema"
)

// PutBackupHeaderParams contains header parameters for PutBackupHandler.
type PutBackupHeaderParams struct {
	ContentType string `header:"Content-Type" binding:"required,eq=application/octet-stream"`
}

// PutBackupPathParams contains path parameters for PutBackupHandler.
type PutBackupPathParams struct {
	ThreemaSafeBackupID string `uri:"threemaSafeBackupID" binding:"required,hexadecimal,len=64"`
}

// PutBackupHandler Gin route
func (a *AppContext) PutBackupHandler(c *gin.Context) {
	if err := c.ShouldBindHeader(&PutBackupHeaderParams{}); err != nil {
		log.Println(err)
		c.Data(http.StatusBadRequest, "", []byte{})

		return
	}

	var pathParams PutBackupPathParams
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

	var threemaSafeEncryptedBackup threema.EncryptedBackup
	threemaSafeEncryptedBackup, _ = c.GetRawData()

	if err := threemaSafeEncryptedBackup.Validate(a.Config.Server.Backups.MaxBackupBytes); err != nil {
		log.Println(err)
		c.Data(http.StatusBadRequest, "", []byte{})

		return
	}

	backupInput := &threema.BackupInput{
		BackupID:        threemaSafeBackupID,
		EncryptedBackup: threemaSafeEncryptedBackup,
		RetentionDays:   a.Config.Server.Backups.RetentionDays,
	}

	if err := a.StorageBackend.PutBackup(c.Request.Context(), backupInput); err != nil {
		log.Println(err)
		c.Data(http.StatusInternalServerError, "", []byte{})

		return
	}

	c.Data(http.StatusOK, "", []byte{})
}

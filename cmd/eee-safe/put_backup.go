package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joeig/eee-safe/pkg/threema"
)

// PutBackupHandler Gin route
func (a *AppCtx) PutBackupHandler(c *gin.Context) {
	if c.GetHeader("Content-Type") != "application/octet-stream" {
		log.Println(&contentTypeInvalid{})

		c.Data(http.StatusBadRequest, "", []byte{})

		return
	}

	threemaSafeBackupID, err := threema.ConvertToBackupID(c.Param("threemaSafeBackupID"))
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

	if err := storageBackend.PutBackup(backupInput); err != nil {
		log.Println(err)

		c.Data(http.StatusInternalServerError, "", []byte{})

		return
	}

	c.Data(http.StatusOK, "", []byte{})
}

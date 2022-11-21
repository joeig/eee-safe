package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joeig/eee-safe/pkg/storage"
	"github.com/joeig/eee-safe/pkg/threema"
)

// DeleteBackupPathParams contains path parameters for DeleteBackupHandler.
type DeleteBackupPathParams struct {
	ThreemaSafeBackupID string `uri:"threemaSafeBackupID" binding:"required,hexadecimal,len=64"`
}

// DeleteBackupHandler Gin route
func (a *AppContext) DeleteBackupHandler(c *gin.Context) {
	var pathParams DeleteBackupPathParams
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

	if err := a.StorageBackend.DeleteBackup(c.Request.Context(), threemaSafeBackupID); err != nil {
		log.Println(err)

		switch err.(type) {
		case *storage.ErrBackupIDNotFound:
			c.Data(http.StatusNotFound, "", []byte{})
		default:
			c.Data(http.StatusInternalServerError, "", []byte{})
		}

		return
	}

	c.Data(http.StatusOK, "", []byte{})
}

package common

import (
	"fmt"
	"github.com/joeig/eee-safe/threema"
)

type StorageBackendError struct {
	APIError interface{}
}

func (e *StorageBackendError) Error() string {
	return fmt.Sprintf("Storage Backend API error: %v", e.APIError)
}

type BackupIDNotFoundError struct {
	BackupID threema.BackupID
}

func (e *BackupIDNotFoundError) Error() string {
	return fmt.Sprintf("BackupInput ID \"%s\" not found", e.BackupID.String())
}

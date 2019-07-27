package storage

import (
	"fmt"
	"github.com/joeig/eee-safe/threema"
)

// StorageBackendError is being used whenever a storage backend error occurred
type StorageBackendError struct {
	APIError interface{}
}

func (e *StorageBackendError) Error() string {
	return fmt.Sprintf("Storage Backend API error: %v", e.APIError)
}

// BackupIDNotFoundError is being used whenever a backup ID was not found
type BackupIDNotFoundError struct {
	BackupID threema.BackupID
}

func (e *BackupIDNotFoundError) Error() string {
	return fmt.Sprintf("BackupInput ID \"%s\" not found", e.BackupID.String())
}

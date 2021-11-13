package storage

import (
	"fmt"

	"github.com/joeig/eee-safe/pkg/threema"
)

// ErrUnknown is being used whenever a storage backend error occurred
type ErrUnknown struct {
	APIError interface{}
}

func (e *ErrUnknown) Error() string {
	return fmt.Sprintf("storage backend API error: %v", e.APIError)
}

// ErrBackupIDNotFound is being used whenever a backup ID was not found
type ErrBackupIDNotFound struct {
	BackupID threema.BackupID
}

func (e *ErrBackupIDNotFound) Error() string {
	return fmt.Sprintf("backup ID %q not found", e.BackupID.String())
}

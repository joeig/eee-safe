package storage

import (
	"github.com/joeig/eee-safe/threema"
)

// StorageBackendType sets the storage backend type
type StorageBackendType string

// StorageBackend is an interface for basic storage operations
type StorageBackend interface {
	PutBackup(backupInput *threema.BackupInput) error
	GetBackup(backupID threema.BackupID) (*threema.BackupOutput, error)
	DeleteBackup(backupID threema.BackupID) error
}

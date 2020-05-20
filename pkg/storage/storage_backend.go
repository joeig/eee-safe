package storage

import (
	"github.com/joeig/eee-safe/pkg/threema"
)

// BackendType sets the storage backend type
type BackendType string

// Backend is an interface for basic storage operations
type Backend interface {
	PutBackup(backupInput *threema.BackupInput) error
	GetBackup(backupID threema.BackupID) (*threema.BackupOutput, error)
	DeleteBackup(backupID threema.BackupID) error
}

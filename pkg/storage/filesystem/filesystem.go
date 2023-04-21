package filesystem

import (
	"context"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/joeig/eee-safe/pkg/storage"
	"github.com/joeig/eee-safe/pkg/threema"
)

// File defines file system operations.
type File interface {
	// WriteFile writes a file to a file system.
	WriteFile(name string, data []byte, perm os.FileMode) error
	// ReadFile reads a file from a file system.
	ReadFile(name string) ([]byte, error)
	// Remove deletes a file from a file system.
	Remove(name string) error
	// Stat returns infos about a file on a file system.
	Stat(name string) (os.FileInfo, error)
}

// Filesystem defines the configuration of the filesystem storage backend type.
type Filesystem struct {
	mu   sync.RWMutex
	file File

	// Directory is the path to the directory which contains the backups.
	Directory string `mapstructure:"directory"`
	// Permissions contains the desired permission for the files being written to the file system.
	Permissions os.FileMode `mapstructure:"permissions"`
}

// SetFile sets the file instance.
func (f *Filesystem) SetFile(file File) {
	f.file = file
}

// PutBackup stores a backup to the filesystem.
func (f *Filesystem) PutBackup(_ context.Context, backupInput *threema.BackupInput) error {
	fileName, err := f.generateFileName(backupInput.BackupID)
	if err != nil {
		return err
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	if err := f.file.WriteFile(fileName, backupInput.EncryptedBackup, f.Permissions); err != nil {
		return &ErrFileNotWritable{FileName: fileName, UpstreamError: err}
	}

	return nil
}

// GetBackup returns a backup from the filesystem.
func (f *Filesystem) GetBackup(_ context.Context, backupID threema.BackupID) (*threema.BackupOutput, error) {
	fileName, err := f.generateFileName(backupID)
	if err != nil {
		return &threema.BackupOutput{}, err
	}

	f.mu.RLock()
	defer f.mu.RUnlock()

	info, err := f.file.Stat(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			return &threema.BackupOutput{}, &storage.ErrBackupIDNotFound{BackupID: backupID}
		}

		return &threema.BackupOutput{}, &ErrFileNotReadable{FileName: fileName, UpstreamError: err}
	}

	data, err := f.file.ReadFile(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			return &threema.BackupOutput{}, &storage.ErrBackupIDNotFound{BackupID: backupID}
		}

		return &threema.BackupOutput{}, &ErrFileNotReadable{FileName: fileName, UpstreamError: err}
	}

	dataString := threema.EncryptedBackup(data)

	return &threema.BackupOutput{
		BackupID:        backupID,
		EncryptedBackup: dataString,
		CreationTime:    info.ModTime(),
	}, nil
}

// DeleteBackup deletes a backup from the filesystem.
func (f *Filesystem) DeleteBackup(_ context.Context, backupID threema.BackupID) error {
	fileName, err := f.generateFileName(backupID)
	if err != nil {
		return err
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	if err := f.file.Remove(fileName); err != nil {
		if os.IsNotExist(err) {
			return &storage.ErrBackupIDNotFound{BackupID: backupID}
		}

		return &ErrFileNotRemovable{FileName: fileName, UpstreamError: err}
	}

	return nil
}

func (f *Filesystem) generateFileName(backupID threema.BackupID) (string, error) {
	if f.Directory == "" {
		return "", &ErrDirectoryInvalid{Directory: f.Directory}
	}

	return path.Join(f.Directory, strings.ToLower(backupID.String())), nil
}

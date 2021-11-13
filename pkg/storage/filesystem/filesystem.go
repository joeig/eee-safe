package filesystem

import (
	"context"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/joeig/eee-safe/pkg/storage"
	"github.com/joeig/eee-safe/pkg/threema"
)

// Filesystem defines the configuration of the filesystem storage backend type
type Filesystem struct {
	mu sync.RWMutex

	Directory   string      `mapstructure:"directory"`
	Permissions os.FileMode `mapstructure:"permissions"`
}

// PutBackup stores a backup to the filesystem
func (f *Filesystem) PutBackup(_ context.Context, backupInput *threema.BackupInput) error {
	fileName, err := f.generateFileName(backupInput.BackupID)
	if err != nil {
		return err
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	if err := ioutil.WriteFile(fileName, backupInput.EncryptedBackup, f.Permissions); err != nil {
		return &FileNotWritableError{FileName: fileName, UpstreamError: err}
	}

	return nil
}

// GetBackup returns a backup from the filesystem
func (f *Filesystem) GetBackup(_ context.Context, backupID threema.BackupID) (*threema.BackupOutput, error) {
	fileName, err := f.generateFileName(backupID)
	if err != nil {
		return &threema.BackupOutput{}, err
	}

	f.mu.RLock()
	defer f.mu.RUnlock()

	info, err := os.Stat(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			return &threema.BackupOutput{}, &storage.BackupIDNotFoundError{BackupID: backupID}
		}

		return &threema.BackupOutput{}, &FileNotReadableError{FileName: fileName, UpstreamError: err}
	}

	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			return &threema.BackupOutput{}, &storage.BackupIDNotFoundError{BackupID: backupID}
		}

		return &threema.BackupOutput{}, &FileNotReadableError{FileName: fileName, UpstreamError: err}
	}

	dataString := threema.EncryptedBackup(data)

	return &threema.BackupOutput{
		BackupID:        backupID,
		EncryptedBackup: dataString,
		CreationTime:    info.ModTime(),
	}, nil
}

// DeleteBackup deletes a backup from the filesystem
func (f *Filesystem) DeleteBackup(_ context.Context, backupID threema.BackupID) error {
	fileName, err := f.generateFileName(backupID)
	if err != nil {
		return err
	}

	f.mu.Lock()
	defer f.mu.Unlock()

	if err := os.Remove(fileName); err != nil {
		if os.IsNotExist(err) {
			return &storage.BackupIDNotFoundError{BackupID: backupID}
		}

		return &FileNotRemovableError{FileName: fileName, UpstreamError: err}
	}

	return nil
}

func (f *Filesystem) generateFileName(backupID threema.BackupID) (string, error) {
	if f.Directory == "" {
		return "", &DirectoryInvalidError{Directory: f.Directory}
	}

	return path.Join(f.Directory, strings.ToLower(backupID.String())), nil
}

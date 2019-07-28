package filesystem

import (
	"github.com/joeig/eee-safe/pkg/debug"
	"github.com/joeig/eee-safe/pkg/storage"
	"github.com/joeig/eee-safe/pkg/threema"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

// PutBackup stores a backup to the filesystem
func (f *Filesystem) PutBackup(backupInput *threema.BackupInput) error {
	if backupInput.RetentionDays != 0 {
		debug.Printf("Filesystem storage backend does currently not support backup retention.")
	}
	fileName, err := f.generateFileName(backupInput.BackupID)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(fileName, backupInput.EncryptedBackup, f.Permissions); err != nil {
		return &FileNotWritableError{FileName: fileName, UpstreamError: err}
	}
	debug.Printf("Writing file \"%s\" (%d bytes)", fileName, len(backupInput.EncryptedBackup))
	return nil
}

// GetBackup returns a backup from the filesystem
func (f *Filesystem) GetBackup(backupID threema.BackupID) (*threema.BackupOutput, error) {
	fileName, err := f.generateFileName(backupID)
	if err != nil {
		return &threema.BackupOutput{}, err
	}
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
	debug.Printf("Reading file \"%s\" (%d bytes)", fileName, len(dataString))
	return &threema.BackupOutput{
		BackupID:        backupID,
		EncryptedBackup: dataString,
		CreationTime:    info.ModTime(),
	}, nil
}

// DeleteBackup deletes a backup from the filesystem
func (f *Filesystem) DeleteBackup(backupID threema.BackupID) error {
	fileName, err := f.generateFileName(backupID)
	if err != nil {
		return err
	}
	if err := os.Remove(fileName); err != nil {
		if os.IsNotExist(err) {
			return &storage.BackupIDNotFoundError{BackupID: backupID}
		}
		return &FileNotRemovableError{FileName: fileName, UpstreamError: err}
	}
	debug.Printf("Deleting file \"%s\"", fileName)
	return nil
}

func (f *Filesystem) generateFileName(backupID threema.BackupID) (string, error) {
	if f.Directory == "" {
		return "", &DirectoryInvalidError{Directory: f.Directory}
	}
	return path.Join(f.Directory, strings.ToLower(backupID.String())), nil
}

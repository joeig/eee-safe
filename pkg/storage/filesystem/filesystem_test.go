package filesystem

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/joeig/eee-safe/pkg/threema"
)

type mockFileInfo struct {
	ModTimeOutput time.Time
}

func (m *mockFileInfo) Name() string {
	return ""
}

func (m *mockFileInfo) Size() int64 {
	return 0
}

func (m *mockFileInfo) Mode() os.FileMode {
	return 0
}

func (m *mockFileInfo) ModTime() time.Time {
	return m.ModTimeOutput
}

func (m *mockFileInfo) IsDir() bool {
	return false
}

func (m *mockFileInfo) Sys() any {
	return nil
}

type mockFile struct {
	WriteFileErr   error
	ReadFileOutput []byte
	ReadFileErr    error
	RemoveErr      error
	StatOutput     os.FileInfo
	StatErr        error
}

func (m *mockFile) WriteFile(_ string, _ []byte, _ os.FileMode) error {
	return m.WriteFileErr
}

func (m *mockFile) ReadFile(_ string) ([]byte, error) {
	return m.ReadFileOutput, m.ReadFileErr
}

func (m *mockFile) Remove(_ string) error {
	return m.RemoveErr
}

func (m *mockFile) Stat(_ string) (os.FileInfo, error) {
	return m.StatOutput, m.StatErr
}

func TestFilesystem_SetFile(t *testing.T) {
	file := &mockFile{}
	filesystem := &Filesystem{}

	filesystem.SetFile(file)

	if filesystem.file != file {
		t.Error("invalid file")
	}
}

func TestFilesystem_PutBackup(t *testing.T) {
	filesystem := &Filesystem{file: &mockFile{}, Directory: "mock"}

	if err := filesystem.PutBackup(context.Background(), &threema.BackupInput{BackupID: threema.BackupID{'b'}, EncryptedBackup: []byte{'e'}}); err != nil {
		t.Error("unexpected error")
	}
}

func TestFilesystem_PutBackup_error(t *testing.T) {
	filesystem := &Filesystem{file: &mockFile{WriteFileErr: errors.New("mock")}, Directory: "mock"}

	if err := filesystem.PutBackup(context.Background(), &threema.BackupInput{BackupID: threema.BackupID{'b'}, EncryptedBackup: []byte{'e'}}); err == nil {
		t.Error("no error")
	}
}

func TestFilesystem_GetBackup(t *testing.T) {
	backupID := threema.BackupID{'b'}
	encryptedBackup := []byte{'r'}
	creationTime := time.UnixMicro(1)
	filesystem := &Filesystem{file: &mockFile{ReadFileOutput: encryptedBackup, StatOutput: &mockFileInfo{ModTimeOutput: creationTime}}, Directory: "mock"}

	output, err := filesystem.GetBackup(context.Background(), backupID)

	if output == nil {
		t.Error("no output")
		t.FailNow()
	}

	if output.BackupID != backupID {
		t.Error("invalid backup ID")
	}

	if !bytes.Equal(output.EncryptedBackup, encryptedBackup) {
		t.Error("invalid encrypted backup")
	}

	if output.CreationTime != creationTime {
		t.Error("invalid creation time")
	}

	if err != nil {
		t.Error("unexpected error")
	}
}

func TestFilesystem_GetBackup_statError(t *testing.T) {
	filesystem := &Filesystem{file: &mockFile{StatErr: errors.New("mock")}, Directory: "mock"}

	output, err := filesystem.GetBackup(context.Background(), threema.BackupID{'b'})

	if output == nil {
		t.Error("no output")
	}

	if err == nil {
		t.Error("no error")
	}
}

func TestFilesystem_GetBackup_readFileError(t *testing.T) {
	filesystem := &Filesystem{file: &mockFile{ReadFileErr: errors.New("mock"), StatOutput: &mockFileInfo{ModTimeOutput: time.UnixMicro(1)}}, Directory: "mock"}

	output, err := filesystem.GetBackup(context.Background(), threema.BackupID{'b'})

	if output == nil {
		t.Error("no output")
	}

	if err == nil {
		t.Error("no error")
	}
}

func TestFilesystem_DeleteBackup(t *testing.T) {
	filesystem := &Filesystem{file: &mockFile{}, Directory: "mock"}

	if err := filesystem.DeleteBackup(context.Background(), threema.BackupID{'b'}); err != nil {
		t.Error("unexpected error")
	}
}

func TestFilesystem_DeleteBackup_error(t *testing.T) {
	filesystem := &Filesystem{file: &mockFile{RemoveErr: errors.New("mock")}, Directory: "mock"}

	if err := filesystem.DeleteBackup(context.Background(), threema.BackupID{'b'}); err == nil {
		t.Error("no error")
	}
}

func TestGenerateFileName(t *testing.T) {
	filesystem := Filesystem{
		Directory:   "foo",
		Permissions: 0600,
	}

	backupID, _ := threema.ConvertToBackupID("c8435aaa32e3de72426e04e845d1251d87df5aaa32e3de72426e04e845d1251d")

	fileName, err := filesystem.generateFileName(backupID)
	if err != nil {
		t.Error("File name generation failed, should be OK")
	}

	if fileName != fmt.Sprintf("%s/%s", filesystem.Directory, backupID.String()) {
		t.Error("Generated file name is wrong")
	}
}

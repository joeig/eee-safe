package threema

import (
	"testing"
)

func TestConvertStringToBackupID(t *testing.T) {
	backupIDTemplate := [32]byte{135, 223, 90, 170, 50, 227, 222, 114, 66, 110, 4, 232, 69, 209, 37, 29, 135, 223, 90, 170, 50, 227, 222, 114, 66, 110, 4, 232, 69, 209, 37, 29}
	backupIDTemplateString := "87df5aaa32e3de72426e04e845d1251d87df5aaa32e3de72426e04e845d1251d"

	t.Run("TestBackupIDOK", func(t *testing.T) {
		b, err := ConvertToBackupID(backupIDTemplateString)
		if b != backupIDTemplate {
			t.Errorf("Backup ID conversion mismatch: %v != %v", b, backupIDTemplate)
		}
		if err != nil {
			t.Errorf("Backup ID conversion failed, should be OK: %v", err)
		}
	})
	t.Run("TestTooShortBackupID", func(t *testing.T) {
		if _, err := ConvertToBackupID("abcdef"); err == nil {
			t.Error("Backup ID is too short, but validation returns OK")
		}
	})
	t.Run("TestTooLongBackupID", func(t *testing.T) {
		if _, err := ConvertToBackupID(backupIDTemplateString + "ab"); err == nil {
			t.Error("Backup ID is too long, but validation returns OK")
		}
	})
}

func TestConvertBackupIDToString(t *testing.T) {
	backupIDTemplate := BackupID{135, 223, 90, 170, 50, 227, 222, 114, 66, 110, 4, 232, 69, 209, 37, 29, 135, 223, 90, 170, 50, 227, 222, 114, 66, 110, 4, 232, 69, 209, 37, 29}
	backupIDTemplateString := "87df5aaa32e3de72426e04e845d1251d87df5aaa32e3de72426e04e845d1251d"

	t.Run("TestBackupIDOK", func(t *testing.T) {
		b := backupIDTemplate.String()
		if b != backupIDTemplateString {
			t.Errorf("Backup ID conversion mismatch: %v != %v", b, backupIDTemplate)
		}
	})
}

func TestValidateEncryptedBackup(t *testing.T) {
	t.Run("TestEncryptedBackupOK", func(t *testing.T) {
		b := EncryptedBackup{'x', 'y', 'z'}
		if b.Validate(8) != nil {
			t.Error("Encrypted backup validation failed, should be OK")
		}
	})
	t.Run("TestEmptyEncryptedBackup", func(t *testing.T) {
		b := EncryptedBackup{}
		if b.Validate(8) == nil {
			t.Error("Backup ID is too short, but validation returns OK")
		}
	})
	t.Run("TestTooLongEncryptedBackup", func(t *testing.T) {
		b := EncryptedBackup{'x', 'y', 'z'}
		if b.Validate(2) == nil {
			t.Error("Backup ID is too long, but validation returns OK")
		}
	})
}

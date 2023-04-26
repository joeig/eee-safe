package threema

import (
	"errors"
	"testing"
)

func TestErrBackupIDInvalid_Error(t *testing.T) {
	e := &ErrBackupIDInvalid{BackupID: BackupID{'f', 'o', 'o'}}

	if e.Error() != "backup ID invalid (32 bytes)" {
		t.Error("invalid error representation")
	}
}

func TestErrBackupIDStringInvalid_Error(t *testing.T) {
	e := &ErrBackupIDStringInvalid{BackupIDString: "foo", UpstreamError: errors.New("upstream-error")}

	if e.Error() != "backup ID string invalid (3 bytes): upstream-error" {
		t.Error("invalid error representation")
	}
}

func TestErrBackupIDLength_Error(t *testing.T) {
	e := &ErrBackupIDLength{BackupID: BackupID{'f', 'o', 'o'}, DesiredLength: 16}

	if e.Error() != "wrong backup ID length (32 bytes instead of 16 bytes)" {
		t.Error("invalid error representation")
	}
}

func TestErrBackupIDStringLength_Error(t *testing.T) {
	e := &ErrBackupIDStringLength{BackupIDString: "foo", DesiredLength: 32}

	if e.Error() != "wrong backup ID string length (3 bytes instead of 32 bytes)" {
		t.Error("invalid error representation")
	}
}

func TestErrEncryptedBackupInvalid_Error(t *testing.T) {
	e := &ErrEncryptedBackupInvalid{EncryptedBackup: EncryptedBackup{'f', 'o', 'o'}}

	if e.Error() != "encrypted backup invalid (3 bytes)" {
		t.Error("invalid error representation")
	}
}

func TestErrEncryptedBackupLength_Error(t *testing.T) {
	e := &ErrEncryptedBackupLength{EncryptedBackup: EncryptedBackup{'f', 'o', 'o'}, MaxLength: 1}

	if e.Error() != "wrong encrypted backup string length (3 bytes is greater than 1 bytes)" {
		t.Error("invalid error representation")
	}
}

func TestErrInvalidUserAgent_Error(t *testing.T) {
	e := &ErrInvalidUserAgent{UserAgent: "foo"}

	if e.Error() != "user agent invalid: foo" {
		t.Error("invalid error representation")
	}
}

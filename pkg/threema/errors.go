package threema

import "fmt"

// ErrBackupIDInvalid is being used whenever a backup ID is invalid
type ErrBackupIDInvalid struct {
	BackupID BackupID
}

func (e *ErrBackupIDInvalid) Error() string {
	return fmt.Sprintf("backup ID invalid (%d bytes)", len(e.BackupID))
}

// ErrBackupIDStringInvalid is being used whenever a backup ID string is invalid
type ErrBackupIDStringInvalid struct {
	BackupIDString string
	UpstreamError  error
}

func (e *ErrBackupIDStringInvalid) Error() string {
	return fmt.Sprintf("backup ID string invalid (%d bytes): %v", len(e.BackupIDString), e.UpstreamError)
}

// ErrBackupIDLength is being used whenever the length of a backup ID is invalid
type ErrBackupIDLength struct {
	BackupID      BackupID
	DesiredLength uint
}

func (e *ErrBackupIDLength) Error() string {
	return fmt.Sprintf("wrong backup ID length (%d bytes instead of %d bytes)", len(e.BackupID), e.DesiredLength)
}

// ErrBackupIDStringLength is being used whenever the length of a backup ID string is invalid
type ErrBackupIDStringLength struct {
	BackupIDString string
	DesiredLength  uint
}

func (e *ErrBackupIDStringLength) Error() string {
	return fmt.Sprintf("wrong backup ID string length (%d bytes instead of %d bytes)", len(e.BackupIDString), e.DesiredLength)
}

// ErrEncryptedBackupInvalid is being used whenever an encrypted backup is invalid
type ErrEncryptedBackupInvalid struct {
	EncryptedBackup EncryptedBackup
}

func (e *ErrEncryptedBackupInvalid) Error() string {
	return fmt.Sprintf("encrypted backup invalid (%d bytes)", len(e.EncryptedBackup))
}

// ErrEncryptedBackupLength is being used whenever the length of an encrypted backup is invalid
type ErrEncryptedBackupLength struct {
	EncryptedBackup EncryptedBackup
	MaxLength       uint
}

func (e *ErrEncryptedBackupLength) Error() string {
	return fmt.Sprintf("wrong encrypted backup string length (%d bytes is greater than %d bytes)", len(e.EncryptedBackup), e.MaxLength)
}

// ErrInvalidUserAgent is being used whenever a user agent is invalid
type ErrInvalidUserAgent struct {
	UserAgent string
}

func (e *ErrInvalidUserAgent) Error() string {
	return fmt.Sprintf("user agent invalid: %s", e.UserAgent)
}

package threema

import "fmt"

// BackupIDInvalidError is being used whenever a backup ID is invalid
type BackupIDInvalidError struct {
	BackupID BackupID
}

func (e *BackupIDInvalidError) Error() string {
	return fmt.Sprintf("BackupInput ID invalid (%d bytes)", len(e.BackupID))
}

// BackupIDStringInvalidError is being used whenever a backup ID string is invalid
type BackupIDStringInvalidError struct {
	BackupIDString string
	UpstreamError  error
}

func (e *BackupIDStringInvalidError) Error() string {
	return fmt.Sprintf("BackupInput ID string invalid (%d bytes): %v", len(e.BackupIDString), e.UpstreamError)
}

// BackupIDLengthError is being used whenever the length of a backup ID is invalid
type BackupIDLengthError struct {
	BackupID      BackupID
	DesiredLength uint
}

func (e *BackupIDLengthError) Error() string {
	return fmt.Sprintf("Wrong backup ID length (%d bytes instead of %d bytes)", len(e.BackupID), e.DesiredLength)
}

// BackupIDStringLengthError is being used whenever the length of a backup ID string is invalid
type BackupIDStringLengthError struct {
	BackupIDString string
	DesiredLength  uint
}

func (e *BackupIDStringLengthError) Error() string {
	return fmt.Sprintf("Wrong backup ID string length (%d bytes instead of %d bytes)", len(e.BackupIDString), e.DesiredLength)
}

// EncryptedBackupInvalidError is being used whenever an encrypted backup is invalid
type EncryptedBackupInvalidError struct {
	EncryptedBackup EncryptedBackup
}

func (e *EncryptedBackupInvalidError) Error() string {
	return fmt.Sprintf("Encrypted backup invalid (%d bytes)", len(e.EncryptedBackup))
}

// EncryptedBackupLengthError is being used whenever the length of an encrypted backup is invalid
type EncryptedBackupLengthError struct {
	EncryptedBackup EncryptedBackup
	MaxLength       uint
}

func (e *EncryptedBackupLengthError) Error() string {
	return fmt.Sprintf("Wrong encrypted backup string length (%d bytes is greater than %d bytes)", len(e.EncryptedBackup), e.MaxLength)
}

// InvalidUserAgentError is being used whenever a user agent is invalid
type InvalidUserAgentError struct {
	UserAgent string
}

func (e *InvalidUserAgentError) Error() string {
	return fmt.Sprintf("User agent invalid: %s", e.UserAgent)
}

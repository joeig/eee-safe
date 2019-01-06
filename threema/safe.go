package threema

import (
	"encoding/hex"
	"time"
)

// BackupIDLength is the length in bytes of the backup ID
const BackupIDLength = 32

// BackupID is the backup ID type
type BackupID [BackupIDLength]byte

// EncryptedBackup is the encrypted backup type
type EncryptedBackup []byte

// BackupInput is a struct containing the backup
type BackupInput struct {
	BackupID        BackupID
	EncryptedBackup EncryptedBackup
	RetentionDays   uint
}

// BackupOutput is a struct containing the backup
type BackupOutput struct {
	BackupID        BackupID
	EncryptedBackup EncryptedBackup
	CreationTime    time.Time
	ExpirationTime  time.Time
}

// ConvertToBackupID converts string to BackupID
func ConvertToBackupID(s string) (BackupID, error) {
	backupIDHex, err := hex.DecodeString(s)
	if err != nil {
		return BackupID{}, &BackupIDStringInvalidError{BackupIDString: s, UpstreamError: err}
	}
	sLength := len(backupIDHex)
	if sLength != BackupIDLength {
		return BackupID{}, &BackupIDStringLengthError{BackupIDString: s, DesiredLength: BackupIDLength * 2}
	}
	var threemaSafeBackupID BackupID
	copy(threemaSafeBackupID[:], backupIDHex)
	return threemaSafeBackupID, nil
}

// String converts BackupID to string
func (b *BackupID) String() string {
	var backupID []byte
	for _, b := range b {
		backupID = append(backupID, b)
	}
	return hex.EncodeToString(backupID)
}

// ValidateEncryptedBackup validates the EncryptedBackup
func (t *EncryptedBackup) Validate(maxBackupBytes uint) error {
	tLength := len(*t)
	if tLength == 0 || uint(tLength) > maxBackupBytes {
		return &EncryptedBackupLengthError{EncryptedBackup: *t, MaxLength: maxBackupBytes}
	}
	return nil
}

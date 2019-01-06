package filesystem

import (
	"fmt"
	"github.com/joeig/eee-safe/threema"
	"testing"
)

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

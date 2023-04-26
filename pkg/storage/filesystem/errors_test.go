package filesystem

import (
	"errors"
	"testing"
)

func TestErrDirectoryInvalid_Error(t *testing.T) {
	e := &ErrDirectoryInvalid{Directory: "foo"}

	if e.Error() != "directory invalid: \"foo\"" {
		t.Error("invalid error representation")
	}
}

func TestErrFileNotWritable_Error(t *testing.T) {
	e := &ErrFileNotWritable{FileName: "foo", UpstreamError: errors.New("upstream-error")}

	if e.Error() != "file not writable (\"foo\"): upstream-error" {
		t.Error("invalid error representation")
	}
}

func TestErrFileNotReadable_Error(t *testing.T) {
	e := &ErrFileNotReadable{FileName: "foo", UpstreamError: errors.New("upstream-error")}

	if e.Error() != "file not readable (\"foo\"): upstream-error" {
		t.Error("invalid error representation")
	}
}

func TestErrFileNotRemovable_Error(t *testing.T) {
	e := &ErrFileNotRemovable{FileName: "foo", UpstreamError: errors.New("upstream-error")}

	if e.Error() != "file not removable (\"foo\"): upstream-error" {
		t.Error("invalid error representation")
	}
}

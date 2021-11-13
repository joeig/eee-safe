package filesystem

import "fmt"

// ErrDirectoryInvalid is being used whenever a given directory is invalid
type ErrDirectoryInvalid struct {
	Directory string
}

func (e *ErrDirectoryInvalid) Error() string {
	return fmt.Sprintf("directory invalid: %q", e.Directory)
}

// ErrFileNotWritable is being used whenever a file handler is not writable
type ErrFileNotWritable struct {
	FileName      string
	UpstreamError error
}

func (e *ErrFileNotWritable) Error() string {
	return fmt.Sprintf("file not writable (%q): %v", e.FileName, e.UpstreamError)
}

// ErrFileNotReadable is being used whenever a file handler is not readable
type ErrFileNotReadable struct {
	FileName      string
	UpstreamError error
}

func (e *ErrFileNotReadable) Error() string {
	return fmt.Sprintf("file not readable (%q): %v", e.FileName, e.UpstreamError)
}

// ErrFileNotRemovable is being used whenever a file handler is not removable
type ErrFileNotRemovable struct {
	FileName      string
	UpstreamError error
}

func (e *ErrFileNotRemovable) Error() string {
	return fmt.Sprintf("file not removable (%q): %v", e.FileName, e.UpstreamError)
}

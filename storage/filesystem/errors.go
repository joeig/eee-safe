package filesystem

import "fmt"

// DirectoryInvalidError is being used whenever a given directory is invalid
type DirectoryInvalidError struct {
	Directory string
}

func (e *DirectoryInvalidError) Error() string {
	return fmt.Sprintf("Directory invalid: %s", e.Directory)
}

// FileNotWritableError is being used whenever a file handler is not writable
type FileNotWritableError struct {
	FileName      string
	UpstreamError error
}

func (e *FileNotWritableError) Error() string {
	return fmt.Sprintf("File not writable (\"%s\"): %v", e.FileName, e.UpstreamError)
}

// FileNotReadableError is being used whenever a file handler is not readable
type FileNotReadableError struct {
	FileName      string
	UpstreamError error
}

func (e *FileNotReadableError) Error() string {
	return fmt.Sprintf("File not readable (\"%s\"): %v", e.FileName, e.UpstreamError)
}

// FileNotRemovableError is being used whenever a file handler is not removable
type FileNotRemovableError struct {
	FileName      string
	UpstreamError error
}

func (e *FileNotRemovableError) Error() string {
	return fmt.Sprintf("File not removable (\"%s\"): %v", e.FileName, e.UpstreamError)
}

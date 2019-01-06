package filesystem

import "fmt"

type DirectoryInvalidError struct {
	Directory string
}

func (e *DirectoryInvalidError) Error() string {
	return fmt.Sprintf("Directory invalid: %s", e.Directory)
}

type FileNotWritableError struct {
	FileName      string
	UpstreamError error
}

func (e *FileNotWritableError) Error() string {
	return fmt.Sprintf("File not writable (\"%s\"): %v", e.FileName, e.UpstreamError)
}

type FileNotReadableError struct {
	FileName      string
	UpstreamError error
}

func (e *FileNotReadableError) Error() string {
	return fmt.Sprintf("File not readable (\"%s\"): %v", e.FileName, e.UpstreamError)
}

type FileNotRemovableError struct {
	FileName      string
	UpstreamError error
}

func (e *FileNotRemovableError) Error() string {
	return fmt.Sprintf("File not removable (\"%s\"): %v", e.FileName, e.UpstreamError)
}

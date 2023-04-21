package filesystem

import "os"

// OSFile implements a File using operating system handlers.
type OSFile struct{}

// WriteFile writes a file to the local file system.
func (o *OSFile) WriteFile(name string, data []byte, perm os.FileMode) error {
	return os.WriteFile(name, data, perm)
}

// ReadFile reads a file from the local file system.
func (o *OSFile) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}

// Remove deletes a file from the local file system.
func (o *OSFile) Remove(name string) error {
	return os.Remove(name)
}

// Stat returns infos about a file on the local file system.
func (o *OSFile) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

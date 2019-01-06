package filesystem

import (
	"os"
)

// Filesystem defines the configuration of the filesystem storage backend type
type Filesystem struct {
	Directory   string      `mapstructure:"directory"`
	Permissions os.FileMode `mapstructure:"permissions"`
}

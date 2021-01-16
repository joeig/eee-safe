package filesystem

import (
	"os"
	"sync"
)

// Filesystem defines the configuration of the filesystem storage backend type
type Filesystem struct {
	mu sync.RWMutex

	Directory   string      `mapstructure:"directory"`
	Permissions os.FileMode `mapstructure:"permissions"`
}

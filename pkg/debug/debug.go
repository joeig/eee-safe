package debug

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

// Debug contains the current state of debug operation
var Debug bool

// PrintPrefix contains the prefix of debug output lines
const PrintPrefix = "[eee-safe-debug] "

// Printf wraps fmt.Printf and displays debug messages unless Debug is set to false
func Printf(format string, values ...interface{}) {
	if Debug || flag.Lookup("test.v") != nil {
		if !strings.HasSuffix(format, "\n") {
			format += "\n"
		}
		_, _ = fmt.Fprintf(os.Stderr, PrintPrefix+format, values...)
	}
}

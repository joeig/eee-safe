package debug

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

var Debug bool

func Printf(format string, values ...interface{}) {
	if Debug || flag.Lookup("test.v") != nil {
		if !strings.HasSuffix(format, "\n") {
			format += "\n"
		}
		fmt.Fprintf(os.Stderr, "[eee-safe-debug] "+format, values...)
	}
}

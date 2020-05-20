package debug

import (
	"bytes"
	"flag"
	"io"
	"os"
	"testing"
)

func assertDebugOutput(t *testing.T, debugMode bool, input string, outputExpected bool) {
	oldStderrHandle := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	Debug = debugMode
	flag.CommandLine = flag.NewFlagSet("", flag.ExitOnError)

	Printf(input)

	outChannel := make(chan string)

	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		outChannel <- buf.String()
	}()

	_ = w.Close()
	os.Stderr = oldStderrHandle
	out := <-outChannel

	expectedOutputWithPrefix := ""
	if outputExpected {
		expectedOutputWithPrefix = PrintPrefix + input + "\n"
	}

	if out != expectedOutputWithPrefix {
		t.Errorf("Debug output does not match: %s", out)
	}
}

func TestDebugPrintf(t *testing.T) {
	t.Run("DebugTrueOutput", func(t *testing.T) {
		assertDebugOutput(t, true, "foo", true)
	})
	t.Run("DebugFalseNoOutput", func(t *testing.T) {
		assertDebugOutput(t, false, "foo", false)
	})
}

package app

import (
	"bytes"
	"log"
	"os"
	"testing"
)

func TestSanityCheck(t *testing.T) {
	// Capture log output
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()

	// Unset environment variable to trigger log message
	os.Unsetenv("SERVER_ADDRESS")

	// Run the sanity check
	sanityCheck()

	// Check the log output
	if !bytes.Contains(buf.Bytes(), []byte("environment variable SERVER_ADDRESS not defined")) {
		t.Error("Expected log message not found")
	}
}

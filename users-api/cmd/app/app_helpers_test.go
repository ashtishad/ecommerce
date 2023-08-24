package app

import (
	"bytes"
	"log"
	"os"
	"testing"
)

func TestSanityCheck(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()

	_ = os.Unsetenv("SERVER_ADDRESS")

	sanityCheck()

	if !bytes.Contains(buf.Bytes(), []byte("environment variable SERVER_ADDRESS not defined")) {
		t.Error("Expected log message not found")
	}
}

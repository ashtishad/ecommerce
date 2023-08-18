package app

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
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

func TestWriteResponse(t *testing.T) {
	tests := []struct {
		code     int
		data     interface{}
		expected string
	}{
		{http.StatusOK, map[string]string{"message": "success"}, `{"message":"success"}`},
		{http.StatusBadRequest, map[string]string{"error": "bad request"}, `{"error":"bad request"}`},
		{http.StatusInternalServerError, map[string]string{"error": "internal server error"}, `{"error":"internal server error"}`},
	}

	for _, test := range tests {
		rr := httptest.NewRecorder()
		writeResponse(rr, test.code, test.data)

		if status := rr.Code; status != test.code {
			t.Errorf("handler returned wrong status code: got %v want %v", status, test.code)
		}

		expected := test.expected + "\n"
		if rr.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
		}
	}
}

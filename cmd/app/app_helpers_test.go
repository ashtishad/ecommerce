package app

import (
	"bytes"
	"github.com/ashtishad/ecommerce/domain"
	"log"
	"net/http"
	"net/http/httptest"
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

func TestValidateCreateUserInput(t *testing.T) {
	tests := []struct {
		name    string
		input   domain.NewUserRequestDTO
		wantErr bool
	}{
		{
			name: "valid input",
			input: domain.NewUserRequestDTO{
				Email:        "test@example.com",
				Password:     "password123",
				FullName:     "John Doe",
				Phone:        "1234567890",
				SignUpOption: "general",
			},
			wantErr: false,
		},
		{
			name: "invalid email",
			input: domain.NewUserRequestDTO{
				Email:        "invalid-email",
				Password:     "password123",
				FullName:     "John Doe",
				Phone:        "1234567890",
				SignUpOption: "general",
			},
			wantErr: true,
		},
		{
			name: "invalid password",
			input: domain.NewUserRequestDTO{
				Email:        "test@example.com",
				Password:     "pass",
				FullName:     "John Doe",
				Phone:        "1234567890",
				SignUpOption: "general",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateCreateUserInput(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateCreateUserInput() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

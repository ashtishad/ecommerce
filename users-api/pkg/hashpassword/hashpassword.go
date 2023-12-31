package hashpassword

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// GenerateSalt makes 16 bytes random salt, converts salt to hexadecimal string to avoid database encoding problems
func GenerateSalt() (string, error) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)

	if err != nil {
		return "", fmt.Errorf("could not generate salt: %w", err)
	}

	return hex.EncodeToString(salt), nil
}

// HashPassword takes password and salt,
// Decodes salt from hexadecimal string, then returns hashedPassword
func HashPassword(password string, saltHex string) string {
	passwordBytes := []byte(password)
	salt, _ := hex.DecodeString(saltHex)

	passwordBytes = append(passwordBytes, salt...)

	hash := sha256.Sum256(passwordBytes)

	hashedPassword := hex.EncodeToString(hash[:])

	return hashedPassword
}

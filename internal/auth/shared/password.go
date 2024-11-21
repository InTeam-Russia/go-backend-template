package shared

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"runtime"

	"golang.org/x/crypto/argon2"
)

const HASH_PASSWORD_TIME = 1
const HASH_PASSWORD_MEMORY = 64 * 1024
const HASH_PASSWORD_LEN = 32

func GenerateSalt(size int) ([]byte, error) {
	salt := make([]byte, size)
	if _, err := rand.Read(salt); err != nil {
		return nil, fmt.Errorf("Error generating salt: %w", err)
	}
	return salt, nil
}

func HashPassword(password string, salt []byte) []byte {
	hash := argon2.IDKey(
		[]byte(password),
		salt,
		HASH_PASSWORD_TIME,
		HASH_PASSWORD_MEMORY,
		(uint8)(runtime.NumCPU()),
		HASH_PASSWORD_LEN,
	)
	return hash
}

func ValidPassword(enteredPassword string, passwordHash []byte, passwordSalt []byte) bool {
	fmt.Println(passwordSalt)
	enteredPasswordHash := HashPassword(enteredPassword, passwordSalt)
	return bytes.Equal(enteredPasswordHash, passwordHash)
}

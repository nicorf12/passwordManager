package security

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"golang.org/x/crypto/pbkdf2"
)

// Generar un hash con PBKDF2
func GenerateHash(password string, salt []byte) string {
	iterations := 100_000
	hash := pbkdf2.Key([]byte(password), salt, iterations, 32, sha256.New)
	return base64.StdEncoding.EncodeToString(hash)
}

// Generar un salt aleatorio
func GenerateSalt() ([]byte, error) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	return salt, err
}

// Verificar una contraseña
func VerifyPassword(password string, storedHash string, salt []byte) bool {
	iterations := 100_000
	hash := pbkdf2.Key([]byte(password), salt, iterations, 32, sha256.New)
	return base64.StdEncoding.EncodeToString(hash) == storedHash
}

package test

import (
	"password_manager/security"
	"testing"
)

func TestEncryptDecryptAES(t *testing.T) {
	password := "testpassword"
	originalText := "Hello, World!"

	encrypted, err := security.EncryptAES([]byte(originalText), password)
	if err != nil {
		t.Fatalf("Error al encriptar: %v", err)
	}

	decrypted, err := security.DecryptAES(encrypted, password)
	if err != nil {
		t.Fatalf("Error al desencriptar: %v", err)
	}

	if decrypted != originalText {
		t.Fatalf("Texto desencriptado no coincide con el original. Esperado: %s, Obtenido: %s", originalText, decrypted)
	}
}

func TestEncryptDecryptXChaCha20Poly1305(t *testing.T) {
	key := "thisisaverysecurekey123456"
	originalText := "Hello, World!"

	encrypted, err := security.EncryptXChaCha20Poly1305([]byte(originalText), key)
	if err != nil {
		t.Fatalf("Error al encriptar: %v", err)
	}

	decrypted, err := security.DecryptXChaCha20Poly1305(encrypted, key)
	if err != nil {
		t.Fatalf("Error al desencriptar: %v", err)
	}

	if decrypted != originalText {
		t.Fatalf("Texto desencriptado no coincide con el original. Esperado: %s, Obtenido: %s", originalText, decrypted)
	}
}

func TestEncryptDecryptDES(t *testing.T) {
	key := "12345678"
	originalText := "Hello, World!"

	encrypted, err := security.EncryptDES([]byte(originalText), key)
	if err != nil {
		t.Fatalf("Error al encriptar: %v", err)
	}

	decrypted, err := security.DecryptDES(encrypted, key)
	if err != nil {
		t.Fatalf("Error al desencriptar: %v", err)
	}

	if decrypted != originalText {
		t.Fatalf("Texto desencriptado no coincide con el original. Esperado: %s, Obtenido: %s", originalText, decrypted)
	}
}

func TestInvalidKeyLengthDES(t *testing.T) {
	invalidKey := "short"
	originalText := "Hello, World!"

	_, err := security.EncryptDES([]byte(originalText), invalidKey)
	if err == nil {
		t.Fatal("Se esperaba un error debido a la longitud incorrecta de la clave")
	}
}

func TestGenerateHash(t *testing.T) {
	password := "testpassword"
	salt, err := security.GenerateSalt()
	if err != nil {
		t.Fatalf("Error al generar salt: %v", err)
	}

	hash := security.GenerateHash(password, salt)

	if hash == "" {
		t.Fatal("El hash generado está vacío")
	}
}

func TestVerifyPassword(t *testing.T) {
	password := "testpassword"
	salt, err := security.GenerateSalt()
	if err != nil {
		t.Fatalf("Error al generar salt: %v", err)
	}

	storedHash := security.GenerateHash(password, salt)

	if !security.VerifyPassword(password, storedHash, salt) {
		t.Fatal("La verificación de la contraseña falló")
	}

	if security.VerifyPassword("wrongpassword", storedHash, salt) {
		t.Fatal("La verificación de la contraseña incorrecta debería haber fallado")
	}
}

func TestGenerateSecurePassword(t *testing.T) {
	var (
		upperSet   = "ABCDEFGHIJKLMNÑOPQRSTUVWXYZ"
		lowerSet   = "abcdefghijklmnñopqrstuvwxyz"
		numberSet  = "0123456789"
		specialSet = "!@#$%^&*()-_=+[]{}<>?/|~"
	)
	length := 12
	useUpper := true
	useLower := true
	useNumbers := true
	useSpecials := true

	password, err := security.GenerateSecurePassword(length, useUpper, useLower, useNumbers, useSpecials)
	if err != nil {
		t.Fatalf("Error al generar la contraseña segura: %v", err)
	}

	if len(password) != length {
		t.Fatalf("La longitud de la contraseña debería ser %d, pero obtuvo %d", length, len(password))
	}

	// Verificar que la contraseña contiene caracteres de diferentes conjuntos
	if !containsCharacterSet(password, useUpper, upperSet) ||
		!containsCharacterSet(password, useLower, lowerSet) ||
		!containsCharacterSet(password, useNumbers, numberSet) ||
		!containsCharacterSet(password, useSpecials, specialSet) {
		t.Error("La contraseña generada no contiene los conjuntos de caracteres esperados")
	}
}

// Función auxiliar para verificar si un conjunto de caracteres está presente en la contraseña
func containsCharacterSet(password string, includeSet bool, charSet string) bool {
	if !includeSet {
		return true
	}
	for _, char := range password {
		if contains(charSet, char) {
			return true
		}
	}
	return false
}

// Función auxiliar para verificar si un conjunto de caracteres contiene un caracter
func contains(set string, char rune) bool {
	for _, c := range set {
		if c == char {
			return true
		}
	}
	return false
}

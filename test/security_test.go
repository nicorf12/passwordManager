package test

import (
	"password_manager/security"
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	password := "testpassword"
	originalText := "Hello, World!"

	// Encriptar
	encrypted, err := security.Encrypt([]byte(originalText), password)
	if err != nil {
		t.Fatalf("Error al encriptar: %v", err)
	}

	// Desencriptar
	decrypted, err := security.Decrypt(encrypted, password)
	if err != nil {
		t.Fatalf("Error al desencriptar: %v", err)
	}

	// Verificar que el texto original y el desencriptado coinciden
	if decrypted != originalText {
		t.Fatalf("Texto desencriptado no coincide con el original. Esperado: %s, Obtenido: %s", originalText, decrypted)
	}
}

func TestGenerateHash(t *testing.T) {
	password := "testpassword"
	salt, err := security.GenerateSalt()
	if err != nil {
		t.Fatalf("Error al generar salt: %v", err)
	}

	// Generar el hash
	hash := security.GenerateHash(password, salt)

	// Verificar que el hash no sea vacío
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

	// Generar el hash del password
	storedHash := security.GenerateHash(password, salt)

	// Verificar que la contraseña sea correcta
	if !security.VerifyPassword(password, storedHash, salt) {
		t.Fatal("La verificación de la contraseña falló")
	}

	// Verificar una contraseña incorrecta
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

	// Generar la contraseña segura
	password, err := security.GenerateSecurePassword(length, useUpper, useLower, useNumbers, useSpecials)
	if err != nil {
		t.Fatalf("Error al generar la contraseña segura: %v", err)
	}

	// Verificar que la contraseña tenga la longitud correcta
	if len(password) != length {
		t.Fatalf("La longitud de la contraseña debería ser %d, pero obtuvo %d", length, len(password))
	}

	// Verificar que la contraseña contiene caracteres de diferentes conjuntos
	if !containsCharacterSet(password, useUpper, upperSet) ||
		!containsCharacterSet(password, useLower, lowerSet) ||
		!containsCharacterSet(password, useNumbers, numberSet) ||
		!containsCharacterSet(password, useSpecials, specialSet) {
		t.Fatal("La contraseña generada no contiene los conjuntos de caracteres esperados")
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

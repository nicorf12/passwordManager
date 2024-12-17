package security

import (
	"strings"
	"unicode"
)

// CalculatePasswordSecurity evalúa la contraseña y devuelve un valor entre 0 y 100
func CalculatePasswordSecurity(password string) float64 {
	var puntos float64
	longitud := len(password)
	tieneMayuscula, tieneMinuscula := false, false
	tieneDigito, tieneEspecial := false, false
	tienePatronSimple := containsSimplePattern(password)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			tieneMayuscula = true
		case unicode.IsLower(char):
			tieneMinuscula = true
		case unicode.IsDigit(char):
			tieneDigito = true
		case strings.ContainsRune("!@#$%^&*()_+[]{}|;:'\",.<>?/`~", char):
			tieneEspecial = true
		}
	}

	if longitud >= 8 {
		puntos += 25
	}
	if tieneMayuscula {
		puntos += 20
	}
	if tieneMinuscula {
		puntos += 20
	}
	if tieneDigito {
		puntos += 20
	}
	if tieneEspecial {
		puntos += 15
	}
	if tienePatronSimple {
		puntos = 0
	}
	if puntos > 100 {
		return 100
	}
	return puntos
}

// Verifica patrones comunes
func containsSimplePattern(password string) bool {
	simplePatterns := []string{
		"123456", "abcdef", "qwerty", "password",
		"111111", "letmein", "iloveyou", "admin",
	}

	for _, pattern := range simplePatterns {
		if strings.Contains(strings.ToLower(password), pattern) {
			return true
		}
	}

	return containsSequential(password)
}

// Verifica secuencias numéricas y de letras
func containsSequential(password string) bool {
	for i := 0; i < len(password)-2; i++ {
		if password[i]+1 == password[i+1] && password[i]+2 == password[i+2] {
			return true
		}
		if password[i]-1 == password[i+1] && password[i]-2 == password[i+2] {
			return true
		}
	}
	return false
}

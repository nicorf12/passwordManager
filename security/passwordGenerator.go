package security

import (
	"crypto/rand"
	"errors"
	"math/big"
)

// Configuración interna no exportable
type config struct {
	length          int
	includeUpper    bool
	includeLower    bool
	includeNumbers  bool
	includeSpecials bool
}

// Caracteres disponibles
var (
	upperSet   = "ABCDEFGHIJKLMNÑOPQRSTUVWXYZ"
	lowerSet   = "abcdefghijklmnñopqrstuvwxyz"
	numberSet  = "0123456789"
	specialSet = "!@#$%^&*()-_=+[]{}<>?/|~"
)

// GenerateSecurePassword Función pública para generar contraseñas
func GenerateSecurePassword(length int, useUpper, useLower, useNumbers, useSpecials bool) (string, error) {
	if length <= 0 {
		return "", errors.New("la longitud de la contraseña debe ser mayor que cero")
	}

	cfg := config{
		length:          length,
		includeUpper:    useUpper,
		includeLower:    useLower,
		includeNumbers:  useNumbers,
		includeSpecials: useSpecials,
	}

	charPool := buildCharPool(cfg)
	if charPool == "" {
		return "", errors.New("no se seleccionaron conjuntos de caracteres válidos")
	}

	password := make([]byte, cfg.length)
	for i := range password {
		randIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charPool))))
		if err != nil {
			return "", err
		}
		password[i] = charPool[randIndex.Int64()]
	}

	return string(password), nil
}

// Construye el conjunto de caracteres según la configuración
func buildCharPool(cfg config) string {
	charPool := ""

	if cfg.includeUpper {
		charPool += upperSet
	}
	if cfg.includeLower {
		charPool += lowerSet
	}
	if cfg.includeNumbers {
		charPool += numberSet
	}
	if cfg.includeSpecials {
		charPool += specialSet
	}

	return charPool
}

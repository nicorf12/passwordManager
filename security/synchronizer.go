package security

import (
	"fmt"
)

// Synchronizer encripta una cadena 3 veces usando claves específicas
func Synchronizer(data string, key1, key2, key3 string) (string, error) {
	enc1, err := EncryptAES([]byte(data), key1)
	if err != nil {
		return "", fmt.Errorf("error en primera encriptación: %v", err)
	}

	enc2, err := EncryptAES([]byte(enc1), key2)
	if err != nil {
		return "", fmt.Errorf("error en segunda encriptación: %v", err)
	}

	enc3, err := EncryptAES([]byte(enc2), key3)
	if err != nil {
		return "", fmt.Errorf("error en tercera encriptación: %v", err)
	}

	return enc3, nil
}

// Desynchronizer desencripta una cadena 3 veces usando claves específicas
func Desynchronizer(data string, key1, key2, key3 string) (string, error) {
	dec1, err := DecryptAES(data, key3)
	if err != nil {
		return "", fmt.Errorf("error en primera desencriptación: %v", err)
	}

	dec2, err := DecryptAES(dec1, key2)
	if err != nil {
		return "", fmt.Errorf("error en segunda desencriptación: %v", err)
	}

	dec3, err := DecryptAES(dec2, key1)
	if err != nil {
		return "", fmt.Errorf("error en tercera desencriptación: %v", err)
	}

	return dec3, nil
}

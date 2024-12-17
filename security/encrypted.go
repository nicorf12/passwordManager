package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	_ "encoding/hex"
	"fmt"
	"golang.org/x/crypto/chacha20poly1305"
	"io"
)

// Genera una clave de 32 bytes usando la contraseña del usuario
func generateKeyFromPassword(password string) ([]byte, error) {
	hash := sha256.New()
	_, err := hash.Write([]byte(password))
	if err != nil {
		return nil, fmt.Errorf("error generando el hash: %v", err)
	}
	return hash.Sum(nil), nil
}

// EncryptAES encripta datos usando AES
func EncryptAES(data []byte, password string) (string, error) {
	key, err := generateKeyFromPassword(password)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(data))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], data)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptAES desencripta datos usando AES
func DecryptAES(encoded string, password string) (string, error) {
	key, err := generateKeyFromPassword(password)
	if err != nil {
		return "", err
	}

	ciphertext, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext), nil
}

// EncryptXChaCha20Poly1305 encripta datos usando XChaCha20-Poly1305
func EncryptXChaCha20Poly1305(data []byte, key string) (string, error) {
	keyBytes, _ := generateKeyFromPassword(key)
	aead, err := chacha20poly1305.NewX(keyBytes)
	if err != nil {
		return "", fmt.Errorf("error creando cifrador: %v", err)
	}

	nonce := make([]byte, chacha20poly1305.NonceSizeX)
	if _, err := rand.Read(nonce); err != nil {
		return "", fmt.Errorf("error generando nonce: %v", err)
	}

	ciphertext := aead.Seal(nil, nonce, data, nil)
	return base64.StdEncoding.EncodeToString(append(nonce, ciphertext...)), nil
}

// DecryptXChaCha20Poly1305 desencripta datos usando XChaCha20-Poly1305
func DecryptXChaCha20Poly1305(encoded string, key string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", fmt.Errorf("error decodificando: %v", err)
	}

	keyBytes, _ := generateKeyFromPassword(key)

	aead, err := chacha20poly1305.NewX(keyBytes)
	if err != nil {
		return "", fmt.Errorf("error creando cifrador: %v", err)
	}

	nonce := ciphertext[:chacha20poly1305.NonceSizeX]
	ciphertext = ciphertext[chacha20poly1305.NonceSizeX:]

	plaintext, err := aead.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("error desencriptando: %v", err)
	}

	return string(plaintext), nil
}

// EncryptDES encripta usando DES
func EncryptDES(data []byte, key string) (string, error) {
	keyBytes, _ := generateKeyFromPassword(key)
	if len(keyBytes) < des.BlockSize {
		return "", fmt.Errorf("clave demasiado corta")
	}
	
	keyBytes = keyBytes[:des.BlockSize]

	block, err := des.NewCipher(keyBytes)
	if err != nil {
		return "", fmt.Errorf("error creando cifrador DES: %v", err)
	}

	ciphertext := make([]byte, des.BlockSize+len(data))
	iv := ciphertext[:des.BlockSize]
	if _, err := rand.Read(iv); err != nil {
		return "", fmt.Errorf("error generando IV: %v", err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[des.BlockSize:], data)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptDES desencripta usando DES
func DecryptDES(encoded string, key string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", fmt.Errorf("error decodificando: %v", err)
	}

	keyBytes, _ := generateKeyFromPassword(key)
	if len(keyBytes) < des.BlockSize {
		return "", fmt.Errorf("clave demasiado corta")
	}

	keyBytes = keyBytes[:des.BlockSize]

	block, err := des.NewCipher(keyBytes)
	if err != nil {
		return "", fmt.Errorf("error creando cifrador DES: %v", err)
	}

	iv := ciphertext[:des.BlockSize]
	ciphertext = ciphertext[des.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)
	return string(ciphertext), nil
}

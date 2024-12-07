package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	_ "encoding/hex"
	"fmt"
	"io"
)

// Genera una clave de 32 bytes usando la contraseña del usuario
func GenerateKeyFromPassword(password string) ([]byte, error) {
	// Utilizamos SHA-256 para generar un hash de 32 bytes
	hash := sha256.New()
	_, err := hash.Write([]byte(password))
	if err != nil {
		return nil, fmt.Errorf("error generando el hash: %v", err)
	}
	return hash.Sum(nil), nil
}

func Encrypt(data []byte, password string) (string, error) {
	// Generamos la clave de 32 bytes a partir de la contraseña
	key, err := GenerateKeyFromPassword(password)
	if err != nil {
		return "", err
	}

	// Usamos AES con la clave generada
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Creamos un vector de inicialización aleatorio
	ciphertext := make([]byte, aes.BlockSize+len(data))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	// Encriptamos los datos
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], data)

	// Devolvemos el texto cifrado en base64
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func Decrypt(encoded string, password string) (string, error) {
	// Generamos la clave de 32 bytes a partir de la contraseña
	key, err := GenerateKeyFromPassword(password)
	if err != nil {
		return "", err
	}

	// Decodificamos el texto cifrado en base64
	ciphertext, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", err
	}

	// Usamos AES con la clave generada
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Extraemos el IV (primer bloque) y los datos cifrados
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	// Desencriptamos los datos
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext), nil
}

package security

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

type SessionData struct {
	UserID         int64     `json:"user_id"`
	UserMail       string    `json:"user_mail"`
	HashedPassword string    `json:"hashed_password"`
	LastAccessed   time.Time `json:"last_accessed"`
}

func SaveSession(data SessionData) error {
	data.LastAccessed = time.Now()

	err := os.MkdirAll("tmp", os.ModePerm)
	if err != nil {
		return err
	}

	file, err := os.Create("tmp/session")
	if err != nil {
		return err
	}
	defer file.Close()

	sessionJSON, err := json.Marshal(data)
	if err != nil {
		return err
	}

	encryptedData, err := EncryptAES(sessionJSON, data.HashedPassword)
	if err != nil {
		return err
	}

	_, err = file.Write([]byte(encryptedData + data.HashedPassword))
	return err
}

func LoadSession() (*SessionData, error) {
	file, err := os.Open("tmp/session")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	encryptedData, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	// Extrae la contraseña hash que está al final del archivo (últimos 44 caracteres)
	hashedPasswordLength := 44
	encryptedSession := encryptedData[:len(encryptedData)-hashedPasswordLength]
	hashedPassword := string(encryptedData[len(encryptedData)-hashedPasswordLength:])

	decryptedData, err := DecryptAES(string(encryptedSession), hashedPassword)
	if err != nil {
		return nil, err
	}

	var data SessionData
	if err := json.Unmarshal([]byte(decryptedData), &data); err != nil {
		return nil, err
	}

	if time.Since(data.LastAccessed) > 2*24*time.Hour {
		err := ClearSession()
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("la sesión ha expirado, el archivo de sesión ha sido eliminado")
	}

	data.LastAccessed = time.Now()

	err = SaveSession(data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func OnLoginSuccess(userID int64, userMail string, hashedPassword string) error {
	sessionData := SessionData{
		UserID:         userID,
		UserMail:       userMail,
		HashedPassword: hashedPassword,
	}
	err := SaveSession(sessionData)
	if err != nil {
		return err
	}
	return nil
}

func ClearSession() error {
	return os.Remove("tmp/session")
}

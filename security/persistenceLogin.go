package security

import (
	"encoding/json"
	"fmt"
	"os"
)

type SessionData struct {
	UserID         int64  `json:"user_id"`
	UserMail       string `json:"user_mail"`
	HashedPassword string `json:"hashed_password"`
}

func SaveSession(data SessionData) error {
	err := os.MkdirAll("tmp", os.ModePerm)
	if err != nil {
		return err
	}

	file, err := os.Create("tmp/session.json")
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(data)
}

func LoadSession() (*SessionData, error) {
	file, err := os.Open("tmp/session.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var data SessionData
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		return nil, err
	}

	fmt.Println(data)

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
	return os.Remove("tmp/session.json")
}

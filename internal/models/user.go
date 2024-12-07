package models

import (
	"errors"
	"regexp"
)

// User representa un usuario como un TDA
type User struct {
	id       int64
	email    string
	password string
}

// NewUser crea un nuevo usuario y realiza validaciones iniciales
func NewUser(email, password string) (*User, error) {
	if !isValidEmail(email) {
		return nil, errors.New("Invalid email")
	}
	/*
		if len(password) < 8 {
			return nil, errors.New("Short password")
		}
	*/
	return &User{email: email, password: password}, nil
}

// GetID devuelve el ID del usuario
func (u *User) GetID() int64 {
	return u.id
}

// GetEmail devuelve el email del usuario
func (u *User) GetEmail() string {
	return u.email
}

// SetEmail actualiza el email del usuario después de validarlo
func (u *User) SetEmail(email string) error {
	if !isValidEmail(email) {
		return errors.New("email inválido")
	}
	u.email = email
	return nil
}

// GetPassword devuelve la contraseña cifrada del usuario
func (u *User) GetPassword() string {
	return u.password
}

// SetPassword permite cambiar la contraseña después de validarla
func (u *User) SetPassword(password string) error {
	if len(password) < 8 {
		return errors.New("la contraseña debe tener al menos 8 caracteres")
	}
	u.password = password
	return nil
}

// EncryptPassword cifra la contraseña del usuario utilizando una función hash externa
func (u *User) EncryptPassword(hashFunc func(string) string) {
	u.password = hashFunc(u.password)
}

// Helper para validar emails
func isValidEmail(email string) bool {
	regex := `^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`
	match, _ := regexp.MatchString(regex, email)
	return match
}

package controllers

import (
	"errors"
	"password_manager/internal/models"
	"password_manager/security"
)

// ControllerUser es el TDA que gestiona las operaciones relacionadas con usuarios.
type ControllerUser struct {
	currentUser  *models.User  // Usuario actual logueado
	dbController *DBController // Controlador de la base de datos
}

// NewControllerUser crea e inicializa una instancia de ControllerUser.
func NewControllerUser(dbC *DBController) *ControllerUser {
	return &ControllerUser{
		dbController: dbC,
		currentUser:  nil,
	}
}

// Login valida las credenciales y establece el usuario actual si son correctas.
func (c *ControllerUser) Login(email, password string) error {
	storedPassword, salt, err := c.dbController.GetUserByEmail(email)
	if err != nil {
		return err
	}

	if storedPassword != security.GenerateHash(password, salt) {
		return errors.New("Incorrect credentials")
	}

	c.currentUser, err = models.NewUser(email, password)
	if err != nil {
		return err
	}
	return nil
}

// IsLoggedIn verifica si hay un usuario actualmente logueado.
func (c *ControllerUser) IsLoggedIn() bool {
	return c.currentUser != nil
}

// Logout cierra sesión, limpiando la información del usuario actual.
func (c *ControllerUser) Logout() {
	c.currentUser = nil
}

// GetCurrentUser devuelve el correo del usuario actual, si está logueado.
func (c *ControllerUser) GetCurrentUserEmail() string {
	if c.currentUser == nil {
		return ""
	}
	return c.currentUser.GetEmail()
}

func (c *ControllerUser) GetCurrentUserId() int64 {
	if c.currentUser == nil {
		return 0
	}
	return c.currentUser.GetID()
}

func (c *ControllerUser) GetCurrentUserPassword() string {
	if c.currentUser == nil {
		return ""
	}
	return c.currentUser.GetPassword()
}

func (c *ControllerUser) GenerateNewPasswordSafe(length int, useUpper, useLower, useNumbers, useSpecials bool) (string, error) {
	password, err := security.GenerateSecurePassword(length,
		useUpper,
		useLower,
		useNumbers,
		useSpecials)
	if err != nil {
		return "", err
	}
	return password, nil
}

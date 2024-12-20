package controllers

import (
	"errors"
	"fmt"
	"log"
	"password_manager/internal/models"
	"password_manager/security"
	"strconv"
	"strings"
)

// ControllerUser es el TDA que gestiona las operaciones relacionadas con usuarios.
type ControllerUser struct {
	currentUser  *models.User  // Usuario actual logueado
	dbController *DBController // Controlador de la base de datos
	config       *security.Config
}

// NewControllerUser crea e inicializa una instancia de ControllerUser.
func NewControllerUser(config *security.Config, dbC *DBController) *ControllerUser {
	return &ControllerUser{
		dbController: dbC,
		currentUser:  nil,
		config:       config,
	}
}

// NewControllerUser crea e inicializa una instancia de ControllerUser.
func NewControllerUserWithSession(config *security.Config, dbC *DBController, id int64, email string, password string) (*ControllerUser, error) {
	storedEmail, storedPassword, _, err := dbC.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	if storedEmail != email {
		return nil, errors.New("invalid email")
	}

	if storedPassword != password {
		return nil, errors.New("invalid password")
	}

	user, err := models.NewUser(email, password)
	if err != nil {
		return nil, fmt.Errorf("error creando el usuario: %v", err)
	}

	user.SetID(id)
	return &ControllerUser{
		dbController: dbC,
		currentUser:  user,
		config:       config,
	}, nil
}

// Login valida las credenciales y establece el usuario actual si son correctas.
func (c *ControllerUser) Login(email, password string) error {
	id, storedPassword, salt, err := c.dbController.GetUserByEmail(email)
	if err != nil {
		return err
	}

	passwordHashed := security.GenerateHash(password, salt)
	if storedPassword != passwordHashed {
		return errors.New("incorrect credentials")
	}

	c.currentUser, err = models.NewUser(email, passwordHashed)
	if err != nil {
		return err
	}
	c.currentUser.SetID(id)

	security.OnLoginSuccess(id, email, passwordHashed)

	return nil
}

// IsLoggedIn verifica si hay un usuario actualmente logueado.
func (c *ControllerUser) IsLoggedIn() bool {
	return c.currentUser != nil
}

// Logout cierra sesión, limpiando la información del usuario actual.
func (c *ControllerUser) Logout() {
	c.currentUser = nil
	security.ClearSession()
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

// GetPasswordSecurityLevel calcula el nivel de seguridad de una contraseña del 1 al 100
func (c *ControllerUser) GetPasswordSecurityLevel(password string) float64 {
	return security.CalculatePasswordSecurity(password)
}

func (c *ControllerUser) GetConfig() (string, string) {
	return c.config.Lang, c.config.Theme
}

func (c *ControllerUser) SetConfig(lang, theme string) {
	c.config.Lang = lang
	c.config.Theme = theme

	err := security.SaveConfig(c.config)
	if err != nil {
		return
	}
}

func (c *ControllerUser) SomeoneLoggedIn() bool {
	if c.currentUser == nil {
		return false
	}
	return true
}

func (c *ControllerUser) EncryptToExport(data, answer0, answer1, answer2 string) (string, error) {
	enc, err := security.Synchronizer("successful\n"+data, answer0, answer1, answer2)
	if err != nil {
		return "", err
	}
	return enc, nil
}

func (c *ControllerUser) DecryptToImport(data, answer0, answer1, answer2 string) (string, error) {
	enc, err := security.Desynchronizer(data, answer0, answer1, answer2)
	if err != nil {
		return "", err
	}

	lines := strings.Split(enc, "\n")
	if len(lines) < 2 || lines[0] != "successful" {
		return "", fmt.Errorf("incorrect keys")
	}

	dataUser := strings.Split(lines[1], ";")
	if len(dataUser) != 3 {
		return "", fmt.Errorf("insufficient user data")
	}
	userID, err := c.dbController.EnterUserToImport(dataUser[0], dataUser[1], dataUser[2])
	if err != nil {
		return "", err
	}

	if len(lines) < 3 {
		return "", nil
	}

	for i, line := range lines[2:] {
		datos := strings.Split(line, ";")
		if len(datos) != 8 {
			log.Printf("Línea %d inválida: %s\n", i+2, line)
			continue
		}

		encryptedID, err := strconv.ParseInt(datos[0], 10, 64)
		if err != nil {
			log.Printf("Error al parsear encryptedID en línea %d: %v\n", i+2, err)
			continue
		}

		isFavorite, err := strconv.Atoi(datos[7])
		if err != nil {
			log.Printf("Error al parsear isFavorite en línea %d: %v\n", i+2, err)
			continue
		}

		err = c.dbController.EnterPasswordToImport(
			0, userID, encryptedID,
			datos[1], datos[2], datos[3], datos[4], datos[5], datos[6], isFavorite)
		if err != nil {
			log.Printf("Error al insertar línea %d: %v\n", i+2, err)
		}
	}

	return enc, nil
}

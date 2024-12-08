package controllers

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	_ "modernc.org/sqlite" // Importa el driver de SQLite
	"password_manager/internal/models"
	"password_manager/security"
)

// Estructura para el controlador de la base de datos
type DBController struct {
	DB *sql.DB
}

// Inicializa la base de datos y crea las tablas necesarias
func NewDBController() (*DBController, error) {
	db, err := sql.Open("sqlite", "db/contraseñas.db")
	if err != nil {
		return nil, fmt.Errorf("Error connecting to database: %v", err)
	}

	// Crear las tablas si no existen
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			email TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL,
			salt TEXT NOT NULL
		);
		
		CREATE TABLE IF NOT EXISTS passwords (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			label TEXT NOT NULL,
			password TEXT NOT NULL,
			FOREIGN KEY(user_id) REFERENCES users(id)
		);
	`)
	if err != nil {
		return nil, fmt.Errorf("Error creating tables: %v", err)
	}

	return &DBController{DB: db}, nil
}

// Cierra la conexión con la base de datos
func (controller *DBController) Close() error {
	return controller.DB.Close()
}

// Inserta un nuevo usuario con su hash y salt en la base de datos
func (controller *DBController) InsertUser(email, password string) (int64, error) {
	_, err := models.NewUser(email, password)
	if err != nil {
		return 0, err
	}

	salt, err := security.GenerateSalt()
	if err != nil {
		return 0, fmt.Errorf("Error generating the salt: %v", err)
	}
	hash := security.GenerateHash(password, salt)

	fmt.Println("Inserting...")
	result, err := controller.DB.Exec("INSERT INTO users (email, password, salt) VALUES (?, ?, ?)", email, hash, base64.StdEncoding.EncodeToString(salt))
	if err != nil {
		return 0, fmt.Errorf("Error inserting user: %v", err)
	}
	fmt.Println("Inserted successfully")
	return result.LastInsertId()
}

// Inserta una nueva contraseña para un usuario
func (controller *DBController) InsertPassword(userID int64, label, password string, userPassword string) (int64, error) {
	if len(password) < 8 {
		return 0, fmt.Errorf("Password must be at least 8 characters")
	}

	encryptedPassword, err := security.Encrypt([]byte(password), userPassword)
	if err != nil {
		return 0, fmt.Errorf("Error encrypting password: %v", err)
	}

	fmt.Println("Saving password")
	result, err := controller.DB.Exec("INSERT INTO passwords (user_id, label, password) VALUES (?, ?, ?)", userID, label, encryptedPassword)
	if err != nil {
		return 0, fmt.Errorf("Error inserting password: %v", err)
	}

	fmt.Println("Password saved")
	passwordID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("Error getting password id: %v", err)
	}
	return passwordID, nil
}

// Obtiene un usuario por su ID
func (controller *DBController) GetUserByID(userID int64) (string, string, error) {
	row := controller.DB.QueryRow("SELECT email, password FROM users WHERE id = ?", userID)
	var email, password string
	if err := row.Scan(&email, &password); err != nil {
		return "", "", fmt.Errorf("Error getting user: %v", err)
	}
	return email, password, nil
}

// Obtiene todas las contraseñas de un usuario por su ID
func (controller *DBController) GetPasswordsByUserID(userID int64, userPassword string) ([]map[string]string, error) {
	rows, err := controller.DB.Query("SELECT id, label, password FROM passwords WHERE user_id = ?", userID)
	if err != nil {
		return nil, fmt.Errorf("Error getting passwords: %v", err)
	}
	defer rows.Close()

	var passwords []map[string]string
	for rows.Next() {
		var id int64
		var label, password string
		if err := rows.Scan(&id, &label, &password); err != nil {
			return nil, fmt.Errorf("Error scanning row: %v", err)
		}
		decryptPassword, _ := security.Decrypt(password, userPassword)
		passwords = append(passwords, map[string]string{
			"id":       fmt.Sprintf("%d", id),
			"label":    label,
			"password": decryptPassword,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Error iterating over rows: %v", err)
	}
	return passwords, nil
}

// GetUserByEmail devuelve los datos de un usuario dado su email
func (controller *DBController) GetUserByEmail(email string) (int64, string, []byte, error) {
	row := controller.DB.QueryRow("SELECT id, password, salt FROM users WHERE email = ?", email)
	var idUser int64
	var hash, saltBase64 string

	if err := row.Scan(&idUser, &hash, &saltBase64); err != nil {
		return 0, "", nil, fmt.Errorf("Error getting user data: %v", err)
	}

	salt, err := base64.StdEncoding.DecodeString(saltBase64)
	if err != nil {
		return 0, "", nil, fmt.Errorf("Error decoding the salt: %v", err)
	}
	return idUser, hash, salt, nil
}

// Elimina un usuario y todas sus contraseñas
func (controller *DBController) DeleteUser(userID int64) error {
	_, err := controller.DB.Exec("DELETE FROM passwords WHERE user_id = ?", userID)
	if err != nil {
		return fmt.Errorf("Error deleting user passwords: %v", err)
	}

	_, err = controller.DB.Exec("DELETE FROM users WHERE id = ?", userID)
	if err != nil {
		return fmt.Errorf("Error deleting user: %v", err)
	}

	return nil
}

// Elimina una contraseña por su ID
func (controller *DBController) DeletePassword(passwordID int64) error {
	_, err := controller.DB.Exec("DELETE FROM passwords WHERE id = ?", passwordID)
	if err != nil {
		return fmt.Errorf("Error deleting user: %v", err)
	}
	return nil
}

// Actualiza la contraseña en la base de datos
func (controller *DBController) EditPassword(passwordID int64, newPassword string, userPassword string) error {
	if len(newPassword) < 8 {
		return fmt.Errorf("Password must be at least 8 characters")
	}

	encryptedPassword, err := security.Encrypt([]byte(newPassword), userPassword)
	if err != nil {
		return fmt.Errorf("Error encrypting password: %v", err)
	}

	_, err = controller.DB.Exec("UPDATE passwords SET password = ? WHERE id = ?", encryptedPassword, passwordID)
	if err != nil {
		return fmt.Errorf("Error updating password: %v", err)
	}
	return nil
}

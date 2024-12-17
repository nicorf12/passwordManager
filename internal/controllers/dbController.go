package controllers

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	_ "modernc.org/sqlite" // Importa el driver de SQLite
	"password_manager/internal/models"
	"password_manager/security"
	"strconv"
	"strings"
)

var encryptionMethods = map[int64]struct {
	Encrypt func([]byte, string) (string, error)
	Decrypt func(string, string) (string, error)
}{
	1: { // AES-256
		Encrypt: security.EncryptAES,
		Decrypt: security.DecryptAES,
	},
	2: { // XChaCha20-Poly1305
		Encrypt: security.EncryptXChaCha20Poly1305,
		Decrypt: security.DecryptXChaCha20Poly1305,
	},
	3: { // DES
		Encrypt: security.EncryptDES,
		Decrypt: security.DecryptDES,
	},
}

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
			folder_id INTEGER NOT NULL,
			label TEXT NOT NULL,
			name TEXT NOT NULL,
			password TEXT NOT NULL,
			website TEXT,
			note TEXT,
			encrypted_id INTEGER NOT NULL,
			last_update DATETIME DEFAULT CURRENT_TIMESTAMP,
			is_favorite INTEGER DEFAULT 0,
			FOREIGN KEY(user_id) REFERENCES users(id),
			FOREIGN KEY(folder_id) REFERENCES folders(id),
		    FOREIGN KEY(encrypted_id) REFERENCES encrypted(id)
		);
	
		CREATE TABLE IF NOT EXISTS folders (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL
		);

		CREATE TABLE IF NOT EXISTS encrypted (
		    id INTEGER PRIMARY KEY AUTOINCREMENT,
		    name TEXT NOT NULL
		);

		INSERT INTO encrypted (name)
		SELECT 'AES-256' WHERE NOT EXISTS (SELECT 1 FROM encrypted WHERE name = 'AES-256');
		INSERT INTO encrypted (name)
		SELECT 'XChaCha20-Poly1305' WHERE NOT EXISTS (SELECT 1 FROM encrypted WHERE name = 'XChaCha20-Poly1305');
		INSERT INTO encrypted (name)
		SELECT 'DES' WHERE NOT EXISTS (SELECT 1 FROM encrypted WHERE name = 'DES');
	`)

	if err != nil {
		return nil, fmt.Errorf("Error creating tables: %v", err)
	}

	return &DBController{DB: db}, nil
}

// Cierra la conexión con la base de datos
func (db *DBController) Close() error {
	return db.DB.Close()
}

// Inserta un nuevo usuario con su hash y salt en la base de datos
func (db *DBController) InsertUser(email, password string) (int64, error) {
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
	result, err := db.DB.Exec("INSERT INTO users (email, password, salt) VALUES (?, ?, ?)", email, hash, base64.StdEncoding.EncodeToString(salt))
	if err != nil {
		return 0, fmt.Errorf("Error inserting user: %v", err)
	}
	fmt.Println("Inserted successfully")
	return result.LastInsertId()
}

// Inserta una nueva contraseña para un usuario
func (db *DBController) InsertPassword(userID int64, folderID int64, label, name, password, website, note string, encryptedID int64, userPassword string) (int64, error) {
	if len(password) < 8 {
		return 0, fmt.Errorf("Password must be at least 8 characters")
	}

	encryptionMethod, exists := encryptionMethods[encryptedID]
	if !exists {
		return 0, fmt.Errorf("Unsupported encryption method")
	}

	encryptedPassword, err := encryptionMethod.Encrypt([]byte(password), userPassword)
	if err != nil {
		return 0, fmt.Errorf("Error encrypting password: %v", err)
	}

	fmt.Println("Saving password")
	result, err := db.DB.Exec("INSERT INTO passwords (user_id, folder_id,label, name, password, website, note, encrypted_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", userID, folderID, label, name, encryptedPassword, website, note, encryptedID)
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

// Crear una nueva carpeta
func (db *DBController) InsertFolder(name string) (int64, error) {
	result, err := db.DB.Exec("INSERT INTO folders (name) VALUES (?)", name)
	if err != nil {
		return 0, fmt.Errorf("error al crear la carpeta: %v", err)
	}
	folderID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error al obtener el ID de la carpeta creada: %v", err)
	}
	return folderID, nil
}

// Obtiene un usuario por su ID
func (db *DBController) GetUserByID(userID int64) (string, string, error) {
	row := db.DB.QueryRow("SELECT email, password FROM users WHERE id = ?", userID)
	var email, password string
	if err := row.Scan(&email, &password); err != nil {
		return "", "", fmt.Errorf("error getting user: %v", err)
	}
	return email, password, nil
}

// Obtiene todas las contraseñas de un usuario por su ID
func (db *DBController) GetPasswordsByUserID(userID int64, userPassword string) ([]map[string]string, error) {
	rows, err := db.DB.Query("SELECT id,folder_id,label,name,password,website,note,encrypted_id,last_update,is_favorite FROM passwords WHERE user_id = ?", userID)
	if err != nil {
		return nil, fmt.Errorf("error getting passwords: %v", err)
	}
	defer rows.Close()

	var passwords []map[string]string
	for rows.Next() {
		var id, folderId, encryptedId int64
		var label, name, password, website, note, lastUpdate string
		var isFavorite int
		if err := rows.Scan(&id, &folderId, &label, &name, &password, &website, &note, &encryptedId, &lastUpdate, &isFavorite); err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		encryptionMethod, exists := encryptionMethods[encryptedId]
		if !exists {
			return nil, fmt.Errorf("unsupported encryption method")
		}

		decryptedPassword, err := encryptionMethod.Decrypt(password, userPassword)
		if err != nil {
			return nil, fmt.Errorf("error decrypting password: %v", err)
		}
		passwords = append(passwords, map[string]string{
			"id":           fmt.Sprintf("%d", id),
			"folder_id":    fmt.Sprintf("%d", folderId),
			"label":        label,
			"name":         name,
			"password":     decryptedPassword,
			"website":      website,
			"note":         note,
			"encrypted_id": strconv.FormatInt(encryptedId, 10),
			"last_update":  lastUpdate,
			"is_favorite":  fmt.Sprintf("%d", isFavorite),
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %v", err)
	}
	return passwords, nil
}

// GetUserByEmail devuelve los datos de un usuario dado su email
func (db *DBController) GetUserByEmail(email string) (int64, string, []byte, error) {
	row := db.DB.QueryRow("SELECT id, password, salt FROM users WHERE email = ?", email)
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

// GetAllFolders devuelve todas las carpetas como un mapa de nombre a ID.
func (db *DBController) GetAllFolders() (map[string]int64, error) {
	rows, err := db.DB.Query("SELECT id, name FROM folders")
	if err != nil {
		return nil, fmt.Errorf("error fetching folders: %v", err)
	}
	defer rows.Close()

	folders := make(map[string]int64)
	for rows.Next() {
		var id int64
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			return nil, fmt.Errorf("error scanning folder: %v", err)
		}
		folders[name] = id
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %v", err)
	}

	return folders, nil
}

func (db *DBController) GetAllEncrypted() (map[string]int64, error) {
	rows, err := db.DB.Query("SELECT id, name FROM encrypted")
	if err != nil {
		return nil, fmt.Errorf("error fetching encrypted types: %v", err)
	}
	defer rows.Close()

	encrypted := make(map[string]int64)

	for rows.Next() {
		var id int64
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			return nil, fmt.Errorf("error scanning encrypted record: %v", err)
		}
		encrypted[name] = id
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %v", err)
	}

	return encrypted, nil
}

// Obtiene todas las contraseñas de un usuario para una carpeta específica
func (db *DBController) GetPasswordsByFolderAndUserID(userID int64, folderID int64, userPassword string) ([]map[string]string, error) {
	rows, err := db.DB.Query("SELECT id, folder_id, label, name, password, website, note, encrypted_id, last_update, is_favorite FROM passwords WHERE user_id = ? AND folder_id = ?", userID, folderID)
	if err != nil {
		return nil, fmt.Errorf("error getting passwords by folder and user: %v", err)
	}
	defer rows.Close()

	var passwords []map[string]string
	for rows.Next() {
		var id, folderId, encryptedId int64
		var label, name, password, website, note, lastUpdate string
		var isFavorite int
		if err := rows.Scan(&id, &folderId, &label, &name, &password, &website, &note, &encryptedId, &lastUpdate, &isFavorite); err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		encryptionMethod, exists := encryptionMethods[encryptedId]
		if !exists {
			return nil, fmt.Errorf("unsupported encryption method")
		}

		decryptedPassword, _ := encryptionMethod.Decrypt(password, userPassword)
		passwords = append(passwords, map[string]string{
			"id":           fmt.Sprintf("%d", id),
			"folder_id":    fmt.Sprintf("%d", folderId),
			"label":        label,
			"name":         name,
			"password":     decryptedPassword,
			"website":      website,
			"note":         note,
			"encrypted_id": fmt.Sprintf("%d", encryptedId),
			"last_update":  lastUpdate,
			"is_favorite":  fmt.Sprintf("%d", isFavorite),
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %v", err)
	}
	return passwords, nil
}

// Obtiene todas las contraseñas marcadas como favoritas de un usuario por su ID
func (db *DBController) GetPasswordsByFavoriteAndUserID(userID int64, userPassword string) ([]map[string]string, error) {
	rows, err := db.DB.Query("SELECT id, folder_id, label, name, password, website, note, encrypted_id, last_update, is_favorite FROM passwords WHERE user_id = ? AND is_favorite = 1", userID)
	if err != nil {
		return nil, fmt.Errorf("error getting favorite passwords by user: %v", err)
	}
	defer rows.Close()

	var passwords []map[string]string
	for rows.Next() {
		var id, folderId, encryptedId int64
		var label, name, password, website, note, lastUpdate string
		var isFavorite int
		if err := rows.Scan(&id, &folderId, &label, &name, &password, &website, &note, &encryptedId, &lastUpdate, &isFavorite); err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		encryptionMethod, exists := encryptionMethods[encryptedId]
		if !exists {
			return nil, fmt.Errorf("unsupported encryption method")
		}

		decryptedPassword, _ := encryptionMethod.Decrypt(password, userPassword)
		passwords = append(passwords, map[string]string{
			"id":          fmt.Sprintf("%d", id),
			"folder_id":   fmt.Sprintf("%d", folderId),
			"label":       label,
			"name":        name,
			"password":    decryptedPassword,
			"website":     website,
			"note":        note,
			"encrypted":   fmt.Sprintf("%d", encryptedId),
			"last_update": lastUpdate,
			"is_favorite": fmt.Sprintf("%d", isFavorite),
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %v", err)
	}
	return passwords, nil
}

// Elimina un usuario y todas sus contraseñas
func (db *DBController) DeleteUser(userID int64) error {
	_, err := db.DB.Exec("DELETE FROM passwords WHERE user_id = ?", userID)
	if err != nil {
		return fmt.Errorf("Error deleting user passwords: %v", err)
	}

	_, err = db.DB.Exec("DELETE FROM users WHERE id = ?", userID)
	if err != nil {
		return fmt.Errorf("Error deleting user: %v", err)
	}

	return nil
}

// Elimina una contraseña por su ID
func (db *DBController) DeletePassword(passwordID int64) error {
	_, err := db.DB.Exec("DELETE FROM passwords WHERE id = ?", passwordID)
	if err != nil {
		return fmt.Errorf("Error deleting user: %v", err)
	}
	return nil
}

// Borrar una carpeta
func (db *DBController) DeleteFolder(folderID int64) error {
	result, err := db.DB.Exec("DELETE FROM folders WHERE id = ?", folderID)
	if err != nil {
		return fmt.Errorf("Error al borrar la carpeta: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Error al verificar filas afectadas: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("No se encontró una carpeta con ID %d", folderID)
	}
	return nil
}

// Actualiza la contraseña en la base de datos
func (db *DBController) EditPassword(passwordID int64, updates map[string]interface{}, userPassword string) error {
	var setClauses []string
	var args []interface{}

	if newPassword, ok := updates["password"]; ok {
		if len(newPassword.(string)) < 8 {
			return fmt.Errorf("password must be at least 8 characters")
		}

		encryptionMethod, exists := encryptionMethods[updates["encrypted_id"].(int64)]
		if !exists {
			return fmt.Errorf("unsupported encryption method")
		}

		encryptedPassword, err := encryptionMethod.Encrypt([]byte(newPassword.(string)), userPassword)
		if err != nil {
			return fmt.Errorf("error encrypting password: %v", err)
		}

		setClauses = append(setClauses, "password = ?")
		args = append(args, encryptedPassword)
	}

	for column, value := range updates {
		if column != "password" {
			setClauses = append(setClauses, fmt.Sprintf("%s = ?", column))
			args = append(args, value)
		}
	}

	if len(setClauses) == 0 {
		return fmt.Errorf("no fields to update")
	}

	setClause := "SET " + strings.Join(setClauses, ", ")
	sqlQuery := fmt.Sprintf("UPDATE passwords %s WHERE id = ?", setClause)

	args = append(args, passwordID)
	_, err := db.DB.Exec(sqlQuery, args...)
	if err != nil {
		return fmt.Errorf("error updating password: %v", err)
	}

	return nil
}

// Actualizar el nombre de una carpeta
func (db *DBController) EditFolder(folderID int64, newName string) error {
	result, err := db.DB.Exec("UPDATE folders SET name = ? WHERE id = ?", newName, folderID)
	if err != nil {
		return fmt.Errorf("Error al actualizar la carpeta: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Error al verificar filas afectadas: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("No se encontró una carpeta con ID %d", folderID)
	}
	return nil
}

// Cambia el estado de favorito de una contraseña por su ID
func (db *DBController) EditFavoritePassword(passwordID int64) error {
	result, err := db.DB.Exec("UPDATE passwords SET is_favorite = CASE WHEN is_favorite = 1 THEN 0 ELSE 1 END WHERE id = ?", passwordID)
	if err != nil {
		return fmt.Errorf("error updating favorite status: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no password found with the specified ID")
	}

	return nil
}

package test

import (
	"database/sql"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"password_manager/internal/controllers"
	"testing"

	_ "modernc.org/sqlite"
)

func setupDBController(t *testing.T) *controllers.DBController {
	dbController, err := controllers.NewDBController()
	if err != nil {
		t.Fatalf("Error al crear el DBController: %v", err)
	}
	t.Cleanup(func() {
		dbController.Close()
		os.RemoveAll("db")
	})
	return dbController
}

func TestCreateDBController(t *testing.T) {
	dbController := setupDBController(t)

	tableColumns := map[string][]string{
		"users": {
			"id", "email", "password", "salt",
		},
		"passwords": {
			"id", "user_id", "folder_id", "label", "name", "password", "website", "note", "encrypted_id", "last_update", "is_favorite",
		},
		"folders": {
			"id", "name",
		},
		"encrypted": {
			"id", "name",
		},
	}

	for table, expectedColumns := range tableColumns {
		query := `SELECT name FROM sqlite_master WHERE type='table' AND name=?`
		row := dbController.DB.QueryRow(query, table)

		var name string
		err := row.Scan(&name)
		if err != nil {
			t.Errorf("La tabla %s no existe: %v", table, err)
			continue
		}

		columnsQuery := `PRAGMA table_info(` + table + `)`
		rows, err := dbController.DB.Query(columnsQuery)
		if err != nil {
			t.Errorf("Error al consultar las columnas de la tabla %s: %v", table, err)
			continue
		}
		defer rows.Close()

		var columns []string
		for rows.Next() {
			var cid int
			var name, ctype string
			var dfltValue sql.NullString
			var pk, notnull bool
			if err := rows.Scan(&cid, &name, &ctype, &notnull, &dfltValue, &pk); err != nil {
				t.Errorf("Error al leer columnas de la tabla %s: %v", table, err)
			} else {
				columns = append(columns, name)
			}
		}
		if len(columns) == 0 {
			t.Errorf("La tabla %s no tiene columnas", table)
		} else {
			for _, expectedColumn := range expectedColumns {
				found := false
				for _, column := range columns {
					if column == expectedColumn {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("La tabla %s no tiene la columna esperada: %s", table, expectedColumn)
				}
			}
		}
	}
}

func TestInsertUser(t *testing.T) {
	dbController := setupDBController(t)

	email := "testuser@example.com"
	password := "testpassword123"

	userID, err := dbController.InsertUser(email, password)
	if err != nil {
		t.Fatalf("Error al insertar el usuario: %v", err)
	}

	if userID <= 0 {
		t.Fatalf("ID de usuario insertado no es válido: %d", userID)
	}

	userIDInserted, _, _, err := dbController.GetUserByEmail(email)
	if err != nil {
		t.Fatalf("Error al recuperar el usuario por correo: %v", err)
	}
	assert.Equal(t, userIDInserted, userID, "El id del usuario insertado no coincide")

	emailInserted, _, _, err := dbController.GetUserByID(userID)
	if err != nil {
		t.Fatalf("Error al recuperar el usuario por id: %v", err)
	}
	assert.Equal(t, emailInserted, email, "El email del usuario insertado no coincide")
}

func TestInsertPassword(t *testing.T) {
	dbController := setupDBController(t)

	userEmail := "testuser@example.com"
	userPassword := "testpw12345678"
	userID, err := dbController.InsertUser(userEmail, userPassword)
	require.NoError(t, err, "Error al insertar el usuario")

	folderID := int64(0)
	passwordData := map[string]interface{}{
		"label":       "test",
		"name":        "test@gmail.com",
		"password":    "testpassword123",
		"website":     "www.test.com",
		"note":        "test",
		"encryptedID": int64(1),
	}

	_, err = dbController.InsertPassword(
		userID, folderID,
		passwordData["label"].(string),
		passwordData["name"].(string),
		passwordData["password"].(string),
		passwordData["website"].(string),
		passwordData["note"].(string),
		passwordData["encryptedID"].(int64),
		userPassword,
	)
	require.NoError(t, err, "Error al insertar la contraseña")
}

func TestGetPasswordsByUserID(t *testing.T) {
	dbController := setupDBController(t)

	userEmail := "testuser@example.com"
	userPassword := "testpw12345678"
	userID, err := dbController.InsertUser(userEmail, userPassword)

	var folderID int64 = 0
	label := "test"
	name := "test@gmail.com"
	password := "testpassword123"
	website := "www.test.com"
	note := "test"
	var encryptedID int64 = 1
	_, err = dbController.InsertPassword(userID, folderID, label, name, password, website, note, encryptedID, userPassword)
	if err != nil {
		t.Fatalf("Error al insertar la contraseña: %v", err)
	}

	passwords, err := dbController.GetPasswordsByUserID(userID, userPassword)
	if err != nil {
		t.Fatalf("Error al obtener las contraseñas: %v", err)
	}

	if len(passwords) != 1 {
		t.Fatalf("Contraseña no insertada correctamente")
		return
	}

	correcto := true
	if passwords[0]["folder_id"] != fmt.Sprintf("%d", folderID) {
		correcto = false
	}
	if passwords[0]["label"] != label {
		correcto = false
	}
	if passwords[0]["name"] != name {
		correcto = false
	}
	if passwords[0]["password"] != password {
		correcto = false
	}
	if passwords[0]["website"] != website {
		correcto = false
	}
	if passwords[0]["note"] != note {
		correcto = false
	}
	if passwords[0]["encrypted_id"] != fmt.Sprintf("%d", encryptedID) {
		correcto = false
	}

	if !correcto {
		t.Fatalf("Datos no insertados correctamente")
	}
}

func TestGetPasswordsByFolderAndUserID(t *testing.T) {
	dbController := setupDBController(t)

	userEmail := "testuser@example.com"
	userPassword := "testpw12345678"
	userID, err := dbController.InsertUser(userEmail, userPassword)

	folderID, err := dbController.InsertFolder("carpeta test")
	if err != nil {
		t.Fatalf("Error al insertar la carpeta: %v", err)
		return
	}

	label := "test"
	name := "test@gmail.com"
	password := "testpassword123"
	website := "www.test.com"
	note := "test"
	var encryptedID int64 = 1
	_, err = dbController.InsertPassword(userID, folderID, label, name, password, website, note, encryptedID, userPassword)
	if err != nil {
		t.Fatalf("Error al insertar la contraseña: %v", err)
		return
	}

	passwords, err := dbController.GetPasswordsByFolderAndUserID(folderID, userID, userPassword)
	if err != nil {
		t.Fatalf("Error al obtener las contraseñas: %v", err)
		return
	}

	if len(passwords) != 1 {
		t.Fatalf("Contraseña no insertada correctamente")
		return
	}

	correcto := true
	if passwords[0]["folder_id"] != fmt.Sprintf("%d", folderID) {
		correcto = false
	}
	if passwords[0]["label"] != label {
		correcto = false
	}
	if passwords[0]["name"] != name {
		correcto = false
	}
	if passwords[0]["password"] != password {
		correcto = false
	}
	if passwords[0]["website"] != website {
		correcto = false
	}
	if passwords[0]["note"] != note {
		correcto = false
	}
	if passwords[0]["encrypted_id"] != fmt.Sprintf("%d", encryptedID) {
		correcto = false
	}

	if !correcto {
		t.Fatalf("Datos no insertados correctamente")
	}
}

func TestGetPasswordsByFavoriteAndUserID(t *testing.T) {
	dbController := setupDBController(t)

	userEmail := "testuser@example.com"
	userPassword := "testpw12345678"
	userID, err := dbController.InsertUser(userEmail, userPassword)

	var (
		folderID    int64 = 0
		label             = "test"
		name              = "test@gmail.com"
		password          = "testpassword123"
		website           = "www.test.com"
		note              = "test"
		encryptedID int64 = 1
	)

	passwordID, err := dbController.InsertPassword(userID, folderID, label, name, password, website, note, encryptedID, userPassword)
	if err != nil {
		t.Fatalf("Error al insertar la contraseña: %v", err)
		return
	}

	passwords, err := dbController.GetPasswordsByFavoriteAndUserID(userID, userPassword)
	if len(passwords) != 0 {
		t.Fatalf("Contraseña no obtenida correctamente")
		return
	}

	err = dbController.EditFavoritePassword(passwordID)
	if err != nil {
		t.Fatalf("Error en la actualizacion a favorito: %v", err)
		return
	}

	passwords, err = dbController.GetPasswordsByFavoriteAndUserID(userID, userPassword)
	if len(passwords) != 1 {
		t.Fatalf("Contraseña no obtenida correctamente")
		return
	}

	if passwords[0]["is_favorite"] != fmt.Sprintf("%d", 1) {
		t.Fatalf("Error en la actualizacion a favorito")
	}
}

func TestGetFolders(t *testing.T) {
	dbController := setupDBController(t)

	dbController.InsertFolder("carpeta test1")
	dbController.InsertFolder("carpeta test2")
	dbController.InsertFolder("carpeta test3")

	folders, err := dbController.GetAllFolders()
	if err != nil {
		t.Fatalf("Error al obtener carpetas %v", err)
		return
	}

	if len(folders) != 3 {
		t.Fatalf("Error al obtener carpetas %v", err)
		return
	}

	var ok bool
	var id int64
	id, ok = folders["carpeta test1"]
	if !ok || id != 1 {
		t.Fatalf("Carpeta test0 folder not correct")
	}
	id, ok = folders["carpeta test2"]
	if !ok || id != 2 {
		t.Fatalf("Carpeta test0 folder not correct")
	}
	id, ok = folders["carpeta test3"]
	if !ok || id != 3 {
		t.Fatalf("Carpeta test0 folder not correct")
	}

}

func TestGetEncrypted(t *testing.T) {
	dbController := setupDBController(t)

	encrypted, err := dbController.GetAllEncrypted()
	if err != nil {
		t.Fatalf("Error al obtener encriptados %v", err)
		return
	}
	if len(encrypted) != 3 { // actualmente hay 3 tipos de encriptado
		t.Fatalf("Error al obtener encriptados %v", err)
	}
}

func TestDeletePassword(t *testing.T) {
	dbController := setupDBController(t)

	userEmail := "testuser@example.com"
	userPassword := "testpw12345678"
	userID, err := dbController.InsertUser(userEmail, userPassword)
	require.NoError(t, err, "Error al insertar el usuario")

	folderID := int64(0)
	passwordData := map[string]interface{}{
		"label":       "test",
		"name":        "test@gmail.com",
		"password":    "testpassword123",
		"website":     "www.test.com",
		"note":        "test",
		"encryptedID": int64(1),
	}

	passwordID, err := dbController.InsertPassword(
		userID, folderID,
		passwordData["label"].(string),
		passwordData["name"].(string),
		passwordData["password"].(string),
		passwordData["website"].(string),
		passwordData["note"].(string),
		passwordData["encryptedID"].(int64),
		userPassword,
	)
	require.NoError(t, err, "Error al insertar la contraseña")

	err = dbController.DeletePassword(passwordID)
	require.NoError(t, err, "Error al eliminar la contraseña")
}

func TestEditPassword(t *testing.T) {
	dbController := setupDBController(t)

	userEmail := "testuser@example.com"
	userPassword := "testpw12345678"
	userID, err := dbController.InsertUser(userEmail, userPassword)
	require.NoError(t, err, "Error al insertar el usuario")

	folderID := int64(0)
	passwordData := map[string]interface{}{
		"label":       "test",
		"name":        "test@gmail.com",
		"password":    "testpassword123",
		"website":     "www.test.com",
		"note":        "test",
		"encryptedID": int64(1),
	}

	passwordID, err := dbController.InsertPassword(
		userID, folderID,
		passwordData["label"].(string),
		passwordData["name"].(string),
		passwordData["password"].(string),
		passwordData["website"].(string),
		passwordData["note"].(string),
		passwordData["encryptedID"].(int64),
		userPassword,
	)
	require.NoError(t, err, "Error al insertar la contraseña")

	passwordData = map[string]interface{}{
		"label":        "test editado",
		"name":         "test@gmail.com",
		"password":     "testpassword456",
		"website":      "www.test.com",
		"note":         "test editado",
		"encrypted_id": int64(2),
	}
	err = dbController.EditPassword(passwordID, passwordData, userPassword)
	require.NoError(t, err, "Error al editar la contraseña")
}

func TestDeleteFolder(t *testing.T) {
	dbController := setupDBController(t)

	folderID, err := dbController.InsertFolder("carpeta test")
	require.NoError(t, err, "Error al insertar la carpeta")

	err = dbController.DeleteFolder(folderID)
	require.NoError(t, err, "Error al eliminar la carpeta")
}

func TestGetDataToExport(t *testing.T) {
	dbController := setupDBController(t)

	userEmail := "testuser@example.com"
	userPassword := "testpw12345678"
	userID, err := dbController.InsertUser(userEmail, userPassword)
	require.NoError(t, err, "Error al insertar el usuario")

	folderID := int64(0)
	passwordData := map[string]interface{}{
		"label":       "test",
		"name":        "test@gmail.com",
		"password":    "testpassword123",
		"website":     "www.test.com",
		"note":        "test",
		"encryptedID": int64(1),
	}

	_, err = dbController.InsertPassword(
		userID, folderID,
		passwordData["label"].(string),
		passwordData["name"].(string),
		passwordData["password"].(string),
		passwordData["website"].(string),
		passwordData["note"].(string),
		passwordData["encryptedID"].(int64),
		userPassword,
	)
	require.NoError(t, err, "Error al insertar la contraseña")

	export, err := dbController.GetDataToExport(userID)
	require.NoError(t, err, "Error al exportar")
	require.NotEqual(t, export, "")
}

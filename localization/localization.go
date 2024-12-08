package localization

import (
	"encoding/json"
	"fmt"
	"os"
)

// Localizer contiene los textos traducidos para un idioma
type Localizer struct {
	Translations map[string]string
}

// NewLocalizer carga un archivo JSON de idioma
func NewLocalizer(languageFile string) (*Localizer, error) {
	filePath := "localization/" + languageFile + ".json"
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("Error opening language file: %w", err)
	}
	defer file.Close()

	var translations map[string]string
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&translations); err != nil {
		return nil, fmt.Errorf("Error decoding JSON file: %w", err)
	}

	return &Localizer{Translations: translations}, nil
}

// Get obtiene la traducción para una clave dada
func (l *Localizer) Get(key string) string {
	if value, exists := l.Translations[key]; exists {
		return value
	}
	return fmt.Sprintf("{{%s}}", key) // Devuelve la clave como placeholder si no existe
}

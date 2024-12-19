package security

import (
	"encoding/json"
	"os"
)

type Config struct {
	Lang  string `json:"languaje"`
	Theme string `json:"theme"`
}

func SaveConfig(data *Config) error {
	err := os.MkdirAll("tmp", os.ModePerm)
	if err != nil {
		return err
	}

	file, err := os.Create("tmp/config.json")
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(*data)
}

func LoadConfig() (*Config, error) {
	file, err := os.Open("tmp/config.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var data Config
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		return nil, err
	}

	return &data, nil
}

func LoadConfigDefault() *Config {
	return &Config{
		Lang:  "en",
		Theme: "Dark",
	}
}

func ClearConfig() error {
	return os.Remove("tmp/config.json")
}

package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DBurl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read() (Config, error) {
	configPath, err := getConfigFilePath()
	if err != nil {
		return Config{}, fmt.Errorf("error getting config file path: %v", err)
	}

	content, err := os.ReadFile(configPath)
	if err != nil {
		return Config{}, fmt.Errorf("error reading file: %v", err)
	}

	var data Config

	if err := json.Unmarshal(content, &data); err != nil {
		return Config{}, fmt.Errorf("error unmarshalling json: %v", err)
	}

	return data, nil
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("error finding home directory: %v", err)
	}
	return homeDir + "/" + configFileName, nil
}

func (c Config) SetUser(username string) error {
	c.CurrentUserName = username
	write(c)
	return nil
}

func write(cfg Config) error {
	jsonData, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	cfgPath, err := getConfigFilePath()
	if err != nil {
		return fmt.Errorf("error getting config file name during write: %v", err)
	}

	if err := os.WriteFile(cfgPath, jsonData, 0644); err != nil {
		return err
	}

	return nil
}

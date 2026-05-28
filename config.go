package gniphyl

import (
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

//go:embed extensions.json
var embeddedConfig embed.FS

// Config represents the file extension categories configuration.
type Config struct {
	Extensions map[string][]string `json:"extensions"`
}

// PathsConfig represents the user's configured paths.
type PathsConfig struct {
	Paths []string `json:"paths"`
}

const appName = "gniphyl"

// getConfigFolder returns the platform-specific configuration folder path.
func getConfigFolder() (string, error) {
	var configPath string

	switch runtime.GOOS {
	case "windows":
		localAppData := os.Getenv("LOCALAPPDATA")
		if localAppData == "" {
			return "", fmt.Errorf("LOCALAPPDATA environment variable not set")
		}
		configPath = filepath.Join(localAppData, appName)
	case "linux", "darwin":
		home := os.Getenv("HOME")
		if home == "" {
			return "", fmt.Errorf("HOME environment variable not set")
		}
		configPath = filepath.Join(home, ".config", appName)
	default:
		return "", fmt.Errorf("unsupported system: %s", runtime.GOOS)
	}

	if err := os.MkdirAll(configPath, 0755); err != nil {
		return "", err
	}

	return configPath, nil
}

// getConfigFilePath returns the full path to the user's paths config file.
func getConfigFilePath() (string, error) {
	configFolder, err := getConfigFolder()
	if err != nil {
		return "", err
	}
	return filepath.Join(configFolder, "config.json"), nil
}

// migrateOldConfig migrates from the legacy config.toml to config.json.
func migrateOldConfig() error {
	configFolder, err := getConfigFolder()
	if err != nil {
		return err
	}
	oldPath := filepath.Join(configFolder, "config.toml")
	newPath := filepath.Join(configFolder, "config.json")

	if _, err := os.Stat(oldPath); os.IsNotExist(err) {
		return nil // no old config to migrate
	}

	// If new config already exists, just remove the old one
	if _, err := os.Stat(newPath); err == nil {
		return os.Remove(oldPath)
	}

	// Read old, write new, remove old
	data, err := os.ReadFile(oldPath)
	if err != nil {
		return err
	}
	if err := os.WriteFile(newPath, data, 0644); err != nil {
		return err
	}
	return os.Remove(oldPath)
}

// LoadPathsConfig loads the user's configured paths from the configuration file.
func LoadPathsConfig() (*PathsConfig, error) {
	if err := migrateOldConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "warning: config migration failed: %v\n", err)
	}

	configFile, err := getConfigFilePath()
	if err != nil {
		return nil, err
	}

	config := &PathsConfig{Paths: []string{}}

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return config, nil
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, config); err != nil {
		return nil, err
	}

	return config, nil
}

// SavePathsConfig saves the user's configured paths to the configuration file.
func SavePathsConfig(config *PathsConfig) error {
	configFile, err := getConfigFilePath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configFile, data, 0644)
}

// LoadExtensionsConfig loads the file extension categories from the embedded
// extensions.json file.
func LoadExtensionsConfig() (*Config, error) {
	data, err := embeddedConfig.ReadFile("extensions.json")
	if err != nil {
		return nil, err
	}

	config := &Config{}
	if err := json.Unmarshal(data, config); err != nil {
		return nil, err
	}

	return config, nil
}

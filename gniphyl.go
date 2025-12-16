// Package gniphyl provides file organization functionality for categorizing
// files based on their extensions. It can be used both as a CLI tool and as a
// Go package in other applications.

package gniphyl

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

//go:embed config.json
var embeddedConfig embed.FS

// Config represents the file extension categories configuration
type Config struct {
	Extensions map[string][]string `json:"extensions"`
}

// PathsConfig represents the user's configured paths
type PathsConfig struct {
	Paths []string `json:"paths"`
}

const appName = "gniphyl"

// getConfigFolder returns the platform-specific configuration folder path
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

// getConfigFilePath returns the full path to the configuration file
func getConfigFilePath() (string, error) {
	configFolder, err := getConfigFolder()
	if err != nil {
		return "", err
	}
	return filepath.Join(configFolder, "config.toml"), nil
}

// LoadPathsConfig loads the user's configured paths from the configuration file
func LoadPathsConfig() (*PathsConfig, error) {
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

// SavePathsConfig saves the user's configured paths to the configuration file
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

// LoadExtensionsConfig loads the file extension categories from the embedded configuration
func LoadExtensionsConfig() (*Config, error) {
	data, err := embeddedConfig.ReadFile("config.json")
	if err != nil {
		return nil, err
	}

	config := &Config{}
	if err := json.Unmarshal(data, config); err != nil {
		return nil, err
	}

	return config, nil
}

// Organize organizes files in the specified directory into categorized folders
// based on their file extensions. This is the main function that can be used
// programmatically by other Go applications.
func Organize(path string) error {
	config, err := LoadExtensionsConfig()
	if err != nil {
		return fmt.Errorf("failed to load extensions config: %w", err)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("path does not exist: %s", path)
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	fileCount := 0
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		itemPath := filepath.Join(path, entry.Name())
		ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(entry.Name())), ".")

		category := "others"
		for cat, exts := range config.Extensions {
			for _, e := range exts {
				if ext == e {
					category = cat
					break
				}
			}
			if category != "others" {
				break
			}
		}

		folderPath := filepath.Join(path, category)
		if err := os.MkdirAll(folderPath, 0755); err != nil {
			continue
		}

		destPath := filepath.Join(folderPath, entry.Name())
		counter := 1
		baseName := strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))
		extension := filepath.Ext(entry.Name())

		for {
			if _, err := os.Stat(destPath); os.IsNotExist(err) {
				break
			}
			destPath = filepath.Join(folderPath, fmt.Sprintf("%s_%d%s", baseName, counter, extension))
			counter++
		}

		if err := os.Rename(itemPath, destPath); err != nil {
			continue
		}

		fileCount++
	}

	return nil
}

// GetFileCategory returns the category for a given file extension
func GetFileCategory(ext string, config *Config) string {
	ext = strings.TrimPrefix(strings.ToLower(ext), ".")
	category := "others"

	for cat, exts := range config.Extensions {
		for _, e := range exts {
			if ext == e {
				return cat
			}
		}
	}

	return category
}

// GetSupportedCategories returns a list of all supported file categories
func GetSupportedCategories(config *Config) []string {
	categories := make([]string, 0, len(config.Extensions))
	for cat := range config.Extensions {
		categories = append(categories, cat)
	}
	return categories
}

// GetExtensionsForCategory returns the file extensions for a specific category
func GetExtensionsForCategory(category string, config *Config) ([]string, error) {
	exts, ok := config.Extensions[category]
	if !ok {
		return nil, errors.New("category not found")
	}
	return exts, nil
}
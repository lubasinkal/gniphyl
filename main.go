package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

//go:embed config.json
var embeddedConfig embed.FS

type Config struct {
	Extensions map[string][]string `json:"extensions"`
}

type PathsConfig struct {
	Paths []string `json:"paths"`
}

const appName = "gniphyl"

const (
	colorReset   = "\033[0m"
	colorRed     = "\033[31;1m"
	colorGreen   = "\033[32;1m"
	colorYellow  = "\033[33;1m"
	colorCyan    = "\033[36;1m"
	colorMagenta = "\033[35;1m"
)

func green(s string) string   { return colorGreen + s + colorReset }
func red(s string) string     { return colorRed + s + colorReset }
func yellow(s string) string  { return colorYellow + s + colorReset }
func cyan(s string) string    { return colorCyan + s + colorReset }
func magenta(s string) string { return colorMagenta + s + colorReset }

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

func getConfigFilePath() (string, error) {
	configFolder, err := getConfigFolder()
	if err != nil {
		return "", err
	}
	return filepath.Join(configFolder, "config.toml"), nil
}

func loadPathsConfig() (*PathsConfig, error) {
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

func savePathsConfig(config *PathsConfig) error {
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

func loadExtensionsConfig() (*Config, error) {
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

func organize(path string) error {
	config, err := loadExtensionsConfig()
	if err != nil {
		return fmt.Errorf("failed to load extensions config: %w", err)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Printf("%s %s does not exist.\n", red("Error:"), path)
		return err
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	fmt.Printf("%s Organizing files in directory: %s\n", green("✓"), cyan(path))

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
			fmt.Printf("%s Failed to create folder %s: %v\n", yellow("Warning:"), category, err)
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
			fmt.Printf("%s Skipping %s due to error: %v\n", yellow("Warning:"), entry.Name(), err)
			continue
		}

		fileCount++
	}

	fmt.Printf("%s Gniphyl run completed (%d files organized)\n", yellow("✓"), fileCount)
	return nil
}

func showHelp() {
	fmt.Printf(`%s

This tool is designed to help you organize your files efficiently.
You can add paths, list them, and perform other operations.

%s

%s
  gniphyl [command] [args]
  gniphyl --help

%s
  add [path]    Add a new path to the configuration
  rm [path]     Remove a path from the configuration
  list          List all configured paths
  run           Run the organization process on the configured paths
  --help, -h    Show this help message

%s
  gniphyl add /path/to/folder
  gniphyl list
  gniphyl run
`,
		magenta("Gniphyl CLI - File organization tool"),
		yellow("Made by lubasinkal ;) https://lubasiverse.pages.dev"),
		cyan("Usage:"),
		cyan("Commands:"),
		cyan("Examples:"),
	)
}

func cmdAdd(args []string) {
	if len(args) == 0 {
		fmt.Printf("%s Please provide a path to add.\n", red("Error:"))
		fmt.Println("Usage: gniphyl add [path]")
		return
	}

	path := args[0]

	config, err := loadPathsConfig()
	if err != nil {
		fmt.Printf("%s Failed to load config: %v\n", red("Error:"), err)
		return
	}

	for _, p := range config.Paths {
		if p == path {
			fmt.Printf("%s Path %s is already in the configuration.\n", yellow("Warning:"), path)
			return
		}
	}

	config.Paths = append(config.Paths, path)

	if err := savePathsConfig(config); err != nil {
		fmt.Printf("%s Failed to save config: %v\n", red("Error:"), err)
		return
	}

	fmt.Printf("%s Added %s\n", green("✓"), path)
}

func cmdRm(args []string) {
	if len(args) == 0 {
		fmt.Printf("%s Please provide a path to remove.\n", red("Error:"))
		fmt.Println("Usage: gniphyl rm [path]")
		return
	}

	path := args[0]

	config, err := loadPathsConfig()
	if err != nil {
		fmt.Printf("%s Failed to load config: %v\n", red("Error:"), err)
		return
	}

	found := false
	newPaths := []string{}
	for _, p := range config.Paths {
		if p != path {
			newPaths = append(newPaths, p)
		} else {
			found = true
		}
	}

	if !found {
		fmt.Printf("%s Path %s not found in the configuration.\n", yellow("Warning:"), path)
		return
	}

	config.Paths = newPaths

	if err := savePathsConfig(config); err != nil {
		fmt.Printf("%s Failed to save config: %v\n", red("Error:"), err)
		return
	}

	fmt.Printf("%s Removed %s\n", green("✓"), path)
}

func cmdList(args []string) {
	config, err := loadPathsConfig()
	if err != nil {
		fmt.Printf("%s Failed to load config: %v\n", red("Error:"), err)
		return
	}

	if len(config.Paths) == 0 {
		fmt.Printf("%s No paths configured.\n", yellow("Info:"))
		return
	}

	fmt.Printf("\n%s\n", magenta("Configured Paths:"))
	fmt.Println(strings.Repeat("-", 50))
	for i, path := range config.Paths {
		fmt.Printf("%s %s\n", cyan(fmt.Sprintf("%d.", i+1)), path)
	}
	fmt.Println(strings.Repeat("-", 50))
}

func cmdRun(args []string) {
	config, err := loadPathsConfig()
	if err != nil {
		fmt.Printf("%s Failed to load config: %v\n", red("Error:"), err)
		return
	}

	if len(config.Paths) == 0 {
		fmt.Printf("%s No paths configured. Please add paths first.\n", red("Error:"))
		return
	}

	fmt.Printf("%s Organizing the following paths:\n", green("✓"))
	for _, path := range config.Paths {
		fmt.Printf(" - %s\n", cyan(path))
		if err := organize(path); err != nil {
			fmt.Printf("%s Failed to organize %s: %v\n", red("Error:"), path, err)
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		showHelp()
		return
	}

	command := os.Args[1]
	args := []string{}
	if len(os.Args) > 2 {
		args = os.Args[2:]
	}

	switch command {
	case "add":
		cmdAdd(args)
	case "rm":
		cmdRm(args)
	case "list":
		cmdList(args)
	case "run":
		cmdRun(args)
	case "--help", "-h", "help":
		showHelp()
	default:
		fmt.Printf("%s Unknown command: %s\n", red("Error:"), command)
		fmt.Println("Use 'gniphyl --help' for available commands.")
		os.Exit(1)
	}
}

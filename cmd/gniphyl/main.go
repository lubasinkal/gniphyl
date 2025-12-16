package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/lubasinkal/gniphyl"
)

const (
	colorReset   = "\u001B[0m"
	colorRed     = "\u001B[31;1m"
	colorGreen   = "\u001B[32;1m"
	colorYellow  = "\u001B[33;1m"
	colorCyan    = "\u001B[36;1m"
	colorMagenta = "\u001B[35;1m"
)

func green(s string) string   { return colorGreen + s + colorReset }
func red(s string) string     { return colorRed + s + colorReset }
func yellow(s string) string  { return colorYellow + s + colorReset }
func cyan(s string) string    { return colorCyan + s + colorReset }
func magenta(s string) string { return colorMagenta + s + colorReset }

func organizeWithUI(path string) error {
	fmt.Printf("%s Organizing files in directory: %s\n", green("✓"), cyan(path))

	err := gniphyl.Organize(path)
	if err != nil {
		return err
	}

	// Count organized files for display
	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	fileCount := 0
	for _, entry := range entries {
		if !entry.IsDir() {
			fileCount++
		}
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

	config, err := gniphyl.LoadPathsConfig()
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

	if err := gniphyl.SavePathsConfig(config); err != nil {
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

	config, err := gniphyl.LoadPathsConfig()
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

	if err := gniphyl.SavePathsConfig(config); err != nil {
		fmt.Printf("%s Failed to save config: %v\n", red("Error:"), err)
		return
	}

	fmt.Printf("%s Removed %s\n", green("✓"), path)
}

func cmdList(args []string) {
	config, err := gniphyl.LoadPathsConfig()
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
	config, err := gniphyl.LoadPathsConfig()
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
		if err := organizeWithUI(path); err != nil {
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
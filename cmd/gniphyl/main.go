package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"

	"github.com/lubasinkal/gniphyl"
)

// Overridable for testing
var osExit = os.Exit

// Version information, set at build time via -ldflags
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// Terminal color support
const (
	colorReset   = "\u001B[0m"
	colorRed     = "\u001B[31;1m"
	colorGreen   = "\u001B[32;1m"
	colorYellow  = "\u001B[33;1m"
	colorCyan    = "\u001B[36;1m"
	colorMagenta = "\u001B[35;1m"
)

var useColor = true

func init() {
	// Disable colors on Windows unless ANSI support is confirmed
	if runtime.GOOS == "windows" {
		// Modern Windows Terminal (10+) supports ANSI, but we can't easily
		// detect it without syscalls. Default to no colors on Windows.
		useColor = os.Getenv("GNIPHYL_COLOR") == "1" || os.Getenv("TERM_PROGRAM") != ""
	}
}

func green(s string) string {
	if !useColor {
		return s
	}
	return colorGreen + s + colorReset
}
func red(s string) string {
	if !useColor {
		return s
	}
	return colorRed + s + colorReset
}
func yellow(s string) string {
	if !useColor {
		return s
	}
	return colorYellow + s + colorReset
}
func cyan(s string) string {
	if !useColor {
		return s
	}
	return colorCyan + s + colorReset
}
func magenta(s string) string {
	if !useColor {
		return s
	}
	return colorMagenta + s + colorReset
}

func showVersion() {
	fmt.Printf("gniphyl %s (commit: %s, built: %s, %s/%s)\n", version, commit, date, runtime.GOOS, runtime.GOARCH)
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
  run [path]    Run organization. If a path is provided, organizes it directly
                without adding to the config. Otherwise, organizes all
                configured paths.
  --version, -v Show version information
  --help, -h    Show this help message

%s
  gniphyl add /path/to/folder
  gniphyl list
  gniphyl run
  gniphyl run /path/to/folder
`,
		magenta("Gniphyl CLI - File organization tool"),
		yellow("Made by lubasinkal ;) https://lubasiverse.com"),
		cyan("Usage:"),
		cyan("Commands:"),
		cyan("Examples:"),
	)
}

func printError(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "%s %s\n", red("Error:"), fmt.Sprintf(format, args...))
}

func printWarning(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "%s %s\n", yellow("Warning:"), fmt.Sprintf(format, args...))
}

func cmdAdd(args []string) {
	if len(args) == 0 {
		printError("Please provide a path to add.")
		fmt.Println("Usage: gniphyl add [path]")
		return
	}

	path := args[0]

	config, err := gniphyl.LoadPathsConfig()
	if err != nil {
		printError("Failed to load config: %v", err)
		return
	}

	if slices.Contains(config.Paths, path) {
		printWarning("Path %s is already in the configuration.", path)
		return
	}

	config.Paths = append(config.Paths, path)

	if err := gniphyl.SavePathsConfig(config); err != nil {
		printError("Failed to save config: %v", err)
		return
	}

	fmt.Printf("%s Added %s\n", green("✓"), path)
}

func cmdRm(args []string) {
	if len(args) == 0 {
		printError("Please provide a path to remove.")
		fmt.Println("Usage: gniphyl rm [path]")
		return
	}

	path := args[0]

	config, err := gniphyl.LoadPathsConfig()
	if err != nil {
		printError("Failed to load config: %v", err)
		return
	}

	i := slices.Index(config.Paths, path)
	if i < 0 {
		printWarning("Path %s not found in the configuration.", path)
		return
	}

	config.Paths = slices.Delete(config.Paths, i, i+1)

	if err := gniphyl.SavePathsConfig(config); err != nil {
		printError("Failed to save config: %v", err)
		return
	}

	fmt.Printf("%s Removed %s\n", green("✓"), path)
}

func cmdList(args []string) {
	config, err := gniphyl.LoadPathsConfig()
	if err != nil {
		printError("Failed to load config: %v", err)
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

func organizeSinglePath(path string) {
	fmt.Printf("%s Organizing: %s\n", green("✓"), cyan(path))

	err := gniphyl.Organize(path)
	if err != nil {
		if orgErr, ok := err.(*gniphyl.OrganizeError); ok {
			for _, f := range orgErr.Failures {
				fmt.Fprintf(os.Stderr, "  %s %s -> %s: %v\n", red("✗"), f.Name, f.Category, f.Error)
			}
			return
		}
		printError("Failed to organize %s: %v", path, err)
		return
	}

	// Count subdirectories (organized categories)
	entries, err := os.ReadDir(path)
	if err != nil {
		return
	}

	catCount := 0
	totalFiles := 0
	for _, entry := range entries {
		if entry.IsDir() {
			catCount++
			catFiles, _ := countFilesInDir(filepath.Join(path, entry.Name()))
			totalFiles += catFiles
		}
	}

	fmt.Printf("%s Done — %d files organized into %d categories\n", yellow("✓"), totalFiles, catCount)
}

func countFilesInDir(dir string) (int, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0, err
	}
	count := 0
	for _, e := range entries {
		if !e.IsDir() {
			count++
		}
	}
	return count, nil
}

func cmdRun(args []string) {
	// If a path argument is provided, organize it directly
	if len(args) > 0 {
		path := args[0]
		organizeSinglePath(path)
		return
	}

	// Otherwise, use configured paths
	config, err := gniphyl.LoadPathsConfig()
	if err != nil {
		printError("Failed to load config: %v", err)
		return
	}

	if len(config.Paths) == 0 {
		printError("No paths configured. Please add paths first, or use: gniphyl run /path/to/dir")
		return
	}

	for _, path := range config.Paths {
		organizeSinglePath(path)
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
	case "--version", "-v", "version":
		showVersion()
	case "--help", "-h", "help":
		showHelp()
	default:
		// Check if first arg looks like a path (contains slash/backslash)
		if strings.ContainsAny(command, "/\\") {
			cmdRun(os.Args[1:])
		} else {
			printError("Unknown command: %s", command)
			fmt.Println("Use 'gniphyl --help' for available commands.")
			osExit(1)
		}
	}
}

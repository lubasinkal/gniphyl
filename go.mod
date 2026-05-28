module github.com/lubasinkal/gniphyl

go 1.21

// Gniphyl is a cross-platform file organization tool that helps you organize files
// into categorized folders based on their extensions. It provides both CLI and
// programmatic Go package functionality for file organization tasks.
//
// Features:
//   - Organize files by extension categories
//   - Cross-platform support (Windows, macOS, Linux)
//   - No external dependencies (pure Go standard library)
//   - CLI and Go package usage
//
// Example usage:
//   go install github.com/lubasinkal/gniphyl@latest
//   gniphyl add /path/to/directory
//   gniphyl run
//
// As a Go package:
//   import "github.com/lubasinkal/gniphyl"
//   err := gniphyl.Organize("/path/to/directory")

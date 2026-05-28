// Package gniphyl provides file organization functionality for categorizing
// files based on their extensions. It can be used both as a CLI tool and as a
// Go package in other applications.
//
// # Features
//
//   - Organize files by extension categories (images, videos, documents, etc.)
//   - Cross-platform support (Windows, macOS, Linux)
//   - No external dependencies (pure Go standard library)
//   - CLI and Go package usage
//
// # Example
//
//	import "github.com/lubasinkal/gniphyl"
//
//	func main() {
//	    err := gniphyl.Organize("/path/to/directory")
//	    if err != nil {
//	        log.Fatal(err)
//	    }
//	}
package gniphyl

# Gniphyl

[![Go Reference](https://pkg.go.dev/badge/github.com/lubasinkal/gniphyl.svg)](https://pkg.go.dev/github.com/lubasinkal/gniphyl)
[![Go Report Card](https://goreportcard.com/badge/github.com/lubasinkal/gniphyl)](https://goreportcard.com/report/github.com/lubasinkal/gniphyl)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

**Gniphyl** is a cross-platform command-line interface (CLI) tool written in **Go** that helps you organize files within specified directories. It enables you to add, delete, list, and organize directories, ensuring that files are neatly sorted based on their type (e.g., `images`, `videos`, `documents`).

**Built with Go's standard library only** - no external dependencies required! The tool is fast, lightweight, and produces standalone executables for **Windows**, **macOS**, and **Linux**.

---
Running the tool
![start](assets/image.png "starting from terminal")
___
Output after running the command:  **gniphyl run**

![output](assets/output.png "Organized Downloads folder")
___
## Features

- ✅ **Add directory paths** to a configuration file
- ✅ **Delete paths** by their index
- ✅ **List all configured paths**
- ✅ **Organize files** into categorized folders based on their extensions
- ✅ **Cross-platform** support (Windows, macOS, Linux)
- ✅ **No external dependencies** (pure Go standard library)
- ✅ **Lightweight** (~3MB binary)

---

## Installation

### Option 1: Install with Go (Recommended)

If you have Go installed, this is the easiest method. The binary will automatically be added to your PATH:

```bash
go install github.com/lubasinkal/gniphyl@latest
```

**Note:** Make sure `$GOPATH/bin` (or `$HOME/go/bin`) is in your PATH. To add it:

```bash
# Add to ~/.bashrc, ~/.zshrc, or equivalent
export PATH=$PATH:$(go env GOPATH)/bin

# On Windows (PowerShell)
$env:Path += ";$(go env GOPATH)\bin"
```

After installation, you can run `gniphyl` from anywhere!

### Option 2: Download Precompiled Binaries

Download the precompiled executables from the [Releases](https://github.com/lubasinkal/gniphyl/releases) section:

- **Windows**: `gniphyl-windows-amd64.exe` or `gniphyl-windows-arm64.exe`
- **macOS**: `gniphyl-macos-amd64` (Intel) or `gniphyl-macos-arm64` (Apple Silicon)
- **Linux**: `gniphyl-linux-amd64` or `gniphyl-linux-arm64`

### Option 3: Build from Source

**Prerequisites:**
- Go 1.21 or higher ([Download Go](https://go.dev/dl/))

**Build steps:**

```bash
# Clone the repository
git clone https://github.com/lubasinkal/gniphyl.git
cd gniphyl

# Build the CLI tool
go build -o gniphyl ./cmd/gniphyl/

# Or build for specific platforms
# Windows:
GOOS=windows GOARCH=amd64 go build -o gniphyl.exe ./cmd/gniphyl/

# macOS (Intel):
GOOS=darwin GOARCH=amd64 go build -o gniphyl ./cmd/gniphyl/

# macOS (Apple Silicon):
GOOS=darwin GOARCH=arm64 go build -o gniphyl ./cmd/gniphyl/

# Linux:
GOOS=linux GOARCH=amd64 go build -o gniphyl ./cmd/gniphyl/
```

### Option 4: Use as a Go Package

Add Gniphyl as a dependency to your Go project:

```bash
go get github.com/lubasinkal/gniphyl
```

Then import and use it in your code:

```go
import "github.com/lubasinkal/gniphyl"

func main() {
    err := gniphyl.Organize("/path/to/directory")
    if err != nil {
        log.Fatal(err)
    }
}
```

**See the [example](example/main.go) for more usage examples.**

### Adding to PATH

To run `gniphyl` from anywhere, add it to your system's PATH:

   To make it easier to run the `gniphyl` commands from any terminal session, move the downloaded executable to a directory that is part of your system's `PATH`.

   #### On **Windows**:
   - Move `gniphyl.exe` to a directory like `C:\Windows\System32` or another directory already in your system’s `PATH`.
   - Alternatively, add the directory where `gniphyl.exe` is located to your system’s `PATH` environment variable.

     To add a directory to the `PATH`:
     1. Right-click on **This PC** or **Computer** and select **Properties**.
     2. Select **Advanced system settings**.
     3. Click the **Environment Variables** button.
     4. In the **System variables** section, scroll to find the `Path` variable and click **Edit**.
     5. Add the path to the directory containing `gniphyl.exe` and click **OK**.

   #### On **macOS** and **Linux**:
   - Move the executable to `/usr/local/bin` (or another directory in your system’s `PATH`).
     ```bash
     sudo mv gniphyl /usr/local/bin/
     ```
     For **Linux**:
     ```bash
     sudo mv gniphyl-linux /usr/local/bin/
     ```
   - Ensure the executable is accessible by running:
     ```bash
     sudo chmod +x /usr/local/bin/gniphyl
     ```

     For **Linux**:
     ```bash
     sudo chmod +x /usr/local/bin/gniphyl-linux
     ```

---

## Usage

### Commands

```bash
# Show help
gniphyl --help

# Show version
gniphyl --version

# Add a directory path to organize
gniphyl add <directory_path>

# Remove a path from configuration
gniphyl rm <directory_path>

# List all configured paths
gniphyl list

# Organize files in all configured paths
gniphyl run

# Organize a specific directory directly (no config needed)
gniphyl run /path/to/directory
```

### Go Package Usage

Gniphyl can also be used as a Go package in your projects:

```go
import "github.com/lubasinkal/gniphyl"

func main() {
    // Load configuration
    config, err := gniphyl.LoadExtensionsConfig()
    if err != nil {
        log.Fatal(err)
    }

    // Organize a specific directory
    err = gniphyl.Organize("/path/to/directory")
    if err != nil {
        log.Fatal(err)
    }

    // Get file category for a specific extension
    category := gniphyl.GetFileCategory(".jpg", config)
    fmt.Println("Category:", category) // Output: images

    // Get all supported categories
    categories := gniphyl.GetSupportedCategories(config)
    fmt.Println("Supported categories:", categories)
}
```

---

## Example Workflow

```bash
# Add your Downloads folder
gniphyl add /users/name/downloads

# List configured paths
gniphyl list

# Output:
# Configured Paths:
# --------------------------------------------------
# 1. /users/name/downloads
# --------------------------------------------------

# Organize all configured paths
gniphyl run

# Or organize a directory directly (without adding it to config)
gniphyl run /path/to/downloads

# Files will be sorted into folders:
# - images/     (jpg, png, gif, etc.)
# - videos/     (mp4, mkv, avi, etc.)
# - documents/  (pdf, docx, txt, etc.)
# - compressed/ (zip, rar, tar, etc.)
# - audio/      (mp3, wav, flac, etc.)
# - code/       (html, css, js, py, etc.)
# - executables/ (exe, msi)
# - others/     (unrecognized extensions)
```

## API Documentation

### Package Functions

#### `Organize(path string) error`
Organizes files in the specified directory into categorized folders based on their file extensions.

**Parameters:**
- `path`: The directory path to organize

**Returns:**
- `error`: Any error encountered during organization

**Example:**
```go
err := gniphyl.Organize("/path/to/directory")
if err != nil {
    log.Fatal(err)
}
```

#### `LoadExtensionsConfig() (*Config, error)`
Loads the file extension categories from the embedded configuration.

**Returns:**
- `*Config`: Configuration containing extension mappings
- `error`: Any error encountered during loading

**Example:**
```go
config, err := gniphyl.LoadExtensionsConfig()
if err != nil {
    log.Fatal(err)
}
```

#### `LoadPathsConfig() (*PathsConfig, error)`
Loads the user's configured paths from the configuration file.

**Returns:**
- `*PathsConfig`: User's configured paths
- `error`: Any error encountered during loading

**Example:**
```go
pathsConfig, err := gniphyl.LoadPathsConfig()
if err != nil {
    log.Fatal(err)
}
```

#### `SavePathsConfig(config *PathsConfig) error`
Saves the user's configured paths to the configuration file.

**Parameters:**
- `config`: PathsConfig containing paths to save

**Returns:**
- `error`: Any error encountered during saving

**Example:**
```go
config := &gniphyl.PathsConfig{Paths: []string{"/path/to/dir"}}
err := gniphyl.SavePathsConfig(config)
if err != nil {
    log.Fatal(err)
}
```

#### `GetFileCategory(ext string, config *Config) string`
Returns the category for a given file extension.

**Parameters:**
- `ext`: File extension (e.g., ".jpg", "jpg")
- `config`: Configuration containing extension mappings

**Returns:**
- `string`: The category name

**Example:**
```go
config, _ := gniphyl.LoadExtensionsConfig()
category := gniphyl.GetFileCategory(".jpg", config)
fmt.Println(category) // Output: images
```

#### `GetSupportedCategories(config *Config) []string`
Returns a list of all supported file categories.

**Parameters:**
- `config`: Configuration containing extension mappings

**Returns:**
- `[]string`: List of category names

**Example:**
```go
config, _ := gniphyl.LoadExtensionsConfig()
categories := gniphyl.GetSupportedCategories(config)
fmt.Println(categories) // Output: [images videos documents compressed executables audio code]
```

#### `GetExtensionsForCategory(category string, config *Config) ([]string, error)`
Returns the file extensions for a specific category.

**Parameters:**
- `category`: Category name (e.g., "images")
- `config`: Configuration containing extension mappings

**Returns:**
- `[]string`: List of file extensions
- `error`: Error if category not found

**Example:**
```go
config, _ := gniphyl.LoadExtensionsConfig()
exts, err := gniphyl.GetExtensionsForCategory("images", config)
if err != nil {
    log.Fatal(err)
}
fmt.Println(exts) // Output: [jpg jpeg png gif bmp webp]
```

## File Categories

The tool organizes files based on their extensions into the following categories:

- **images**: jpg, jpeg, png, gif, bmp, webp
- **videos**: mp4, mkv, webm, flv, avi, mov
- **documents**: pdf, doc, docx, xls, xlsx, ppt, pptx, txt, csv
- **compressed**: zip, rar, tar, gz, 7z
- **executables**: exe, msi
- **audio**: mp3, wav, flac, m4a, aac
- **code**: html, css, js, py, java, c, cpp, h, hpp, php, sql
- **others**: any other file type

You can customize these categories by modifying `extensions.json` in the source and rebuilding.

## Technical Details

- **Language**: Go (using only standard library)
- **Module Path**: `github.com/lubasinkal/gniphyl`
- **Dependencies**: None (stdlib only)
- **Config Storage**: 
  - Windows: `%LOCALAPPDATA%\gniphyl\config.json`
  - macOS/Linux: `~/.config/gniphyl/config.json`
- **Config Format**: JSON
  - *Migrated automatically from legacy `config.toml` on first run*
- **Binary Size**: ~3MB (standalone, no runtime required)
- **Supported Platforms**: Windows, macOS (Intel & Apple Silicon), Linux (amd64 & arm64)
- **Go Version**: 1.21+
- **License**: MIT
- **Package Documentation**: [pkg.go.dev](https://pkg.go.dev/github.com/lubasinkal/gniphyl)

---

## Contributing

Contributions are welcome! Please follow these guidelines:

### Development Setup

```bash
# Fork the repository
# Clone your fork
git clone https://github.com/yourusername/gniphyl.git
cd gniphyl

# Install dependencies (none required, but you may want to install these tools)
go install golang.org/x/tools/gopls@latest          # Go language server
go install honnef.co/go/tools/cmd/staticcheck@latest # Static analysis
```

### Development Workflow

1. **Create a feature branch:**
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes:**
   - Follow Go conventions and best practices
   - Add appropriate tests
   - Update documentation as needed

3. **Run tests and checks:**
   ```bash
   go test ./...
   go vet ./...
   staticcheck ./...
   ```

4. **Commit your changes:**
   ```bash
   git commit -m "feat: add your feature description"
   ```
   - Use conventional commit messages
   - Keep commits focused and atomic

5. **Push to your fork:**
   ```bash
   git push origin feature/your-feature-name
   ```

6. **Create a Pull Request:**
   - Target the `main` branch
   - Provide a clear description of your changes
   - Reference any related issues

### Code Style Guidelines

- Follow [Effective Go](https://go.dev/doc/effective_go) guidelines
- Use `gofmt` for code formatting
- Write clear, descriptive commit messages
- Include appropriate documentation
- Add tests for new functionality

### Testing

Run the test suite:
```bash
go test -v ./...
```

Run static analysis:
```bash
go vet ./...
staticcheck ./...
```

---

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

## Contact

For any questions or issues, please contact:
- **Name:** Lubasi Nkalolang
- **Email:** lubasinkal@outlook.com
- **GitHub:** [lubasinkal](https://github.com/lubasinkal)
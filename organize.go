package gniphyl

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

// FileResult records the outcome of organizing a single file.
type FileResult struct {
	Name     string `json:"name"`
	Category string `json:"category"`
	Error    error  `json:"-"`
}

// OrganizeError aggregates failures that occurred during file organization.
type OrganizeError struct {
	Path     string       `json:"path"`
	Failures []FileResult `json:"failures"`
}

func (e *OrganizeError) Error() string {
	return fmt.Sprintf("organized %s with %d failure(s)", e.Path, len(e.Failures))
}

// Organize organizes files in the specified directory into categorized folders
// based on their file extensions. This is the main function that can be used
// programmatically by other Go applications.
//
// If any files fail to move, the returned error will be of type *OrganizeError
// and will contain details about each failure.
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

	var failures []FileResult
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		itemPath := filepath.Join(path, entry.Name())
		ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(entry.Name())), ".")

		category := "others"
		for cat, exts := range config.Extensions {
			if slices.Contains(exts, ext) {
				category = cat
			}
			if category != "others" {
				break
			}
		}

		folderPath := filepath.Join(path, category)
		if err := os.MkdirAll(folderPath, 0755); err != nil {
			failures = append(failures, FileResult{Name: entry.Name(), Category: category, Error: err})
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
			failures = append(failures, FileResult{Name: entry.Name(), Category: category, Error: err})
			continue
		}
	}

	if len(failures) > 0 {
		return &OrganizeError{Path: path, Failures: failures}
	}

	return nil
}

// GetFileCategory returns the category for a given file extension.
// The extension may include a leading dot (e.g., ".jpg") or not (e.g., "jpg").
// Unrecognized extensions return "others".
func GetFileCategory(ext string, config *Config) string {
	ext = strings.TrimPrefix(strings.ToLower(ext), ".")

	for cat, exts := range config.Extensions {
		if slices.Contains(exts, ext) {
			return cat
		}
	}

	return "others"
}

// GetSupportedCategories returns a list of all supported file categories.
func GetSupportedCategories(config *Config) []string {
	categories := make([]string, 0, len(config.Extensions))
	for cat := range config.Extensions {
		categories = append(categories, cat)
	}
	return categories
}

// GetExtensionsForCategory returns the file extensions for a specific category.
// Returns an error if the category is not found.
func GetExtensionsForCategory(category string, config *Config) ([]string, error) {
	exts, ok := config.Extensions[category]
	if !ok {
		return nil, errors.New("category not found")
	}
	return exts, nil
}

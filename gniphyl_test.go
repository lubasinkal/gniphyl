package gniphyl

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadExtensionsConfig(t *testing.T) {
	config, err := LoadExtensionsConfig()
	if err != nil {
		t.Fatalf("LoadExtensionsConfig() returned error: %v", err)
	}

	if config == nil {
		t.Fatal("LoadExtensionsConfig() returned nil config")
	}

	// Verify known categories exist
	expectedCategories := []string{"images", "videos", "documents", "compressed", "executables", "audio", "code"}
	for _, cat := range expectedCategories {
		if _, ok := config.Extensions[cat]; !ok {
			t.Errorf("expected category %q not found in config", cat)
		}
	}

	// Verify images has common extensions
	imageExts := config.Extensions["images"]
	if len(imageExts) == 0 {
		t.Error("images category has no extensions")
	}

	foundJpg := false
	for _, ext := range imageExts {
		if ext == "jpg" {
			foundJpg = true
			break
		}
	}
	if !foundJpg {
		t.Error("images category should contain 'jpg'")
	}
}

func TestGetFileCategory(t *testing.T) {
	config := &Config{
		Extensions: map[string][]string{
			"images":      {"jpg", "jpeg", "png"},
			"documents":   {"pdf", "txt"},
			"executables": {"exe"},
		},
	}

	tests := []struct {
		name     string
		ext      string
		expected string
	}{
		{"jpg extension", ".jpg", "images"},
		{"jpeg without dot", "jpeg", "images"},
		{"PNG uppercase", ".PNG", "images"},
		{"pdf with dot", ".pdf", "documents"},
		{"txt without dot", "txt", "documents"},
		{"exe with dot", ".exe", "executables"},
		{"unknown ext", ".xyz", "others"},
		{"empty string", "", "others"},
		{"dot only", ".", "others"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetFileCategory(tt.ext, config)
			if got != tt.expected {
				t.Errorf("GetFileCategory(%q) = %q, want %q", tt.ext, got, tt.expected)
			}
		})
	}
}

func TestGetSupportedCategories(t *testing.T) {
	config := &Config{
		Extensions: map[string][]string{
			"images":    {"jpg"},
			"documents": {"pdf"},
			"audio":     {"mp3"},
		},
	}

	categories := GetSupportedCategories(config)
	if len(categories) != 3 {
		t.Errorf("expected 3 categories, got %d: %v", len(categories), categories)
	}

	// Check all expected categories are present
	expected := map[string]bool{"images": true, "documents": true, "audio": true}
	for _, cat := range categories {
		delete(expected, cat)
	}
	if len(expected) > 0 {
		t.Errorf("missing categories: %v", expected)
	}
}

func TestGetExtensionsForCategory(t *testing.T) {
	config := &Config{
		Extensions: map[string][]string{
			"images": {"jpg", "png"},
		},
	}

	t.Run("existing category", func(t *testing.T) {
		exts, err := GetExtensionsForCategory("images", config)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(exts) != 2 || exts[0] != "jpg" || exts[1] != "png" {
			t.Errorf("GetExtensionsForCategory('images') = %v, want [jpg png]", exts)
		}
	})

	t.Run("non-existing category", func(t *testing.T) {
		_, err := GetExtensionsForCategory("nonexistent", config)
		if err == nil {
			t.Error("expected error for nonexistent category, got nil")
		}
	})
}

func TestOrganize(t *testing.T) {
	// Create a temporary directory structure
	tmpDir := t.TempDir()

	// Create test files with various extensions
	testFiles := []string{
		"photo.jpg",
		"document.pdf",
		"archive.zip",
		"script.js",
		"music.mp3",
		"readme.txt",
		"unknown.xyz",
	}

	for _, f := range testFiles {
		path := filepath.Join(tmpDir, f)
		if err := os.WriteFile(path, []byte("test content"), 0644); err != nil {
			t.Fatalf("failed to create test file %s: %v", f, err)
		}
	}

	// Run Organize
	if err := Organize(tmpDir); err != nil {
		t.Fatalf("Organize() returned error: %v", err)
	}

	// Verify files were moved to correct categories
	expected := map[string][]string{
		"images":     {"photo.jpg"},
		"documents":  {"document.pdf", "readme.txt"},
		"compressed": {"archive.zip"},
		"code":       {"script.js"},
		"audio":      {"music.mp3"},
		"others":     {"unknown.xyz"},
	}

	for category, files := range expected {
		catDir := filepath.Join(tmpDir, category)
		dirEntries, err := os.ReadDir(catDir)
		if err != nil {
			t.Errorf("failed to read category dir %s: %v", category, err)
			continue
		}

		if len(dirEntries) != len(files) {
			t.Errorf("category %q has %d files, want %d", category, len(dirEntries), len(files))
		}

		for _, f := range files {
			path := filepath.Join(catDir, f)
			if _, err := os.Stat(path); os.IsNotExist(err) {
				t.Errorf("expected file %s in category %q, but not found", f, category)
			}
		}
	}

	// Verify original directory has no leftover files (only category dirs)
	entries, err := os.ReadDir(tmpDir)
	if err != nil {
		t.Fatalf("failed to read tmp dir: %v", err)
	}

	for _, e := range entries {
		if !e.IsDir() {
			t.Errorf("unexpected leftover file in root: %s", e.Name())
		}
	}
}

func TestOrganizeWithDuplicateNames(t *testing.T) {
	tmpDir := t.TempDir()

	// Create same-named files in different categories
	files := []string{"readme.txt", "notes.txt"}
	for _, f := range files {
		if err := os.WriteFile(filepath.Join(tmpDir, f), []byte("test"), 0644); err != nil {
			t.Fatalf("failed to create %s: %v", f, err)
		}
	}

	if err := Organize(tmpDir); err != nil {
		t.Fatalf("Organize() returned error: %v", err)
	}

	// Both should be in documents
	docDir := filepath.Join(tmpDir, "documents")
	if _, err := os.Stat(filepath.Join(docDir, "readme.txt")); os.IsNotExist(err) {
		t.Error("expected readme.txt in documents")
	}
	if _, err := os.Stat(filepath.Join(docDir, "notes.txt")); os.IsNotExist(err) {
		t.Error("expected notes.txt in documents")
	}
}

func TestOrganizeNonexistentPath(t *testing.T) {
	err := Organize("/nonexistent/path/that/does/not/exist")
	if err == nil {
		t.Error("expected error for nonexistent path, got nil")
	}
}

func TestOrganizeEmptyDirectory(t *testing.T) {
	tmpDir := t.TempDir()

	err := Organize(tmpDir)
	if err != nil {
		t.Fatalf("Organize() on empty dir returned error: %v", err)
	}

	// Should still be empty
	entries, err := os.ReadDir(tmpDir)
	if err != nil {
		t.Fatalf("failed to read dir: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected empty directory, found %d entries", len(entries))
	}
}

func TestOrganizeIdempotent(t *testing.T) {
	tmpDir := t.TempDir()

	testFiles := []string{"test.jpg", "test.pdf"}
	for _, f := range testFiles {
		if err := os.WriteFile(filepath.Join(tmpDir, f), []byte("test"), 0644); err != nil {
			t.Fatalf("failed to create %s: %v", f, err)
		}
	}

	// Run twice
	if err := Organize(tmpDir); err != nil {
		t.Fatalf("first Organize() returned error: %v", err)
	}
	if err := Organize(tmpDir); err != nil {
		t.Fatalf("second Organize() returned error: %v", err)
	}

	// Should still be organized correctly
	imgDir := filepath.Join(tmpDir, "images")
	if _, err := os.Stat(filepath.Join(imgDir, "test.jpg")); os.IsNotExist(err) {
		t.Error("expected test.jpg in images after second run")
	}
}

func TestPathsConfigRoundTrip(t *testing.T) {
	// Save and load paths config
	original := &PathsConfig{
		Paths: []string{"/path/one", "/path/two", "/path/three"},
	}

	if err := SavePathsConfig(original); err != nil {
		t.Fatalf("SavePathsConfig() returned error: %v", err)
	}

	loaded, err := LoadPathsConfig()
	if err != nil {
		t.Fatalf("LoadPathsConfig() returned error: %v", err)
	}

	if len(loaded.Paths) != len(original.Paths) {
		t.Fatalf("loaded %d paths, want %d", len(loaded.Paths), len(original.Paths))
	}

	for i, p := range original.Paths {
		if loaded.Paths[i] != p {
			t.Errorf("path[%d] = %q, want %q", i, loaded.Paths[i], p)
		}
	}
}

func TestOrganizeErrorImplementsError(t *testing.T) {
	err := &OrganizeError{
		Path: "/tmp",
		Failures: []FileResult{
			{Name: "bad.exe", Category: "executables", Error: os.ErrPermission},
		},
	}
	if err.Error() == "" {
		t.Error("OrganizeError.Error() should not be empty")
	}

	// Verify it satisfies the error interface
	var e error = err
	if e == nil {
		t.Error("OrganizeError should implement error")
	}
}

func TestPathsConfigEmpty(t *testing.T) {
	config := &PathsConfig{Paths: []string{}}
	if err := SavePathsConfig(config); err != nil {
		t.Fatalf("SavePathsConfig() returned error: %v", err)
	}

	loaded, err := LoadPathsConfig()
	if err != nil {
		t.Fatalf("LoadPathsConfig() returned error: %v", err)
	}

	if loaded.Paths == nil {
		t.Error("expected non-nil Paths slice")
	}
	if len(loaded.Paths) != 0 {
		t.Errorf("expected empty paths, got %v", loaded.Paths)
	}
}

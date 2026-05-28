package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

// TestMainHelp verifies the help output doesn't crash
func TestMainHelp(t *testing.T) {
	// Just verify that the binary builds and runs --help
	// We use os/exec equivalent by capturing stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Simulate --help
	os.Args = []string{"gniphyl", "--help"}
	main()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	if !contains(output, "Gniphyl CLI") {
		t.Error("help output should contain 'Gniphyl CLI'")
	}
	if !contains(output, "Commands:") {
		t.Error("help output should contain 'Commands:'")
	}
}

func TestMainVersion(t *testing.T) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	os.Args = []string{"gniphyl", "--version"}
	main()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	if !contains(output, "gniphyl") {
		t.Error("version output should contain 'gniphyl'")
	}
}

func TestRunPathDirectly(t *testing.T) {
	// Create a temp dir with a test file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.jpg")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	// Run gniphyl run <tmpDir>
	os.Args = []string{"gniphyl", "run", tmpDir}

	// Capture stderr for errors
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	main()

	w.Close()
	os.Stderr = oldStderr

	var errBuf bytes.Buffer
	errBuf.ReadFrom(r)
	if errBuf.Len() > 0 {
		t.Errorf("unexpected stderr output: %s", errBuf.String())
	}

	// Check file was organized
	imgDir := filepath.Join(tmpDir, "images")
	if _, err := os.Stat(imgDir); os.IsNotExist(err) {
		t.Errorf("expected images directory to be created")
	}
	if _, err := os.Stat(filepath.Join(imgDir, "test.jpg")); os.IsNotExist(err) {
		t.Errorf("expected test.jpg to be moved to images/")
	}
}

func TestUnknownCommandShowsError(t *testing.T) {
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	// Capture exit
	oldExit := osExit
	defer func() { osExit = oldExit }()
	osExit = func(code int) {}

	os.Args = []string{"gniphyl", "unknowncommand"}
	main()

	w.Close()
	os.Stderr = oldStderr

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	if !contains(output, "Error") {
		t.Error("expected error output for unknown command")
	}
}

// contains is a simple strings.Contains wrapper for testing
func contains(s, substr string) bool {
	return len(s) >= len(substr) && containsStr(s, substr)
}

func containsStr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

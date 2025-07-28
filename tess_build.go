//go:build cgo
// +build cgo

package tess

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func init() {
	// Check if we're in a module context and need to build the library
	if isModuleDependency() {
		if err := ensureLibraryBuilt(); err != nil {
			panic(fmt.Sprintf("Failed to build libtess2: %v", err))
		}
	}
}

// isModuleDependency checks if this package is being used as a dependency
func isModuleDependency() bool {
	// Check if we're in a Go module cache or vendor directory
	wd, err := os.Getwd()
	if err != nil {
		return false
	}
	
	// If we're in a module cache or vendor directory, we're a dependency
	return contains(wd, "pkg/mod") || contains(wd, "vendor")
}

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || 
		(len(s) > len(substr) && (s[:len(substr)] == substr || 
		s[len(s)-len(substr):] == substr || 
		func() bool {
			for i := 1; i <= len(s)-len(substr); i++ {
				if s[i:i+len(substr)] == substr {
					return true
				}
			}
			return false
		}())))
}

// ensureLibraryBuilt ensures the libtess2 library is built
func ensureLibraryBuilt() error {
	// Get the package directory
	pkgDir, err := getPackageDir()
	if err != nil {
		return fmt.Errorf("failed to get package directory: %w", err)
	}
	
	// Check if library already exists
	libPath := filepath.Join(pkgDir, "libtess2", "libtess2.a")
	if _, err := os.Stat(libPath); err == nil {
		return nil // Library already exists
	}
	
	// Try to build the library
	return buildLibrary(pkgDir)
}

// getPackageDir returns the directory containing this package
func getPackageDir() (string, error) {
	// Get the current file's directory
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("failed to get caller info")
	}
	
	return filepath.Dir(filename), nil
}

// buildLibrary builds the libtess2 library
func buildLibrary(pkgDir string) error {
	// Check if make is available
	if _, err := exec.LookPath("make"); err != nil {
		return fmt.Errorf("make not found: %w", err)
	}
	
	// Check if gcc is available
	if _, err := exec.LookPath("gcc"); err != nil {
		return fmt.Errorf("gcc not found: %w", err)
	}
	
	// Run make to build the library
	cmd := exec.Command("make", "libtess2/libtess2.a")
	cmd.Dir = pkgDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("make failed: %w", err)
	}
	
	return nil
} 
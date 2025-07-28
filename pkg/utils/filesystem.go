package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

// FileExists checks if a file exists
func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

// DirExists checks if a directory exists
func DirExists(dirPath string) bool {
	info, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

// EnsureDir creates a directory if it doesn't exist
func EnsureDir(dirPath string) error {
	if dirPath == "" {
		return nil
	}
	
	if DirExists(dirPath) {
		return nil
	}
	
	return os.MkdirAll(dirPath, 0755)
}

// CreateFileIfNotExists creates a file if it doesn't exist
func CreateFileIfNotExists(filePath string, content []byte) error {
	if FileExists(filePath) {
		return fmt.Errorf("file already exists: %s", filePath)
	}
	
	if err := EnsureDir(filepath.Dir(filePath)); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	
	if err := os.WriteFile(filePath, content, 0644); err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	
	return nil
}

// IsEmptyDir checks if a directory is empty
func IsEmptyDir(dirPath string) (bool, error) {
	if !DirExists(dirPath) {
		return false, fmt.Errorf("directory does not exist: %s", dirPath)
	}
	
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return false, fmt.Errorf("failed to read directory: %w", err)
	}
	
	return len(entries) == 0, nil
} 
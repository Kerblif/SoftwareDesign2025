package local

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"kr-02/internal/pkg/file_analysis/storage"
)

// LocalStorage implements the WordCloudStorage interface using the local filesystem
type LocalStorage struct {
	basePath string
}

// NewLocalStorage creates a new LocalStorage instance
func NewLocalStorage(basePath string) (storage.WordCloudStorage, error) {
	// Create the base directory if it doesn't exist
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}
	return &LocalStorage{basePath: basePath}, nil
}

// SaveWordCloud saves a word cloud image to the local filesystem
func (s *LocalStorage) SaveWordCloud(ctx context.Context, location string, image []byte) error {
	fullPath := filepath.Join(s.basePath, location)
	
	// Create the directory if it doesn't exist
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	
	// Write the file
	if err := os.WriteFile(fullPath, image, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	
	return nil
}

// GetWordCloud retrieves a word cloud image from the local filesystem
func (s *LocalStorage) GetWordCloud(ctx context.Context, location string) ([]byte, error) {
	fullPath := filepath.Join(s.basePath, location)
	
	// Read the file
	image, err := os.ReadFile(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("word cloud image not found at location %s", location)
		}
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	
	return image, nil
}
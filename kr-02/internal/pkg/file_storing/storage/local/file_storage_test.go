package local_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"kr-02/internal/pkg/file_storing/storage/local"
)

func TestNewLocalStorage(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "file_storage_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Test successful creation
	storage, err := local.NewLocalStorage(tempDir)
	assert.NoError(t, err)
	assert.NotNil(t, storage)

	// Test creation with non-existent directory (should create it)
	nonExistentDir := filepath.Join(tempDir, "non_existent")
	storage, err = local.NewLocalStorage(nonExistentDir)
	assert.NoError(t, err)
	assert.NotNil(t, storage)

	// Verify the directory was created
	_, err = os.Stat(nonExistentDir)
	assert.NoError(t, err)
}

func TestSaveFile(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "file_storage_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a storage instance
	storage, err := local.NewLocalStorage(tempDir)
	require.NoError(t, err)

	// Test saving a file
	testContent := []byte("test file content")
	location := "test.txt"
	err = storage.SaveFile(context.Background(), location, testContent)
	assert.NoError(t, err)

	// Verify the file was created
	fullPath := filepath.Join(tempDir, location)
	savedContent, err := os.ReadFile(fullPath)
	assert.NoError(t, err)
	assert.Equal(t, testContent, savedContent)

	// Test saving to a subdirectory
	subDirLocation := filepath.Join("subdir", "test.txt")
	err = storage.SaveFile(context.Background(), subDirLocation, testContent)
	assert.NoError(t, err)

	// Verify the file was created in the subdirectory
	fullSubDirPath := filepath.Join(tempDir, subDirLocation)
	savedSubDirContent, err := os.ReadFile(fullSubDirPath)
	assert.NoError(t, err)
	assert.Equal(t, testContent, savedSubDirContent)
}

func TestGetFile(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "file_storage_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a storage instance
	storage, err := local.NewLocalStorage(tempDir)
	require.NoError(t, err)

	// Create a test file
	testContent := []byte("test file content")
	location := "test.txt"
	fullPath := filepath.Join(tempDir, location)
	err = os.WriteFile(fullPath, testContent, 0644)
	require.NoError(t, err)

	// Test retrieving the file
	retrievedContent, err := storage.GetFile(context.Background(), location)
	assert.NoError(t, err)
	assert.Equal(t, testContent, retrievedContent)

	// Test retrieving a non-existent file
	_, err = storage.GetFile(context.Background(), "non_existent.txt")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestSaveAndGetFile(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "file_storage_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a storage instance
	storage, err := local.NewLocalStorage(tempDir)
	require.NoError(t, err)

	// Test saving and then retrieving a file
	testContent := []byte("test file content")
	location := "test.txt"
	
	// Save the file
	err = storage.SaveFile(context.Background(), location, testContent)
	assert.NoError(t, err)
	
	// Retrieve the file
	retrievedContent, err := storage.GetFile(context.Background(), location)
	assert.NoError(t, err)
	assert.Equal(t, testContent, retrievedContent)
}
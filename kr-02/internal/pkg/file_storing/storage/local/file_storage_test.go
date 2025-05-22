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

	// Test creation with empty path (should use current directory)
	t.Run("Empty path", func(t *testing.T) {
		storage, err := local.NewLocalStorage("")
		assert.NoError(t, err)
		assert.NotNil(t, storage)
	})

	// Test creation with relative path
	t.Run("Relative path", func(t *testing.T) {
		relativeDir := "./test_relative_dir"
		defer os.RemoveAll(relativeDir)

		storage, err := local.NewLocalStorage(relativeDir)
		assert.NoError(t, err)
		assert.NotNil(t, storage)

		// Verify the directory was created
		_, err = os.Stat(relativeDir)
		assert.NoError(t, err)
	})
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

	// Test saving with empty location
	t.Run("Empty location", func(t *testing.T) {
		err = storage.SaveFile(context.Background(), "", testContent)
		assert.NoError(t, err)

		// Verify the file was created with the default filename
		defaultPath := filepath.Join(tempDir, "default.txt")
		savedEmptyLocContent, err := os.ReadFile(defaultPath)
		assert.NoError(t, err)
		assert.Equal(t, testContent, savedEmptyLocContent)
	})

	// Test saving with empty content
	t.Run("Empty content", func(t *testing.T) {
		emptyContent := []byte{}
		emptyContentLoc := "empty.txt"
		err = storage.SaveFile(context.Background(), emptyContentLoc, emptyContent)
		assert.NoError(t, err)

		// Verify the empty file was created
		emptyContentPath := filepath.Join(tempDir, emptyContentLoc)
		savedEmptyContent, err := os.ReadFile(emptyContentPath)
		assert.NoError(t, err)
		assert.Empty(t, savedEmptyContent)
	})

	// Test overwriting existing file
	t.Run("Overwrite existing file", func(t *testing.T) {
		// First save
		firstContent := []byte("first file content")
		overwriteLoc := "overwrite.txt"
		err = storage.SaveFile(context.Background(), overwriteLoc, firstContent)
		assert.NoError(t, err)

		// Verify first save
		overwritePath := filepath.Join(tempDir, overwriteLoc)
		savedFirstContent, err := os.ReadFile(overwritePath)
		assert.NoError(t, err)
		assert.Equal(t, firstContent, savedFirstContent)

		// Second save (overwrite)
		secondContent := []byte("second file content")
		err = storage.SaveFile(context.Background(), overwriteLoc, secondContent)
		assert.NoError(t, err)

		// Verify overwrite
		savedSecondContent, err := os.ReadFile(overwritePath)
		assert.NoError(t, err)
		assert.Equal(t, secondContent, savedSecondContent)
	})
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

	// Test retrieving with empty location
	t.Run("Empty location", func(t *testing.T) {
		// Create a file with default filename
		defaultPath := filepath.Join(tempDir, "default.txt")
		defaultContent := []byte("default content")
		err = os.WriteFile(defaultPath, defaultContent, 0644)
		require.NoError(t, err)

		// Retrieve the file with empty location (should use default filename)
		retrievedDefaultContent, err := storage.GetFile(context.Background(), "")
		assert.NoError(t, err)
		assert.Equal(t, defaultContent, retrievedDefaultContent)
	})

	// Test retrieving an empty file
	t.Run("Empty file", func(t *testing.T) {
		// Create an empty file
		emptyFilePath := filepath.Join(tempDir, "empty.txt")
		err = os.WriteFile(emptyFilePath, []byte{}, 0644)
		require.NoError(t, err)

		// Retrieve the empty file
		retrievedEmptyFile, err := storage.GetFile(context.Background(), "empty.txt")
		assert.NoError(t, err)
		assert.Empty(t, retrievedEmptyFile)
	})
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

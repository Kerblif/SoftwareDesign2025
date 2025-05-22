package local_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"kr-02/internal/pkg/file_analysis/storage/local"
)

func TestNewLocalStorage(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "wordcloud_test")
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

func TestSaveWordCloud(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "wordcloud_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a storage instance
	storage, err := local.NewLocalStorage(tempDir)
	require.NoError(t, err)

	// Test saving a word cloud
	testImage := []byte("test image data")
	location := "test.png"
	err = storage.SaveWordCloud(context.Background(), location, testImage)
	assert.NoError(t, err)

	// Verify the file was created
	fullPath := filepath.Join(tempDir, location)
	savedImage, err := os.ReadFile(fullPath)
	assert.NoError(t, err)
	assert.Equal(t, testImage, savedImage)

	// Test saving to a subdirectory
	subDirLocation := filepath.Join("subdir", "test.png")
	err = storage.SaveWordCloud(context.Background(), subDirLocation, testImage)
	assert.NoError(t, err)

	// Verify the file was created in the subdirectory
	fullSubDirPath := filepath.Join(tempDir, subDirLocation)
	savedSubDirImage, err := os.ReadFile(fullSubDirPath)
	assert.NoError(t, err)
	assert.Equal(t, testImage, savedSubDirImage)
}

func TestGetWordCloud(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "wordcloud_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a storage instance
	storage, err := local.NewLocalStorage(tempDir)
	require.NoError(t, err)

	// Create a test image file
	testImage := []byte("test image data")
	location := "test.png"
	fullPath := filepath.Join(tempDir, location)
	err = os.WriteFile(fullPath, testImage, 0644)
	require.NoError(t, err)

	// Test retrieving the word cloud
	retrievedImage, err := storage.GetWordCloud(context.Background(), location)
	assert.NoError(t, err)
	assert.Equal(t, testImage, retrievedImage)

	// Test retrieving a non-existent word cloud
	_, err = storage.GetWordCloud(context.Background(), "non_existent.png")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestSaveAndGetWordCloud(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "wordcloud_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a storage instance
	storage, err := local.NewLocalStorage(tempDir)
	require.NoError(t, err)

	// Test saving and then retrieving a word cloud
	testImage := []byte("test image data")
	location := "test.png"
	
	// Save the word cloud
	err = storage.SaveWordCloud(context.Background(), location, testImage)
	assert.NoError(t, err)
	
	// Retrieve the word cloud
	retrievedImage, err := storage.GetWordCloud(context.Background(), location)
	assert.NoError(t, err)
	assert.Equal(t, testImage, retrievedImage)
}
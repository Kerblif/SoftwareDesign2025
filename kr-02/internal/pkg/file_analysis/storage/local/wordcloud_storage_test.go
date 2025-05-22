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

	// Test saving with empty location
	t.Run("Empty location", func(t *testing.T) {
		err = storage.SaveWordCloud(context.Background(), "", testImage)
		assert.NoError(t, err)

		// Verify the file was created with the default filename
		defaultPath := filepath.Join(tempDir, "default.png")
		savedEmptyLocImage, err := os.ReadFile(defaultPath)
		assert.NoError(t, err)
		assert.Equal(t, testImage, savedEmptyLocImage)
	})

	// Test saving with empty content
	t.Run("Empty content", func(t *testing.T) {
		emptyImage := []byte{}
		emptyContentLoc := "empty.png"
		err = storage.SaveWordCloud(context.Background(), emptyContentLoc, emptyImage)
		assert.NoError(t, err)

		// Verify the empty file was created
		emptyContentPath := filepath.Join(tempDir, emptyContentLoc)
		savedEmptyImage, err := os.ReadFile(emptyContentPath)
		assert.NoError(t, err)
		assert.Empty(t, savedEmptyImage)
	})

	// Test overwriting existing file
	t.Run("Overwrite existing file", func(t *testing.T) {
		// First save
		firstImage := []byte("first image data")
		overwriteLoc := "overwrite.png"
		err = storage.SaveWordCloud(context.Background(), overwriteLoc, firstImage)
		assert.NoError(t, err)

		// Verify first save
		overwritePath := filepath.Join(tempDir, overwriteLoc)
		savedFirstImage, err := os.ReadFile(overwritePath)
		assert.NoError(t, err)
		assert.Equal(t, firstImage, savedFirstImage)

		// Second save (overwrite)
		secondImage := []byte("second image data")
		err = storage.SaveWordCloud(context.Background(), overwriteLoc, secondImage)
		assert.NoError(t, err)

		// Verify overwrite
		savedSecondImage, err := os.ReadFile(overwritePath)
		assert.NoError(t, err)
		assert.Equal(t, secondImage, savedSecondImage)
	})
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

	// Test retrieving with empty location
	t.Run("Empty location", func(t *testing.T) {
		// Create a file with default filename
		defaultPath := filepath.Join(tempDir, "default.png")
		defaultImage := []byte("default image")
		err = os.WriteFile(defaultPath, defaultImage, 0644)
		require.NoError(t, err)

		// Retrieve the file with empty location (should use default filename)
		retrievedDefaultImage, err := storage.GetWordCloud(context.Background(), "")
		assert.NoError(t, err)
		assert.Equal(t, defaultImage, retrievedDefaultImage)
	})

	// Test retrieving an empty file
	t.Run("Empty file", func(t *testing.T) {
		// Create an empty file
		emptyFilePath := filepath.Join(tempDir, "empty.png")
		err = os.WriteFile(emptyFilePath, []byte{}, 0644)
		require.NoError(t, err)

		// Retrieve the empty file
		retrievedEmptyFile, err := storage.GetWordCloud(context.Background(), "empty.png")
		assert.NoError(t, err)
		assert.Empty(t, retrievedEmptyFile)
	})
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

package file_storing_service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"kr-02/internal/pkg/file_storing/repository"
	"kr-02/internal/pkg/file_storing/service"
	"kr-02/internal/pkg/file_storing/storage"
)

// MockFileRepository is a mock implementation of the FileRepository interface
type MockFileRepository struct {
	mock.Mock
}

func (m *MockFileRepository) SaveFile(ctx context.Context, id, name, hash, location string) error {
	args := m.Called(ctx, id, name, hash, location)
	return args.Error(0)
}

func (m *MockFileRepository) GetFileByID(ctx context.Context, id string) (string, string, error) {
	args := m.Called(ctx, id)
	return args.String(0), args.String(1), args.Error(2)
}

func (m *MockFileRepository) GetFileByHash(ctx context.Context, hash string) (string, error) {
	args := m.Called(ctx, hash)
	return args.String(0), args.Error(1)
}

// MockFileStorage is a mock implementation of the FileStorage interface
type MockFileStorage struct {
	mock.Mock
}

func (m *MockFileStorage) SaveFile(ctx context.Context, location string, content []byte) error {
	args := m.Called(ctx, location, content)
	return args.Error(0)
}

func (m *MockFileStorage) GetFile(ctx context.Context, location string) ([]byte, error) {
	args := m.Called(ctx, location)
	return args.Get(0).([]byte), args.Error(1)
}

func TestFileService_UploadFile_NewFile(t *testing.T) {
	// Create mocks
	mockRepo := new(MockFileRepository)
	mockStorage := new(MockFileStorage)

	// Create service
	fileService := service.NewFileService(mockRepo, mockStorage)

	// Test data
	ctx := context.Background()
	fileName := "test.txt"
	content := []byte("test content")

	// Set up expectations
	mockRepo.On("GetFileByHash", ctx, mock.AnythingOfType("string")).Return("", nil)
	mockStorage.On("SaveFile", ctx, mock.AnythingOfType("string"), content).Return(nil)
	mockRepo.On("SaveFile", ctx, mock.AnythingOfType("string"), fileName, mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(nil)

	// Call the method
	fileID, err := fileService.UploadFile(ctx, fileName, content)

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, fileID)
	mockRepo.AssertExpectations(t)
	mockStorage.AssertExpectations(t)
}

func TestFileService_UploadFile_ExistingFile(t *testing.T) {
	// Create mocks
	mockRepo := new(MockFileRepository)
	mockStorage := new(MockFileStorage)

	// Create service
	fileService := service.NewFileService(mockRepo, mockStorage)

	// Test data
	ctx := context.Background()
	fileName := "test.txt"
	content := []byte("test content")
	existingFileID := "existing-file-id"

	// Set up expectations
	mockRepo.On("GetFileByHash", ctx, mock.AnythingOfType("string")).Return(existingFileID, nil)

	// Call the method
	fileID, err := fileService.UploadFile(ctx, fileName, content)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, existingFileID, fileID)
	mockRepo.AssertExpectations(t)
	mockStorage.AssertNotCalled(t, "SaveFile")
	mockRepo.AssertNotCalled(t, "SaveFile")
}

func TestFileService_GetFile(t *testing.T) {
	// Create mocks
	mockRepo := new(MockFileRepository)
	mockStorage := new(MockFileStorage)

	// Create service
	fileService := service.NewFileService(mockRepo, mockStorage)

	// Test data
	ctx := context.Background()
	fileID := "test-file-id"
	fileName := "test.txt"
	location := "test-location"
	content := []byte("test content")

	// Set up expectations
	mockRepo.On("GetFileByID", ctx, fileID).Return(fileName, location, nil)
	mockStorage.On("GetFile", ctx, location).Return(content, nil)

	// Call the method
	gotFileName, gotContent, err := fileService.GetFile(ctx, fileID)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, fileName, gotFileName)
	assert.Equal(t, content, gotContent)
	mockRepo.AssertExpectations(t)
	mockStorage.AssertExpectations(t)
}
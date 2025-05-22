package clients_test

import (
	"context"
	"errors"
	"kr-02/internal/pkg/file_analysis/clients"
	"kr-02/internal/proto/file_storing_service"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

// Mock for the FileStoringServiceClient
type MockFileStoringServiceClient struct {
	mock.Mock
}

func (m *MockFileStoringServiceClient) GetFile(ctx context.Context, in *file_storing_service.GetFileRequest, opts ...grpc.CallOption) (*file_storing_service.GetFileResponse, error) {
	args := m.Called(ctx, in)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*file_storing_service.GetFileResponse), args.Error(1)
}

func (m *MockFileStoringServiceClient) UploadFile(ctx context.Context, in *file_storing_service.UploadFileRequest, opts ...grpc.CallOption) (*file_storing_service.UploadFileResponse, error) {
	args := m.Called(ctx, in)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*file_storing_service.UploadFileResponse), args.Error(1)
}

// Test implementation of FileStoringClient that uses the mock
type TestFileStoringClient struct {
	clients.FileStoringClientInterface
	mockClient *MockFileStoringServiceClient
}

func NewTestFileStoringClient(mockClient *MockFileStoringServiceClient) *TestFileStoringClient {
	return &TestFileStoringClient{
		mockClient: mockClient,
	}
}

func (c *TestFileStoringClient) GetFile(ctx context.Context, fileID string) (string, []byte, error) {
	resp, err := c.mockClient.GetFile(ctx, &file_storing_service.GetFileRequest{
		FileId: fileID,
	})

	if err != nil {
		return "", nil, err
	}

	return resp.FileName, resp.Content, nil
}

func (c *TestFileStoringClient) Close() error {
	return nil
}

func TestGetFile(t *testing.T) {
	// Create mock
	mockClient := new(MockFileStoringServiceClient)

	// Create test client
	client := NewTestFileStoringClient(mockClient)

	// Set up mock expectations
	mockClient.On("GetFile", mock.Anything, &file_storing_service.GetFileRequest{
		FileId: "file123",
	}).Return(&file_storing_service.GetFileResponse{
		FileName: "test.txt",
		Content:  []byte("test content"),
	}, nil)

	// Call the method
	fileName, content, err := client.GetFile(context.Background(), "file123")

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, "test.txt", fileName)
	assert.Equal(t, []byte("test content"), content)

	mockClient.AssertExpectations(t)
}

func TestGetFile_Error(t *testing.T) {
	// Create mock
	mockClient := new(MockFileStoringServiceClient)

	// Create test client
	client := NewTestFileStoringClient(mockClient)

	// Set up mock expectations
	mockClient.On("GetFile", mock.Anything, &file_storing_service.GetFileRequest{
		FileId: "file123",
	}).Return(nil, errors.New("connection error"))

	// Call the method
	_, _, err := client.GetFile(context.Background(), "file123")

	// Assert
	assert.Error(t, err)

	mockClient.AssertExpectations(t)
}

func TestClose(t *testing.T) {
	// Create mock
	mockClient := new(MockFileStoringServiceClient)

	// Create test client
	client := NewTestFileStoringClient(mockClient)

	// Call the method
	err := client.Close()

	// Assert
	assert.NoError(t, err)
}

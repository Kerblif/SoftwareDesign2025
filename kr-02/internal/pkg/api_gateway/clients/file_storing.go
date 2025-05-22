package clients

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	pb "kr-02/internal/proto/file_storing_service"
)

// FileStoringClient provides methods for interacting with the File Storing Service
type FileStoringClient struct {
	client pb.FileStoringServiceClient
	conn   *grpc.ClientConn
}

// NewFileStoringClient creates a new FileStoringClient instance
func NewFileStoringClient(address string) (*FileStoringClient, error) {
	// Set up a connection to the server with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to File Storing Service: %w", err)
	}

	client := pb.NewFileStoringServiceClient(conn)

	return &FileStoringClient{
		client: client,
		conn:   conn,
	}, nil
}

// Close closes the client connection
func (c *FileStoringClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// UploadFile uploads a file to the File Storing Service
func (c *FileStoringClient) UploadFile(ctx context.Context, fileName string, content []byte) (string, error) {
	// Set a timeout for the request
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Define retry parameters
	maxRetries := 3
	retryDelay := 1 * time.Second

	var resp *pb.UploadFileResponse
	var err error

	// Try the request with retries
	for attempt := 0; attempt < maxRetries; attempt++ {
		resp, err = c.client.UploadFile(ctx, &pb.UploadFileRequest{
			FileName: fileName,
			Content:  content,
		})

		// If successful or not a retryable error, break
		if err == nil {
			break
		}

		// Check if the error is retryable
		s, ok := status.FromError(err)
		if !ok || (s.Code() != codes.Unavailable && s.Code() != codes.DeadlineExceeded) {
			return "", fmt.Errorf("failed to upload file: %w", err)
		}

		// If this was the last attempt, return the error
		if attempt == maxRetries-1 {
			return "", fmt.Errorf("failed to upload file after %d attempts: %w", maxRetries, err)
		}

		// Wait before retrying
		time.Sleep(retryDelay)
		// Increase delay for next retry (exponential backoff)
		retryDelay *= 2
	}

	return resp.FileId, nil
}

// GetFile retrieves a file from the File Storing Service
func (c *FileStoringClient) GetFile(ctx context.Context, fileID string) (string, []byte, error) {
	// Set a timeout for the request
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Define retry parameters
	maxRetries := 3
	retryDelay := 1 * time.Second

	var resp *pb.GetFileResponse
	var err error

	// Try the request with retries
	for attempt := 0; attempt < maxRetries; attempt++ {
		resp, err = c.client.GetFile(ctx, &pb.GetFileRequest{
			FileId: fileID,
		})

		// If successful or not a retryable error, break
		if err == nil {
			break
		}

		// Check if the error is retryable
		s, ok := status.FromError(err)
		if !ok || (s.Code() != codes.Unavailable && s.Code() != codes.DeadlineExceeded) {
			return "", nil, fmt.Errorf("failed to get file: %w", err)
		}

		// If this was the last attempt, return the error
		if attempt == maxRetries-1 {
			return "", nil, fmt.Errorf("failed to get file after %d attempts: %w", maxRetries, err)
		}

		// Wait before retrying
		time.Sleep(retryDelay)
		// Increase delay for next retry (exponential backoff)
		retryDelay *= 2
	}

	return resp.FileName, resp.Content, nil
}

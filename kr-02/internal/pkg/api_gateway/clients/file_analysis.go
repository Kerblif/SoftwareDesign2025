package clients

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "kr-02/internal/proto/file_analysis_service"
)

// FileAnalysisClient provides methods for interacting with the File Analysis Service
type FileAnalysisClient struct {
	client pb.FileAnalysisServiceClient
	conn   *grpc.ClientConn
}

// NewFileAnalysisClient creates a new FileAnalysisClient instance
func NewFileAnalysisClient(address string) (*FileAnalysisClient, error) {
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
		return nil, fmt.Errorf("failed to connect to File Analysis Service: %w", err)
	}
	
	client := pb.NewFileAnalysisServiceClient(conn)
	
	return &FileAnalysisClient{
		client: client,
		conn:   conn,
	}, nil
}

// Close closes the client connection
func (c *FileAnalysisClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// AnalyzeFile sends a request to analyze a file
func (c *FileAnalysisClient) AnalyzeFile(ctx context.Context, fileID string, generateWordCloud bool) (*pb.AnalyzeFileResponse, error) {
	// Set a timeout for the request
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second) // Analysis might take longer
	defer cancel()
	
	// Make the request
	resp, err := c.client.AnalyzeFile(ctx, &pb.AnalyzeFileRequest{
		FileId:            fileID,
		GenerateWordCloud: generateWordCloud,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to analyze file: %w", err)
	}
	
	return resp, nil
}

// GetWordCloud retrieves a word cloud image
func (c *FileAnalysisClient) GetWordCloud(ctx context.Context, location string) ([]byte, error) {
	// Set a timeout for the request
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	
	// Make the request
	resp, err := c.client.GetWordCloud(ctx, &pb.GetWordCloudRequest{
		Location: location,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get word cloud: %w", err)
	}
	
	return resp.Image, nil
}
package server

import (
	"context"
	"log"

	pb "kr-02/internal/proto/api_gateway"
	fs_pb "kr-02/internal/proto/file_storing_service"
	fa_pb "kr-02/internal/proto/file_analysis_service"
	"kr-02/internal/pkg/api_gateway/clients"
)

// Server implements the APIGatewayServer interface
type Server struct {
	pb.UnimplementedAPIGatewayServer
	fileStoringClient *clients.FileStoringClient
	fileAnalysisClient *clients.FileAnalysisClient
}

// NewServer creates a new Server instance
func NewServer(fileStoringClient *clients.FileStoringClient, fileAnalysisClient *clients.FileAnalysisClient) *Server {
	return &Server{
		fileStoringClient: fileStoringClient,
		fileAnalysisClient: fileAnalysisClient,
	}
}

// UploadFile handles file upload requests
func (s *Server) UploadFile(ctx context.Context, req *fs_pb.UploadFileRequest) (*fs_pb.UploadFileResponse, error) {
	log.Printf("API Gateway: Received upload request for file: %s", req.FileName)
	
	fileID, err := s.fileStoringClient.UploadFile(ctx, req.FileName, req.Content)
	if err != nil {
		log.Printf("API Gateway: Failed to upload file: %v", err)
		return nil, err
	}
	
	log.Printf("API Gateway: File uploaded successfully with ID: %s", fileID)
	return &fs_pb.UploadFileResponse{
		FileId: fileID,
	}, nil
}

// GetFile handles file retrieval requests
func (s *Server) GetFile(ctx context.Context, req *fs_pb.GetFileRequest) (*fs_pb.GetFileResponse, error) {
	log.Printf("API Gateway: Received get file request for ID: %s", req.FileId)
	
	fileName, content, err := s.fileStoringClient.GetFile(ctx, req.FileId)
	if err != nil {
		log.Printf("API Gateway: Failed to get file: %v", err)
		return nil, err
	}
	
	log.Printf("API Gateway: File retrieved successfully: %s", fileName)
	return &fs_pb.GetFileResponse{
		FileName: fileName,
		Content:  content,
	}, nil
}

// AnalyzeFile handles file analysis requests
func (s *Server) AnalyzeFile(ctx context.Context, req *fa_pb.AnalyzeFileRequest) (*fa_pb.AnalyzeFileResponse, error) {
	log.Printf("API Gateway: Received analysis request for file ID: %s", req.FileId)
	
	resp, err := s.fileAnalysisClient.AnalyzeFile(ctx, req.FileId, req.GenerateWordCloud)
	if err != nil {
		log.Printf("API Gateway: Failed to analyze file: %v", err)
		return nil, err
	}
	
	log.Printf("API Gateway: File analyzed successfully: %s", req.FileId)
	return resp, nil
}

// GetWordCloud handles word cloud retrieval requests
func (s *Server) GetWordCloud(ctx context.Context, req *fa_pb.GetWordCloudRequest) (*fa_pb.GetWordCloudResponse, error) {
	log.Printf("API Gateway: Received word cloud request for location: %s", req.Location)
	
	image, err := s.fileAnalysisClient.GetWordCloud(ctx, req.Location)
	if err != nil {
		log.Printf("API Gateway: Failed to get word cloud: %v", err)
		return nil, err
	}
	
	log.Printf("API Gateway: Word cloud retrieved successfully: %s", req.Location)
	return &fa_pb.GetWordCloudResponse{
		Image: image,
	}, nil
}
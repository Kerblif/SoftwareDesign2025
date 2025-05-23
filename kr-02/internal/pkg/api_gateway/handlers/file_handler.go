package handlers

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	"kr-02/internal/pkg/api_gateway/clients"
)

// FileHandler handles file-related operations
type FileHandler struct {
	client *clients.FileStoringClient
}

// NewFileHandler creates a new FileHandler instance
func NewFileHandler(client *clients.FileStoringClient) *FileHandler {
	return &FileHandler{client: client}
}

// UploadFile godoc
// @Summary Upload a file
// @Description Upload a file to the storage
// @Tags files
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File to upload"
// @Success 200 {object} map[string]string "Returns the file ID"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/files [post]
func (h *FileHandler) UploadFile(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file provided"})
		return
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}

	fileID, err := h.client.UploadFile(c.Request.Context(), header.Filename, content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"file_id": fileID})
}

// GetFile godoc
// @Summary Get a file
// @Description Get a file by its ID
// @Tags files
// @Produce octet-stream
// @Param file_id path string true "File ID"
// @Success 200 {file} binary "File content"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 404 {object} map[string]string "File not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/files/{file_id} [get]
func (h *FileHandler) GetFile(c *gin.Context) {
	fileID := c.Param("file_id")
	if fileID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File ID is required"})
		return
	}

	fileName, content, err := h.client.GetFile(c.Request.Context(), fileID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Data(http.StatusOK, "application/octet-stream", content)
}

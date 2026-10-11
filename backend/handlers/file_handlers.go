package handlers

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type FileUploadResponse struct {
	ID       string `json:"id"`
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
	URL      string `json:"url"`
	Type     string `json:"type"`
}

type FileHandler struct {
	uploadDir     string
	maxFileSize   int64
	allowedTypes  []string
}

// NewFileHandler creates a new file handler
func NewFileHandler() *FileHandler {
	return &FileHandler{
		uploadDir:    "./uploads",
		maxFileSize:  10 << 20, // 10MB
		allowedTypes: []string{".jpg", ".jpeg", ".png", ".gif", ".pdf", ".doc", ".docx", ".txt", ".csv"},
	}
}

// UploadFile handles file uploads
func (h *FileHandler) UploadFile(c *gin.Context) {
	// Check file size
	if c.Request.ContentLength > h.maxFileSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File too large"})
		return
	}

	// Parse multipart form
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}
	defer file.Close()

	// Validate file type
	if !h.isAllowedFileType(header.Filename) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File type not allowed"})
		return
	}

	// Get upload context
	uploadType := c.PostForm("type") // image, document, receipt
	entityType := c.PostForm("entity_type") // campaign, milestone, expense
	entityID := c.PostForm("entity_id")

	// Generate unique filename
	ext := filepath.Ext(header.Filename)
	filename := fmt.Sprintf("%s_%d%s", uuid.New().String(), time.Now().Unix(), ext)

	// Create upload directory if it doesn't exist
	uploadPath := filepath.Join(h.uploadDir, uploadType, entityType)
	if err := os.MkdirAll(uploadPath, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create upload directory"})
		return
	}

	// Save file
	fullPath := filepath.Join(uploadPath, filename)
	dst, err := os.Create(fullPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save file"})
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not write file"})
		return
	}

	// Create response
	response := FileUploadResponse{
		ID:       filename,
		Filename: header.Filename,
		Size:     header.Size,
		URL:      fmt.Sprintf("/api/files/%s/%s/%s", uploadType, entityType, filename),
		Type:     uploadType,
	}

	// Here you would typically save file metadata to database
	// Example: h.saveFileMetadata(response, entityType, entityID)

	c.JSON(http.StatusOK, response)
}

// UploadMultipleFiles handles multiple file uploads
func (h *FileHandler) UploadMultipleFiles(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not parse form"})
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No files uploaded"})
		return
	}

	var responses []FileUploadResponse
	uploadType := c.PostForm("type")
	entityType := c.PostForm("entity_type")

	for _, fileHeader := range files {
		if fileHeader.Size > h.maxFileSize {
			continue // Skip files that are too large
		}

		if !h.isAllowedFileType(fileHeader.Filename) {
			continue // Skip disallowed file types
		}

		file, err := fileHeader.Open()
		if err != nil {
			continue // Skip files that can't be opened
		}

		// Generate unique filename
		ext := filepath.Ext(fileHeader.Filename)
		filename := fmt.Sprintf("%s_%d%s", uuid.New().String(), time.Now().Unix(), ext)

		// Create upload directory
		uploadPath := filepath.Join(h.uploadDir, uploadType, entityType)
		if err := os.MkdirAll(uploadPath, 0755); err != nil {
			file.Close()
			continue
		}

		// Save file
		fullPath := filepath.Join(uploadPath, filename)
		dst, err := os.Create(fullPath)
		if err != nil {
			file.Close()
			continue
		}

		_, err = io.Copy(dst, file)
		file.Close()
		dst.Close()

		if err != nil {
			continue
		}

		response := FileUploadResponse{
			ID:       filename,
			Filename: fileHeader.Filename,
			Size:     fileHeader.Size,
			URL:      fmt.Sprintf("/api/files/%s/%s/%s", uploadType, entityType, filename),
			Type:     uploadType,
		}

		responses = append(responses, response)
	}

	c.JSON(http.StatusOK, gin.H{"files": responses})
}

// ServeFile serves uploaded files
func (h *FileHandler) ServeFile(c *gin.Context) {
	uploadType := c.Param("type")
	entityType := c.Param("entity_type")
	filename := c.Param("filename")

	filePath := filepath.Join(h.uploadDir, uploadType, entityType, filename)

	// Security check - ensure file path is within upload directory
	absUploadDir, _ := filepath.Abs(h.uploadDir)
	absFilePath, _ := filepath.Abs(filePath)
	if !strings.HasPrefix(absFilePath, absUploadDir) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	c.File(filePath)
}

// DeleteFile removes a file
func (h *FileHandler) DeleteFile(c *gin.Context) {
	uploadType := c.Param("type")
	entityType := c.Param("entity_type")
	filename := c.Param("filename")

	filePath := filepath.Join(h.uploadDir, uploadType, entityType, filename)

	// Security check
	absUploadDir, _ := filepath.Abs(h.uploadDir)
	absFilePath, _ := filepath.Abs(filePath)
	if !strings.HasPrefix(absFilePath, absUploadDir) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	if err := os.Remove(filePath); err != nil {
		if os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete file"})
		return
	}

	// Here you would typically remove file metadata from database

	c.JSON(http.StatusOK, gin.H{"message": "File deleted successfully"})
}

// GetFileInfo returns file metadata
func (h *FileHandler) GetFileInfo(c *gin.Context) {
	uploadType := c.Param("type")
	entityType := c.Param("entity_type")
	filename := c.Param("filename")

	filePath := filepath.Join(h.uploadDir, uploadType, entityType, filename)

	info, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get file info"})
		return
	}

	response := gin.H{
		"filename":    filename,
		"size":        info.Size(),
		"modified_at": info.ModTime(),
		"url":         fmt.Sprintf("/api/files/%s/%s/%s", uploadType, entityType, filename),
	}

	c.JSON(http.StatusOK, response)
}

// isAllowedFileType checks if file extension is allowed
func (h *FileHandler) isAllowedFileType(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	for _, allowedExt := range h.allowedTypes {
		if ext == allowedExt {
			return true
		}
	}
	return false
}

// ValidateImageFile validates image files specifically
func (h *FileHandler) ValidateImageFile(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}
	defer file.Close()

	// Check file size (max 5MB for images)
	if header.Size > 5<<20 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image file too large (max 5MB)"})
		return
	}

	// Check file type
	ext := strings.ToLower(filepath.Ext(header.Filename))
	imageTypes := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
	allowed := false
	for _, allowedExt := range imageTypes {
		if ext == allowedExt {
			allowed = true
			break
		}
	}

	if !allowed {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image format"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"valid":    true,
		"filename": header.Filename,
		"size":     header.Size,
		"type":     ext,
	})
}

// GetUploadProgress returns upload progress for large files
func (h *FileHandler) GetUploadProgress(c *gin.Context) {
	uploadID := c.Param("upload_id")
	
	// This would typically track progress in Redis or database
	// For now, return mock progress
	progress := gin.H{
		"upload_id":       uploadID,
		"progress":        85,
		"bytes_uploaded":  "8.5MB",
		"total_bytes":     "10MB",
		"status":          "uploading",
	}

	c.JSON(http.StatusOK, progress)
}
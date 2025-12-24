package file

import (
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service   *FileService
	uploadDir string
}

func NewHandler(service *FileService, uploadDir string) *Handler {
	return &Handler{
		service:   service,
		uploadDir: uploadDir,
	}
}

// @Summary Upload a file
// @Description Upload a file with optional TTL (in hours). Returns a download URL.
// @Tags files
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File to upload"
// @Param ttl query int false "Time-to-live in hours (max 24, default 24)"
// @Success 200 {object} UploadResponse
// @Failure 400 {object} ErrorResponse "Invalid input (e.g. file missing, size too large)"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /upload [post]
func (h *Handler) UploadHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Error: "file is required"})
		return
	}
	ttlStr := c.Query("ttl")
	ttlHours := 0
	if ttlStr != "" {
		if t, err := strconv.Atoi(ttlStr); err == nil {
			ttlHours = t
		}
	}
	src, err := file.Open()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to open file"})
		return
	}
	defer src.Close()
	id, err := h.service.Upload(c.Request.Context(), file.Filename, src, file.Size, ttlHours)
	if err != nil {
		if errors.Is(err, ErrInvalidFilename) || errors.Is(err, ErrInvalidFileSize) {
			c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Error: fmt.Sprintf("%v or %v", ErrInvalidFilename, ErrInvalidFileSize)})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Error: "upload failed"})
		return
	}
	c.JSON(http.StatusOK, UploadResponse{URL: "/f/" + id})
}

// @Summary Download a file
// @Description Download a file using its unique ID. File must not be expired.
// @Tags files
// @Produce octet-stream
// @Param id path string true "File ID"
// @Success 200 {file} file "File content"
// @Failure 404 {object} ErrorResponse "File not found or expired"
// @Failure 500 {object} ErrorResponse "Internal server error"
// @Router /f/{id} [get]
func (h *Handler) DownloadHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.AbortWithStatusJSON(http.StatusNotFound, ErrorResponse{Error: "invalid file id"})
		return
	}
	path, err := h.service.Download(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, ErrFileNotFound) || errors.Is(err, ErrFileExpired) {
			c.AbortWithStatusJSON(http.StatusNotFound, ErrorResponse{Error: fmt.Sprintf("%v or %v", ErrFileNotFound, ErrFileExpired)})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Error: "download failed"})
		return
	}
	fullPath := filepath.Join(h.uploadDir, path)
	c.File(fullPath)
}

package helper

import (
	"errors"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ValidateUploadedFile(c *gin.Context, file *multipart.FileHeader, maxSize int64, allowedMimeTypes map[string]bool) bool {
	if file.Size > maxSize {
		ErrorResponse(c, http.StatusRequestEntityTooLarge, "File size exceeds the allowed limit", nil)
		return false
	}

	src, err := file.Open()
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to open uploaded file", err)
		return false
	}
	defer src.Close()

	buffer := make([]byte, 512)
	_, err = src.Read(buffer)
	if err != nil {
		ErrorResponse(c, http.StatusInternalServerError, "Failed to read uploaded file", err)
		return false
	}

	contentType := http.DetectContentType(buffer)
	if !allowedMimeTypes[contentType] {
		ErrorResponse(c, http.StatusBadRequest, "Invalid file type", errors.New("unsupported file type"))
		return false
	}

	return true
}

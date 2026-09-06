package services

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/deepakgudla/bookvault/internal/interfaces"
	"github.com/google/uuid"
)

var _ UploadServiceInterface = (*UploadService)(nil)

// UploadService validates and stores uploaded product images.
type UploadService struct {
	provider    interfaces.UploadProvider
	maxFileSize int64
}

// NewUploadService creates an upload service.
func NewUploadService(provider interfaces.UploadProvider, maxFileSize int64) *UploadService {
	return &UploadService{provider: provider, maxFileSize: maxFileSize}
}

// UploadProductImage validates and stores a product image.
func (s *UploadService) UploadProductImage(productID uint, file *multipart.FileHeader) (string, error) {
	if s.maxFileSize > 0 && file.Size > s.maxFileSize {
		return "", fmt.Errorf("file exceeds maximum size of %d bytes", s.maxFileSize)
	}

	ext := strings.ToLower(filepath.Ext(file.Filename)) // image extensions ()jpg, jpeg etc...
	if !isValidImageExt(ext) {
		return "", fmt.Errorf("invalid file type: %s", ext)
	}

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer func() { _ = src.Close() }()

	buffer := make([]byte, 512)
	read, err := src.Read(buffer)
	if err != nil && read == 0 {
		return "", err
	}
	contentType := http.DetectContentType(buffer[:read])
	if !isValidImageContentType(contentType, ext) {
		return "", fmt.Errorf("file content does not match extension")
	}

	newFileName := uuid.New().String()

	path := fmt.Sprintf("products/%d/%s%s", productID, newFileName, ext)

	return s.provider.UploadFile(file, path)
}

func isValidImageContentType(contentType, ext string) bool {
	allowed := map[string]string{
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".gif":  "image/gif",
		".webp": "image/webp",
	}
	return allowed[ext] == contentType
}

func isValidImageExt(ext string) bool {
	validExt := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
	for _, validExt := range validExt {
		if ext == validExt {
			return true
		}
	}

	return false
}

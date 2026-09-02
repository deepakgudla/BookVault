package providers

import (
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
)

// LocalUploadProvider stores uploaded files on the local filesystem.
type LocalUploadProvider struct {
	basePath string
}

// NewLocalUploadProvider creates a local filesystem upload provider.
func NewLocalUploadProvider(basePath string) *LocalUploadProvider {
	return &LocalUploadProvider{basePath: basePath}
}

// UploadFile stores a multipart file at the requested path.
func (p *LocalUploadProvider) UploadFile(file *multipart.FileHeader, path string) (string, error) {
	fullPath := filepath.Join(p.basePath, path)

	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		return "", err
	}

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer func() {
		if err := src.Close(); err != nil {
			log.Printf("failed to close source file: %v", err)
		}
	}()

	dst, err := os.Create(fullPath)
	if err != nil {
		return "", err
	}
	defer func() {
		if err := dst.Close(); err != nil {
			log.Printf("failed to close destination file: %v", err)
		}
	}()

	if _, err := dst.ReadFrom(src); err != nil {
		return "", err
	}

	return fmt.Sprintf("/uploads/%s", path), nil
}

// DeleteFile removes an uploaded file from the local filesystem.
func (p *LocalUploadProvider) DeleteFile(path string) error {
	fullPath := filepath.Join(p.basePath, path)
	return os.Remove(fullPath)
}

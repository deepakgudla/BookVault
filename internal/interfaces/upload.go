package interfaces

import "mime/multipart"

// UploadProvider uploads and deletes files in a storage backend.
type UploadProvider interface {
	UploadFile(file *multipart.FileHeader, path string) (string, error)
	DeleteFile(path string) error
}

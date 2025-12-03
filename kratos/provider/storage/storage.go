// Package storage provides object storage abstraction for various providers.
package storage

import (
	"context"
	"io"
)

// Storage defines the interface for object storage operations.
// It abstracts different storage providers (S3, OSS, COS, MinIO) behind a common interface.
type Storage interface {
	// PutObject uploads an object to storage.
	PutObject(ctx context.Context, key string, data []byte) error

	// PutObjectFromReader uploads an object from a reader.
	PutObjectFromReader(ctx context.Context, key string, reader io.Reader, size int64) error

	// GetObject retrieves an object from storage.
	GetObject(ctx context.Context, key string) ([]byte, error)

	// GetObjectReader returns a reader for the object.
	GetObjectReader(ctx context.Context, key string) (io.ReadCloser, error)

	// DeleteObject deletes an object from storage.
	DeleteObject(ctx context.Context, key string) error

	// Exists checks if an object exists in storage.
	Exists(ctx context.Context, key string) (bool, error)

	// GetObjectURL returns a URL for accessing the object (optional, may return empty string).
	GetObjectURL(ctx context.Context, key string, expiresIn int64) (string, error)
}

// NoOpStorage is a no-op implementation of Storage interface.
// It can be used when object storage is disabled.
type NoOpStorage struct{}

func (n *NoOpStorage) PutObject(ctx context.Context, key string, data []byte) error {
	return nil
}

func (n *NoOpStorage) PutObjectFromReader(ctx context.Context, key string, reader io.Reader, size int64) error {
	return nil
}

func (n *NoOpStorage) GetObject(ctx context.Context, key string) ([]byte, error) {
	return nil, ErrObjectNotFound
}

func (n *NoOpStorage) GetObjectReader(ctx context.Context, key string) (io.ReadCloser, error) {
	return nil, ErrObjectNotFound
}

func (n *NoOpStorage) DeleteObject(ctx context.Context, key string) error {
	return nil
}

func (n *NoOpStorage) Exists(ctx context.Context, key string) (bool, error) {
	return false, nil
}

func (n *NoOpStorage) GetObjectURL(ctx context.Context, key string, expiresIn int64) (string, error) {
	return "", nil
}

// Storage operation errors.
var (
	ErrObjectNotFound = &StorageError{Message: "object not found"}
	ErrInvalidConfig = &StorageError{Message: "invalid storage configuration"}
)

// StorageError represents a storage operation error.
type StorageError struct {
	Message string
	Err     error
}

func (e *StorageError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

func (e *StorageError) Unwrap() error {
	return e.Err
}


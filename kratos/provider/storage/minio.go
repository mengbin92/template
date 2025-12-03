// Package storage provides MinIO implementation for object storage.
package storage

import (
	"bytes"
	"context"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/pkg/errors"
)

// MinIOStorage implements Storage interface using MinIO client.
type MinIOStorage struct {
	client     *minio.Client
	bucketName string
	pathPrefix string
}

// NewMinIOStorage creates a new MinIO storage instance.
func NewMinIOStorage(endpoint, accessKeyID, secretAccessKey, bucketName, region string, useSSL bool, pathPrefix string) (*MinIOStorage, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
		Region: region,
	})
	if err != nil {
		return nil, errors.Wrap(err, "create minio client")
	}

	storage := &MinIOStorage{
		client:     client,
		bucketName: bucketName,
		pathPrefix: pathPrefix,
	}

	ctx := context.Background()
	exists, err := client.BucketExists(ctx, bucketName)
	if err != nil {
		return nil, errors.Wrapf(err, "check bucket existence: %s", bucketName)
	}

	if !exists {
		err = client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: region})
		if err != nil {
			return nil, errors.Wrapf(err, "create bucket: %s", bucketName)
		}
	}

	return storage, nil
}

func (m *MinIOStorage) buildKey(key string) string {
	if m.pathPrefix == "" {
		return key
	}
	if m.pathPrefix[len(m.pathPrefix)-1] == '/' {
		return m.pathPrefix + key
	}
	return m.pathPrefix + "/" + key
}

func (m *MinIOStorage) PutObject(ctx context.Context, key string, data []byte) error {
	key = m.buildKey(key)
	_, err := m.client.PutObject(ctx, m.bucketName, key, bytes.NewReader(data), int64(len(data)), minio.PutObjectOptions{})
	return errors.Wrapf(err, "put object: %s", key)
}

func (m *MinIOStorage) PutObjectFromReader(ctx context.Context, key string, reader io.Reader, size int64) error {
	key = m.buildKey(key)
	_, err := m.client.PutObject(ctx, m.bucketName, key, reader, size, minio.PutObjectOptions{})
	return errors.Wrapf(err, "put object from reader: %s", key)
}

func (m *MinIOStorage) GetObject(ctx context.Context, key string) ([]byte, error) {
	key = m.buildKey(key)
	obj, err := m.client.GetObject(ctx, m.bucketName, key, minio.GetObjectOptions{})
	if err != nil {
		if minioErr, ok := err.(minio.ErrorResponse); ok && minioErr.Code == "NoSuchKey" {
			return nil, ErrObjectNotFound
		}
		return nil, errors.Wrapf(err, "get object: %s", key)
	}
	defer obj.Close()

	data, err := io.ReadAll(obj)
	if err != nil {
		return nil, errors.Wrapf(err, "read object: %s", key)
	}
	return data, nil
}

func (m *MinIOStorage) GetObjectReader(ctx context.Context, key string) (io.ReadCloser, error) {
	key = m.buildKey(key)
	obj, err := m.client.GetObject(ctx, m.bucketName, key, minio.GetObjectOptions{})
	if err != nil {
		if minioErr, ok := err.(minio.ErrorResponse); ok && minioErr.Code == "NoSuchKey" {
			return nil, ErrObjectNotFound
		}
		return nil, errors.Wrapf(err, "get object reader: %s", key)
	}
	return obj, nil
}

func (m *MinIOStorage) DeleteObject(ctx context.Context, key string) error {
	key = m.buildKey(key)
	err := m.client.RemoveObject(ctx, m.bucketName, key, minio.RemoveObjectOptions{})
	return errors.Wrapf(err, "delete object: %s", key)
}

func (m *MinIOStorage) Exists(ctx context.Context, key string) (bool, error) {
	key = m.buildKey(key)
	_, err := m.client.StatObject(ctx, m.bucketName, key, minio.StatObjectOptions{})
	if err != nil {
		if minio.ToErrorResponse(err).Code == "NoSuchKey" {
			return false, nil
		}
		return false, errors.Wrapf(err, "stat object: %s", key)
	}
	return true, nil
}

func (m *MinIOStorage) GetObjectURL(ctx context.Context, key string, expiresIn int64) (string, error) {
	key = m.buildKey(key)
	expires := time.Duration(expiresIn) * time.Second
	if expiresIn == 0 {
		expires = 7 * 24 * time.Hour // Default 7 days
	}
	url, err := m.client.PresignedGetObject(ctx, m.bucketName, key, expires, nil)
	if err != nil {
		return "", errors.Wrapf(err, "get presigned url: %s", key)
	}
	return url.String(), nil
}


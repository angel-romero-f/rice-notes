// Package storage provides file upload and download services
package storage

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// Uploader defines the interface for file upload operations
type Uploader interface {
	Upload(ctx context.Context, key string, body io.Reader, contentType string, size int64) error
	GetPresignedURL(ctx context.Context, key string, expiration time.Duration) (string, error)
	Delete(ctx context.Context, key string) error
}

// S3Uploader implements Uploader interface using AWS S3
type S3Uploader struct {
	client *s3.Client
	bucket string
}

// NewS3Uploader creates a new S3 uploader instance
func NewS3Uploader(ctx context.Context, bucket string, region string) (*S3Uploader, error) {
	// Load AWS config
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		slog.Error("Failed to load AWS config", "error", err)
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	client := s3.NewFromConfig(cfg)

	// Note: Skipping HeadBucket check as it requires additional permissions
	// The bucket access will be validated on first upload operation

	slog.Info("S3 uploader initialized successfully", "bucket", bucket, "region", region)

	return &S3Uploader{
		client: client,
		bucket: bucket,
	}, nil
}

// Upload uploads a file to S3
func (s *S3Uploader) Upload(ctx context.Context, key string, body io.Reader, contentType string, size int64) error {
	slog.Debug("Starting S3 upload", "key", key, "contentType", contentType, "size", size)

	input := &s3.PutObjectInput{
		Bucket:        aws.String(s.bucket),
		Key:           aws.String(key),
		Body:          body,
		ContentType:   aws.String(contentType),
		ContentLength: aws.Int64(size),
		ServerSideEncryption: types.ServerSideEncryptionAes256,
	}

	_, err := s.client.PutObject(ctx, input)
	if err != nil {
		slog.Error("Failed to upload to S3", "error", err, "key", key)
		return fmt.Errorf("failed to upload file to S3: %w", err)
	}

	slog.Info("File uploaded to S3 successfully", "key", key, "bucket", s.bucket)
	return nil
}

// GetPresignedURL generates a presigned URL for downloading a file
func (s *S3Uploader) GetPresignedURL(ctx context.Context, key string, expiration time.Duration) (string, error) {
	slog.Debug("Generating presigned URL", "key", key, "expiration", expiration)

	presignClient := s3.NewPresignClient(s.client)

	request, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = expiration
	})

	if err != nil {
		slog.Error("Failed to generate presigned URL", "error", err, "key", key)
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	slog.Debug("Presigned URL generated successfully", "key", key)
	return request.URL, nil
}

// Delete removes a file from S3
func (s *S3Uploader) Delete(ctx context.Context, key string) error {
	slog.Debug("Deleting file from S3", "key", key)

	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		slog.Error("Failed to delete from S3", "error", err, "key", key)
		return fmt.Errorf("failed to delete file from S3: %w", err)
	}

	slog.Info("File deleted from S3 successfully", "key", key, "bucket", s.bucket)
	return nil
}

// GenerateFileKey creates a structured S3 key for a file
func GenerateFileKey(userEmail, noteID, fileName string) string {
	return fmt.Sprintf("notes/%s/%s/%s", userEmail, noteID, fileName)
}

// MockUploader is a mock implementation of Uploader for testing
type MockUploader struct {
	files map[string][]byte
}

// NewMockUploader creates a new mock uploader for testing
func NewMockUploader() *MockUploader {
	return &MockUploader{
		files: make(map[string][]byte),
	}
}

// Upload simulates file upload by storing in memory
func (m *MockUploader) Upload(ctx context.Context, key string, body io.Reader, contentType string, size int64) error {
	data, err := io.ReadAll(body)
	if err != nil {
		return fmt.Errorf("failed to read body: %w", err)
	}
	
	m.files[key] = data
	slog.Debug("Mock upload successful", "key", key, "size", len(data))
	return nil
}

// GetPresignedURL returns a mock URL
func (m *MockUploader) GetPresignedURL(ctx context.Context, key string, expiration time.Duration) (string, error) {
	if _, exists := m.files[key]; !exists {
		return "", fmt.Errorf("file not found: %s", key)
	}
	return fmt.Sprintf("https://mock-bucket.s3.amazonaws.com/%s?expires=%d", key, time.Now().Add(expiration).Unix()), nil
}

// Delete removes file from mock storage
func (m *MockUploader) Delete(ctx context.Context, key string) error {
	if _, exists := m.files[key]; !exists {
		return fmt.Errorf("file not found: %s", key)
	}
	delete(m.files, key)
	slog.Debug("Mock delete successful", "key", key)
	return nil
}
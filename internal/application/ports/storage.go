package ports

import (
	"context"
	"io"
)

// StorageAdapter is the port that storage backends (S3, SIO, etc.) must implement.
type StorageAdapter interface {
	// Init is called once before any operations. Use it to create clients, etc.
	Init(cfg map[string]string) error

	// Dispose is called once after all operations. Use it to clean up resources.
	Dispose() error

	// Bucket operations
	CreateBucket(ctx context.Context, bucket string) error
	DeleteBucket(ctx context.Context, bucket string) error

	// Object operations
	PutObject(ctx context.Context, bucket, key string, body io.Reader, size int64) error
	GetObject(ctx context.Context, bucket, key string) (io.ReadCloser, error)
	DeleteObject(ctx context.Context, bucket, key string) error
	HeadObject(ctx context.Context, bucket, key string) (ObjectMeta, error)
	ListObjects(ctx context.Context, bucket, prefix string, maxKeys int) ([]ObjectEntry, error)

	// Multipart operations (SIO-specific but defined at port level)
	MultipartPut(ctx context.Context, bucket, key string, body io.Reader, size int64, partSize int64) error

	// Restore (SIO-specific)
	RestoreObject(ctx context.Context, bucket, key string, days int) error
}

// ObjectMeta holds metadata returned from HeadObject.
type ObjectMeta struct {
	ContentLength int64
	ContentType   string
	ETag          string
}

// ObjectEntry holds a single item in a list-objects response.
type ObjectEntry struct {
	Key  string
	Size int64
}

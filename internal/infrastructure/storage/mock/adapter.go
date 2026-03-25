package mock

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"sort"
	"strings"
	"sync"

	"github.com/sine-io/cosbench-go/internal/application/ports"
)

type Adapter struct {
	mu      sync.RWMutex
	buckets map[string]map[string][]byte
}

var _ ports.StorageAdapter = (*Adapter)(nil)

func New() *Adapter {
	return &Adapter{buckets: map[string]map[string][]byte{}}
}

func (a *Adapter) Init(cfg map[string]string) error { return nil }
func (a *Adapter) Dispose() error                   { return nil }

func (a *Adapter) CreateBucket(ctx context.Context, bucket string) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if _, ok := a.buckets[bucket]; !ok {
		a.buckets[bucket] = map[string][]byte{}
	}
	return nil
}

func (a *Adapter) DeleteBucket(ctx context.Context, bucket string) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	delete(a.buckets, bucket)
	return nil
}

func (a *Adapter) PutObject(ctx context.Context, bucket, key string, body io.Reader, size int64) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if _, ok := a.buckets[bucket]; !ok {
		a.buckets[bucket] = map[string][]byte{}
	}
	data, err := io.ReadAll(body)
	if err != nil {
		return err
	}
	a.buckets[bucket][key] = data
	return nil
}

func (a *Adapter) GetObject(ctx context.Context, bucket, key string) (io.ReadCloser, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	obj, ok := a.buckets[bucket][key]
	if !ok {
		return nil, fmt.Errorf("object not found: %s/%s", bucket, key)
	}
	return io.NopCloser(bytes.NewReader(obj)), nil
}

func (a *Adapter) DeleteObject(ctx context.Context, bucket, key string) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if _, ok := a.buckets[bucket]; ok {
		delete(a.buckets[bucket], key)
	}
	return nil
}

func (a *Adapter) HeadObject(ctx context.Context, bucket, key string) (ports.ObjectMeta, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	obj, ok := a.buckets[bucket][key]
	if !ok {
		return ports.ObjectMeta{}, fmt.Errorf("object not found: %s/%s", bucket, key)
	}
	return ports.ObjectMeta{ContentLength: int64(len(obj))}, nil
}

func (a *Adapter) ListObjects(ctx context.Context, bucket, prefix string, maxKeys int) ([]ports.ObjectEntry, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	entries := []ports.ObjectEntry{}
	for key, data := range a.buckets[bucket] {
		if prefix == "" || strings.HasPrefix(key, prefix) {
			entries = append(entries, ports.ObjectEntry{Key: key, Size: int64(len(data))})
		}
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Key < entries[j].Key })
	if maxKeys > 0 && len(entries) > maxKeys {
		entries = entries[:maxKeys]
	}
	return entries, nil
}

func (a *Adapter) MultipartPut(ctx context.Context, bucket, key string, body io.Reader, size int64, partSize int64) error {
	return a.PutObject(ctx, bucket, key, body, size)
}

func (a *Adapter) RestoreObject(ctx context.Context, bucket, key string, days int) error {
	_, err := a.HeadObject(ctx, bucket, key)
	return err
}

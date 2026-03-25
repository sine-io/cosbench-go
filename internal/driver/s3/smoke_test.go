package s3

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	kvconfig "github.com/sine-io/cosbench-go/internal/infrastructure/config"
)

type smokeConfig struct {
	Backend      string
	Endpoint     string
	AccessKey    string
	SecretKey    string
	Region       string
	PathStyleRaw string
	BucketPrefix string
}

func TestLoadLiveConfigRequiresEndpointAndCredentials(t *testing.T) {
	t.Setenv("COSBENCH_SMOKE_ENDPOINT", "")
	t.Setenv("COSBENCH_SMOKE_ACCESS_KEY", "")
	t.Setenv("COSBENCH_SMOKE_SECRET_KEY", "")

	if _, ok := loadSmokeConfigFromEnv(); ok {
		t.Fatal("expected smoke config to be disabled without required env")
	}
}

func TestLoadLiveConfigAppliesDefaults(t *testing.T) {
	t.Setenv("COSBENCH_SMOKE_ENDPOINT", "http://127.0.0.1:9000")
	t.Setenv("COSBENCH_SMOKE_ACCESS_KEY", "ak")
	t.Setenv("COSBENCH_SMOKE_SECRET_KEY", "sk")
	t.Setenv("COSBENCH_SMOKE_BACKEND", "")
	t.Setenv("COSBENCH_SMOKE_REGION", "")
	t.Setenv("COSBENCH_SMOKE_PATH_STYLE", "")
	t.Setenv("COSBENCH_SMOKE_BUCKET_PREFIX", "")

	cfg, ok := loadSmokeConfigFromEnv()
	if !ok {
		t.Fatal("expected smoke config to be enabled")
	}
	if cfg.Backend != "s3" {
		t.Fatalf("backend = %q", cfg.Backend)
	}
	if cfg.Region != "us-east-1" {
		t.Fatalf("region = %q", cfg.Region)
	}
	if cfg.BucketPrefix != "cosbench-go-smoke" {
		t.Fatalf("bucket prefix = %q", cfg.BucketPrefix)
	}
}

func TestSmokeObjectLifecycle(t *testing.T) {
	cfg, ok := loadSmokeConfigFromEnv()
	if !ok {
		t.Skip("set COSBENCH_SMOKE_ENDPOINT, COSBENCH_SMOKE_ACCESS_KEY, and COSBENCH_SMOKE_SECRET_KEY to run live smoke tests")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	adapter := newSmokeAdapter(t, cfg)
	bucket := smokeBucketName(cfg.BucketPrefix)
	key := "smoke-object.txt"
	body := []byte("cosbench-go smoke")

	t.Cleanup(func() {
		_ = adapter.DeleteObject(context.Background(), bucket, key)
		_ = adapter.DeleteObject(context.Background(), bucket, "smoke-multipart.bin")
		_ = adapter.DeleteBucket(context.Background(), bucket)
		_ = adapter.Dispose()
	})

	if err := adapter.CreateBucket(ctx, bucket); err != nil {
		t.Fatalf("CreateBucket(): %v", err)
	}
	if err := adapter.PutObject(ctx, bucket, key, bytes.NewReader(body), int64(len(body))); err != nil {
		t.Fatalf("PutObject(): %v", err)
	}
	rc, err := adapter.GetObject(ctx, bucket, key)
	if err != nil {
		t.Fatalf("GetObject(): %v", err)
	}
	defer rc.Close()
	got, err := io.ReadAll(rc)
	if err != nil {
		t.Fatalf("ReadAll(): %v", err)
	}
	if !bytes.Equal(got, body) {
		t.Fatalf("payload = %q", got)
	}
	meta, err := adapter.HeadObject(ctx, bucket, key)
	if err != nil {
		t.Fatalf("HeadObject(): %v", err)
	}
	if meta.ContentLength != int64(len(body)) {
		t.Fatalf("content length = %d", meta.ContentLength)
	}
	items, err := adapter.ListObjects(ctx, bucket, "smoke-", 1000)
	if err != nil {
		t.Fatalf("ListObjects(): %v", err)
	}
	found := false
	for _, item := range items {
		if item.Key == key {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("uploaded key %q not found in %#v", key, items)
	}
	if err := adapter.DeleteObject(ctx, bucket, key); err != nil {
		t.Fatalf("DeleteObject(): %v", err)
	}
}

func TestSmokeSIOMultipartLifecycle(t *testing.T) {
	cfg, ok := loadSmokeConfigFromEnv()
	if !ok {
		t.Skip("set COSBENCH_SMOKE_ENDPOINT, COSBENCH_SMOKE_ACCESS_KEY, and COSBENCH_SMOKE_SECRET_KEY to run live smoke tests")
	}
	if cfg.Backend != "sio" {
		t.Skip("set COSBENCH_SMOKE_BACKEND=sio to run SIO multipart smoke coverage")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	adapter := newSmokeAdapter(t, cfg)
	bucket := smokeBucketName(cfg.BucketPrefix)
	key := "smoke-multipart.bin"
	payload := bytes.Repeat([]byte("m"), 5*1024*1024+1)

	t.Cleanup(func() {
		_ = adapter.DeleteObject(context.Background(), bucket, key)
		_ = adapter.DeleteBucket(context.Background(), bucket)
		_ = adapter.Dispose()
	})

	if err := adapter.CreateBucket(ctx, bucket); err != nil {
		t.Fatalf("CreateBucket(): %v", err)
	}
	if err := adapter.MultipartPut(ctx, bucket, key, bytes.NewReader(payload), int64(len(payload)), 5*1024*1024); err != nil {
		t.Fatalf("MultipartPut(): %v", err)
	}
	meta, err := adapter.HeadObject(ctx, bucket, key)
	if err != nil {
		t.Fatalf("HeadObject(): %v", err)
	}
	if meta.ContentLength != int64(len(payload)) {
		t.Fatalf("content length = %d", meta.ContentLength)
	}
}

func loadSmokeConfigFromEnv() (smokeConfig, bool) {
	cfg := smokeConfig{
		Backend:      strings.TrimSpace(getenv("COSBENCH_SMOKE_BACKEND", "s3")),
		Endpoint:     strings.TrimSpace(getenv("COSBENCH_SMOKE_ENDPOINT", "")),
		AccessKey:    strings.TrimSpace(getenv("COSBENCH_SMOKE_ACCESS_KEY", "")),
		SecretKey:    strings.TrimSpace(getenv("COSBENCH_SMOKE_SECRET_KEY", "")),
		Region:       strings.TrimSpace(getenv("COSBENCH_SMOKE_REGION", "us-east-1")),
		PathStyleRaw: strings.TrimSpace(getenv("COSBENCH_SMOKE_PATH_STYLE", "")),
		BucketPrefix: strings.TrimSpace(getenv("COSBENCH_SMOKE_BUCKET_PREFIX", "cosbench-go-smoke")),
	}
	if cfg.Endpoint == "" || cfg.AccessKey == "" || cfg.SecretKey == "" {
		return smokeConfig{}, false
	}
	if cfg.Backend == "" {
		cfg.Backend = "s3"
	}
	if cfg.Region == "" {
		cfg.Region = "us-east-1"
	}
	if cfg.BucketPrefix == "" {
		cfg.BucketPrefix = "cosbench-go-smoke"
	}
	return cfg, true
}

func newSmokeAdapter(t *testing.T, cfg smokeConfig) *Adapter {
	t.Helper()
	raw := []string{
		"endpoint=" + cfg.Endpoint,
		"accesskey=" + cfg.AccessKey,
		"secretkey=" + cfg.SecretKey,
		"region=" + cfg.Region,
	}
	if cfg.PathStyleRaw != "" {
		raw = append(raw, "path_style_access="+cfg.PathStyleRaw)
	}

	adapter := NewAdapter(cfg.Backend, strings.Join(raw, ";"))
	if err := adapter.Init(kvconfig.ParseKVConfig(strings.Join(raw, ";"))); err != nil {
		t.Fatalf("Init(): %v", err)
	}
	return adapter
}

func smokeBucketName(prefix string) string {
	safe := strings.ToLower(prefix)
	safe = strings.ReplaceAll(safe, "_", "-")
	safe = strings.ReplaceAll(safe, " ", "-")
	return fmt.Sprintf("%s-%d", safe, time.Now().UnixNano())
}

func getenv(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

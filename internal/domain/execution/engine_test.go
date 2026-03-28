package execution

import (
	"bytes"
	"context"
	"io"
	"math/rand"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/sine-io/cosbench-go/internal/application/ports"
	"github.com/sine-io/cosbench-go/internal/domain/workload"
)

type stubStorage struct {
	putPayload       []byte
	putSize          int64
	multipartPayload []byte
	multipartSize    int64
	partSize         int64
	restoreDays      int
	deleteCalls      []string
	listCalls        []string
}

func (s *stubStorage) Init(map[string]string) error                                      { return nil }
func (s *stubStorage) Dispose() error                                                    { return nil }
func (s *stubStorage) CreateBucket(context.Context, string) error                        { return nil }
func (s *stubStorage) DeleteBucket(context.Context, string) error                        { return nil }
func (s *stubStorage) PutObject(ctx context.Context, bucket, key string, body io.Reader, size int64) error {
	data, err := io.ReadAll(body)
	if err != nil {
		return err
	}
	s.putPayload = append([]byte(nil), data...)
	s.putSize = size
	return nil
}
func (s *stubStorage) GetObject(context.Context, string, string) (io.ReadCloser, error) {
	return io.NopCloser(strings.NewReader("abc")), nil
}
func (s *stubStorage) DeleteObject(ctx context.Context, bucket, key string) error {
	s.deleteCalls = append(s.deleteCalls, bucket+"/"+key)
	return nil
}
func (s *stubStorage) HeadObject(context.Context, string, string) (ports.ObjectMeta, error) {
	return ports.ObjectMeta{ContentLength: 3}, nil
}
func (s *stubStorage) ListObjects(ctx context.Context, bucket, prefix string, maxKeys int) ([]ports.ObjectEntry, error) {
	s.listCalls = append(s.listCalls, bucket+"/"+prefix)
	return []ports.ObjectEntry{{Key: prefix + "-1"}, {Key: prefix + "-2"}}, nil
}
func (s *stubStorage) MultipartPut(ctx context.Context, bucket, key string, body io.Reader, size int64, partSize int64) error {
	data, err := io.ReadAll(body)
	if err != nil {
		return err
	}
	s.multipartPayload = append([]byte(nil), data...)
	s.multipartSize = size
	s.partSize = partSize
	return nil
}
func (s *stubStorage) RestoreObject(ctx context.Context, bucket, key string, days int) error {
	s.restoreDays = days
	return nil
}

func TestEngineRunTotalOps(t *testing.T) {
	e := Engine{Work: workload.Work{Workers: 2, TotalOps: 4, Operations: []workload.Operation{{Type: "read", Ratio: 100, Config: "containers=c(1);objects=c(1)"}}}, Storage: &stubStorage{}}
	res := e.Run(context.Background())
	if res.Err != nil {
		t.Fatal(res.Err)
	}
	if len(res.Samples) != 4 {
		t.Fatalf("samples = %d", len(res.Samples))
	}
}

func TestExecuteOpMFileWriteUsesLocalFileContents(t *testing.T) {
	storage := &stubStorage{}
	path := t.TempDir() + "/payload.bin"
	if err := os.WriteFile(path, []byte("payload"), 0o644); err != nil {
		t.Fatal(err)
	}

	n, err := executeOp(context.Background(), storage, "", workload.Operation{
		Type:   "mfilewrite",
		Ratio:  100,
		Config: "containers=c(1);objects=c(1);files=" + path + ";part_size=4",
	}, 1, 1, nil)
	if err != nil {
		t.Fatal(err)
	}
	if n != int64(len("payload")) {
		t.Fatalf("bytes = %d", n)
	}
	if !bytes.Equal(storage.multipartPayload, []byte("payload")) {
		t.Fatalf("multipart payload = %q", storage.multipartPayload)
	}
	if storage.multipartSize != int64(len("payload")) || storage.partSize != 4 {
		t.Fatalf("multipart size=%d partSize=%d", storage.multipartSize, storage.partSize)
	}
}

func TestExecuteOpFileWriteUsesLocalFileContents(t *testing.T) {
	storage := &stubStorage{}
	path := t.TempDir() + "/payload.bin"
	if err := os.WriteFile(path, []byte("payload"), 0o644); err != nil {
		t.Fatal(err)
	}

	n, err := executeOp(context.Background(), storage, "", workload.Operation{
		Type:   "filewrite",
		Ratio:  100,
		Config: "containers=c(1);objects=c(1);files=" + path,
	}, 1, 1, nil)
	if err != nil {
		t.Fatal(err)
	}
	if n != int64(len("payload")) {
		t.Fatalf("bytes = %d", n)
	}
	if !bytes.Equal(storage.putPayload, []byte("payload")) {
		t.Fatalf("put payload = %q", storage.putPayload)
	}
	if storage.putSize != int64(len("payload")) {
		t.Fatalf("put size = %d", storage.putSize)
	}
}

func TestExecuteOpDelayWaits(t *testing.T) {
	start := time.Now()
	if _, err := executeOp(context.Background(), &stubStorage{}, "", workload.Operation{
		Type:   "delay",
		Ratio:  100,
		Config: "duration=25ms",
	}, 1, 1, nil); err != nil {
		t.Fatal(err)
	}
	if elapsed := time.Since(start); elapsed < 20*time.Millisecond {
		t.Fatalf("delay elapsed too short: %s", elapsed)
	}
}

func TestExecuteOpCleanupScansTargetsSequentially(t *testing.T) {
	storage := &stubStorage{}
	if _, err := executeOp(context.Background(), storage, "", workload.Operation{
		Type:   "cleanup",
		Ratio:  100,
		Config: "cprefix=b;oprefix=o;containers=r(1,2);objects=r(1,3)",
	}, 1, 1, rand.New(rand.NewSource(1))); err != nil {
		t.Fatal(err)
	}
	want := []string{"b1/o1", "b1/o2", "b1/o3", "b2/o1", "b2/o2", "b2/o3"}
	if !reflect.DeepEqual(storage.deleteCalls, want) {
		t.Fatalf("delete calls = %#v want %#v", storage.deleteCalls, want)
	}
}

func TestExecuteOpListScansTargetsSequentially(t *testing.T) {
	storage := &stubStorage{}
	n, err := executeOp(context.Background(), storage, "", workload.Operation{
		Type:   "list",
		Ratio:  100,
		Config: "cprefix=b;oprefix=o;containers=r(1,2);objects=r(1,2)",
	}, 1, 1, rand.New(rand.NewSource(1)))
	if err != nil {
		t.Fatal(err)
	}
	wantCalls := []string{"b1/o1", "b1/o2", "b2/o1", "b2/o2"}
	if !reflect.DeepEqual(storage.listCalls, wantCalls) {
		t.Fatalf("list calls = %#v want %#v", storage.listCalls, wantCalls)
	}
	if n != 8 {
		t.Fatalf("listed count = %d", n)
	}
}

func TestExecuteOpMWriteUsesStorageLevelPartSizeByDefault(t *testing.T) {
	storage := &stubStorage{}
	n, err := executeOp(context.Background(), storage, "part_size=7340032", workload.Operation{
		Type:   "mwrite",
		Ratio:  100,
		Config: "containers=c(1);objects=c(1);sizes=c(1)KB",
	}, 1, 1, rand.New(rand.NewSource(1)))
	if err != nil {
		t.Fatal(err)
	}
	if n != 1000 {
		t.Fatalf("bytes = %d", n)
	}
	if storage.partSize != 7340032 {
		t.Fatalf("partSize = %d", storage.partSize)
	}
}

func TestExecuteOpRestoreUsesStorageLevelRestoreDaysByDefault(t *testing.T) {
	storage := &stubStorage{}
	if _, err := executeOp(context.Background(), storage, "restore_days=7", workload.Operation{
		Type:   "restore",
		Ratio:  100,
		Config: "containers=c(1);objects=c(1)",
	}, 1, 1, rand.New(rand.NewSource(1))); err != nil {
		t.Fatal(err)
	}
	if storage.restoreDays != 7 {
		t.Fatalf("restoreDays = %d", storage.restoreDays)
	}
}

func TestResolvedStorageConfigIncludesAuthAndStorageConfig(t *testing.T) {
	raw := storageConfig(&workload.StorageSpec{Type: "sio", Config: "endpoint=http://storage"}, &workload.AuthSpec{Type: "basic", Config: "username=work;password=secret"})
	if !strings.Contains(raw, "endpoint=http://storage") {
		t.Fatalf("missing storage config in %q", raw)
	}
	if !strings.Contains(raw, "username=work;password=secret") {
		t.Fatalf("missing auth config in %q", raw)
	}
}

func TestEngineRunReturnsContextCanceledOnExternalCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(20 * time.Millisecond)
		cancel()
	}()

	e := Engine{Work: workload.Work{
		Workers: 1,
		Operations: []workload.Operation{{
			Type:   "write",
			Ratio:  100,
			Config: "containers=c(1);objects=u(1,100);sizes=c(1)KB",
		}},
	}, Storage: &stubStorage{}}
	res := e.Run(ctx)
	if res.Err != context.Canceled {
		t.Fatalf("err = %v", res.Err)
	}
	if len(res.Samples) == 0 {
		t.Fatal("expected partial samples before cancellation")
	}
}

func TestEngineRunTreatsRuntimeDeadlineAsNormalCompletion(t *testing.T) {
	e := Engine{Work: workload.Work{
		Workers: 1,
		Runtime: 1,
		Operations: []workload.Operation{{
			Type:   "write",
			Ratio:  100,
			Config: "containers=c(1);objects=u(1,100);sizes=c(1)KB",
		}},
	}, Storage: &stubStorage{}}
	res := e.Run(context.Background())
	if res.Err != nil {
		t.Fatalf("err = %v", res.Err)
	}
	if len(res.Samples) == 0 {
		t.Fatal("expected samples from runtime-bounded execution")
	}
}

package execution

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"math/rand"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/sine-io/cosbench-go/internal/application/ports"
	"github.com/sine-io/cosbench-go/internal/domain/workload"
)

type Engine struct {
	Work    workload.Work
	Storage ports.StorageAdapter
}

type Result struct {
	Samples []Sample
	Err     error
}

type Sample struct {
	Timestamp   time.Time
	OpType      string
	OpCount     int64
	ByteCount   int64
	ErrorCount  int64
	TotalTimeMs int64
}

func (e *Engine) Run(ctx context.Context) Result {
	workers := e.Work.Workers
	if workers <= 0 {
		workers = 1
	}
	picker, err := NewWeightedOperationPicker(e.Work.Operations)
	if err != nil {
		return Result{Err: err}
	}

	runCtx := ctx
	var cancel context.CancelFunc
	if e.Work.Runtime > 0 {
		runCtx, cancel = context.WithTimeout(ctx, time.Duration(e.Work.Runtime)*time.Second)
	} else {
		runCtx, cancel = context.WithCancel(ctx)
	}
	defer cancel()

	var globalOps atomic.Int64
	var wg sync.WaitGroup
	var mu sync.Mutex
	samples := make([]Sample, 0, max(1, workers))

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			r := rand.New(rand.NewSource(time.Now().UnixNano() + int64(workerID)))
			for {
				select {
				case <-runCtx.Done():
					return
				default:
				}
				if e.Work.TotalOps > 0 {
					n := globalOps.Add(1)
					if n > int64(e.Work.TotalOps) {
						return
					}
				}
				op := picker.Pick(r)
				start := time.Now()
				bytesN, err := executeOp(runCtx, e.Storage, storageConfig(e.Work.Storage), op, workerID+1, workers, r)
				elapsed := time.Since(start)
				s := Sample{Timestamp: start, OpType: op.Type, OpCount: 1, ByteCount: bytesN, TotalTimeMs: elapsed.Milliseconds()}
				if err != nil {
					s.ErrorCount = 1
				}
				mu.Lock()
				samples = append(samples, s)
				mu.Unlock()
			}
		}(i)
	}

	wg.Wait()
	if err := runCtx.Err(); err != nil {
		if err == context.DeadlineExceeded {
			return Result{Samples: samples}
		}
		if err == context.Canceled {
			return Result{Samples: samples, Err: err}
		}
		return Result{Samples: samples, Err: err}
	}
	return Result{Samples: samples}
}

func executeOp(ctx context.Context, sa ports.StorageAdapter, storageRaw string, op workload.Operation, idx, all int, r *rand.Rand) (int64, error) {
	pc, err := ParseOpConfigWithStorage(storageRaw, op.Config)
	if err != nil {
		return 0, err
	}
	t := pc.NextTarget(r, idx, all)
	switch op.Type {
	case "init":
		return 0, sa.CreateBucket(ctx, t.Bucket)
	case "dispose":
		return 0, sa.DeleteBucket(ctx, t.Bucket)
	case "prepare", "write":
		body := bytes.NewReader(make([]byte, t.Size))
		return t.Size, sa.PutObject(ctx, t.Bucket, t.Key, body, t.Size)
	case "mprepare", "mwrite":
		body := bytes.NewReader(make([]byte, t.Size))
		return t.Size, sa.MultipartPut(ctx, t.Bucket, t.Key, body, t.Size, pc.PartSize)
	case "read":
		rc, err := sa.GetObject(ctx, t.Bucket, t.Key)
		if err != nil {
			return 0, err
		}
		defer rc.Close()
		n, err := io.Copy(io.Discard, rc)
		return n, err
	case "delete", "cleanup":
		if op.Type == "cleanup" {
			for _, target := range pc.ScanTargets() {
				if err := sa.DeleteObject(ctx, target.Bucket, target.Key); err != nil {
					return 0, err
				}
			}
			return 0, nil
		}
		return 0, sa.DeleteObject(ctx, t.Bucket, t.Key)
	case "head":
		meta, err := sa.HeadObject(ctx, t.Bucket, t.Key)
		return meta.ContentLength, err
	case "list":
		var total int64
		for _, target := range pc.ScanTargets() {
			items, err := sa.ListObjects(ctx, target.Bucket, target.Key, 1000)
			if err != nil {
				return total, err
			}
			total += int64(len(items))
		}
		return total, nil
	case "restore":
		return 0, sa.RestoreObject(ctx, t.Bucket, t.Key, pc.RestoreDays)
	case "localwrite":
		body, size, err := openInputFile(t.File)
		if err != nil {
			return 0, err
		}
		defer body.Close()
		return size, sa.PutObject(ctx, t.Bucket, t.Key, body, size)
	case "mfilewrite":
		body, size, err := openInputFile(t.File)
		if err != nil {
			return 0, err
		}
		defer body.Close()
		return size, sa.MultipartPut(ctx, t.Bucket, t.Key, body, size, pc.PartSize)
	case "delay":
		d := pc.Delay
		if d <= 0 {
			d = time.Millisecond
		}
		timer := time.NewTimer(d)
		defer timer.Stop()
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		case <-timer.C:
			return 0, nil
		}
	default:
		return 0, fmt.Errorf("unsupported operation: %s", op.Type)
	}
}

func ValidateOperation(op workload.Operation, storageRaw string) error {
	pc, err := ParseOpConfigWithStorage(storageRaw, op.Config)
	if err != nil {
		return err
	}
	switch op.Type {
	case "init", "dispose", "prepare", "write", "mprepare", "mwrite", "read", "delete", "cleanup", "head", "list", "restore", "delay":
		return nil
	case "localwrite", "mfilewrite":
		target := pc.NextTarget(rand.New(rand.NewSource(1)), 1, 1)
		file, _, err := openInputFile(target.File)
		if err != nil {
			return err
		}
		return file.Close()
	default:
		return fmt.Errorf("unsupported operation: %s", op.Type)
	}
}

func storageConfig(spec *workload.StorageSpec) string {
	if spec == nil {
		return ""
	}
	return spec.Config
}

func openInputFile(path string) (*os.File, int64, error) {
	if path == "" {
		return nil, 0, fmt.Errorf("missing files config for file-backed operation")
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, 0, err
	}
	info, err := file.Stat()
	if err != nil {
		_ = file.Close()
		return nil, 0, err
	}
	if info.IsDir() {
		_ = file.Close()
		return nil, 0, fmt.Errorf("file-backed operation requires a file, got directory: %s", path)
	}
	return file, info.Size(), nil
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

package controlplane

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/sine-io/cosbench-go/internal/domain"
	"github.com/sine-io/cosbench-go/internal/snapshot"
)

func TestRepresentativeFixtureProducesReportableResult(t *testing.T) {
	store, err := snapshot.New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	mgr, err := New(store)
	if err != nil {
		t.Fatal(err)
	}
	endpoint, err := mgr.CreateEndpoint(domain.EndpointConfig{Name: "mock", Type: domain.EndpointTypeMock})
	if err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(filepath.Join("..", "..", "testdata", "workloads", "s3-active-subset.xml"))
	if err != nil {
		t.Fatal(err)
	}
	job, err := mgr.CreateJobFromXML(raw, endpoint.ID)
	if err != nil {
		t.Fatal(err)
	}
	if err := mgr.StartJob(context.Background(), job.ID); err != nil {
		t.Fatal(err)
	}
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		loaded, ok := mgr.GetJob(job.ID)
		if !ok {
			t.Fatal("job disappeared")
		}
		if loaded.Status == domain.JobStatusSucceeded {
			result, ok := mgr.GetJobResult(job.ID)
			if !ok {
				t.Fatal("missing result")
			}
			if result.Metrics.OperationCount == 0 || len(result.StageTotals) == 0 || len(result.Metrics.ByOperation) == 0 {
				t.Fatalf("unexpected result summary: %#v", result)
			}
			return
		}
		time.Sleep(20 * time.Millisecond)
	}
	loaded, _ := mgr.GetJob(job.ID)
	t.Fatalf("job did not finish: %#v", loaded)
}

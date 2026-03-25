package snapshot

import (
	"testing"
	"time"

	"github.com/sine-io/cosbench-go/internal/domain"
)

func TestStoreRoundTrip(t *testing.T) {
	store, err := New(t.TempDir())
	if err != nil {
		t.Fatalf("New(): %v", err)
	}
	now := time.Now().UTC()
	endpoint := domain.EndpointConfig{ID: "ep-1", Name: "mock", Type: domain.EndpointTypeMock, CreatedAt: now, UpdatedAt: now}
	job := domain.Job{ID: "job-1", Name: "job", Status: domain.JobStatusCreated, CreatedAt: now, Workload: domain.Workload{Name: "wl"}}
	result := domain.JobResult{JobID: job.ID, UpdatedAt: now}
	events := []domain.JobEvent{{JobID: job.ID, OccurredAt: now, Level: domain.EventLevelInfo, Message: "created"}}

	if err := store.SaveEndpoint(endpoint); err != nil {
		t.Fatalf("SaveEndpoint(): %v", err)
	}
	if err := store.SaveJob(job); err != nil {
		t.Fatalf("SaveJob(): %v", err)
	}
	if err := store.SaveResult(result); err != nil {
		t.Fatalf("SaveResult(): %v", err)
	}
	if err := store.SaveEvents(job.ID, events); err != nil {
		t.Fatalf("SaveEvents(): %v", err)
	}

	jobs, err := store.LoadJobs()
	if err != nil {
		t.Fatalf("LoadJobs(): %v", err)
	}
	if len(jobs) != 1 || jobs[0].ID != job.ID {
		t.Fatalf("jobs = %#v", jobs)
	}
	endpoints, err := store.LoadEndpoints()
	if err != nil {
		t.Fatalf("LoadEndpoints(): %v", err)
	}
	if len(endpoints) != 1 || endpoints[0].ID != endpoint.ID {
		t.Fatalf("endpoints = %#v", endpoints)
	}
	loadedResult, err := store.LoadResult(job.ID)
	if err != nil {
		t.Fatalf("LoadResult(): %v", err)
	}
	if loadedResult.JobID != job.ID {
		t.Fatalf("result = %#v", loadedResult)
	}
	loadedEvents, err := store.LoadEvents(job.ID)
	if err != nil {
		t.Fatalf("LoadEvents(): %v", err)
	}
	if len(loadedEvents) != 1 || loadedEvents[0].Message != "created" {
		t.Fatalf("events = %#v", loadedEvents)
	}
}

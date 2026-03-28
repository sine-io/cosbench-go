package reporting

import (
	"testing"
	"time"

	legacyexec "github.com/sine-io/cosbench-go/internal/domain/execution"
)

func TestBuildTimelineAggregatesSamplesIntoBuckets(t *testing.T) {
	base := time.Unix(0, 0).UTC()
	timeline := BuildTimeline([]legacyexec.Sample{
		{Timestamp: base.Add(100 * time.Millisecond), OpCount: 1, ByteCount: 100, ErrorCount: 0, TotalTimeMs: 10},
		{Timestamp: base.Add(900 * time.Millisecond), OpCount: 1, ByteCount: 200, ErrorCount: 1, TotalTimeMs: 30},
		{Timestamp: base.Add(1100 * time.Millisecond), OpCount: 1, ByteCount: 300, ErrorCount: 0, TotalTimeMs: 20},
	}, time.Second)
	if len(timeline) != 2 {
		t.Fatalf("bucket count = %d", len(timeline))
	}
	first := timeline[0]
	if first.OperationCount != 2 || first.ByteCount != 300 || first.ErrorCount != 1 {
		t.Fatalf("first bucket = %#v", first)
	}
	if first.AvgLatencyMs != 20 {
		t.Fatalf("first avg latency = %f", first.AvgLatencyMs)
	}
	second := timeline[1]
	if second.OperationCount != 1 || second.ByteCount != 300 || second.ErrorCount != 0 {
		t.Fatalf("second bucket = %#v", second)
	}
	if !second.Timestamp.After(first.Timestamp) {
		t.Fatalf("unexpected bucket order: %#v", timeline)
	}
}

package reporting

import (
	"testing"
	"time"

	legacyexec "github.com/sine-io/cosbench-go/internal/domain/execution"
)

func TestSummarizeIncludesPercentiles(t *testing.T) {
	now := time.Now()
	summary := Summarize([]legacyexec.Sample{
		{Timestamp: now, OpType: "read", OpCount: 1, TotalTimeMs: 10},
		{Timestamp: now.Add(time.Millisecond), OpType: "read", OpCount: 1, TotalTimeMs: 20},
		{Timestamp: now.Add(2 * time.Millisecond), OpType: "write", OpCount: 1, TotalTimeMs: 30},
		{Timestamp: now.Add(3 * time.Millisecond), OpType: "write", OpCount: 1, TotalTimeMs: 40},
	})
	if summary.P50LatencyMs <= 0 || summary.P95LatencyMs <= 0 || summary.P99LatencyMs <= 0 {
		t.Fatalf("missing percentile values: %#v", summary)
	}
	if len(summary.ByOperation) != 2 {
		t.Fatalf("by operation = %#v", summary.ByOperation)
	}
}

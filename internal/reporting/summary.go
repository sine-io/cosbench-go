package reporting

import (
	"sort"

	"github.com/sine-io/cosbench-go/internal/domain"
	legacyexec "github.com/sine-io/cosbench-go/internal/domain/execution"
)

func Summarize(samples []legacyexec.Sample) domain.MetricsSummary {
	if len(samples) == 0 {
		return domain.MetricsSummary{}
	}
	perOp := map[string]*domain.OperationMetrics{}
	perOpLatencies := map[string][]int64{}
	latencies := make([]int64, 0, len(samples))
	var firstTs, lastTs int64
	firstTs = samples[0].Timestamp.UnixMilli()
	lastTs = firstTs
	var summary domain.MetricsSummary
	for _, sample := range samples {
		ts := sample.Timestamp.UnixMilli()
		if ts < firstTs {
			firstTs = ts
		}
		if end := ts + sample.TotalTimeMs; end > lastTs {
			lastTs = end
		}
		summary.OperationCount += sample.OpCount
		summary.ByteCount += sample.ByteCount
		summary.ErrorCount += sample.ErrorCount
		summary.TotalLatencyMs += sample.TotalTimeMs
		latencies = append(latencies, sample.TotalTimeMs)
		item := perOp[sample.OpType]
		if item == nil {
			item = &domain.OperationMetrics{Operation: sample.OpType}
			perOp[sample.OpType] = item
		}
		item.OperationCount += sample.OpCount
		item.ByteCount += sample.ByteCount
		item.ErrorCount += sample.ErrorCount
		item.AvgLatencyMs += float64(sample.TotalTimeMs)
		perOpLatencies[sample.OpType] = append(perOpLatencies[sample.OpType], sample.TotalTimeMs)
	}
	summary.DurationMs = lastTs - firstTs
	if summary.OperationCount > 0 {
		summary.AvgLatencyMs = float64(summary.TotalLatencyMs) / float64(summary.OperationCount)
	}
	if summary.DurationMs > 0 {
		summary.OpsPerSecond = float64(summary.OperationCount) / (float64(summary.DurationMs) / 1000)
	}
	summary.P50LatencyMs = percentile(latencies, 0.50)
	summary.P95LatencyMs = percentile(latencies, 0.95)
	summary.P99LatencyMs = percentile(latencies, 0.99)
	for _, item := range perOp {
		if item.OperationCount > 0 {
			item.AvgLatencyMs /= float64(item.OperationCount)
		}
		item.P50LatencyMs = percentile(perOpLatencies[item.Operation], 0.50)
		item.P95LatencyMs = percentile(perOpLatencies[item.Operation], 0.95)
		item.P99LatencyMs = percentile(perOpLatencies[item.Operation], 0.99)
		summary.ByOperation = append(summary.ByOperation, *item)
	}
	sort.Slice(summary.ByOperation, func(i, j int) bool { return summary.ByOperation[i].Operation < summary.ByOperation[j].Operation })
	return summary
}

func Merge(parts ...domain.MetricsSummary) domain.MetricsSummary {
	merged := domain.MetricsSummary{}
	perOp := map[string]*domain.OperationMetrics{}
	weightedP50 := 0.0
	weightedP95 := 0.0
	weightedP99 := 0.0
	for _, part := range parts {
		merged.OperationCount += part.OperationCount
		merged.ByteCount += part.ByteCount
		merged.ErrorCount += part.ErrorCount
		merged.TotalLatencyMs += part.TotalLatencyMs
		merged.DurationMs += part.DurationMs
		weightedP50 += part.P50LatencyMs * float64(part.OperationCount)
		weightedP95 += part.P95LatencyMs * float64(part.OperationCount)
		weightedP99 += part.P99LatencyMs * float64(part.OperationCount)
		for _, item := range part.ByOperation {
			target := perOp[item.Operation]
			if target == nil {
				target = &domain.OperationMetrics{Operation: item.Operation}
				perOp[item.Operation] = target
			}
			target.OperationCount += item.OperationCount
			target.ByteCount += item.ByteCount
			target.ErrorCount += item.ErrorCount
			target.AvgLatencyMs += item.AvgLatencyMs * float64(item.OperationCount)
		}
	}
	if merged.OperationCount > 0 {
		merged.AvgLatencyMs = float64(merged.TotalLatencyMs) / float64(merged.OperationCount)
		merged.P50LatencyMs = weightedP50 / float64(merged.OperationCount)
		merged.P95LatencyMs = weightedP95 / float64(merged.OperationCount)
		merged.P99LatencyMs = weightedP99 / float64(merged.OperationCount)
	}
	if merged.DurationMs > 0 {
		merged.OpsPerSecond = float64(merged.OperationCount) / (float64(merged.DurationMs) / 1000)
	}
	for _, item := range perOp {
		if item.OperationCount > 0 {
			item.AvgLatencyMs /= float64(item.OperationCount)
		}
		merged.ByOperation = append(merged.ByOperation, *item)
	}
	sort.Slice(merged.ByOperation, func(i, j int) bool { return merged.ByOperation[i].Operation < merged.ByOperation[j].Operation })
	return merged
}

func percentile(values []int64, p float64) float64 {
	if len(values) == 0 {
		return 0
	}
	cp := append([]int64(nil), values...)
	sort.Slice(cp, func(i, j int) bool { return cp[i] < cp[j] })
	idx := int(float64(len(cp)-1) * p)
	if idx < 0 {
		idx = 0
	}
	if idx >= len(cp) {
		idx = len(cp) - 1
	}
	return float64(cp[idx])
}

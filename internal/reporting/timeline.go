package reporting

import (
	"sort"
	"time"

	"github.com/sine-io/cosbench-go/internal/domain"
	legacyexec "github.com/sine-io/cosbench-go/internal/domain/execution"
)

func BuildTimeline(samples []legacyexec.Sample, bucketSize time.Duration) []domain.TimelinePoint {
	if len(samples) == 0 {
		return nil
	}
	if bucketSize <= 0 {
		bucketSize = time.Second
	}
	type accum struct {
		point domain.TimelinePoint
		total int64
	}
	buckets := map[int64]*accum{}
	order := make([]int64, 0)
	for _, sample := range samples {
		ts := sample.Timestamp.UTC()
		key := ts.UnixNano() / bucketSize.Nanoseconds()
		item := buckets[key]
		if item == nil {
			item = &accum{point: domain.TimelinePoint{Timestamp: time.Unix(0, key*bucketSize.Nanoseconds()).UTC()}}
			buckets[key] = item
			order = append(order, key)
		}
		item.point.OperationCount += sample.OpCount
		item.point.ByteCount += sample.ByteCount
		item.point.ErrorCount += sample.ErrorCount
		item.total += sample.TotalTimeMs
	}
	sort.Slice(order, func(i, j int) bool { return order[i] < order[j] })
	points := make([]domain.TimelinePoint, 0, len(order))
	for _, key := range order {
		item := buckets[key]
		if item.point.OperationCount > 0 {
			item.point.AvgLatencyMs = float64(item.total) / float64(item.point.OperationCount)
		}
		points = append(points, item.point)
	}
	return points
}

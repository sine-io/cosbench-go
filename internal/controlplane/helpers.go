package controlplane

import (
	"fmt"
	"sort"
	"time"

	"github.com/sine-io/cosbench-go/internal/domain"
)

func newID(prefix string) string {
	return fmt.Sprintf("%s-%d", prefix, time.Now().UTC().UnixNano())
}

func sortJobs(items []domain.Job) {
	sort.Slice(items, func(i, j int) bool { return items[i].CreatedAt.After(items[j].CreatedAt) })
}

func sortEndpoints(items []domain.EndpointConfig) {
	sort.Slice(items, func(i, j int) bool { return items[i].UpdatedAt.After(items[j].UpdatedAt) })
}

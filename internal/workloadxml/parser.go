package workloadxml

import (
	"fmt"
	"os"

	"github.com/sine-io/cosbench-go/internal/domain"
	legacyxml "github.com/sine-io/cosbench-go/internal/infrastructure/xml"
)

func Parse(data []byte) (domain.Workload, error) {
	parsed, err := legacyxml.ParseWorkload(data)
	if err != nil {
		return domain.Workload{}, fmt.Errorf("parse workload xml: %w", err)
	}
	return domain.WorkloadFromLegacy(parsed), nil
}

func ParseFile(path string) (domain.Workload, []byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return domain.Workload{}, nil, err
	}
	parsed, err := Parse(data)
	if err != nil {
		return domain.Workload{}, nil, err
	}
	return parsed, data, nil
}

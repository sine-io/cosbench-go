package storage

import (
	"fmt"
	"strings"

	newdriver "github.com/sine-io/cosbench-go/internal/driver/s3"
	"github.com/sine-io/cosbench-go/internal/application/ports"
	"github.com/sine-io/cosbench-go/internal/infrastructure/storage/mock"
)

func NewAdapter(storageType, rawConfig string) (ports.StorageAdapter, error) {
	switch strings.ToLower(strings.TrimSpace(storageType)) {
	case "mock":
		return mock.New(), nil
	case "s3":
		return newdriver.NewAdapter("s3", rawConfig), nil
	case "sio", "siov1", "gdas":
		return newdriver.NewAdapter(strings.ToLower(strings.TrimSpace(storageType)), rawConfig), nil
	default:
		return nil, fmt.Errorf("unsupported storage type: %q", storageType)
	}
}

package storage

import (
	"strings"

	"github.com/sine-io/cosbench-go/internal/application/ports"
	"github.com/sine-io/cosbench-go/internal/infrastructure/config"
)

type RunAdapters struct {
	mock ports.StorageAdapter
}

func NewRunAdapters() *RunAdapters {
	return &RunAdapters{}
}

func (r *RunAdapters) Acquire(storageType, rawConfig string) (ports.StorageAdapter, bool, error) {
	if strings.EqualFold(strings.TrimSpace(storageType), "mock") {
		if r.mock == nil {
			adapter, err := NewAdapter(storageType, rawConfig)
			if err != nil {
				return nil, false, err
			}
			if err := adapter.Init(config.ParseKVConfig(rawConfig)); err != nil {
				return nil, false, err
			}
			r.mock = adapter
		}
		return r.mock, true, nil
	}

	adapter, err := NewAdapter(storageType, rawConfig)
	if err != nil {
		return nil, false, err
	}
	if err := adapter.Init(config.ParseKVConfig(rawConfig)); err != nil {
		return nil, false, err
	}
	return adapter, false, nil
}

func (r *RunAdapters) Close() error {
	if r.mock != nil {
		err := r.mock.Dispose()
		r.mock = nil
		return err
	}
	return nil
}

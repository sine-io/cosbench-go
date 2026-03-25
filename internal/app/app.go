package app

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/sine-io/cosbench-go/internal/controlplane"
	"github.com/sine-io/cosbench-go/internal/snapshot"
	"github.com/sine-io/cosbench-go/internal/web"
)

type Config struct {
	DataDir string
	ViewDir string
}

type App struct {
	Manager *controlplane.Manager
	Handler http.Handler
}

func New(cfg Config) (*App, error) {
	dataDir := cfg.DataDir
	if dataDir == "" {
		dataDir = filepath.Join("data")
	}
	viewDir := cfg.ViewDir
	if viewDir == "" {
		viewDir = filepath.Join("web", "templates")
	}
	store, err := snapshot.New(dataDir)
	if err != nil {
		return nil, fmt.Errorf("snapshot store: %w", err)
	}
	manager, err := controlplane.New(store)
	if err != nil {
		return nil, fmt.Errorf("controlplane: %w", err)
	}
	handler, err := web.NewHandler(manager, viewDir)
	if err != nil {
		return nil, err
	}
	return &App{Manager: manager, Handler: handler}, nil
}

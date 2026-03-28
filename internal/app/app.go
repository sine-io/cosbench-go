package app

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/sine-io/cosbench-go/internal/controlplane"
	driveragent "github.com/sine-io/cosbench-go/internal/driver/agent"
	"github.com/sine-io/cosbench-go/internal/snapshot"
	"github.com/sine-io/cosbench-go/internal/web"
)

type Config struct {
	DataDir string
	ViewDir string
	Mode    Mode
}

type App struct {
	Mode         Mode
	Manager      *controlplane.Manager
	Handler      http.Handler
	loopbackAgent *driveragent.Agent
}

func New(cfg Config) (*App, error) {
	mode, err := normalizeMode(cfg.Mode)
	if err != nil {
		return nil, err
	}
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
	return &App{Mode: mode, Manager: manager, Handler: handler}, nil
}

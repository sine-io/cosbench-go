package app

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/sine-io/cosbench-go/internal/controlplane"
	driveragent "github.com/sine-io/cosbench-go/internal/driver/agent"
	"github.com/sine-io/cosbench-go/internal/snapshot"
	"github.com/sine-io/cosbench-go/internal/web"
)

type Config struct {
	DataDir string
	ViewDir string
	Mode    Mode
	DriverSharedToken string
}

type App struct {
	Mode         Mode
	Manager      *controlplane.Manager
	Handler      http.Handler
	DriverSharedToken string
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
	driverSharedToken := strings.TrimSpace(cfg.DriverSharedToken)
	if driverSharedToken == "" {
		driverSharedToken = strings.TrimSpace(os.Getenv("COSBENCH_DRIVER_SHARED_TOKEN"))
	}
	store, err := snapshot.New(dataDir)
	if err != nil {
		return nil, fmt.Errorf("snapshot store: %w", err)
	}
	manager, err := controlplane.New(store)
	if err != nil {
		return nil, fmt.Errorf("controlplane: %w", err)
	}
	handler, err := web.NewHandler(manager, viewDir, driverSharedToken)
	if err != nil {
		return nil, err
	}
	return &App{Mode: mode, Manager: manager, Handler: handler, DriverSharedToken: driverSharedToken}, nil
}

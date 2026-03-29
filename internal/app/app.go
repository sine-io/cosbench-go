package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/sine-io/cosbench-go/internal/controlplane"
	"github.com/sine-io/cosbench-go/internal/domain"
	driveragent "github.com/sine-io/cosbench-go/internal/driver/agent"
	"github.com/sine-io/cosbench-go/internal/snapshot"
	"github.com/sine-io/cosbench-go/internal/web"
)

type Config struct {
	DataDir string
	ViewDir string
	Mode    Mode
	DriverSharedToken string
	ControllerURL string
	DriverName string
	DriverPollInterval time.Duration
	LeaseSweepInterval time.Duration
}

type App struct {
	Mode         Mode
	Manager      *controlplane.Manager
	Handler      http.Handler
	DriverSharedToken string
	ControllerURL string
	DriverName string
	DriverPollInterval time.Duration
	LeaseSweepInterval time.Duration
	loopbackAgent *driveragent.Agent
	backgroundOnce sync.Once
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
	controllerURL := strings.TrimSpace(cfg.ControllerURL)
	if controllerURL == "" {
		controllerURL = strings.TrimSpace(os.Getenv("COSBENCH_CONTROLLER_URL"))
	}
	driverName := strings.TrimSpace(cfg.DriverName)
	if driverName == "" {
		driverName = strings.TrimSpace(os.Getenv("COSBENCH_DRIVER_NAME"))
	}
	driverPollInterval := cfg.DriverPollInterval
	if driverPollInterval <= 0 {
		driverPollInterval = 100 * time.Millisecond
	}
	leaseSweepInterval := cfg.LeaseSweepInterval
	if leaseSweepInterval <= 0 {
		leaseSweepInterval = 100 * time.Millisecond
	}
	store, err := snapshot.New(dataDir)
	if err != nil {
		return nil, fmt.Errorf("snapshot store: %w", err)
	}
	manager, err := controlplane.New(store)
	if err != nil {
		return nil, fmt.Errorf("controlplane: %w", err)
	}
	if mode == ModeControllerOnly {
		manager.SetRemoteScheduling(true)
	}
	handler, err := web.NewHandler(manager, viewDir, driverSharedToken)
	if err != nil {
		return nil, err
	}
	return &App{
		Mode: mode,
		Manager: manager,
		Handler: handler,
		DriverSharedToken: driverSharedToken,
		ControllerURL: controllerURL,
		DriverName: driverName,
		DriverPollInterval: driverPollInterval,
		LeaseSweepInterval: leaseSweepInterval,
	}, nil
}

func (a *App) StartBackground(ctx context.Context) error {
	if a.Mode == ModeDriverOnly && strings.TrimSpace(a.ControllerURL) == "" {
		return fmt.Errorf("driver-only mode requires controller url")
	}
	a.backgroundOnce.Do(func() {
		switch a.Mode {
		case ModeDriverOnly:
			agent := &driveragent.Agent{
				Client: &driveragent.HTTPClient{
					BaseURL: a.ControllerURL,
					SharedToken: a.DriverSharedToken,
				},
				Name:  a.DriverName,
				Mode:  domain.DriverModeDriver,
				Mirror: a.Manager,
			}
			go func() {
				ticker := time.NewTicker(a.DriverPollInterval)
				defer ticker.Stop()
				for {
					select {
					case <-ctx.Done():
						return
					default:
					}
					_, _ = agent.ProcessOne(ctx)
					select {
					case <-ctx.Done():
						return
					case <-ticker.C:
					}
				}
			}()
		case ModeControllerOnly, ModeCombined:
			go func() {
				ticker := time.NewTicker(a.LeaseSweepInterval)
				defer ticker.Stop()
				for {
					select {
					case <-ctx.Done():
						return
					case <-ticker.C:
						a.Manager.SweepExpiredLeases(time.Now().UTC())
					}
				}
			}()
		}
	})
	return nil
}

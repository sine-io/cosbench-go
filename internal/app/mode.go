package app

import (
	"context"
	"fmt"
	"net/http/httptest"

	"github.com/sine-io/cosbench-go/internal/domain"
	driveragent "github.com/sine-io/cosbench-go/internal/driver/agent"
)

type Mode string

const (
	ModeControllerOnly Mode = "controller-only"
	ModeDriverOnly     Mode = "driver-only"
	ModeCombined       Mode = "combined"
)

func normalizeMode(mode Mode) (Mode, error) {
	if mode == "" {
		return ModeCombined, nil
	}
	switch mode {
	case ModeControllerOnly, ModeDriverOnly, ModeCombined:
		return mode, nil
	default:
		return "", fmt.Errorf("unsupported app mode: %q", mode)
	}
}

func (a *App) ProcessCombinedMission(ctx context.Context) (bool, error) {
	if a.Mode != ModeCombined {
		return false, fmt.Errorf("combined mission processing requires %q mode", ModeCombined)
	}
	server := httptest.NewServer(a.Handler)
	defer server.Close()
	if a.loopbackAgent == nil {
		a.loopbackAgent = &driveragent.Agent{
			Client: &driveragent.HTTPClient{BaseURL: server.URL, SharedToken: a.DriverSharedToken},
			Name:   "combined-loopback",
			Mode:   domain.DriverModeCombined,
		}
	} else {
		a.loopbackAgent.Client = &driveragent.HTTPClient{BaseURL: server.URL, SharedToken: a.DriverSharedToken}
	}
	return a.loopbackAgent.ProcessOne(ctx)
}

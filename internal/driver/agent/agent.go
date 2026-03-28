package agent

import (
	"context"
	"fmt"
	"time"

	"github.com/sine-io/cosbench-go/internal/domain"
	"github.com/sine-io/cosbench-go/internal/executor"
	"github.com/sine-io/cosbench-go/internal/infrastructure/config"
	storagefactory "github.com/sine-io/cosbench-go/internal/infrastructure/storage"
)

type Agent struct {
	Client   *HTTPClient
	DriverID string
	Name     string
	Mode     domain.DriverMode
}

func (a *Agent) ProcessOne(ctx context.Context) (bool, error) {
	if a.Client == nil {
		return false, fmt.Errorf("missing driver client")
	}
	if a.Name == "" {
		a.Name = "driver"
	}
	if a.Mode == "" {
		a.Mode = domain.DriverModeDriver
	}
	if a.DriverID == "" {
		driver, err := a.Client.RegisterDriver(a.Name, a.Mode)
		if err != nil {
			return false, err
		}
		a.DriverID = driver.ID
	}
	if _, err := a.Client.Heartbeat(a.DriverID, time.Now().UTC()); err != nil {
		return false, err
	}
	mission, ok, err := a.Client.ClaimMission(a.DriverID, 30*time.Second)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}

	startEvent := domain.JobEvent{OccurredAt: time.Now().UTC(), Level: domain.EventLevelInfo, Message: "mission started"}
	if err := a.Client.UploadEvents(mission.ID, []domain.JobEvent{startEvent}); err != nil {
		return false, err
	}

	adapter, err := storagefactory.NewAdapter(mission.Work.Storage.Type, mission.Work.Storage.Config)
	if err != nil {
		return false, err
	}
	if err := adapter.Init(config.ParseKVConfig(mission.Work.Storage.Config)); err != nil {
		return false, err
	}
	defer adapter.Dispose()

	stageExecutor := executor.StageExecutor{Storage: adapter}
	workResult := stageExecutor.RunWork(ctx, mission.Work)
	if err := a.Client.UploadSamples(mission.ID, workResult.Samples); err != nil {
		return false, err
	}

	status := domain.MissionStatusSucceeded
	errorMessage := ""
	finishLevel := domain.EventLevelInfo
	finishMessage := "mission finished"
	if workResult.Err != nil {
		status = domain.MissionStatusFailed
		errorMessage = workResult.Err.Error()
		finishLevel = domain.EventLevelError
		finishMessage = "mission finished with failure"
	}
	if err := a.Client.UploadEvents(mission.ID, []domain.JobEvent{{
		OccurredAt: time.Now().UTC(),
		Level:      finishLevel,
		Message:    finishMessage,
	}}); err != nil {
		return false, err
	}
	if err := a.Client.CompleteMission(mission.ID, status, errorMessage); err != nil {
		return false, err
	}
	return true, nil
}

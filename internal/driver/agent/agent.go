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
	Mirror   LocalMirror
}

type LocalMirror interface {
	PutDriverNode(domain.DriverNode) error
	PutMission(domain.Mission) error
	AppendMissionEventsBatch(string, string, []domain.JobEvent) error
	CompleteMission(string, domain.MissionStatus, string) error
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
		if a.Mirror != nil {
			_ = a.Mirror.PutDriverNode(driver)
		}
	}
	driver, err := a.Client.Heartbeat(a.DriverID, time.Now().UTC())
	if err != nil {
		return false, err
	}
	if a.Mirror != nil {
		_ = a.Mirror.PutDriverNode(driver)
	}
	mission, ok, err := a.Client.ClaimMission(a.DriverID, 30*time.Second)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	if a.Mirror != nil {
		_ = a.Mirror.PutMission(mission)
	}

	startEvent := domain.JobEvent{OccurredAt: time.Now().UTC(), Level: domain.EventLevelInfo, Message: "mission started"}
	startBatchID := fmt.Sprintf("events-start-%d", time.Now().UnixNano())
	if err := a.Client.UploadEventsBatch(mission.ID, startBatchID, []domain.JobEvent{startEvent}); err != nil {
		return false, err
	}
	if a.Mirror != nil {
		_ = a.Mirror.AppendMissionEventsBatch(mission.ID, startBatchID, []domain.JobEvent{startEvent})
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
	sampleBatchID := fmt.Sprintf("samples-%d", time.Now().UnixNano())
	if err := a.Client.UploadSamplesBatch(mission.ID, sampleBatchID, workResult.Samples); err != nil {
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
	finishBatchID := fmt.Sprintf("events-finish-%d", time.Now().UnixNano())
	finishEvents := []domain.JobEvent{{
		OccurredAt: time.Now().UTC(),
		Level:      finishLevel,
		Message:    finishMessage,
	}}
	if err := a.Client.UploadEventsBatch(mission.ID, finishBatchID, finishEvents); err != nil {
		return false, err
	}
	if a.Mirror != nil {
		_ = a.Mirror.AppendMissionEventsBatch(mission.ID, finishBatchID, finishEvents)
	}
	if err := a.Client.CompleteMission(mission.ID, status, errorMessage); err != nil {
		return false, err
	}
	if a.Mirror != nil {
		_ = a.Mirror.CompleteMission(mission.ID, status, errorMessage)
	}
	return true, nil
}

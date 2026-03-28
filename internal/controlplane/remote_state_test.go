package controlplane

import (
	"testing"
	"time"

	"github.com/sine-io/cosbench-go/internal/domain"
	"github.com/sine-io/cosbench-go/internal/snapshot"
)

func TestManagerPersistsDriverNodesAndMissions(t *testing.T) {
	store, err := snapshot.New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	mgr, err := New(store)
	if err != nil {
		t.Fatal(err)
	}

	now := time.Now().UTC()
	driver := domain.DriverNode{
		ID:              "drv-1",
		Name:            "driver-1",
		Mode:            domain.DriverModeDriver,
		Status:          domain.DriverStatusHealthy,
		RegisteredAt:    now,
		LastHeartbeatAt: &now,
	}
	if err := mgr.PutDriverNode(driver); err != nil {
		t.Fatalf("PutDriverNode(): %v", err)
	}

	leaseClaimedAt := now.Add(10 * time.Second)
	mission := domain.Mission{
		ID:        "mission-1",
		JobID:     "job-1",
		StageName: "main",
		WorkName:  "work-1",
		Status:    domain.MissionStatusClaimed,
		CreatedAt: now,
		UpdatedAt: now,
		Lease: &domain.MissionLease{
			DriverID:  driver.ID,
			ClaimedAt: &leaseClaimedAt,
			ExpiresAt: now.Add(time.Minute),
		},
	}
	if err := mgr.PutMission(mission); err != nil {
		t.Fatalf("PutMission(): %v", err)
	}

	drivers := mgr.ListDriverNodes()
	if len(drivers) != 1 || drivers[0].ID != driver.ID || drivers[0].Status != domain.DriverStatusHealthy {
		t.Fatalf("drivers = %#v", drivers)
	}

	loadedMission, ok := mgr.GetMission(mission.ID)
	if !ok {
		t.Fatal("expected mission")
	}
	if loadedMission.Lease == nil || loadedMission.Lease.DriverID != driver.ID {
		t.Fatalf("loaded mission = %#v", loadedMission)
	}
	if loadedMission.Lease.ClaimedAt == nil || !loadedMission.Lease.ExpiresAt.After(*loadedMission.Lease.ClaimedAt) {
		t.Fatalf("unexpected mission lease = %#v", loadedMission.Lease)
	}

	reloaded, err := New(store)
	if err != nil {
		t.Fatal(err)
	}
	reloadedDrivers := reloaded.ListDriverNodes()
	if len(reloadedDrivers) != 1 || reloadedDrivers[0].ID != driver.ID {
		t.Fatalf("reloaded drivers = %#v", reloadedDrivers)
	}
	reloadedMission, ok := reloaded.GetMission(mission.ID)
	if !ok {
		t.Fatal("expected reloaded mission")
	}
	if reloadedMission.Lease == nil || reloadedMission.Lease.DriverID != driver.ID {
		t.Fatalf("reloaded mission = %#v", reloadedMission)
	}
}

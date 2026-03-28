package controlplane

import (
	"testing"

	"github.com/sine-io/cosbench-go/internal/domain"
	"github.com/sine-io/cosbench-go/internal/snapshot"
)

func TestScheduleJobStageDecomposesWorkIntoUnitsAndAttempts(t *testing.T) {
	store, err := snapshot.New(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	mgr, err := New(store)
	if err != nil {
		t.Fatal(err)
	}
	job, err := mgr.CreateJobFromXML([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<workload name="unit-decompose">
  <storage type="mock" />
  <workflow>
    <workstage name="main">
      <work name="fanout" workers="3" totalOps="3">
        <operation type="write" ratio="100" config="cprefix=t;containers=c(1);objects=s(1,3);sizes=c(1)KB" />
      </work>
    </workstage>
  </workflow>
</workload>`), "")
	if err != nil {
		t.Fatal(err)
	}

	attempts, err := mgr.ScheduleJobStage(job.ID)
	if err != nil {
		t.Fatalf("ScheduleJobStage(): %v", err)
	}
	if len(attempts) != 3 {
		t.Fatalf("attempts = %#v", attempts)
	}

	units := mgr.ListWorkUnits(job.ID, "main", "fanout")
	if len(units) != 3 {
		t.Fatalf("units = %#v", units)
	}
	for index, unit := range units {
		if unit.UnitIndex != index+1 || unit.UnitCount != 3 {
			t.Fatalf("unit = %#v", unit)
		}
		if unit.Slice.WorkerIndex != index+1 || unit.Slice.WorkerCount != 3 {
			t.Fatalf("unit slice = %#v", unit.Slice)
		}
		if unit.Status != domain.WorkUnitStatusPending {
			t.Fatalf("unit status = %#v", unit)
		}
	}

	allAttempts := mgr.ListMissionAttempts()
	if len(allAttempts) != 3 {
		t.Fatalf("all attempts = %#v", allAttempts)
	}
	for _, attempt := range allAttempts {
		if attempt.Attempt != 1 || attempt.Status != domain.MissionAttemptStatusPending {
			t.Fatalf("attempt = %#v", attempt)
		}
		if attempt.WorkUnitID == "" {
			t.Fatalf("attempt missing work unit id: %#v", attempt)
		}
	}
}

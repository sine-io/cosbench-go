package workload

import (
	"strings"
	"testing"
)

func TestNormalizeAndValidateInheritsConfigAndStorage(t *testing.T) {
	got, err := NormalizeAndValidate(Workload{
		Name:    "inherit",
		Storage: &StorageSpec{Type: "mock", Config: "endpoint=top"},
		Workflow: Workflow{
			Config: "wf=1",
			Stages: []Stage{{
				Name:   "main",
				Config: "stage=1",
				Works: []Work{{
					Name:     "main",
					Workers:  1,
					TotalOps: 1,
					Config:   "work=1",
					Operations: []Operation{{
						Type:   "write",
						Ratio:  100,
						Config: "op=1",
					}},
				}},
			}},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	stage := got.Workflow.Stages[0]
	if stage.Config != "wf=1;stage=1" {
		t.Fatalf("stage config = %q", stage.Config)
	}
	work := stage.Works[0]
	if work.Config != "wf=1;stage=1;work=1" {
		t.Fatalf("work config = %q", work.Config)
	}
	if work.Storage == nil || work.Storage.Type != "mock" || work.Storage.Config != "endpoint=top" {
		t.Fatalf("work storage = %#v", work.Storage)
	}
	if work.Interval != 5 || work.Division != "none" || work.Type != "normal" {
		t.Fatalf("defaults not applied: %#v", work)
	}
	if work.Operations[0].Config != "wf=1;stage=1;work=1;op=1" {
		t.Fatalf("op config = %q", work.Operations[0].Config)
	}
	if work.Operations[0].Division != "none" {
		t.Fatalf("op division = %q", work.Operations[0].Division)
	}
}

func TestNormalizeAndValidateSpecialWorks(t *testing.T) {
	got, err := NormalizeAndValidate(Workload{
		Name:    "special",
		Storage: &StorageSpec{Type: "sio"},
		Workflow: Workflow{
			Stages: []Stage{{
				Name: "main",
				Works: []Work{
					{Type: "prepare", Workers: 2, Config: "containers=c(1);objects=s(1,2);sizes=c(1)KB"},
					{Type: "mprepare", Workers: 1, Config: "containers=c(1);objects=s(1,2);sizes=c(1)KB"},
					{Type: "cleanup", Workers: 3, Config: "containers=c(1);objects=r(1,3)"},
					{Type: "init", Workers: 1, Config: "containers=c(1)"},
					{Type: "dispose", Workers: 1, Config: "containers=c(1)"},
					{Type: "delay", Workers: 9},
				},
			}},
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	prepare := got.Workflow.Stages[0].Works[0]
	if prepare.Name != "prepare" || prepare.TotalOps != 2 || prepare.Division != "object" {
		t.Fatalf("prepare = %#v", prepare)
	}
	if prepare.Operations[0].Type != "prepare" || !strings.Contains(prepare.Operations[0].Config, "createContainer=false") {
		t.Fatalf("prepare op = %#v", prepare.Operations[0])
	}

	mprepare := got.Workflow.Stages[0].Works[1]
	if mprepare.Name != "mprepare" || mprepare.Operations[0].Type != "mprepare" {
		t.Fatalf("mprepare = %#v", mprepare)
	}

	cleanup := got.Workflow.Stages[0].Works[2]
	if cleanup.Name != "cleanup" || cleanup.TotalOps != 3 {
		t.Fatalf("cleanup = %#v", cleanup)
	}
	if !strings.Contains(cleanup.Operations[0].Config, "deleteContainer=false") {
		t.Fatalf("cleanup op = %#v", cleanup.Operations[0])
	}

	initWork := got.Workflow.Stages[0].Works[3]
	if initWork.Name != "init" || initWork.Division != "container" || initWork.TotalOps != 1 {
		t.Fatalf("init = %#v", initWork)
	}
	if !strings.Contains(initWork.Operations[0].Config, "objects=r(0,0);sizes=c(0)B") {
		t.Fatalf("init op = %#v", initWork.Operations[0])
	}

	dispose := got.Workflow.Stages[0].Works[4]
	if dispose.Name != "dispose" || dispose.Operations[0].Type != "dispose" {
		t.Fatalf("dispose = %#v", dispose)
	}

	delay := got.Workflow.Stages[0].Works[5]
	if delay.Name != "delay" || delay.Workers != 1 || delay.TotalOps != 1 || delay.Division != "none" {
		t.Fatalf("delay = %#v", delay)
	}
}

func TestNormalizeAndValidateRejectsSIOOnlyOperationOnNonSIOStorage(t *testing.T) {
	_, err := NormalizeAndValidate(Workload{
		Name:    "non-sio",
		Storage: &StorageSpec{Type: "s3"},
		Workflow: Workflow{
			Stages: []Stage{{
				Name: "main",
				Works: []Work{{
					Name:     "main",
					Workers:  1,
					TotalOps: 1,
					Operations: []Operation{{
						Type:   "mwrite",
						Ratio:  100,
						Config: "containers=c(1);objects=c(1);sizes=c(1)KB",
					}},
				}},
			}},
		},
	})
	if err == nil || !strings.Contains(err.Error(), `requires sio-compatible storage`) {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNormalizeAndValidateFiltersZeroRatioOperations(t *testing.T) {
	got, err := NormalizeAndValidate(Workload{
		Name:    "filter-zero",
		Storage: &StorageSpec{Type: "mock"},
		Workflow: Workflow{
			Stages: []Stage{{
				Name: "main",
				Works: []Work{{
					Name:     "main",
					Workers:  1,
					TotalOps: 1,
					Operations: []Operation{
						{Type: "write", Ratio: 100, Config: "containers=c(1);objects=c(1);sizes=c(1)KB"},
						{Type: "read", Ratio: 0, Config: "containers=c(1);objects=c(1)"},
					},
				}},
			}},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	ops := got.Workflow.Stages[0].Works[0].Operations
	if len(ops) != 1 || ops[0].Type != "write" {
		t.Fatalf("ops = %#v", ops)
	}
}

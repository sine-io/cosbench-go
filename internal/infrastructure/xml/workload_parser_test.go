package xml

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestParseS3SampleWorkload(t *testing.T) {
	path := filepath.Clean("../../../testdata/legacy/s3-config-sample.xml")
	wl, err := ParseWorkloadFile(path)
	if err != nil {
		t.Fatalf("ParseWorkloadFile(): %v", err)
	}
	if wl.Name != "s3-sample" {
		t.Fatalf("workload name = %q", wl.Name)
	}
	if wl.Storage == nil || wl.Storage.Type != "s3" {
		t.Fatalf("storage = %#v", wl.Storage)
	}
	if len(wl.Workflow.Stages) != 5 {
		t.Fatalf("stages = %d", len(wl.Workflow.Stages))
	}
	mainStage := wl.Workflow.Stages[2]
	if mainStage.Name != "main" {
		t.Fatalf("stage[2].name = %q", mainStage.Name)
	}
	if len(mainStage.Works) != 1 {
		t.Fatalf("main works = %d", len(mainStage.Works))
	}
	work := mainStage.Works[0]
	if len(work.Operations) != 2 {
		t.Fatalf("main operations = %d", len(work.Operations))
	}
	if work.Operations[0].Type != "read" || work.Operations[0].Ratio != 80 {
		t.Fatalf("op0 = %#v", work.Operations[0])
	}
	if work.Operations[1].Type != "write" || work.Operations[1].Ratio != 20 {
		t.Fatalf("op1 = %#v", work.Operations[1])
	}
}

func TestParseSIOSampleWorkload(t *testing.T) {
	path := filepath.Clean("../../../testdata/legacy/sio-config-sample.xml")
	wl, err := ParseWorkloadFile(path)
	if err != nil {
		t.Fatalf("ParseWorkloadFile(): %v", err)
	}
	if len(wl.Workflow.Stages) == 0 {
		t.Fatal("expected stages")
	}
	stage := wl.Workflow.Stages[0]
	if stage.Works[0].Type != "mprepare" {
		t.Fatalf("first work type = %q", stage.Works[0].Type)
	}
	if stage.Works[0].Storage == nil || stage.Works[0].Storage.Type != "sio" {
		t.Fatalf("first work storage = %#v", stage.Works[0].Storage)
	}
	if len(stage.Works[0].Operations) != 1 || stage.Works[0].Operations[0].Type != "mprepare" {
		t.Fatalf("first work operations = %#v", stage.Works[0].Operations)
	}
}

func TestParseInheritanceSubsetWorkload(t *testing.T) {
	path := filepath.Clean("../../../testdata/workloads/xml-inheritance-subset.xml")
	wl, err := ParseWorkloadFile(path)
	if err != nil {
		t.Fatalf("ParseWorkloadFile(): %v", err)
	}
	stage := wl.Workflow.Stages[0]
	if stage.Config != "wf=1;stage=1" {
		t.Fatalf("stage config = %q", stage.Config)
	}
	if len(stage.Works) != 2 {
		t.Fatalf("works = %d", len(stage.Works))
	}
	first := stage.Works[0]
	if first.Config != "wf=1;stage=1;work=1" {
		t.Fatalf("first work config = %q", first.Config)
	}
	if first.Storage == nil || !strings.Contains(first.Storage.Config, "endpoint=http://stage") {
		t.Fatalf("first work storage = %#v", first.Storage)
	}
	if len(first.Operations) != 1 || first.Operations[0].Ratio != 100 {
		t.Fatalf("first work operations = %#v", first.Operations)
	}
	second := stage.Works[1]
	if second.Storage == nil || !strings.Contains(second.Storage.Config, "endpoint=http://work") {
		t.Fatalf("second work storage = %#v", second.Storage)
	}
}

func TestParseAttributeSubsetWorkload(t *testing.T) {
	path := filepath.Clean("../../../testdata/workloads/xml-attribute-subset.xml")
	wl, err := ParseWorkloadFile(path)
	if err != nil {
		t.Fatalf("ParseWorkloadFile(): %v", err)
	}
	if wl.Trigger != "nightly" {
		t.Fatalf("trigger = %q", wl.Trigger)
	}
	stage := wl.Workflow.Stages[0]
	if stage.ClosureDelay != 7 || stage.Trigger != "after-upload" {
		t.Fatalf("stage attrs = %#v", stage)
	}
	work := stage.Works[0]
	if work.Interval != 9 || work.Division != "object" || work.RampUp != 4 || work.RampDown != 5 || work.Driver != "driver-a" {
		t.Fatalf("work attrs = %#v", work)
	}
	if work.Operations[0].ID != "write-1" || work.Operations[0].Division != "object" {
		t.Fatalf("operation attrs = %#v", work.Operations[0])
	}
}

func TestParseSpecialOpsSubsetWorkload(t *testing.T) {
	path := filepath.Clean("../../../testdata/workloads/xml-special-ops-subset.xml")
	wl, err := ParseWorkloadFile(path)
	if err != nil {
		t.Fatalf("ParseWorkloadFile(): %v", err)
	}
	stage := wl.Workflow.Stages[0]
	if len(stage.Works) != 4 {
		t.Fatalf("works = %d", len(stage.Works))
	}
	delay := stage.Works[0]
	if delay.Name != "delay" || delay.Workers != 1 || delay.TotalOps != 1 || len(delay.Operations) != 1 {
		t.Fatalf("delay = %#v", delay)
	}
	cleanup := stage.Works[1]
	if cleanup.Name != "cleanup" || !strings.Contains(cleanup.Operations[0].Config, "deleteContainer=false") {
		t.Fatalf("cleanup = %#v", cleanup)
	}
	local := stage.Works[2]
	if local.Operations[0].Type != "localwrite" || local.Storage == nil || local.Storage.Type != "sio" {
		t.Fatalf("localwrite = %#v", local)
	}
	mfile := stage.Works[3]
	if mfile.Operations[0].Type != "mfilewrite" || mfile.Storage == nil || mfile.Storage.Type != "sio" {
		t.Fatalf("mfilewrite = %#v", mfile)
	}
}

func TestParseDelayStageSubsetWorkload(t *testing.T) {
	path := filepath.Clean("../../../testdata/workloads/xml-delay-stage-subset.xml")
	wl, err := ParseWorkloadFile(path)
	if err != nil {
		t.Fatalf("ParseWorkloadFile(): %v", err)
	}
	if wl.Trigger != "nightly" {
		t.Fatalf("trigger = %q", wl.Trigger)
	}
	if len(wl.Workflow.Stages) != 6 {
		t.Fatalf("stages = %d", len(wl.Workflow.Stages))
	}
	delayA := wl.Workflow.Stages[1]
	if delayA.ClosureDelay != 3 || delayA.Works[0].Operations[0].Type != "delay" {
		t.Fatalf("delayA = %#v", delayA)
	}
	delayB := wl.Workflow.Stages[3]
	if delayB.ClosureDelay != 4 || delayB.Works[0].Name != "delay" {
		t.Fatalf("delayB = %#v", delayB)
	}
}

func TestParseSplitRWSubsetWorkload(t *testing.T) {
	path := filepath.Clean("../../../testdata/workloads/xml-splitrw-subset.xml")
	wl, err := ParseWorkloadFile(path)
	if err != nil {
		t.Fatalf("ParseWorkloadFile(): %v", err)
	}
	mainWork := wl.Workflow.Stages[2].Works[0]
	if len(mainWork.Operations) != 2 {
		t.Fatalf("operations = %#v", mainWork.Operations)
	}
	if !strings.Contains(mainWork.Operations[0].Config, "containers=u(1,2)") {
		t.Fatalf("read config = %q", mainWork.Operations[0].Config)
	}
	if !strings.Contains(mainWork.Operations[1].Config, "containers=u(3,4)") {
		t.Fatalf("write config = %q", mainWork.Operations[1].Config)
	}
}

func TestParseReuseDataSubsetWorkload(t *testing.T) {
	path := filepath.Clean("../../../testdata/workloads/mock-reusedata-subset.xml")
	wl, err := ParseWorkloadFile(path)
	if err != nil {
		t.Fatalf("ParseWorkloadFile(): %v", err)
	}
	if len(wl.Workflow.Stages) != 6 {
		t.Fatalf("stages = %d", len(wl.Workflow.Stages))
	}
	if wl.Workflow.Stages[2].Name != "main-read" || wl.Workflow.Stages[3].Name != "main-list" {
		t.Fatalf("unexpected main stage layout: %#v", wl.Workflow.Stages)
	}
}

func TestParseCompatStorageSubsetWorkload(t *testing.T) {
	path := filepath.Clean("../../../testdata/workloads/xml-compat-storage-subset.xml")
	wl, err := ParseWorkloadFile(path)
	if err != nil {
		t.Fatalf("ParseWorkloadFile(): %v", err)
	}
	if len(wl.Workflow.Stages) != 2 {
		t.Fatalf("stages = %d", len(wl.Workflow.Stages))
	}
	siov1 := wl.Workflow.Stages[0].Works[0]
	if siov1.Storage == nil || siov1.Storage.Type != "siov1" {
		t.Fatalf("siov1 storage = %#v", siov1.Storage)
	}
	gdas := wl.Workflow.Stages[1].Works[0]
	if gdas.Storage == nil || gdas.Storage.Type != "gdas" {
		t.Fatalf("gdas storage = %#v", gdas.Storage)
	}
}

func TestParseRangePrefetchSubsetWorkload(t *testing.T) {
	path := filepath.Clean("../../../testdata/workloads/xml-range-prefetch-subset.xml")
	wl, err := ParseWorkloadFile(path)
	if err != nil {
		t.Fatalf("ParseWorkloadFile(): %v", err)
	}
	prefetch := wl.Workflow.Stages[0].Works[0]
	if prefetch.Storage == nil || !strings.Contains(prefetch.Storage.Config, "is_prefetch=true") {
		t.Fatalf("prefetch storage = %#v", prefetch.Storage)
	}
	rangeStage := wl.Workflow.Stages[1].Works[0]
	if rangeStage.Storage == nil || !strings.Contains(rangeStage.Storage.Config, "is_range_request=true") || !strings.Contains(rangeStage.Storage.Config, "file_length=15000000") || !strings.Contains(rangeStage.Storage.Config, "chunk_length=5000000") {
		t.Fatalf("range storage = %#v", rangeStage.Storage)
	}
}

func TestParseAuthToleratedSubsetWorkload(t *testing.T) {
	path := filepath.Clean("../../../testdata/workloads/xml-auth-tolerated-subset.xml")
	wl, err := ParseWorkloadFile(path)
	if err != nil {
		t.Fatalf("ParseWorkloadFile(): %v", err)
	}
	if len(wl.Workflow.Stages) != 1 {
		t.Fatalf("stages = %d", len(wl.Workflow.Stages))
	}
	work := wl.Workflow.Stages[0].Works[0]
	if work.Storage == nil || work.Storage.Type != "mock" {
		t.Fatalf("work storage = %#v", work.Storage)
	}
	if work.Operations[0].Type != "write" {
		t.Fatalf("operation = %#v", work.Operations[0])
	}
}

func TestParseAuthNoneSubsetWorkload(t *testing.T) {
	path := filepath.Clean("../../../testdata/workloads/xml-auth-none-subset.xml")
	wl, err := ParseWorkloadFile(path)
	if err != nil {
		t.Fatalf("ParseWorkloadFile(): %v", err)
	}
	if wl.Trigger != "demo" {
		t.Fatalf("trigger = %q", wl.Trigger)
	}
	if len(wl.Workflow.Stages) != 2 {
		t.Fatalf("stages = %d", len(wl.Workflow.Stages))
	}
	if wl.Workflow.Stages[0].Works[0].Type != "init" || wl.Workflow.Stages[1].Works[0].Operations[0].Type != "read" {
		t.Fatalf("unexpected parsed structure: %#v", wl.Workflow.Stages)
	}
}

func TestParseAuthInheritanceSubsetWorkload(t *testing.T) {
	path := filepath.Clean("../../../testdata/workloads/xml-auth-inheritance-subset.xml")
	wl, err := ParseWorkloadFile(path)
	if err != nil {
		t.Fatalf("ParseWorkloadFile(): %v", err)
	}
	if wl.Auth == nil || wl.Auth.Type != "basic" || wl.Auth.Config != "username=workload;password=root" {
		t.Fatalf("workload auth = %#v", wl.Auth)
	}
	if len(wl.Workflow.Stages) != 2 {
		t.Fatalf("stages = %d", len(wl.Workflow.Stages))
	}
	stageAuth := wl.Workflow.Stages[0]
	if stageAuth.Auth == nil || stageAuth.Auth.Config != "username=stage;password=override" {
		t.Fatalf("stage auth = %#v", stageAuth.Auth)
	}
	inheritedWork := stageAuth.Works[0]
	if inheritedWork.Auth == nil || inheritedWork.Auth.Config != "username=stage;password=override" {
		t.Fatalf("inherited work auth = %#v", inheritedWork.Auth)
	}
	explicitWork := stageAuth.Works[1]
	if explicitWork.Auth == nil || explicitWork.Auth.Config != "username=work;password=leaf" {
		t.Fatalf("explicit work auth = %#v", explicitWork.Auth)
	}
	fallbackStage := wl.Workflow.Stages[1]
	if fallbackStage.Auth == nil || fallbackStage.Auth.Config != "username=workload;password=root" {
		t.Fatalf("fallback stage auth = %#v", fallbackStage.Auth)
	}
	fallbackWork := fallbackStage.Works[0]
	if fallbackWork.Auth == nil || fallbackWork.Auth.Config != "username=workload;password=root" {
		t.Fatalf("fallback work auth = %#v", fallbackWork.Auth)
	}
}

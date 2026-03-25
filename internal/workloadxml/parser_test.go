package workloadxml

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestParseFile(t *testing.T) {
	path := filepath.Clean("../../testdata/legacy/s3-config-sample.xml")
	parsed, raw, err := ParseFile(path)
	if err != nil {
		t.Fatalf("ParseFile(): %v", err)
	}
	if parsed.Name != "s3-sample" {
		t.Fatalf("workload name = %q", parsed.Name)
	}
	if !strings.Contains(string(raw), "<workload") {
		t.Fatal("expected raw xml contents")
	}
	if len(parsed.Workflow.Stages) != 5 {
		t.Fatalf("stages = %d", len(parsed.Workflow.Stages))
	}
}

func TestParseRepresentativeFixture(t *testing.T) {
	path := filepath.Clean("../../testdata/workloads/s3-active-subset.xml")
	parsed, _, err := ParseFile(path)
	if err != nil {
		t.Fatalf("ParseFile(): %v", err)
	}
	if parsed.Name != "s3-active-subset" {
		t.Fatalf("workload name = %q", parsed.Name)
	}
	if len(parsed.Workflow.Stages) != 2 {
		t.Fatalf("stages = %d", len(parsed.Workflow.Stages))
	}
	if parsed.Workflow.Stages[1].Works[0].Operations[0].Ratio != 70 {
		t.Fatalf("unexpected read ratio: %#v", parsed.Workflow.Stages[1].Works[0].Operations)
	}
}

func TestParseSIOMultipartFixture(t *testing.T) {
	path := filepath.Clean("../../testdata/workloads/sio-multipart-subset.xml")
	parsed, _, err := ParseFile(path)
	if err != nil {
		t.Fatalf("ParseFile(): %v", err)
	}
	if parsed.Name != "sio-multipart-subset" {
		t.Fatalf("workload name = %q", parsed.Name)
	}
	if len(parsed.Workflow.Stages) != 2 {
		t.Fatalf("stages = %d", len(parsed.Workflow.Stages))
	}
	work := parsed.Workflow.Stages[1].Works[0]
	if work.Operations[0].Type != "mwrite" {
		t.Fatalf("unexpected op: %#v", work.Operations)
	}
	if work.Storage == nil || work.Storage.Type != "sio" {
		t.Fatalf("unexpected storage: %#v", work.Storage)
	}
}

func TestParseInheritanceSubsetFixture(t *testing.T) {
	path := filepath.Clean("../../testdata/workloads/xml-inheritance-subset.xml")
	parsed, _, err := ParseFile(path)
	if err != nil {
		t.Fatalf("ParseFile(): %v", err)
	}
	stage := parsed.Workflow.Stages[0]
	first := stage.Works[0]
	if first.Config != "wf=1;stage=1;work=1" {
		t.Fatalf("work config = %q", first.Config)
	}
	if len(first.Operations) != 1 || first.Operations[0].Ratio != 100 {
		t.Fatalf("operations = %#v", first.Operations)
	}
	if first.Operations[0].Config != "wf=1;stage=1;work=1;op=1;containers=c(1);objects=c(1)" {
		t.Fatalf("op config = %q", first.Operations[0].Config)
	}
	if first.Storage == nil || !strings.Contains(first.Storage.Config, "endpoint=http://stage") {
		t.Fatalf("work storage = %#v", first.Storage)
	}
	second := stage.Works[1]
	if second.Storage == nil || !strings.Contains(second.Storage.Config, "endpoint=http://work") {
		t.Fatalf("work override storage = %#v", second.Storage)
	}
}

func TestParseAttributeSubsetFixture(t *testing.T) {
	path := filepath.Clean("../../testdata/workloads/xml-attribute-subset.xml")
	parsed, _, err := ParseFile(path)
	if err != nil {
		t.Fatalf("ParseFile(): %v", err)
	}
	if parsed.Trigger != "nightly" {
		t.Fatalf("trigger = %q", parsed.Trigger)
	}
	stage := parsed.Workflow.Stages[0]
	if stage.ClosureDelay != 7 || stage.Trigger != "after-upload" {
		t.Fatalf("stage attrs = %#v", stage)
	}
	work := stage.Works[0]
	if work.Interval != 9 || work.Division != "object" || work.RampUp != 4 || work.RampDown != 5 || work.Driver != "driver-a" {
		t.Fatalf("work attrs = %#v", work)
	}
}

func TestParseSpecialOpsSubsetFixture(t *testing.T) {
	path := filepath.Clean("../../testdata/workloads/xml-special-ops-subset.xml")
	parsed, _, err := ParseFile(path)
	if err != nil {
		t.Fatalf("ParseFile(): %v", err)
	}
	stage := parsed.Workflow.Stages[0]
	delay := stage.Works[0]
	if delay.Name != "delay" || delay.Workers != 1 || delay.TotalOps != 1 || delay.Operations[0].Type != "delay" {
		t.Fatalf("delay = %#v", delay)
	}
	cleanup := stage.Works[1]
	if cleanup.Name != "cleanup" || !strings.Contains(cleanup.Operations[0].Config, "deleteContainer=false") {
		t.Fatalf("cleanup = %#v", cleanup)
	}
	if stage.Works[2].Operations[0].Type != "localwrite" || stage.Works[3].Operations[0].Type != "mfilewrite" {
		t.Fatalf("special ops = %#v", stage.Works)
	}
}

func TestParseDelayStageSubsetFixture(t *testing.T) {
	path := filepath.Clean("../../testdata/workloads/xml-delay-stage-subset.xml")
	parsed, _, err := ParseFile(path)
	if err != nil {
		t.Fatalf("ParseFile(): %v", err)
	}
	if len(parsed.Workflow.Stages) != 6 {
		t.Fatalf("stages = %d", len(parsed.Workflow.Stages))
	}
	delayA := parsed.Workflow.Stages[1]
	if delayA.ClosureDelay != 3 || delayA.Works[0].Operations[0].Type != "delay" {
		t.Fatalf("delayA = %#v", delayA)
	}
	prepare := parsed.Workflow.Stages[2]
	if prepare.ClosureDelay != 2 || prepare.Works[0].Name != "prepare" {
		t.Fatalf("prepare = %#v", prepare)
	}
}

func TestParseSplitRWSubsetFixture(t *testing.T) {
	path := filepath.Clean("../../testdata/workloads/xml-splitrw-subset.xml")
	parsed, _, err := ParseFile(path)
	if err != nil {
		t.Fatalf("ParseFile(): %v", err)
	}
	mainWork := parsed.Workflow.Stages[2].Works[0]
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

func TestParseReuseDataSubsetFixture(t *testing.T) {
	path := filepath.Clean("../../testdata/workloads/mock-reusedata-subset.xml")
	parsed, _, err := ParseFile(path)
	if err != nil {
		t.Fatalf("ParseFile(): %v", err)
	}
	if len(parsed.Workflow.Stages) != 6 {
		t.Fatalf("stages = %d", len(parsed.Workflow.Stages))
	}
	if parsed.Workflow.Stages[2].Works[0].Operations[0].Type != "read" {
		t.Fatalf("main-read = %#v", parsed.Workflow.Stages[2].Works[0])
	}
	if parsed.Workflow.Stages[3].Works[0].Operations[0].Type != "list" {
		t.Fatalf("main-list = %#v", parsed.Workflow.Stages[3].Works[0])
	}
}

func TestParseCompatStorageSubsetFixture(t *testing.T) {
	path := filepath.Clean("../../testdata/workloads/xml-compat-storage-subset.xml")
	parsed, _, err := ParseFile(path)
	if err != nil {
		t.Fatalf("ParseFile(): %v", err)
	}
	if parsed.Workflow.Stages[0].Works[0].Storage == nil || parsed.Workflow.Stages[0].Works[0].Storage.Type != "siov1" {
		t.Fatalf("siov1 storage = %#v", parsed.Workflow.Stages[0].Works[0].Storage)
	}
	if parsed.Workflow.Stages[1].Works[0].Storage == nil || parsed.Workflow.Stages[1].Works[0].Storage.Type != "gdas" {
		t.Fatalf("gdas storage = %#v", parsed.Workflow.Stages[1].Works[0].Storage)
	}
}

func TestParseRangePrefetchSubsetFixture(t *testing.T) {
	path := filepath.Clean("../../testdata/workloads/xml-range-prefetch-subset.xml")
	parsed, _, err := ParseFile(path)
	if err != nil {
		t.Fatalf("ParseFile(): %v", err)
	}
	if !strings.Contains(parsed.Workflow.Stages[0].Works[0].Storage.Config, "is_prefetch=true") {
		t.Fatalf("prefetch config = %q", parsed.Workflow.Stages[0].Works[0].Storage.Config)
	}
	rangeConfig := parsed.Workflow.Stages[1].Works[0].Storage.Config
	if !strings.Contains(rangeConfig, "is_range_request=true") || !strings.Contains(rangeConfig, "file_length=15000000") || !strings.Contains(rangeConfig, "chunk_length=5000000") {
		t.Fatalf("range config = %q", rangeConfig)
	}
}

func TestParseAuthToleratedSubsetFixture(t *testing.T) {
	path := filepath.Clean("../../testdata/workloads/xml-auth-tolerated-subset.xml")
	parsed, _, err := ParseFile(path)
	if err != nil {
		t.Fatalf("ParseFile(): %v", err)
	}
	if len(parsed.Workflow.Stages) != 1 {
		t.Fatalf("stages = %d", len(parsed.Workflow.Stages))
	}
	work := parsed.Workflow.Stages[0].Works[0]
	if work.Name != "main" || work.Storage == nil || work.Storage.Type != "mock" {
		t.Fatalf("work = %#v", work)
	}
}

func TestParseAuthNoneSubsetFixture(t *testing.T) {
	path := filepath.Clean("../../testdata/workloads/xml-auth-none-subset.xml")
	parsed, _, err := ParseFile(path)
	if err != nil {
		t.Fatalf("ParseFile(): %v", err)
	}
	if parsed.Trigger != "demo" {
		t.Fatalf("trigger = %q", parsed.Trigger)
	}
	if len(parsed.Workflow.Stages) != 2 {
		t.Fatalf("stages = %d", len(parsed.Workflow.Stages))
	}
	if parsed.Workflow.Stages[0].Works[0].Name != "init" || parsed.Workflow.Stages[1].Works[0].Operations[0].Type != "read" {
		t.Fatalf("unexpected structure: %#v", parsed.Workflow.Stages)
	}
}

func TestParseInvalidXML(t *testing.T) {
	_, err := Parse([]byte("<workload"))
	if err == nil {
		t.Fatal("expected parse error")
	}
	if !strings.Contains(err.Error(), "parse workload xml") {
		t.Fatalf("unexpected error: %v", err)
	}
}

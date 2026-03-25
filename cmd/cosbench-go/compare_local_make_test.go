package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestCompareLocalPrunesStaleOutputs(t *testing.T) {
	makeBin, err := exec.LookPath("make")
	if err != nil {
		t.Fatalf("look path make: %v", err)
	}
	goBin, err := exec.LookPath("go")
	if err != nil {
		t.Fatalf("look path go: %v", err)
	}

	rootDir := filepath.Clean("../..")
	outputDir := filepath.Join(t.TempDir(), "compare-local")
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		t.Fatalf("mkdir output dir: %v", err)
	}
	staleFile := filepath.Join(outputDir, "stale.json")
	if err := os.WriteFile(staleFile, []byte("{}\n"), 0o644); err != nil {
		t.Fatalf("seed stale file: %v", err)
	}

	cmd := exec.Command(makeBin, "compare-local", "GO="+goBin, "COMPARE_LOCAL_OUTPUT_DIR="+outputDir)
	cmd.Dir = rootDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("make compare-local failed: %v\n%s", err, output)
	}

	if _, err := os.Stat(staleFile); err == nil {
		t.Fatalf("expected stale output %s to be removed", staleFile)
	}

	for _, name := range []string{
		"s3-active-subset.json",
		"mock-stage-aware.json",
		"mock-reusedata-subset.json",
		"xml-splitrw-subset.json",
	} {
		if _, err := os.Stat(filepath.Join(outputDir, name)); err != nil {
			t.Fatalf("expected fresh output %s: %v", name, err)
		}
	}

	indexData, err := os.ReadFile(filepath.Join(outputDir, "index.json"))
	if err != nil {
		t.Fatalf("read index: %v", err)
	}
	var payload struct {
		Meta struct {
			Filter       string `json:"filter"`
			FixtureCount int    `json:"fixture_count"`
		} `json:"meta"`
		Fixtures []struct {
			Name     string `json:"name"`
			Workload string `json:"workload"`
			Summary  string `json:"summary"`
			Stages   int    `json:"stages"`
			Works    int    `json:"works"`
			Samples  int    `json:"samples"`
			Errors   int64  `json:"errors"`
		} `json:"fixtures"`
	}
	if err := json.Unmarshal(indexData, &payload); err != nil {
		t.Fatalf("unmarshal index: %v", err)
	}
	if len(payload.Fixtures) != 4 {
		t.Fatalf("fixtures = %#v", payload.Fixtures)
	}
	if payload.Meta.Filter != "" {
		t.Fatalf("meta filter = %q", payload.Meta.Filter)
	}
	if payload.Meta.FixtureCount != 4 {
		t.Fatalf("meta fixture_count = %d", payload.Meta.FixtureCount)
	}
	for _, fixture := range payload.Fixtures {
		if fixture.Stages == 0 {
			t.Fatalf("missing stages in fixture %#v", fixture)
		}
		if fixture.Works == 0 {
			t.Fatalf("missing works in fixture %#v", fixture)
		}
		if fixture.Samples == 0 {
			t.Fatalf("missing samples in fixture %#v", fixture)
		}
		if fixture.Errors < 0 {
			t.Fatalf("unexpected errors in fixture %#v", fixture)
		}
	}
}

func TestCompareLocalRejectsUnsafeOutputDir(t *testing.T) {
	makeBin, err := exec.LookPath("make")
	if err != nil {
		t.Fatalf("look path make: %v", err)
	}
	goBin, err := exec.LookPath("go")
	if err != nil {
		t.Fatalf("look path go: %v", err)
	}

	rootDir := filepath.Clean("../..")
	unsafeDir := t.TempDir()
	cmd := exec.Command(makeBin, "compare-local", "GO="+goBin, "COMPARE_LOCAL_OUTPUT_DIR="+unsafeDir)
	cmd.Dir = rootDir
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("expected compare-local to reject unsafe output dir %s\n%s", unsafeDir, output)
	}
}

func TestCompareLocalFilterRunsSingleFixture(t *testing.T) {
	makeBin, err := exec.LookPath("make")
	if err != nil {
		t.Fatalf("look path make: %v", err)
	}
	goBin, err := exec.LookPath("go")
	if err != nil {
		t.Fatalf("look path go: %v", err)
	}

	rootDir := filepath.Clean("../..")
	outputDir := filepath.Join(t.TempDir(), "compare-local")
	cmd := exec.Command(
		makeBin,
		"compare-local",
		"GO="+goBin,
		"COMPARE_LOCAL_OUTPUT_DIR="+outputDir,
		"COMPARE_LOCAL_FILTER=mock-stage-aware",
	)
	cmd.Dir = rootDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("make compare-local failed: %v\n%s", err, output)
	}

	if _, err := os.Stat(filepath.Join(outputDir, "mock-stage-aware.json")); err != nil {
		t.Fatalf("expected filtered output: %v", err)
	}
	if _, err := os.Stat(filepath.Join(outputDir, "s3-active-subset.json")); err == nil {
		t.Fatal("unexpected unfiltered output")
	}

	indexData, err := os.ReadFile(filepath.Join(outputDir, "index.json"))
	if err != nil {
		t.Fatalf("read index: %v", err)
	}
	var payload struct {
		Meta struct {
			Filter       string `json:"filter"`
			FixtureCount int    `json:"fixture_count"`
		} `json:"meta"`
		Fixtures []struct {
			Name string `json:"name"`
		} `json:"fixtures"`
	}
	if err := json.Unmarshal(indexData, &payload); err != nil {
		t.Fatalf("unmarshal index: %v", err)
	}
	if len(payload.Fixtures) != 1 || payload.Fixtures[0].Name != "mock-stage-aware" {
		t.Fatalf("fixtures = %#v", payload.Fixtures)
	}
	if payload.Meta.Filter != "mock-stage-aware" {
		t.Fatalf("meta filter = %q", payload.Meta.Filter)
	}
	if payload.Meta.FixtureCount != 1 {
		t.Fatalf("meta fixture_count = %d", payload.Meta.FixtureCount)
	}
}

func TestCompareLocalFilterRejectsUnknownFixture(t *testing.T) {
	makeBin, err := exec.LookPath("make")
	if err != nil {
		t.Fatalf("look path make: %v", err)
	}
	goBin, err := exec.LookPath("go")
	if err != nil {
		t.Fatalf("look path go: %v", err)
	}

	rootDir := filepath.Clean("../..")
	outputDir := filepath.Join(t.TempDir(), "compare-local")
	cmd := exec.Command(
		makeBin,
		"compare-local",
		"GO="+goBin,
		"COMPARE_LOCAL_OUTPUT_DIR="+outputDir,
		"COMPARE_LOCAL_FILTER=does-not-exist",
	)
	cmd.Dir = rootDir
	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("expected compare-local to reject unknown fixture\n%s", output)
	}
	if !strings.Contains(string(output), "does-not-exist") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestCompareLocalListShowsFixtureNames(t *testing.T) {
	makeBin, err := exec.LookPath("make")
	if err != nil {
		t.Fatalf("look path make: %v", err)
	}

	rootDir := filepath.Clean("../..")
	cmd := exec.Command(makeBin, "--no-print-directory", "compare-local-list")
	cmd.Dir = rootDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("make compare-local-list failed: %v\n%s", err, output)
	}

	lines := strings.Fields(strings.TrimSpace(string(output)))
	want := []string{
		"s3-active-subset",
		"mock-stage-aware",
		"mock-reusedata-subset",
		"xml-splitrw-subset",
	}
	if len(lines) != len(want) {
		t.Fatalf("lines = %#v", lines)
	}
	for i, name := range want {
		if lines[i] != name {
			t.Fatalf("lines = %#v", lines)
		}
	}
}

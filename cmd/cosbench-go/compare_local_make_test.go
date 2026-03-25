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
	summaryData, err := os.ReadFile(filepath.Join(outputDir, "summary.md"))
	if err != nil {
		t.Fatalf("read summary: %v", err)
	}
	if !strings.Contains(string(summaryData), "## Compare Local") || !strings.Contains(string(summaryData), "| Fixture | Workload |") || !strings.Contains(string(summaryData), "Filter: `all`") {
		t.Fatalf("unexpected summary: %s", summaryData)
	}
	var payload struct {
		Meta struct {
			Filter       string `json:"filter"`
			FixtureCount int    `json:"fixture_count"`
			GeneratedAt  string `json:"generated_at"`
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
	if payload.Meta.Filter != "all" {
		t.Fatalf("meta filter = %q", payload.Meta.Filter)
	}
	if payload.Meta.FixtureCount != 4 {
		t.Fatalf("meta fixture_count = %d", payload.Meta.FixtureCount)
	}
	if payload.Meta.GeneratedAt == "" {
		t.Fatal("missing meta generated_at")
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
			GeneratedAt  string `json:"generated_at"`
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
	if payload.Meta.GeneratedAt == "" {
		t.Fatal("missing meta generated_at")
	}
}

func TestCompareLocalFilterRunsFixtureSubset(t *testing.T) {
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
		"COMPARE_LOCAL_FILTER=mock-stage-aware,xml-splitrw-subset",
	)
	cmd.Dir = rootDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("make compare-local failed: %v\n%s", err, output)
	}

	if _, err := os.Stat(filepath.Join(outputDir, "mock-stage-aware.json")); err != nil {
		t.Fatalf("expected subset output: %v", err)
	}
	if _, err := os.Stat(filepath.Join(outputDir, "xml-splitrw-subset.json")); err != nil {
		t.Fatalf("expected subset output: %v", err)
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
	if payload.Meta.Filter != "mock-stage-aware,xml-splitrw-subset" {
		t.Fatalf("meta filter = %q", payload.Meta.Filter)
	}
	if payload.Meta.FixtureCount != 2 {
		t.Fatalf("meta fixture_count = %d", payload.Meta.FixtureCount)
	}
	if len(payload.Fixtures) != 2 {
		t.Fatalf("fixtures = %#v", payload.Fixtures)
	}
}

func TestCompareLocalFilterAcceptsAllAlias(t *testing.T) {
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
		"COMPARE_LOCAL_FILTER=all",
	)
	cmd.Dir = rootDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("make compare-local failed: %v\n%s", err, output)
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
	}
	if err := json.Unmarshal(indexData, &payload); err != nil {
		t.Fatalf("unmarshal index: %v", err)
	}
	if payload.Meta.Filter != "all" {
		t.Fatalf("meta filter = %q", payload.Meta.Filter)
	}
	if payload.Meta.FixtureCount != 4 {
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

func TestCompareLocalListJSONShowsFixtureMetadata(t *testing.T) {
	makeBin, err := exec.LookPath("make")
	if err != nil {
		t.Fatalf("look path make: %v", err)
	}

	rootDir := filepath.Clean("../..")
	cmd := exec.Command(makeBin, "--no-print-directory", "compare-local-list-json")
	cmd.Dir = rootDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("make compare-local-list-json failed: %v\n%s", err, output)
	}

	var payload []struct {
		Name     string `json:"name"`
		Workload string `json:"workload"`
	}
	if err := json.Unmarshal(output, &payload); err != nil {
		t.Fatalf("unmarshal output: %v\n%s", err, output)
	}
	if len(payload) != 4 {
		t.Fatalf("payload = %#v", payload)
	}
	if payload[0].Name != "s3-active-subset" || payload[0].Workload != "testdata/workloads/s3-active-subset.xml" {
		t.Fatalf("payload = %#v", payload)
	}
}

func TestWorktreeAuditTargetRuns(t *testing.T) {
	makeBin, err := exec.LookPath("make")
	if err != nil {
		t.Fatalf("look path make: %v", err)
	}

	rootDir := filepath.Clean("../..")
	cmd := exec.Command(makeBin, "--no-print-directory", "worktree-audit")
	cmd.Dir = rootDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("make worktree-audit failed: %v\n%s", err, output)
	}
	text := string(output)
	if !strings.Contains(text, "PATH\tBRANCH\tCURRENT\tSTATE\tDETAILS") {
		t.Fatalf("unexpected audit output: %s", text)
	}
	if !strings.Contains(text, "\tyes\t") {
		t.Fatalf("expected current row marker: %s", text)
	}
}

func TestWorktreeAuditTargetSupportsBaseRefOverride(t *testing.T) {
	makeBin, err := exec.LookPath("make")
	if err != nil {
		t.Fatalf("look path make: %v", err)
	}

	rootDir := filepath.Clean("../..")
	cmd := exec.Command(makeBin, "--no-print-directory", "worktree-audit", "WORKTREE_AUDIT_BASE_REF=HEAD")
	cmd.Dir = rootDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("make worktree-audit failed: %v\n%s", err, output)
	}
	text := string(output)
	if !strings.Contains(text, "PATH\tBRANCH\tCURRENT\tSTATE\tDETAILS") {
		t.Fatalf("unexpected output: %s", text)
	}
	if strings.Contains(text, "origin/main") {
		t.Fatalf("unexpected default base ref output: %s", text)
	}
}

func TestWorktreeAuditJSONTargetRuns(t *testing.T) {
	makeBin, err := exec.LookPath("make")
	if err != nil {
		t.Fatalf("look path make: %v", err)
	}

	rootDir := filepath.Clean("../..")
	cmd := exec.Command(makeBin, "--no-print-directory", "worktree-audit-json")
	cmd.Dir = rootDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("make worktree-audit-json failed: %v\n%s", err, output)
	}

	var payload struct {
		Summary map[string]any   `json:"summary"`
		Rows    []map[string]any `json:"rows"`
	}
	if err := json.Unmarshal(output, &payload); err != nil {
		t.Fatalf("unmarshal output: %v\n%s", err, output)
	}
	if payload.Summary == nil {
		t.Fatalf("missing summary: %#v", payload)
	}
	if _, ok := payload.Summary["total"]; !ok {
		t.Fatalf("missing total: %#v", payload.Summary)
	}
	if _, ok := payload.Summary["base_ref"]; !ok {
		t.Fatalf("missing base_ref: %#v", payload.Summary)
	}
	if len(payload.Rows) == 0 {
		t.Fatal("expected at least one worktree entry")
	}
	seenMerged := false
	seenActive := false
	for _, row := range payload.Rows {
		state, _ := row["state"].(string)
		if state == "merged" {
			seenMerged = true
			if seenActive {
				t.Fatalf("expected merged rows before active rows: %#v", payload.Rows)
			}
			continue
		}
		if state == "active" {
			seenActive = true
		}
	}
	if seenMerged && payload.Rows[0]["state"] != "merged" {
		t.Fatalf("expected merged rows first when present: %#v", payload.Rows[:1])
	}
	if payload.Rows[0]["path"] == "" || payload.Rows[0]["branch"] == "" || payload.Rows[0]["state"] == "" {
		t.Fatalf("unexpected payload: %#v", payload.Rows[0])
	}
	if _, ok := payload.Rows[0]["current"]; !ok {
		t.Fatalf("missing current: %#v", payload.Rows[0])
	}
	if _, ok := payload.Rows[0]["ahead"]; !ok {
		t.Fatalf("missing ahead: %#v", payload.Rows[0])
	}
	if _, ok := payload.Rows[0]["behind"]; !ok {
		t.Fatalf("missing behind: %#v", payload.Rows[0])
	}
	ahead, ok := payload.Rows[0]["ahead"].(float64)
	if !ok || ahead < 0 {
		t.Fatalf("unexpected payload: %#v", payload.Rows[0])
	}
	behind, ok := payload.Rows[0]["behind"].(float64)
	if !ok || behind < 0 {
		t.Fatalf("unexpected payload: %#v", payload.Rows[0])
	}
	foundCurrent := false
	for _, row := range payload.Rows {
		current, ok := row["current"].(bool)
		if !ok {
			t.Fatalf("unexpected row: %#v", row)
		}
		if current {
			foundCurrent = true
		}
	}
	if !foundCurrent {
		t.Fatalf("expected one current row: %#v", payload.Rows)
	}
}

func TestWorktreeAuditMergedTargetRuns(t *testing.T) {
	makeBin, err := exec.LookPath("make")
	if err != nil {
		t.Fatalf("look path make: %v", err)
	}

	rootDir := filepath.Clean("../..")
	cmd := exec.Command(makeBin, "--no-print-directory", "worktree-audit-merged")
	cmd.Dir = rootDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("make worktree-audit-merged failed: %v\n%s", err, output)
	}
	text := string(output)
	if !strings.Contains(text, "PATH\tBRANCH\tCURRENT\tSTATE\tDETAILS") {
		t.Fatalf("unexpected audit output: %s", text)
	}
	if strings.Contains(text, "\tactive\t") {
		t.Fatalf("unexpected active row in merged-only output: %s", text)
	}
}

func TestWorktreeAuditStaleTargetRuns(t *testing.T) {
	makeBin, err := exec.LookPath("make")
	if err != nil {
		t.Fatalf("look path make: %v", err)
	}

	rootDir := filepath.Clean("../..")
	cmd := exec.Command(makeBin, "--no-print-directory", "worktree-audit-stale")
	cmd.Dir = rootDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("make worktree-audit-stale failed: %v\n%s", err, output)
	}
	text := string(output)
	if !strings.Contains(text, "PATH\tBRANCH\tCURRENT\tSTATE\tDETAILS") {
		t.Fatalf("unexpected output: %s", text)
	}
	if strings.Contains(text, "\tmerged\t") {
		t.Fatalf("unexpected merged row in stale-only output: %s", text)
	}
}

func TestWorktreeAuditMergedJSONTargetRuns(t *testing.T) {
	makeBin, err := exec.LookPath("make")
	if err != nil {
		t.Fatalf("look path make: %v", err)
	}

	rootDir := filepath.Clean("../..")
	cmd := exec.Command(makeBin, "--no-print-directory", "worktree-audit-merged-json")
	cmd.Dir = rootDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("make worktree-audit-merged-json failed: %v\n%s", err, output)
	}

	var payload struct {
		Rows []struct {
			Path   string `json:"path"`
			Branch string `json:"branch"`
			State  string `json:"state"`
		} `json:"rows"`
	}
	if err := json.Unmarshal(output, &payload); err != nil {
		t.Fatalf("unmarshal output: %v\n%s", err, output)
	}
	for _, row := range payload.Rows {
		if row.State != "merged" {
			t.Fatalf("unexpected row: %#v", row)
		}
	}
}

func TestWorktreePrunePlanTargetRuns(t *testing.T) {
	makeBin, err := exec.LookPath("make")
	if err != nil {
		t.Fatalf("look path make: %v", err)
	}

	rootDir := filepath.Clean("../..")
	cmd := exec.Command(makeBin, "--no-print-directory", "worktree-prune-plan")
	cmd.Dir = rootDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("make worktree-prune-plan failed: %v\n%s", err, output)
	}
	text := string(output)
	if !strings.Contains(text, "# Suggested cleanup commands") {
		t.Fatalf("unexpected output: %s", text)
	}
	if !strings.Contains(text, "git worktree remove") && !strings.Contains(text, "# no merged worktrees to prune") {
		t.Fatalf("unexpected output: %s", text)
	}
	if strings.Contains(text, "worktree-prune-safety") {
		t.Fatalf("unexpected self-removal plan: %s", text)
	}
}

func TestWorktreePrunePlanJSONTargetRuns(t *testing.T) {
	makeBin, err := exec.LookPath("make")
	if err != nil {
		t.Fatalf("look path make: %v", err)
	}

	rootDir := filepath.Clean("../..")
	cmd := exec.Command(makeBin, "--no-print-directory", "worktree-prune-plan-json")
	cmd.Dir = rootDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("make worktree-prune-plan-json failed: %v\n%s", err, output)
	}

	var payload []struct {
		Path     string   `json:"path"`
		Branch   string   `json:"branch"`
		Commands []string `json:"commands"`
	}
	if err := json.Unmarshal(output, &payload); err != nil {
		t.Fatalf("unmarshal output: %v\n%s", err, output)
	}
	for _, row := range payload {
		if row.Path == "" || row.Branch == "" {
			t.Fatalf("unexpected row: %#v", row)
		}
		if len(row.Commands) == 0 {
			t.Fatalf("unexpected row: %#v", row)
		}
	}
}

func TestWorktreeCleanupReportTargetRuns(t *testing.T) {
	makeBin, err := exec.LookPath("make")
	if err != nil {
		t.Fatalf("look path make: %v", err)
	}

	rootDir := filepath.Clean("../..")
	cmd := exec.Command(makeBin, "--no-print-directory", "worktree-cleanup-report")
	cmd.Dir = rootDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("make worktree-cleanup-report failed: %v\n%s", err, output)
	}
	text := string(output)
	if !strings.Contains(text, "# Worktree Cleanup Report") {
		t.Fatalf("unexpected output: %s", text)
	}
	if !strings.Contains(text, "## Prune Plan") {
		t.Fatalf("unexpected output: %s", text)
	}
}

func TestCompareLocalListRespectsFilter(t *testing.T) {
	makeBin, err := exec.LookPath("make")
	if err != nil {
		t.Fatalf("look path make: %v", err)
	}

	rootDir := filepath.Clean("../..")
	cmd := exec.Command(makeBin, "--no-print-directory", "compare-local-list", "COMPARE_LOCAL_FILTER=mock-stage-aware,xml-splitrw-subset")
	cmd.Dir = rootDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("make compare-local-list failed: %v\n%s", err, output)
	}

	lines := strings.Fields(strings.TrimSpace(string(output)))
	want := []string{"mock-stage-aware", "xml-splitrw-subset"}
	if len(lines) != len(want) {
		t.Fatalf("lines = %#v", lines)
	}
	for i, name := range want {
		if lines[i] != name {
			t.Fatalf("lines = %#v", lines)
		}
	}
}

func TestCompareLocalListJSONRespectsFilter(t *testing.T) {
	makeBin, err := exec.LookPath("make")
	if err != nil {
		t.Fatalf("look path make: %v", err)
	}

	rootDir := filepath.Clean("../..")
	cmd := exec.Command(makeBin, "--no-print-directory", "compare-local-list-json", "COMPARE_LOCAL_FILTER=mock-stage-aware,xml-splitrw-subset")
	cmd.Dir = rootDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("make compare-local-list-json failed: %v\n%s", err, output)
	}

	var payload []struct {
		Name string `json:"name"`
	}
	if err := json.Unmarshal(output, &payload); err != nil {
		t.Fatalf("unmarshal output: %v\n%s", err, output)
	}
	if len(payload) != 2 || payload[0].Name != "mock-stage-aware" || payload[1].Name != "xml-splitrw-subset" {
		t.Fatalf("payload = %#v", payload)
	}
}

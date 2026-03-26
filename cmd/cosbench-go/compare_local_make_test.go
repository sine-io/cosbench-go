package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func makeCommand(t *testing.T, args ...string) *exec.Cmd {
	t.Helper()
	cmd := exec.Command(mustLookPath(t, "make"), args...)
	cmd.Dir = repoRootDir()
	return cmd
}

func runMakeSuccess(t *testing.T, args ...string) []byte {
	t.Helper()
	return runCommandSuccess(t, makeCommand(t, args...))
}

func runMakeFailure(t *testing.T, args ...string) []byte {
	t.Helper()
	return runCommandFailure(t, makeCommand(t, args...))
}

func TestCompareLocalPrunesStaleOutputs(t *testing.T) {
	goBin := mustLookPath(t, "go")
	outputDir := filepath.Join(t.TempDir(), "compare-local")
	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		t.Fatalf("mkdir output dir: %v", err)
	}
	staleFile := filepath.Join(outputDir, "stale.json")
	if err := os.WriteFile(staleFile, []byte("{}\n"), 0o644); err != nil {
		t.Fatalf("seed stale file: %v", err)
	}

	runMakeSuccess(t, "compare-local", "GO="+goBin, "COMPARE_LOCAL_OUTPUT_DIR="+outputDir)

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

	indexData := mustReadFile(t, filepath.Join(outputDir, "index.json"))
	summaryData := mustReadFile(t, filepath.Join(outputDir, "summary.md"))
	if !strings.Contains(string(summaryData), "## Compare Local") || !strings.Contains(string(summaryData), "| Fixture | Workload |") || !strings.Contains(string(summaryData), "Filter: `all`") {
		t.Fatalf("unexpected summary: %s", summaryData)
	}
	if !strings.Contains(string(summaryData), "Artifact directory: `"+outputDir+"`") {
		t.Fatalf("unexpected summary artifact directory: %s", summaryData)
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
	mustUnmarshalJSON(t, indexData, &payload)
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
	goBin := mustLookPath(t, "go")
	unsafeDir := t.TempDir()
	output := runMakeFailure(t, "compare-local", "GO="+goBin, "COMPARE_LOCAL_OUTPUT_DIR="+unsafeDir)
	if !strings.Contains(string(output), "COMPARE_LOCAL_OUTPUT_DIR") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestCompareLocalFilterRunsSingleFixture(t *testing.T) {
	goBin := mustLookPath(t, "go")
	outputDir := filepath.Join(t.TempDir(), "compare-local")
	runMakeSuccess(
		t,
		"compare-local",
		"GO="+goBin,
		"COMPARE_LOCAL_OUTPUT_DIR="+outputDir,
		"COMPARE_LOCAL_FILTER=mock-stage-aware",
	)

	if _, err := os.Stat(filepath.Join(outputDir, "mock-stage-aware.json")); err != nil {
		t.Fatalf("expected filtered output: %v", err)
	}
	if _, err := os.Stat(filepath.Join(outputDir, "s3-active-subset.json")); err == nil {
		t.Fatal("unexpected unfiltered output")
	}

	indexData := mustReadFile(t, filepath.Join(outputDir, "index.json"))
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
	mustUnmarshalJSON(t, indexData, &payload)
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
	goBin := mustLookPath(t, "go")
	outputDir := filepath.Join(t.TempDir(), "compare-local")
	runMakeSuccess(
		t,
		"compare-local",
		"GO="+goBin,
		"COMPARE_LOCAL_OUTPUT_DIR="+outputDir,
		"COMPARE_LOCAL_FILTER=mock-stage-aware,xml-splitrw-subset",
	)

	if _, err := os.Stat(filepath.Join(outputDir, "mock-stage-aware.json")); err != nil {
		t.Fatalf("expected subset output: %v", err)
	}
	if _, err := os.Stat(filepath.Join(outputDir, "xml-splitrw-subset.json")); err != nil {
		t.Fatalf("expected subset output: %v", err)
	}
	if _, err := os.Stat(filepath.Join(outputDir, "s3-active-subset.json")); err == nil {
		t.Fatal("unexpected unfiltered output")
	}

	indexData := mustReadFile(t, filepath.Join(outputDir, "index.json"))
	var payload struct {
		Meta struct {
			Filter       string `json:"filter"`
			FixtureCount int    `json:"fixture_count"`
		} `json:"meta"`
		Fixtures []struct {
			Name string `json:"name"`
		} `json:"fixtures"`
	}
	mustUnmarshalJSON(t, indexData, &payload)
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

func TestCompareLocalFilterTrimsWhitespaceAroundFixtureNames(t *testing.T) {
	goBin := mustLookPath(t, "go")
	outputDir := filepath.Join(t.TempDir(), "compare-local")
	runMakeSuccess(
		t,
		"compare-local",
		"GO="+goBin,
		"COMPARE_LOCAL_OUTPUT_DIR="+outputDir,
		"COMPARE_LOCAL_FILTER=mock-stage-aware, xml-splitrw-subset",
	)

	if _, err := os.Stat(filepath.Join(outputDir, "mock-stage-aware.json")); err != nil {
		t.Fatalf("expected subset output: %v", err)
	}
	if _, err := os.Stat(filepath.Join(outputDir, "xml-splitrw-subset.json")); err != nil {
		t.Fatalf("expected subset output: %v", err)
	}
	if _, err := os.Stat(filepath.Join(outputDir, "s3-active-subset.json")); err == nil {
		t.Fatal("unexpected unfiltered output")
	}

	indexData := mustReadFile(t, filepath.Join(outputDir, "index.json"))
	var payload struct {
		Meta struct {
			Filter       string `json:"filter"`
			FixtureCount int    `json:"fixture_count"`
		} `json:"meta"`
		Fixtures []struct {
			Name string `json:"name"`
		} `json:"fixtures"`
	}
	mustUnmarshalJSON(t, indexData, &payload)
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
	goBin := mustLookPath(t, "go")
	outputDir := filepath.Join(t.TempDir(), "compare-local")
	runMakeSuccess(
		t,
		"compare-local",
		"GO="+goBin,
		"COMPARE_LOCAL_OUTPUT_DIR="+outputDir,
		"COMPARE_LOCAL_FILTER=all",
	)

	indexData := mustReadFile(t, filepath.Join(outputDir, "index.json"))
	var payload struct {
		Meta struct {
			Filter       string `json:"filter"`
			FixtureCount int    `json:"fixture_count"`
		} `json:"meta"`
	}
	mustUnmarshalJSON(t, indexData, &payload)
	if payload.Meta.Filter != "all" {
		t.Fatalf("meta filter = %q", payload.Meta.Filter)
	}
	if payload.Meta.FixtureCount != 4 {
		t.Fatalf("meta fixture_count = %d", payload.Meta.FixtureCount)
	}
}

func TestCompareLocalFilterDeduplicatesRepeatedFixtureNames(t *testing.T) {
	goBin := mustLookPath(t, "go")
	outputDir := filepath.Join(t.TempDir(), "compare-local")
	runMakeSuccess(
		t,
		"compare-local",
		"GO="+goBin,
		"COMPARE_LOCAL_OUTPUT_DIR="+outputDir,
		"COMPARE_LOCAL_FILTER=mock-stage-aware,mock-stage-aware",
	)

	indexData := mustReadFile(t, filepath.Join(outputDir, "index.json"))
	summaryData := mustReadFile(t, filepath.Join(outputDir, "summary.md"))
	var payload struct {
		Meta struct {
			Filter       string `json:"filter"`
			FixtureCount int    `json:"fixture_count"`
		} `json:"meta"`
		Fixtures []struct {
			Name string `json:"name"`
		} `json:"fixtures"`
	}
	mustUnmarshalJSON(t, indexData, &payload)
	if payload.Meta.Filter != "mock-stage-aware" {
		t.Fatalf("meta filter = %q", payload.Meta.Filter)
	}
	if payload.Meta.FixtureCount != 1 {
		t.Fatalf("meta fixture_count = %d", payload.Meta.FixtureCount)
	}
	if len(payload.Fixtures) != 1 || payload.Fixtures[0].Name != "mock-stage-aware" {
		t.Fatalf("fixtures = %#v", payload.Fixtures)
	}
	if !strings.Contains(string(summaryData), "Filter: `mock-stage-aware`") {
		t.Fatalf("unexpected summary: %s", summaryData)
	}
}

func TestCompareLocalFilterRejectsUnknownFixture(t *testing.T) {
	goBin := mustLookPath(t, "go")
	outputDir := filepath.Join(t.TempDir(), "compare-local")
	output := runMakeFailure(
		t,
		"compare-local",
		"GO="+goBin,
		"COMPARE_LOCAL_OUTPUT_DIR="+outputDir,
		"COMPARE_LOCAL_FILTER=does-not-exist",
	)
	if !strings.Contains(string(output), "does-not-exist") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestCompareLocalListShowsFixtureNames(t *testing.T) {
	output := runMakeSuccess(t, "--no-print-directory", "compare-local-list")

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
	output := runMakeSuccess(t, "--no-print-directory", "compare-local-list-json")

	var payload []struct {
		Name     string `json:"name"`
		Workload string `json:"workload"`
	}
	mustUnmarshalJSON(t, output, &payload)
	if len(payload) != 4 {
		t.Fatalf("payload = %#v", payload)
	}
	if payload[0].Name != "s3-active-subset" || payload[0].Workload != "testdata/workloads/s3-active-subset.xml" {
		t.Fatalf("payload = %#v", payload)
	}
}

func TestWorktreeAuditTargetRuns(t *testing.T) {
	output := runMakeSuccess(t, "--no-print-directory", "worktree-audit")
	text := string(output)
	for _, expected := range []string{
		"# Generated at:",
		"# Base ref:",
		"# Current worktree:",
	} {
		if !strings.Contains(text, expected) {
			t.Fatalf("missing %q in output: %s", expected, text)
		}
	}
	if !strings.Contains(text, "PATH\tBRANCH\tCURRENT\tSTATE\tDETAILS") {
		t.Fatalf("unexpected audit output: %s", text)
	}
	if !strings.Contains(text, "\tyes\t") {
		t.Fatalf("expected current row marker: %s", text)
	}
}

func TestWorktreeAuditTargetSupportsBaseRefOverride(t *testing.T) {
	output := runMakeSuccess(t, "--no-print-directory", "worktree-audit", "WORKTREE_AUDIT_BASE_REF=HEAD")
	text := string(output)
	if !strings.Contains(text, "PATH\tBRANCH\tCURRENT\tSTATE\tDETAILS") {
		t.Fatalf("unexpected output: %s", text)
	}
	if strings.Contains(text, "origin/main") {
		t.Fatalf("unexpected default base ref output: %s", text)
	}
}

func TestWorktreeAuditJSONTargetRuns(t *testing.T) {
	output := runMakeSuccess(t, "--no-print-directory", "worktree-audit-json")

	var payload struct {
		GeneratedAt string `json:"generated_at"`
		Meta        struct {
			GeneratedAt    string `json:"generated_at"`
			BaseRef        string `json:"base_ref"`
			CurrentWorktree string `json:"current_worktree"`
		} `json:"meta"`
		Summary     map[string]any   `json:"summary"`
		Rows        []map[string]any `json:"rows"`
		Views       map[string]struct {
			Summary map[string]any   `json:"summary"`
			Rows    []map[string]any `json:"rows"`
		} `json:"views"`
	}
	mustUnmarshalJSON(t, output, &payload)
	if payload.GeneratedAt == "" {
		t.Fatalf("missing generated_at: %#v", payload)
	}
	if payload.Meta.GeneratedAt == "" || payload.Meta.BaseRef == "" || payload.Meta.CurrentWorktree == "" {
		t.Fatalf("missing meta: %#v", payload.Meta)
	}
	auditView, ok := payload.Views["audit"]
	if !ok {
		t.Fatalf("missing audit view: %#v", payload.Views)
	}
	if auditView.Summary == nil {
		t.Fatalf("missing audit view summary: %#v", auditView)
	}
	if len(auditView.Rows) == 0 {
		t.Fatalf("missing audit view rows: %#v", auditView)
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
	if _, ok := payload.Summary["stale"]; !ok {
		t.Fatalf("missing stale: %#v", payload.Summary)
	}
	if _, ok := payload.Summary["prune_candidates"]; !ok {
		t.Fatalf("missing prune_candidates: %#v", payload.Summary)
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
	output := runMakeSuccess(t, "--no-print-directory", "worktree-audit-merged")
	text := string(output)
	for _, expected := range []string{
		"# Generated at:",
		"# Base ref:",
		"# Current worktree:",
	} {
		if !strings.Contains(text, expected) {
			t.Fatalf("missing %q in output: %s", expected, text)
		}
	}
	if !strings.Contains(text, "PATH\tBRANCH\tCURRENT\tSTATE\tDETAILS") {
		t.Fatalf("unexpected audit output: %s", text)
	}
	if strings.Contains(text, "\tactive\t") {
		t.Fatalf("unexpected active row in merged-only output: %s", text)
	}
}

func TestWorktreeAuditStaleTargetRuns(t *testing.T) {
	output := runMakeSuccess(t, "--no-print-directory", "worktree-audit-stale")
	text := string(output)
	for _, expected := range []string{
		"# Generated at:",
		"# Base ref:",
		"# Current worktree:",
	} {
		if !strings.Contains(text, expected) {
			t.Fatalf("missing %q in output: %s", expected, text)
		}
	}
	if !strings.Contains(text, "PATH\tBRANCH\tCURRENT\tSTATE\tDETAILS") {
		t.Fatalf("unexpected output: %s", text)
	}
	if strings.Contains(text, "\tmerged\t") {
		t.Fatalf("unexpected merged row in stale-only output: %s", text)
	}
}

func TestWorktreeAuditIntegratedTargetRuns(t *testing.T) {
	output := runMakeSuccess(t, "--no-print-directory", "worktree-audit-integrated")
	text := string(output)
	for _, expected := range []string{
		"# Generated at:",
		"# Base ref:",
		"# Current worktree:",
	} {
		if !strings.Contains(text, expected) {
			t.Fatalf("missing %q in output: %s", expected, text)
		}
	}
	if !strings.Contains(text, "PATH\tBRANCH\tCURRENT\tSTATE\tDETAILS") {
		t.Fatalf("unexpected output: %s", text)
	}
	if strings.Contains(text, "\tactive\t") || strings.Contains(text, "\tmerged\t") {
		t.Fatalf("unexpected non-integrated row in integrated-only output: %s", text)
	}
}

func TestWorktreeAuditMergedJSONTargetRuns(t *testing.T) {
	output := runMakeSuccess(t, "--no-print-directory", "worktree-audit-merged-json")

	var payload struct {
		Rows []struct {
			Path   string `json:"path"`
			Branch string `json:"branch"`
			State  string `json:"state"`
		} `json:"rows"`
	}
	mustUnmarshalJSON(t, output, &payload)
	for _, row := range payload.Rows {
		if row.State != "merged" {
			t.Fatalf("unexpected row: %#v", row)
		}
	}
}

func TestWorktreeAuditIntegratedJSONTargetRuns(t *testing.T) {
	output := runMakeSuccess(t, "--no-print-directory", "worktree-audit-integrated-json")

	var payload struct {
		Rows []struct {
			Path   string `json:"path"`
			Branch string `json:"branch"`
			State  string `json:"state"`
		} `json:"rows"`
	}
	mustUnmarshalJSON(t, output, &payload)
	for _, row := range payload.Rows {
		if row.State != "integrated" {
			t.Fatalf("unexpected row: %#v", row)
		}
	}
}

func TestWorktreeAuditPruneTargetRuns(t *testing.T) {
	output := runMakeSuccess(t, "--no-print-directory", "worktree-audit-prune")
	text := string(output)
	for _, expected := range []string{
		"# Generated at:",
		"# Base ref:",
		"# Current worktree:",
	} {
		if !strings.Contains(text, expected) {
			t.Fatalf("missing %q in output: %s", expected, text)
		}
	}
	if !strings.Contains(text, "PATH\tBRANCH\tCURRENT\tSTATE\tDETAILS") {
		t.Fatalf("unexpected output: %s", text)
	}
	if strings.Contains(text, "\tactive\t") {
		t.Fatalf("unexpected active row in prune-only output: %s", text)
	}
	if strings.Contains(text, "\tyes\t") {
		t.Fatalf("unexpected current row in prune-only output: %s", text)
	}
}

func TestWorktreeAuditPruneJSONTargetRuns(t *testing.T) {
	output := runMakeSuccess(t, "--no-print-directory", "worktree-audit-prune-json")

	var payload struct {
		Rows []struct {
			Path    string `json:"path"`
			Branch  string `json:"branch"`
			State   string `json:"state"`
			Current bool   `json:"current"`
		} `json:"rows"`
	}
	mustUnmarshalJSON(t, output, &payload)
	for _, row := range payload.Rows {
		if row.State != "merged" && row.State != "integrated" {
			t.Fatalf("unexpected row: %#v", row)
		}
		if row.Current {
			t.Fatalf("unexpected current row: %#v", row)
		}
	}
}

func TestWorktreePrunePlanTargetRuns(t *testing.T) {
	repoRoot := repoRootAbs(t)
	output := runMakeSuccess(t, "--no-print-directory", "worktree-prune-plan")
	text := string(output)
	if !strings.Contains(text, "# Suggested cleanup commands") {
		t.Fatalf("unexpected output: %s", text)
	}
	for _, expected := range []string{
		"# Generated at:",
		"# Base ref:",
		"# Current worktree:",
	} {
		if !strings.Contains(text, expected) {
			t.Fatalf("missing %q in output: %s", expected, text)
		}
	}
	if !strings.Contains(text, "git worktree remove") && !strings.Contains(text, "# no prune-candidate worktrees to prune") {
		t.Fatalf("unexpected output: %s", text)
	}
	if strings.Contains(text, "git worktree remove '"+repoRoot+"'") {
		t.Fatalf("unexpected self-removal plan: %s", text)
	}
}

func TestWorktreePrunePlanJSONTargetRuns(t *testing.T) {
	repoRoot := repoRootAbs(t)
	output := runMakeSuccess(t, "--no-print-directory", "worktree-prune-plan-json")

	var payload struct {
		GeneratedAt string `json:"generated_at"`
		Meta        struct {
			GeneratedAt     string `json:"generated_at"`
			BaseRef         string `json:"base_ref"`
			CurrentWorktree string `json:"current_worktree"`
		} `json:"meta"`
		Summary struct {
			BaseRef         string `json:"base_ref"`
			CurrentWorktree string `json:"current_worktree"`
			Total           int    `json:"total"`
			Merged          int    `json:"merged"`
			Integrated      int    `json:"integrated"`
		} `json:"summary"`
		Rows []struct {
			Path     string   `json:"path"`
			Branch   string   `json:"branch"`
			Commands []string `json:"commands"`
		} `json:"rows"`
		Views map[string]struct {
			Summary struct {
				BaseRef         string `json:"base_ref"`
				CurrentWorktree string `json:"current_worktree"`
				Total           int    `json:"total"`
				Merged          int    `json:"merged"`
				Integrated      int    `json:"integrated"`
			} `json:"summary"`
			Rows []struct {
				Path     string   `json:"path"`
				Branch   string   `json:"branch"`
				Commands []string `json:"commands"`
			} `json:"rows"`
		} `json:"views"`
	}
	mustUnmarshalJSON(t, output, &payload)
	if payload.GeneratedAt == "" {
		t.Fatalf("missing generated_at: %#v", payload)
	}
	if payload.Meta.GeneratedAt == "" || payload.Meta.BaseRef == "" || payload.Meta.CurrentWorktree == "" {
		t.Fatalf("missing meta: %#v", payload.Meta)
	}
	pruneView, ok := payload.Views["prune_plan"]
	if !ok {
		t.Fatalf("missing prune_plan view: %#v", payload.Views)
	}
	if pruneView.Summary.CurrentWorktree == "" {
		t.Fatalf("missing prune_plan view summary: %#v", pruneView)
	}
	if pruneView.Rows == nil {
		t.Fatalf("missing prune_plan view rows: %#v", pruneView)
	}
	if payload.Summary.Total < 0 {
		t.Fatalf("unexpected summary: %#v", payload.Summary)
	}
	if payload.Summary.CurrentWorktree == "" {
		t.Fatalf("unexpected summary: %#v", payload.Summary)
	}
	for _, row := range payload.Rows {
		if row.Path == "" || row.Branch == "" {
			t.Fatalf("unexpected row: %#v", row)
		}
		if len(row.Commands) == 0 {
			t.Fatalf("unexpected row: %#v", row)
		}
		if row.Path == repoRoot {
			t.Fatalf("unexpected self-removal row: %#v", row)
		}
	}
}

func TestWorktreeCleanupReportTargetRuns(t *testing.T) {
	reportPath := filepath.Join(t.TempDir(), "worktree-cleanup-report.md")
	output := runMakeSuccess(t, "--no-print-directory", "worktree-cleanup-report", "WORKTREE_CLEANUP_REPORT_PATH="+reportPath)
	text := string(output)
	if !strings.Contains(text, "# Worktree Cleanup Report") {
		t.Fatalf("unexpected output: %s", text)
	}
	for _, expected := range []string{
		"- Generated at:",
		"- Current worktree:",
		"- Integrated:",
		"- Stale:",
		"- Prune candidates:",
	} {
		if !strings.Contains(text, expected) {
			t.Fatalf("missing %q in output: %s", expected, text)
		}
	}
	if !strings.Contains(text, "## Prune Plan") {
		t.Fatalf("unexpected output: %s", text)
	}
	if !strings.Contains(text, "## Integrated") {
		t.Fatalf("unexpected output: %s", text)
	}
	if !strings.Contains(text, "## Prune Candidates") {
		t.Fatalf("unexpected output: %s", text)
	}
	reportData := mustReadFile(t, reportPath)
	if !strings.Contains(string(reportData), "# Worktree Cleanup Report") {
		t.Fatalf("unexpected report file: %s", reportData)
	}
	if !strings.Contains(string(reportData), "## Integrated") {
		t.Fatalf("unexpected report file: %s", reportData)
	}
	if !strings.Contains(string(reportData), "## Prune Candidates") {
		t.Fatalf("unexpected report file: %s", reportData)
	}
	for _, expected := range []string{
		"- Generated at:",
		"- Current worktree:",
		"- Integrated:",
		"- Stale:",
		"- Prune candidates:",
	} {
		if !strings.Contains(string(reportData), expected) {
			t.Fatalf("missing %q in report file: %s", expected, reportData)
		}
	}
}

func TestWorktreeCleanupReportJSONTargetRuns(t *testing.T) {
	output := runMakeSuccess(t, "--no-print-directory", "worktree-cleanup-report-json")

	var payload map[string]any
	mustUnmarshalJSON(t, output, &payload)
	if generatedAt, ok := payload["generated_at"].(string); !ok || generatedAt == "" {
		t.Fatalf("missing generated_at: %#v", payload)
	}
	meta, ok := payload["meta"].(map[string]any)
	if !ok {
		t.Fatalf("missing meta: %#v", payload)
	}
	for _, key := range []string{"generated_at", "base_ref", "current_worktree"} {
		if value, ok := meta[key].(string); !ok || value == "" {
			t.Fatalf("missing meta[%s]: %#v", key, meta)
		}
	}
	for _, key := range []string{"summary", "views", "merged", "integrated", "stale", "prune_candidates", "prune_plan"} {
		if _, ok := payload[key]; !ok {
			t.Fatalf("missing %s: %#v", key, payload)
		}
	}
	views, ok := payload["views"].(map[string]any)
	if !ok {
		t.Fatalf("views is not structured: %#v", payload["views"])
	}
	for _, key := range []string{"merged", "integrated", "stale", "prune_candidates", "prune_plan"} {
		if _, ok := views[key].(map[string]any); !ok {
			t.Fatalf("views[%s] is not structured: %#v", key, views[key])
		}
	}
	if _, ok := payload["merged"].(map[string]any); !ok {
		t.Fatalf("merged is not structured: %#v", payload["merged"])
	}
	if _, ok := payload["integrated"].(map[string]any); !ok {
		t.Fatalf("integrated is not structured: %#v", payload["integrated"])
	}
	if _, ok := payload["stale"].(map[string]any); !ok {
		t.Fatalf("stale is not structured: %#v", payload["stale"])
	}
	if _, ok := payload["prune_candidates"].(map[string]any); !ok {
		t.Fatalf("prune_candidates is not structured: %#v", payload["prune_candidates"])
	}
	prunePlan, ok := payload["prune_plan"].(map[string]any)
	if !ok {
		t.Fatalf("prune_plan is not structured: %#v", payload["prune_plan"])
	}
	if _, ok := prunePlan["summary"].(map[string]any); !ok {
		t.Fatalf("prune_plan summary is not structured: %#v", prunePlan)
	}
	if _, ok := prunePlan["rows"].([]any); !ok {
		t.Fatalf("prune_plan rows are not structured: %#v", prunePlan)
	}
}

func TestWorktreeCleanupReportRespectsBaseRef(t *testing.T) {
	output := runMakeSuccess(t, "--no-print-directory", "worktree-cleanup-report", "WORKTREE_AUDIT_BASE_REF=HEAD")
	text := string(output)
	if !strings.Contains(text, "- Base ref: `HEAD`") {
		t.Fatalf("unexpected output: %s", text)
	}
}

func TestCompareLocalListRespectsFilter(t *testing.T) {
	output := runMakeSuccess(t, "--no-print-directory", "compare-local-list", "COMPARE_LOCAL_FILTER=mock-stage-aware,xml-splitrw-subset")

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

func TestCompareLocalListPreservesFilterOrder(t *testing.T) {
	output := runMakeSuccess(t, "--no-print-directory", "compare-local-list", "COMPARE_LOCAL_FILTER=xml-splitrw-subset,mock-stage-aware")

	lines := strings.Fields(strings.TrimSpace(string(output)))
	want := []string{"xml-splitrw-subset", "mock-stage-aware"}
	if len(lines) != len(want) {
		t.Fatalf("lines = %#v", lines)
	}
	for i, name := range want {
		if lines[i] != name {
			t.Fatalf("lines = %#v", lines)
		}
	}
}

func TestCompareLocalListTrimsFilterWhitespace(t *testing.T) {
	output := runMakeSuccess(t, "--no-print-directory", "compare-local-list", "COMPARE_LOCAL_FILTER=mock-stage-aware, xml-splitrw-subset")

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
	output := runMakeSuccess(t, "--no-print-directory", "compare-local-list-json", "COMPARE_LOCAL_FILTER=mock-stage-aware,xml-splitrw-subset")

	var payload []struct {
		Name string `json:"name"`
	}
	mustUnmarshalJSON(t, output, &payload)
	if len(payload) != 2 || payload[0].Name != "mock-stage-aware" || payload[1].Name != "xml-splitrw-subset" {
		t.Fatalf("payload = %#v", payload)
	}
}

func TestCompareLocalFilterPreservesFixtureOrder(t *testing.T) {
	goBin := mustLookPath(t, "go")
	outputDir := filepath.Join(t.TempDir(), "compare-local")
	runMakeSuccess(
		t,
		"compare-local",
		"GO="+goBin,
		"COMPARE_LOCAL_OUTPUT_DIR="+outputDir,
		"COMPARE_LOCAL_FILTER=xml-splitrw-subset,mock-stage-aware",
	)

	indexData := mustReadFile(t, filepath.Join(outputDir, "index.json"))
	var payload struct {
		Meta struct {
			Filter       string `json:"filter"`
			FixtureCount int    `json:"fixture_count"`
		} `json:"meta"`
		Fixtures []struct {
			Name string `json:"name"`
		} `json:"fixtures"`
	}
	mustUnmarshalJSON(t, indexData, &payload)
	if payload.Meta.Filter != "xml-splitrw-subset,mock-stage-aware" {
		t.Fatalf("meta filter = %q", payload.Meta.Filter)
	}
	if payload.Meta.FixtureCount != 2 {
		t.Fatalf("meta fixture_count = %d", payload.Meta.FixtureCount)
	}
	if len(payload.Fixtures) != 2 {
		t.Fatalf("fixtures = %#v", payload.Fixtures)
	}
	if payload.Fixtures[0].Name != "xml-splitrw-subset" || payload.Fixtures[1].Name != "mock-stage-aware" {
		t.Fatalf("fixtures = %#v", payload.Fixtures)
	}
}

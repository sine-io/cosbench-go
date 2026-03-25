package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestWorktreeAuditJSONMarksPatchEquivalentBranchIntegrated(t *testing.T) {
	gitBin, err := exec.LookPath("git")
	if err != nil {
		t.Fatalf("look path git: %v", err)
	}
	pythonBin, err := exec.LookPath("python3")
	if err != nil {
		t.Fatalf("look path python3: %v", err)
	}

	repoDir := t.TempDir()
	runCmd(t, repoDir, gitBin, "init", "-b", "main")
	runCmd(t, repoDir, gitBin, "config", "user.name", "Test User")
	runCmd(t, repoDir, gitBin, "config", "user.email", "test@example.com")

	filePath := filepath.Join(repoDir, "note.txt")
	if err := os.WriteFile(filePath, []byte("base\n"), 0o644); err != nil {
		t.Fatalf("write base file: %v", err)
	}
	runCmd(t, repoDir, gitBin, "add", "note.txt")
	runCmd(t, repoDir, gitBin, "commit", "-m", "base")

	featureDir := filepath.Join(t.TempDir(), "feature")
	runCmd(t, repoDir, gitBin, "worktree", "add", featureDir, "-b", "feature")

	appendLine(t, filepath.Join(featureDir, "note.txt"), "feature\n")
	runCmd(t, featureDir, gitBin, "add", "note.txt")
	runCmd(t, featureDir, gitBin, "commit", "-m", "feature change")

	appendLine(t, filepath.Join(repoDir, "note.txt"), "feature\n")
	runCmd(t, repoDir, gitBin, "add", "note.txt")
	runCmd(t, repoDir, gitBin, "commit", "-m", "squash-equivalent")

	scriptPath, err := filepath.Abs(filepath.Clean("../../scripts/worktree_audit.py"))
	if err != nil {
		t.Fatalf("abs script path: %v", err)
	}
	cmd := exec.Command(pythonBin, scriptPath, "--json", "main")
	cmd.Dir = repoDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("run script: %v\n%s", err, output)
	}

	var payload struct {
		Summary map[string]any `json:"summary"`
		Rows    []struct {
			Branch string `json:"branch"`
			State  string `json:"state"`
		} `json:"rows"`
	}
	if err := json.Unmarshal(output, &payload); err != nil {
		t.Fatalf("unmarshal output: %v\n%s", err, output)
	}

	found := false
	for _, row := range payload.Rows {
		if row.Branch == "feature" {
			found = true
			if row.State != "integrated" {
				t.Fatalf("feature state = %q", row.State)
			}
		}
	}
	if !found {
		t.Fatalf("missing feature row: %s", output)
	}
	if payload.Summary["integrated"] == nil {
		t.Fatalf("missing integrated count: %#v", payload.Summary)
	}
	if payload.Summary["integrated"].(float64) != 1 {
		t.Fatalf("integrated count = %#v", payload.Summary["integrated"])
	}
	if payload.Summary["prune_candidates"] == nil {
		t.Fatalf("missing prune_candidates count: %#v", payload.Summary)
	}
	if payload.Summary["prune_candidates"].(float64) != 1 {
		t.Fatalf("prune_candidates = %#v", payload.Summary["prune_candidates"])
	}
}

func TestWorktreePrunePlanJSONIncludesBranchContext(t *testing.T) {
	gitBin, err := exec.LookPath("git")
	if err != nil {
		t.Fatalf("look path git: %v", err)
	}
	pythonBin, err := exec.LookPath("python3")
	if err != nil {
		t.Fatalf("look path python3: %v", err)
	}

	repoDir := t.TempDir()
	runCmd(t, repoDir, gitBin, "init", "-b", "main")
	runCmd(t, repoDir, gitBin, "config", "user.name", "Test User")
	runCmd(t, repoDir, gitBin, "config", "user.email", "test@example.com")

	filePath := filepath.Join(repoDir, "note.txt")
	if err := os.WriteFile(filePath, []byte("base\n"), 0o644); err != nil {
		t.Fatalf("write base file: %v", err)
	}
	runCmd(t, repoDir, gitBin, "add", "note.txt")
	runCmd(t, repoDir, gitBin, "commit", "-m", "base")

	featureDir := filepath.Join(t.TempDir(), "feature")
	runCmd(t, repoDir, gitBin, "worktree", "add", featureDir, "-b", "feature")

	appendLine(t, filepath.Join(featureDir, "note.txt"), "feature\n")
	runCmd(t, featureDir, gitBin, "add", "note.txt")
	runCmd(t, featureDir, gitBin, "commit", "-m", "feature change")

	appendLine(t, filepath.Join(repoDir, "note.txt"), "feature\n")
	runCmd(t, repoDir, gitBin, "add", "note.txt")
	runCmd(t, repoDir, gitBin, "commit", "-m", "squash-equivalent")

	scriptPath, err := filepath.Abs(filepath.Clean("../../scripts/worktree_prune_plan.py"))
	if err != nil {
		t.Fatalf("abs script path: %v", err)
	}
	cmd := exec.Command(pythonBin, scriptPath, "--json", "main")
	cmd.Dir = repoDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("run script: %v\n%s", err, output)
	}

	var payload struct {
		Summary struct {
			BaseRef         string `json:"base_ref"`
			CurrentWorktree string `json:"current_worktree"`
			Total           int    `json:"total"`
			Merged          int    `json:"merged"`
			Integrated      int    `json:"integrated"`
		} `json:"summary"`
		Rows []struct {
			Branch   string   `json:"branch"`
			State    string   `json:"state"`
			Details  string   `json:"details"`
			Ahead    int      `json:"ahead"`
			Behind   int      `json:"behind"`
			Commands []string `json:"commands"`
		} `json:"rows"`
	}
	if err := json.Unmarshal(output, &payload); err != nil {
		t.Fatalf("unmarshal output: %v\n%s", err, output)
	}
	if payload.Summary.BaseRef != "main" {
		t.Fatalf("summary = %#v", payload.Summary)
	}
	if payload.Summary.Total != 1 || payload.Summary.Integrated != 1 || payload.Summary.Merged != 0 {
		t.Fatalf("summary = %#v", payload.Summary)
	}
	if len(payload.Rows) != 1 {
		t.Fatalf("payload = %#v", payload)
	}
	if payload.Rows[0].Branch != "feature" || payload.Rows[0].State != "integrated" {
		t.Fatalf("row = %#v", payload.Rows[0])
	}
	if payload.Rows[0].Details == "" {
		t.Fatalf("row = %#v", payload.Rows[0])
	}
	if payload.Rows[0].Ahead < 0 || payload.Rows[0].Behind < 0 {
		t.Fatalf("row = %#v", payload.Rows[0])
	}
	if len(payload.Rows[0].Commands) != 2 {
		t.Fatalf("row = %#v", payload.Rows[0])
	}
}

func runCmd(t *testing.T, dir, bin string, args ...string) {
	t.Helper()
	cmd := exec.Command(bin, args...)
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("%s %v failed: %v\n%s", bin, args, err, output)
	}
}

func appendLine(t *testing.T, path, line string) {
	t.Helper()
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		t.Fatalf("open %s: %v", path, err)
	}
	defer f.Close()
	if _, err := f.WriteString(strings.ReplaceAll(line, "\r\n", "\n")); err != nil {
		t.Fatalf("append %s: %v", path, err)
	}
}

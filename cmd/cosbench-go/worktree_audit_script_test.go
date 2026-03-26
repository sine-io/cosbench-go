package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func setupPatchEquivalentRepo(t *testing.T) (repoDir string, gitBin string, pythonBin string) {
	t.Helper()

	gitBin, err := exec.LookPath("git")
	if err != nil {
		t.Fatalf("look path git: %v", err)
	}
	pythonBin, err = exec.LookPath("python3")
	if err != nil {
		t.Fatalf("look path python3: %v", err)
	}

	repoDir = t.TempDir()
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

	return repoDir, gitBin, pythonBin
}

func setupActiveFeatureRepo(t *testing.T) (repoDir string, gitBin string, pythonBin string) {
	t.Helper()

	gitBin, err := exec.LookPath("git")
	if err != nil {
		t.Fatalf("look path git: %v", err)
	}
	pythonBin, err = exec.LookPath("python3")
	if err != nil {
		t.Fatalf("look path python3: %v", err)
	}

	repoDir = t.TempDir()
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

	return repoDir, gitBin, pythonBin
}

func setupUnicodePatchEquivalentRepo(t *testing.T) (repoDir string, gitBin string, pythonBin string) {
	t.Helper()

	gitBin, err := exec.LookPath("git")
	if err != nil {
		t.Fatalf("look path git: %v", err)
	}
	pythonBin, err = exec.LookPath("python3")
	if err != nil {
		t.Fatalf("look path python3: %v", err)
	}

	repoDir = t.TempDir()
	runCmd(t, repoDir, gitBin, "init", "-b", "main")
	runCmd(t, repoDir, gitBin, "config", "user.name", "Test User")
	runCmd(t, repoDir, gitBin, "config", "user.email", "test@example.com")

	filePath := filepath.Join(repoDir, "note.txt")
	if err := os.WriteFile(filePath, []byte("base\n"), 0o644); err != nil {
		t.Fatalf("write base file: %v", err)
	}
	runCmd(t, repoDir, gitBin, "add", "note.txt")
	runCmd(t, repoDir, gitBin, "commit", "-m", "base")

	featureDir := filepath.Join(t.TempDir(), "工作树")
	runCmd(t, repoDir, gitBin, "worktree", "add", featureDir, "-b", "特性")

	appendLine(t, filepath.Join(featureDir, "note.txt"), "feature\n")
	runCmd(t, featureDir, gitBin, "add", "note.txt")
	runCmd(t, featureDir, gitBin, "commit", "-m", "feature change")

	appendLine(t, filepath.Join(repoDir, "note.txt"), "feature\n")
	runCmd(t, repoDir, gitBin, "add", "note.txt")
	runCmd(t, repoDir, gitBin, "commit", "-m", "squash-equivalent")

	return repoDir, gitBin, pythonBin
}

func setupActiveTrunkRepo(t *testing.T) (repoDir string, gitBin string, pythonBin string) {
	t.Helper()

	gitBin, err := exec.LookPath("git")
	if err != nil {
		t.Fatalf("look path git: %v", err)
	}
	pythonBin, err = exec.LookPath("python3")
	if err != nil {
		t.Fatalf("look path python3: %v", err)
	}

	repoDir = t.TempDir()
	runCmd(t, repoDir, gitBin, "init", "-b", "trunk")
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

	return repoDir, gitBin, pythonBin
}

func setupActiveTrunkRepoWithFeatureWorktree(t *testing.T) (repoDir, featureDir, gitBin, pythonBin string) {
	t.Helper()

	gitBin, err := exec.LookPath("git")
	if err != nil {
		t.Fatalf("look path git: %v", err)
	}
	pythonBin, err = exec.LookPath("python3")
	if err != nil {
		t.Fatalf("look path python3: %v", err)
	}

	repoDir = t.TempDir()
	runCmd(t, repoDir, gitBin, "init", "-b", "trunk")
	runCmd(t, repoDir, gitBin, "config", "user.name", "Test User")
	runCmd(t, repoDir, gitBin, "config", "user.email", "test@example.com")

	filePath := filepath.Join(repoDir, "note.txt")
	if err := os.WriteFile(filePath, []byte("base\n"), 0o644); err != nil {
		t.Fatalf("write base file: %v", err)
	}
	runCmd(t, repoDir, gitBin, "add", "note.txt")
	runCmd(t, repoDir, gitBin, "commit", "-m", "base")

	featureDir = filepath.Join(t.TempDir(), "feature")
	runCmd(t, repoDir, gitBin, "worktree", "add", featureDir, "-b", "feature")

	appendLine(t, filepath.Join(featureDir, "note.txt"), "feature\n")
	runCmd(t, featureDir, gitBin, "add", "note.txt")
	runCmd(t, featureDir, gitBin, "commit", "-m", "feature change")

	return repoDir, featureDir, gitBin, pythonBin
}

func setupDetachedTrunkRepo(t *testing.T) (repoDir string, gitBin string, pythonBin string) {
	t.Helper()

	repoDir, gitBin, pythonBin = setupActiveTrunkRepo(t)
	runCmd(t, repoDir, gitBin, "checkout", "--detach", "HEAD")
	runCmd(t, repoDir, gitBin, "branch", "-D", "trunk")
	return repoDir, gitBin, pythonBin
}

func setupRemoteTrunkOnlyFeatureRepo(t *testing.T) (repoDir, featureDir, gitBin, pythonBin string) {
	t.Helper()

	repoDir, featureDir, gitBin, pythonBin = setupActiveTrunkRepoWithFeatureWorktree(t)
	runCmd(t, repoDir, gitBin, "update-ref", "refs/remotes/origin/trunk", "HEAD")
	runCmd(t, repoDir, gitBin, "checkout", "--detach", "HEAD")
	runCmd(t, repoDir, gitBin, "branch", "-D", "trunk")
	return repoDir, featureDir, gitBin, pythonBin
}

func setupRemoteMasterOnlyFeatureRepo(t *testing.T) (repoDir, featureDir, gitBin, pythonBin string) {
	t.Helper()

	gitBin, err := exec.LookPath("git")
	if err != nil {
		t.Fatalf("look path git: %v", err)
	}
	pythonBin, err = exec.LookPath("python3")
	if err != nil {
		t.Fatalf("look path python3: %v", err)
	}

	repoDir = t.TempDir()
	runCmd(t, repoDir, gitBin, "init", "-b", "master")
	runCmd(t, repoDir, gitBin, "config", "user.name", "Test User")
	runCmd(t, repoDir, gitBin, "config", "user.email", "test@example.com")

	filePath := filepath.Join(repoDir, "note.txt")
	if err := os.WriteFile(filePath, []byte("base\n"), 0o644); err != nil {
		t.Fatalf("write base file: %v", err)
	}
	runCmd(t, repoDir, gitBin, "add", "note.txt")
	runCmd(t, repoDir, gitBin, "commit", "-m", "base")

	runCmd(t, repoDir, gitBin, "update-ref", "refs/remotes/origin/master", "HEAD")

	featureDir = filepath.Join(t.TempDir(), "feature")
	runCmd(t, repoDir, gitBin, "worktree", "add", featureDir, "-b", "feature")

	appendLine(t, filepath.Join(featureDir, "note.txt"), "feature\n")
	runCmd(t, featureDir, gitBin, "add", "note.txt")
	runCmd(t, featureDir, gitBin, "commit", "-m", "feature change")

	runCmd(t, repoDir, gitBin, "checkout", "--detach", "HEAD")
	runCmd(t, repoDir, gitBin, "branch", "-D", "master")
	return repoDir, featureDir, gitBin, pythonBin
}

func setupRemoteHEADOnlyFeatureRepo(t *testing.T) (repoDir, featureDir, gitBin, pythonBin string) {
	t.Helper()

	gitBin, err := exec.LookPath("git")
	if err != nil {
		t.Fatalf("look path git: %v", err)
	}
	pythonBin, err = exec.LookPath("python3")
	if err != nil {
		t.Fatalf("look path python3: %v", err)
	}

	repoDir = t.TempDir()
	runCmd(t, repoDir, gitBin, "init", "-b", "develop")
	runCmd(t, repoDir, gitBin, "config", "user.name", "Test User")
	runCmd(t, repoDir, gitBin, "config", "user.email", "test@example.com")

	filePath := filepath.Join(repoDir, "note.txt")
	if err := os.WriteFile(filePath, []byte("base\n"), 0o644); err != nil {
		t.Fatalf("write base file: %v", err)
	}
	runCmd(t, repoDir, gitBin, "add", "note.txt")
	runCmd(t, repoDir, gitBin, "commit", "-m", "base")

	runCmd(t, repoDir, gitBin, "update-ref", "refs/remotes/origin/develop", "HEAD")
	runCmd(t, repoDir, gitBin, "symbolic-ref", "refs/remotes/origin/HEAD", "refs/remotes/origin/develop")

	featureDir = filepath.Join(t.TempDir(), "feature")
	runCmd(t, repoDir, gitBin, "worktree", "add", featureDir, "-b", "feature")

	appendLine(t, filepath.Join(featureDir, "note.txt"), "feature\n")
	runCmd(t, featureDir, gitBin, "add", "note.txt")
	runCmd(t, featureDir, gitBin, "commit", "-m", "feature change")

	runCmd(t, repoDir, gitBin, "checkout", "--detach", "HEAD")
	runCmd(t, repoDir, gitBin, "branch", "-D", "develop")
	return repoDir, featureDir, gitBin, pythonBin
}

func setupNonOriginRemoteHEADOnlyFeatureRepo(t *testing.T) (repoDir, featureDir, gitBin, pythonBin string) {
	t.Helper()

	gitBin, err := exec.LookPath("git")
	if err != nil {
		t.Fatalf("look path git: %v", err)
	}
	pythonBin, err = exec.LookPath("python3")
	if err != nil {
		t.Fatalf("look path python3: %v", err)
	}

	repoDir = t.TempDir()
	runCmd(t, repoDir, gitBin, "init", "-b", "develop")
	runCmd(t, repoDir, gitBin, "config", "user.name", "Test User")
	runCmd(t, repoDir, gitBin, "config", "user.email", "test@example.com")

	filePath := filepath.Join(repoDir, "note.txt")
	if err := os.WriteFile(filePath, []byte("base\n"), 0o644); err != nil {
		t.Fatalf("write base file: %v", err)
	}
	runCmd(t, repoDir, gitBin, "add", "note.txt")
	runCmd(t, repoDir, gitBin, "commit", "-m", "base")

	runCmd(t, repoDir, gitBin, "update-ref", "refs/remotes/upstream/develop", "HEAD")
	runCmd(t, repoDir, gitBin, "symbolic-ref", "refs/remotes/upstream/HEAD", "refs/remotes/upstream/develop")

	featureDir = filepath.Join(t.TempDir(), "feature")
	runCmd(t, repoDir, gitBin, "worktree", "add", featureDir, "-b", "feature")

	appendLine(t, filepath.Join(featureDir, "note.txt"), "feature\n")
	runCmd(t, featureDir, gitBin, "add", "note.txt")
	runCmd(t, featureDir, gitBin, "commit", "-m", "feature change")

	runCmd(t, repoDir, gitBin, "checkout", "--detach", "HEAD")
	runCmd(t, repoDir, gitBin, "branch", "-D", "develop")
	return repoDir, featureDir, gitBin, pythonBin
}

func setupNonOriginRemoteNamedBranchOnlyFeatureRepo(t *testing.T, branch string) (repoDir, featureDir, gitBin, pythonBin string) {
	t.Helper()

	gitBin, err := exec.LookPath("git")
	if err != nil {
		t.Fatalf("look path git: %v", err)
	}
	pythonBin, err = exec.LookPath("python3")
	if err != nil {
		t.Fatalf("look path python3: %v", err)
	}

	repoDir = t.TempDir()
	runCmd(t, repoDir, gitBin, "init", "-b", branch)
	runCmd(t, repoDir, gitBin, "config", "user.name", "Test User")
	runCmd(t, repoDir, gitBin, "config", "user.email", "test@example.com")

	filePath := filepath.Join(repoDir, "note.txt")
	if err := os.WriteFile(filePath, []byte("base\n"), 0o644); err != nil {
		t.Fatalf("write base file: %v", err)
	}
	runCmd(t, repoDir, gitBin, "add", "note.txt")
	runCmd(t, repoDir, gitBin, "commit", "-m", "base")

	runCmd(t, repoDir, gitBin, "update-ref", "refs/remotes/upstream/"+branch, "HEAD")

	featureDir = filepath.Join(t.TempDir(), "feature")
	runCmd(t, repoDir, gitBin, "worktree", "add", featureDir, "-b", "feature")

	appendLine(t, filepath.Join(featureDir, "note.txt"), "feature\n")
	runCmd(t, featureDir, gitBin, "add", "note.txt")
	runCmd(t, featureDir, gitBin, "commit", "-m", "feature change")

	runCmd(t, repoDir, gitBin, "checkout", "--detach", "HEAD")
	runCmd(t, repoDir, gitBin, "branch", "-D", branch)
	return repoDir, featureDir, gitBin, pythonBin
}

func setupUpstreamAndForkRemoteHEADFeatureRepo(t *testing.T) (repoDir, featureDir, gitBin, pythonBin string) {
	t.Helper()

	gitBin, err := exec.LookPath("git")
	if err != nil {
		t.Fatalf("look path git: %v", err)
	}
	pythonBin, err = exec.LookPath("python3")
	if err != nil {
		t.Fatalf("look path python3: %v", err)
	}

	repoDir = t.TempDir()
	runCmd(t, repoDir, gitBin, "init", "-b", "develop")
	runCmd(t, repoDir, gitBin, "config", "user.name", "Test User")
	runCmd(t, repoDir, gitBin, "config", "user.email", "test@example.com")

	filePath := filepath.Join(repoDir, "note.txt")
	if err := os.WriteFile(filePath, []byte("base\n"), 0o644); err != nil {
		t.Fatalf("write base file: %v", err)
	}
	runCmd(t, repoDir, gitBin, "add", "note.txt")
	runCmd(t, repoDir, gitBin, "commit", "-m", "base")

	runCmd(t, repoDir, gitBin, "update-ref", "refs/remotes/upstream/develop", "HEAD")
	runCmd(t, repoDir, gitBin, "symbolic-ref", "refs/remotes/upstream/HEAD", "refs/remotes/upstream/develop")
	runCmd(t, repoDir, gitBin, "update-ref", "refs/remotes/fork/develop", "HEAD")
	runCmd(t, repoDir, gitBin, "symbolic-ref", "refs/remotes/fork/HEAD", "refs/remotes/fork/develop")

	featureDir = filepath.Join(t.TempDir(), "feature")
	runCmd(t, repoDir, gitBin, "worktree", "add", featureDir, "-b", "feature")

	appendLine(t, filepath.Join(featureDir, "note.txt"), "feature\n")
	runCmd(t, featureDir, gitBin, "add", "note.txt")
	runCmd(t, featureDir, gitBin, "commit", "-m", "feature change")

	runCmd(t, repoDir, gitBin, "checkout", "--detach", "HEAD")
	runCmd(t, repoDir, gitBin, "branch", "-D", "develop")
	return repoDir, featureDir, gitBin, pythonBin
}

func setupUpstreamAndForkRemoteMainFeatureRepo(t *testing.T) (repoDir, featureDir, gitBin, pythonBin string) {
	t.Helper()

	gitBin, err := exec.LookPath("git")
	if err != nil {
		t.Fatalf("look path git: %v", err)
	}
	pythonBin, err = exec.LookPath("python3")
	if err != nil {
		t.Fatalf("look path python3: %v", err)
	}

	repoDir = t.TempDir()
	runCmd(t, repoDir, gitBin, "init", "-b", "main")
	runCmd(t, repoDir, gitBin, "config", "user.name", "Test User")
	runCmd(t, repoDir, gitBin, "config", "user.email", "test@example.com")

	filePath := filepath.Join(repoDir, "note.txt")
	if err := os.WriteFile(filePath, []byte("base\n"), 0o644); err != nil {
		t.Fatalf("write base file: %v", err)
	}
	runCmd(t, repoDir, gitBin, "add", "note.txt")
	runCmd(t, repoDir, gitBin, "commit", "-m", "base")

	runCmd(t, repoDir, gitBin, "update-ref", "refs/remotes/upstream/main", "HEAD")
	runCmd(t, repoDir, gitBin, "update-ref", "refs/remotes/fork/main", "HEAD")

	featureDir = filepath.Join(t.TempDir(), "feature")
	runCmd(t, repoDir, gitBin, "worktree", "add", featureDir, "-b", "feature")

	appendLine(t, filepath.Join(featureDir, "note.txt"), "feature\n")
	runCmd(t, featureDir, gitBin, "add", "note.txt")
	runCmd(t, featureDir, gitBin, "commit", "-m", "feature change")

	runCmd(t, repoDir, gitBin, "checkout", "--detach", "HEAD")
	runCmd(t, repoDir, gitBin, "branch", "-D", "main")
	return repoDir, featureDir, gitBin, pythonBin
}

func setupUpstreamMasterAndForkMainFeatureRepo(t *testing.T) (repoDir, featureDir, gitBin, pythonBin string) {
	t.Helper()

	gitBin, err := exec.LookPath("git")
	if err != nil {
		t.Fatalf("look path git: %v", err)
	}
	pythonBin, err = exec.LookPath("python3")
	if err != nil {
		t.Fatalf("look path python3: %v", err)
	}

	repoDir = t.TempDir()
	runCmd(t, repoDir, gitBin, "init", "-b", "master")
	runCmd(t, repoDir, gitBin, "config", "user.name", "Test User")
	runCmd(t, repoDir, gitBin, "config", "user.email", "test@example.com")

	filePath := filepath.Join(repoDir, "note.txt")
	if err := os.WriteFile(filePath, []byte("base\n"), 0o644); err != nil {
		t.Fatalf("write base file: %v", err)
	}
	runCmd(t, repoDir, gitBin, "add", "note.txt")
	runCmd(t, repoDir, gitBin, "commit", "-m", "base")

	runCmd(t, repoDir, gitBin, "update-ref", "refs/remotes/upstream/master", "HEAD")
	runCmd(t, repoDir, gitBin, "update-ref", "refs/remotes/fork/main", "HEAD")

	featureDir = filepath.Join(t.TempDir(), "feature")
	runCmd(t, repoDir, gitBin, "worktree", "add", featureDir, "-b", "feature")

	appendLine(t, filepath.Join(featureDir, "note.txt"), "feature\n")
	runCmd(t, featureDir, gitBin, "add", "note.txt")
	runCmd(t, featureDir, gitBin, "commit", "-m", "feature change")

	runCmd(t, repoDir, gitBin, "checkout", "--detach", "HEAD")
	runCmd(t, repoDir, gitBin, "branch", "-D", "master")
	return repoDir, featureDir, gitBin, pythonBin
}

func setupRemoteHeadOverridesOriginMainRepo(t *testing.T) (repoDir, featureDir, gitBin, pythonBin string) {
	t.Helper()

	gitBin, err := exec.LookPath("git")
	if err != nil {
		t.Fatalf("look path git: %v", err)
	}
	pythonBin, err = exec.LookPath("python3")
	if err != nil {
		t.Fatalf("look path python3: %v", err)
	}

	repoDir = t.TempDir()
	runCmd(t, repoDir, gitBin, "init", "-b", "develop")
	runCmd(t, repoDir, gitBin, "config", "user.name", "Test User")
	runCmd(t, repoDir, gitBin, "config", "user.email", "test@example.com")

	filePath := filepath.Join(repoDir, "note.txt")
	if err := os.WriteFile(filePath, []byte("base\n"), 0o644); err != nil {
		t.Fatalf("write base file: %v", err)
	}
	runCmd(t, repoDir, gitBin, "add", "note.txt")
	runCmd(t, repoDir, gitBin, "commit", "-m", "base")

	runCmd(t, repoDir, gitBin, "update-ref", "refs/remotes/origin/develop", "HEAD")
	runCmd(t, repoDir, gitBin, "update-ref", "refs/remotes/origin/main", "HEAD")
	runCmd(t, repoDir, gitBin, "symbolic-ref", "refs/remotes/origin/HEAD", "refs/remotes/origin/develop")

	featureDir = filepath.Join(t.TempDir(), "feature")
	runCmd(t, repoDir, gitBin, "worktree", "add", featureDir, "-b", "feature")

	appendLine(t, filepath.Join(featureDir, "note.txt"), "feature\n")
	runCmd(t, featureDir, gitBin, "add", "note.txt")
	runCmd(t, featureDir, gitBin, "commit", "-m", "feature change")

	runCmd(t, repoDir, gitBin, "checkout", "--detach", "HEAD")
	runCmd(t, repoDir, gitBin, "branch", "-D", "develop")
	return repoDir, featureDir, gitBin, pythonBin
}

func runRepoScriptJSON(t *testing.T, repoDir, pythonBin, scriptRel string, target any) {
	t.Helper()

	scriptPath, err := filepath.Abs(filepath.Clean(scriptRel))
	if err != nil {
		t.Fatalf("abs script path: %v", err)
	}
	cmd := exec.Command(pythonBin, scriptPath, "--json", "main")
	cmd.Dir = repoDir
	cmd.Env = append(os.Environ(), "PYTHONDONTWRITEBYTECODE=1")
	output := runCommandSuccess(t, cmd)
	mustUnmarshalJSON(t, output, target)
}

func runRepoScriptTextWithEnv(t *testing.T, repoDir, pythonBin, scriptRel string, env []string, args ...string) string {
	t.Helper()

	scriptPath, err := filepath.Abs(filepath.Clean(scriptRel))
	if err != nil {
		t.Fatalf("abs script path: %v", err)
	}
	cmd := exec.Command(pythonBin, append([]string{scriptPath}, args...)...)
	cmd.Dir = repoDir
	cmd.Env = append(append([]string{}, os.Environ()...), env...)
	return string(runCommandSuccess(t, cmd))
}

func runRepoScriptText(t *testing.T, repoDir, pythonBin, scriptRel string, args ...string) string {
	t.Helper()

	scriptPath, err := filepath.Abs(filepath.Clean(scriptRel))
	if err != nil {
		t.Fatalf("abs script path: %v", err)
	}
	cmd := exec.Command(pythonBin, append([]string{scriptPath}, args...)...)
	cmd.Dir = repoDir
	cmd.Env = append(os.Environ(), "PYTHONDONTWRITEBYTECODE=1")
	return string(runCommandSuccess(t, cmd))
}

func runRepoScriptFailureText(t *testing.T, repoDir, pythonBin, scriptRel string, args ...string) string {
	t.Helper()

	scriptPath, err := filepath.Abs(filepath.Clean(scriptRel))
	if err != nil {
		t.Fatalf("abs script path: %v", err)
	}
	cmd := exec.Command(pythonBin, append([]string{scriptPath}, args...)...)
	cmd.Dir = repoDir
	cmd.Env = append(os.Environ(), "PYTHONDONTWRITEBYTECODE=1")
	return string(runCommandFailure(t, cmd))
}

func runRepoScriptFailureTextWithEnv(t *testing.T, repoDir, pythonBin, scriptRel string, env []string, args ...string) string {
	t.Helper()

	scriptPath, err := filepath.Abs(filepath.Clean(scriptRel))
	if err != nil {
		t.Fatalf("abs script path: %v", err)
	}
	cmd := exec.Command(pythonBin, append([]string{scriptPath}, args...)...)
	cmd.Dir = repoDir
	cmd.Env = append(append([]string{}, os.Environ()...), env...)
	return string(runCommandFailure(t, cmd))
}

func runRepoScriptTextAtDir(t *testing.T, dir, pythonBin, scriptRel string, args ...string) string {
	t.Helper()

	scriptPath, err := filepath.Abs(filepath.Clean(scriptRel))
	if err != nil {
		t.Fatalf("abs script path: %v", err)
	}
	cmd := exec.Command(pythonBin, append([]string{scriptPath}, args...)...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "PYTHONDONTWRITEBYTECODE=1")
	return string(runCommandSuccess(t, cmd))
}

func TestWorktreePrunePlanUsesConfiguredPythonForNestedScripts(t *testing.T) {
	repoDir, _, pythonBin := setupPatchEquivalentRepo(t)

	logPath := filepath.Join(t.TempDir(), "python-wrapper.log")
	wrapperPath := filepath.Join(t.TempDir(), "python-wrapper.sh")
	wrapperScript := "#!/usr/bin/env bash\n" +
		"echo \"$@\" >>\"" + logPath + "\"\n" +
		"exec \"" + pythonBin + "\" \"$@\"\n"
	if err := os.WriteFile(wrapperPath, []byte(wrapperScript), 0o755); err != nil {
		t.Fatalf("write wrapper: %v", err)
	}

	scriptPath, err := filepath.Abs(filepath.Clean("../../scripts/worktree_prune_plan.py"))
	if err != nil {
		t.Fatalf("abs script path: %v", err)
	}
	cmd := exec.Command(pythonBin, scriptPath, "--json", "main")
	cmd.Dir = repoDir
	cmd.Env = append(os.Environ(),
		"PYTHONDONTWRITEBYTECODE=1",
		"PYTHON="+wrapperPath,
	)
	runCommandSuccess(t, cmd)

	logData := mustReadFile(t, logPath)
	if !strings.Contains(string(logData), "worktree_audit.py") {
		t.Fatalf("expected nested script to use configured python, got: %s", logData)
	}
}

func TestWorktreePrunePlanSupportsMultiTokenConfiguredPython(t *testing.T) {
	repoDir, _, pythonBin := setupPatchEquivalentRepo(t)

	scriptPath, err := filepath.Abs(filepath.Clean("../../scripts/worktree_prune_plan.py"))
	if err != nil {
		t.Fatalf("abs script path: %v", err)
	}
	cmd := exec.Command(pythonBin, scriptPath, "--json", "main")
	cmd.Dir = repoDir
	cmd.Env = append(os.Environ(),
		"PYTHONDONTWRITEBYTECODE=1",
		"PYTHON="+pythonBin+" -B",
	)
	output := runCommandSuccess(t, cmd)

	var payload struct {
		Summary struct {
			Total int `json:"total"`
		} `json:"summary"`
		Rows []struct {
			Branch string `json:"branch"`
		} `json:"rows"`
	}
	mustUnmarshalJSON(t, output, &payload)
	if payload.Summary.Total != 1 || len(payload.Rows) != 1 || payload.Rows[0].Branch != "feature" {
		t.Fatalf("unexpected payload: %#v", payload)
	}
}

func TestWorktreePrunePlanRejectsMissingConfiguredPythonGracefully(t *testing.T) {
	repoDir, _, pythonBin := setupPatchEquivalentRepo(t)

	scriptPath, err := filepath.Abs(filepath.Clean("../../scripts/worktree_prune_plan.py"))
	if err != nil {
		t.Fatalf("abs script path: %v", err)
	}
	cmd := exec.Command(pythonBin, scriptPath, "--json", "main")
	cmd.Dir = repoDir
	cmd.Env = append(os.Environ(),
		"PYTHONDONTWRITEBYTECODE=1",
		"PYTHON=/definitely/missing/python3",
	)
	output := string(runCommandFailure(t, cmd))
	if strings.Contains(output, "Traceback") {
		t.Fatalf("unexpected traceback: %s", output)
	}
	if !strings.Contains(output, "unable to execute worktree_audit.py via configured python command") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestWorktreePrunePlanRejectsNonASCIIMissingConfiguredPythonGracefully(t *testing.T) {
	repoDir, _, pythonBin := setupPatchEquivalentRepo(t)

	output := runRepoScriptFailureTextWithEnv(
		t,
		repoDir,
		pythonBin,
		"../../scripts/worktree_prune_plan.py",
		[]string{
			"PYTHONDONTWRITEBYTECODE=1",
			"LC_ALL=C",
			"LANG=C",
			"PYTHONCOERCECLOCALE=0",
			"PYTHONUTF8=0",
			"PYTHON=/缺失/python3",
		},
		"--json",
		"main",
	)
	if strings.Contains(output, "Traceback") {
		t.Fatalf("unexpected traceback: %s", output)
	}
	if !strings.Contains(output, "unable to execute worktree_audit.py via configured python command") {
		t.Fatalf("unexpected output: %s", output)
	}
	if !strings.Contains(output, "/缺失/python3") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestWorktreeAuditJSONMarksPatchEquivalentBranchIntegrated(t *testing.T) {
	repoDir, _, pythonBin := setupPatchEquivalentRepo(t)

	var payload struct {
		Summary map[string]any `json:"summary"`
		Rows    []struct {
			Branch string `json:"branch"`
			State  string `json:"state"`
		} `json:"rows"`
	}
	runRepoScriptJSON(t, repoDir, pythonBin, "../../scripts/worktree_audit.py", &payload)

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
		t.Fatalf("missing feature row: %#v", payload.Rows)
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
	repoDir, _, pythonBin := setupPatchEquivalentRepo(t)

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
	runRepoScriptJSON(t, repoDir, pythonBin, "../../scripts/worktree_prune_plan.py", &payload)
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

func TestWorktreePrunePlanTextUsesPruneCandidateWordingWhenEmpty(t *testing.T) {
	repoDir, _, pythonBin := setupActiveFeatureRepo(t)

	output := runRepoScriptText(t, repoDir, pythonBin, "../../scripts/worktree_prune_plan.py", "main")
	if strings.Contains(output, "# no merged worktrees to prune") {
		t.Fatalf("unexpected merged-only wording: %s", output)
	}
	if !strings.Contains(output, "# no prune-candidate worktrees to prune") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestWorktreeAuditDefaultsToLocalMainWhenOriginMainMissing(t *testing.T) {
	repoDir, _, pythonBin := setupActiveFeatureRepo(t)

	output := runRepoScriptText(t, repoDir, pythonBin, "../../scripts/worktree_audit.py")
	if !strings.Contains(output, "# Base ref: main") {
		t.Fatalf("unexpected output: %s", output)
	}
	if !strings.Contains(output, "\tfeature\t") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestWorktreePrunePlanDefaultsToLocalMainWhenOriginMainMissing(t *testing.T) {
	repoDir, _, pythonBin := setupActiveFeatureRepo(t)

	output := runRepoScriptText(t, repoDir, pythonBin, "../../scripts/worktree_prune_plan.py")
	if !strings.Contains(output, "# Base ref: main") {
		t.Fatalf("unexpected output: %s", output)
	}
	if !strings.Contains(output, "# no prune-candidate worktrees to prune") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestWorktreeCleanupReportDefaultsToLocalMainWhenOriginMainMissing(t *testing.T) {
	repoDir, _, pythonBin := setupActiveFeatureRepo(t)

	output := runRepoScriptText(t, repoDir, pythonBin, "../../scripts/worktree_cleanup_report.py")
	if !strings.Contains(output, "# Worktree Cleanup Report") {
		t.Fatalf("unexpected output: %s", output)
	}
	if !strings.Contains(output, "- Base ref: `main`") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestWorktreeAuditDefaultsToCurrentBranchWhenNoStandardBaseRefExists(t *testing.T) {
	repoDir, _, pythonBin := setupActiveTrunkRepo(t)

	output := runRepoScriptText(t, repoDir, pythonBin, "../../scripts/worktree_audit.py")
	if !strings.Contains(output, "# Base ref: trunk") {
		t.Fatalf("unexpected output: %s", output)
	}
	if !strings.Contains(output, "\tfeature\t") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestWorktreePrunePlanDefaultsToCurrentBranchWhenNoStandardBaseRefExists(t *testing.T) {
	repoDir, _, pythonBin := setupActiveTrunkRepo(t)

	output := runRepoScriptText(t, repoDir, pythonBin, "../../scripts/worktree_prune_plan.py")
	if !strings.Contains(output, "# Base ref: trunk") {
		t.Fatalf("unexpected output: %s", output)
	}
	if !strings.Contains(output, "# no prune-candidate worktrees to prune") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestWorktreeCleanupReportDefaultsToCurrentBranchWhenNoStandardBaseRefExists(t *testing.T) {
	repoDir, _, pythonBin := setupActiveTrunkRepo(t)

	output := runRepoScriptText(t, repoDir, pythonBin, "../../scripts/worktree_cleanup_report.py")
	if !strings.Contains(output, "# Worktree Cleanup Report") {
		t.Fatalf("unexpected output: %s", output)
	}
	if !strings.Contains(output, "- Base ref: `trunk`") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestWorktreeAuditPrefersTrunkOverCurrentBranchWhenRunFromFeatureWorktree(t *testing.T) {
	_, featureDir, _, pythonBin := setupActiveTrunkRepoWithFeatureWorktree(t)

	output := runRepoScriptTextAtDir(t, featureDir, pythonBin, "../../scripts/worktree_audit.py")
	if !strings.Contains(output, "# Base ref: trunk") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestWorktreePrunePlanPrefersTrunkOverCurrentBranchWhenRunFromFeatureWorktree(t *testing.T) {
	_, featureDir, _, pythonBin := setupActiveTrunkRepoWithFeatureWorktree(t)

	output := runRepoScriptTextAtDir(t, featureDir, pythonBin, "../../scripts/worktree_prune_plan.py")
	if !strings.Contains(output, "# Base ref: trunk") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestWorktreeCleanupReportPrefersTrunkOverCurrentBranchWhenRunFromFeatureWorktree(t *testing.T) {
	_, featureDir, _, pythonBin := setupActiveTrunkRepoWithFeatureWorktree(t)

	output := runRepoScriptTextAtDir(t, featureDir, pythonBin, "../../scripts/worktree_cleanup_report.py")
	if !strings.Contains(output, "- Base ref: `trunk`") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestWorktreeAuditPrefersOriginTrunkOverCurrentBranchWhenOnlyRemoteTrunkExists(t *testing.T) {
	_, featureDir, _, pythonBin := setupRemoteTrunkOnlyFeatureRepo(t)

	output := runRepoScriptTextAtDir(t, featureDir, pythonBin, "../../scripts/worktree_audit.py")
	if !strings.Contains(output, "# Base ref: origin/trunk") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestWorktreePrunePlanPrefersOriginTrunkOverCurrentBranchWhenOnlyRemoteTrunkExists(t *testing.T) {
	_, featureDir, _, pythonBin := setupRemoteTrunkOnlyFeatureRepo(t)

	output := runRepoScriptTextAtDir(t, featureDir, pythonBin, "../../scripts/worktree_prune_plan.py")
	if !strings.Contains(output, "# Base ref: origin/trunk") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestWorktreeCleanupReportPrefersOriginTrunkOverCurrentBranchWhenOnlyRemoteTrunkExists(t *testing.T) {
	_, featureDir, _, pythonBin := setupRemoteTrunkOnlyFeatureRepo(t)

	output := runRepoScriptTextAtDir(t, featureDir, pythonBin, "../../scripts/worktree_cleanup_report.py")
	if !strings.Contains(output, "- Base ref: `origin/trunk`") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestWorktreeAuditPrefersOriginMasterOverCurrentBranchWhenOnlyRemoteMasterExists(t *testing.T) {
	_, featureDir, _, pythonBin := setupRemoteMasterOnlyFeatureRepo(t)

	output := runRepoScriptTextAtDir(t, featureDir, pythonBin, "../../scripts/worktree_audit.py")
	if !strings.Contains(output, "# Base ref: origin/master") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestWorktreePrunePlanPrefersOriginMasterOverCurrentBranchWhenOnlyRemoteMasterExists(t *testing.T) {
	_, featureDir, _, pythonBin := setupRemoteMasterOnlyFeatureRepo(t)

	output := runRepoScriptTextAtDir(t, featureDir, pythonBin, "../../scripts/worktree_prune_plan.py")
	if !strings.Contains(output, "# Base ref: origin/master") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestWorktreeCleanupReportPrefersOriginMasterOverCurrentBranchWhenOnlyRemoteMasterExists(t *testing.T) {
	_, featureDir, _, pythonBin := setupRemoteMasterOnlyFeatureRepo(t)

	output := runRepoScriptTextAtDir(t, featureDir, pythonBin, "../../scripts/worktree_cleanup_report.py")
	if !strings.Contains(output, "- Base ref: `origin/master`") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestWorktreeAuditPrefersOriginHeadTargetOverCurrentBranchWhenOnlyRemoteHeadExists(t *testing.T) {
	_, featureDir, _, pythonBin := setupRemoteHEADOnlyFeatureRepo(t)

	output := runRepoScriptTextAtDir(t, featureDir, pythonBin, "../../scripts/worktree_audit.py")
	if !strings.Contains(output, "# Base ref: origin/develop") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestWorktreePrunePlanPrefersOriginHeadTargetOverCurrentBranchWhenOnlyRemoteHeadExists(t *testing.T) {
	_, featureDir, _, pythonBin := setupRemoteHEADOnlyFeatureRepo(t)

	output := runRepoScriptTextAtDir(t, featureDir, pythonBin, "../../scripts/worktree_prune_plan.py")
	if !strings.Contains(output, "# Base ref: origin/develop") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestWorktreeCleanupReportPrefersOriginHeadTargetOverCurrentBranchWhenOnlyRemoteHeadExists(t *testing.T) {
	_, featureDir, _, pythonBin := setupRemoteHEADOnlyFeatureRepo(t)

	output := runRepoScriptTextAtDir(t, featureDir, pythonBin, "../../scripts/worktree_cleanup_report.py")
	if !strings.Contains(output, "- Base ref: `origin/develop`") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestWorktreeAuditPrefersNonOriginHeadTargetOverCurrentBranchWhenOnlyRemoteHeadExists(t *testing.T) {
	_, featureDir, _, pythonBin := setupNonOriginRemoteHEADOnlyFeatureRepo(t)

	output := runRepoScriptTextAtDir(t, featureDir, pythonBin, "../../scripts/worktree_audit.py")
	if !strings.Contains(output, "# Base ref: upstream/develop") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestWorktreePrunePlanPrefersNonOriginHeadTargetOverCurrentBranchWhenOnlyRemoteHeadExists(t *testing.T) {
	_, featureDir, _, pythonBin := setupNonOriginRemoteHEADOnlyFeatureRepo(t)

	output := runRepoScriptTextAtDir(t, featureDir, pythonBin, "../../scripts/worktree_prune_plan.py")
	if !strings.Contains(output, "# Base ref: upstream/develop") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestWorktreeCleanupReportPrefersNonOriginHeadTargetOverCurrentBranchWhenOnlyRemoteHeadExists(t *testing.T) {
	_, featureDir, _, pythonBin := setupNonOriginRemoteHEADOnlyFeatureRepo(t)

	output := runRepoScriptTextAtDir(t, featureDir, pythonBin, "../../scripts/worktree_cleanup_report.py")
	if !strings.Contains(output, "- Base ref: `upstream/develop`") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestWorktreeAuditPrefersNonOriginMainOverCurrentBranchWhenOnlyRemoteMainExists(t *testing.T) {
	_, featureDir, _, pythonBin := setupNonOriginRemoteNamedBranchOnlyFeatureRepo(t, "main")

	output := runRepoScriptTextAtDir(t, featureDir, pythonBin, "../../scripts/worktree_audit.py")
	if !strings.Contains(output, "# Base ref: upstream/main") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestWorktreePrunePlanPrefersNonOriginMainOverCurrentBranchWhenOnlyRemoteMainExists(t *testing.T) {
	_, featureDir, _, pythonBin := setupNonOriginRemoteNamedBranchOnlyFeatureRepo(t, "main")

	output := runRepoScriptTextAtDir(t, featureDir, pythonBin, "../../scripts/worktree_prune_plan.py")
	if !strings.Contains(output, "# Base ref: upstream/main") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestWorktreeCleanupReportPrefersNonOriginMainOverCurrentBranchWhenOnlyRemoteMainExists(t *testing.T) {
	_, featureDir, _, pythonBin := setupNonOriginRemoteNamedBranchOnlyFeatureRepo(t, "main")

	output := runRepoScriptTextAtDir(t, featureDir, pythonBin, "../../scripts/worktree_cleanup_report.py")
	if !strings.Contains(output, "- Base ref: `upstream/main`") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestWorktreeAuditPrefersNonOriginMasterOverCurrentBranchWhenOnlyRemoteMasterExists(t *testing.T) {
	_, featureDir, _, pythonBin := setupNonOriginRemoteNamedBranchOnlyFeatureRepo(t, "master")

	output := runRepoScriptTextAtDir(t, featureDir, pythonBin, "../../scripts/worktree_audit.py")
	if !strings.Contains(output, "# Base ref: upstream/master") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestWorktreePrunePlanPrefersNonOriginMasterOverCurrentBranchWhenOnlyRemoteMasterExists(t *testing.T) {
	_, featureDir, _, pythonBin := setupNonOriginRemoteNamedBranchOnlyFeatureRepo(t, "master")

	output := runRepoScriptTextAtDir(t, featureDir, pythonBin, "../../scripts/worktree_prune_plan.py")
	if !strings.Contains(output, "# Base ref: upstream/master") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestWorktreeCleanupReportPrefersNonOriginMasterOverCurrentBranchWhenOnlyRemoteMasterExists(t *testing.T) {
	_, featureDir, _, pythonBin := setupNonOriginRemoteNamedBranchOnlyFeatureRepo(t, "master")

	output := runRepoScriptTextAtDir(t, featureDir, pythonBin, "../../scripts/worktree_cleanup_report.py")
	if !strings.Contains(output, "- Base ref: `upstream/master`") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestWorktreeAuditPrefersNonOriginTrunkOverCurrentBranchWhenOnlyRemoteTrunkExists(t *testing.T) {
	_, featureDir, _, pythonBin := setupNonOriginRemoteNamedBranchOnlyFeatureRepo(t, "trunk")

	output := runRepoScriptTextAtDir(t, featureDir, pythonBin, "../../scripts/worktree_audit.py")
	if !strings.Contains(output, "# Base ref: upstream/trunk") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestWorktreePrunePlanPrefersNonOriginTrunkOverCurrentBranchWhenOnlyRemoteTrunkExists(t *testing.T) {
	_, featureDir, _, pythonBin := setupNonOriginRemoteNamedBranchOnlyFeatureRepo(t, "trunk")

	output := runRepoScriptTextAtDir(t, featureDir, pythonBin, "../../scripts/worktree_prune_plan.py")
	if !strings.Contains(output, "# Base ref: upstream/trunk") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestWorktreeCleanupReportPrefersNonOriginTrunkOverCurrentBranchWhenOnlyRemoteTrunkExists(t *testing.T) {
	_, featureDir, _, pythonBin := setupNonOriginRemoteNamedBranchOnlyFeatureRepo(t, "trunk")

	output := runRepoScriptTextAtDir(t, featureDir, pythonBin, "../../scripts/worktree_cleanup_report.py")
	if !strings.Contains(output, "- Base ref: `upstream/trunk`") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestWorktreeAuditPrefersUpstreamHeadTargetOverOtherNonOriginHeadTargets(t *testing.T) {
	_, featureDir, _, pythonBin := setupUpstreamAndForkRemoteHEADFeatureRepo(t)

	output := runRepoScriptTextAtDir(t, featureDir, pythonBin, "../../scripts/worktree_audit.py")
	if !strings.Contains(output, "# Base ref: upstream/develop") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestWorktreePrunePlanPrefersUpstreamHeadTargetOverOtherNonOriginHeadTargets(t *testing.T) {
	_, featureDir, _, pythonBin := setupUpstreamAndForkRemoteHEADFeatureRepo(t)

	output := runRepoScriptTextAtDir(t, featureDir, pythonBin, "../../scripts/worktree_prune_plan.py")
	if !strings.Contains(output, "# Base ref: upstream/develop") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestWorktreeCleanupReportPrefersUpstreamHeadTargetOverOtherNonOriginHeadTargets(t *testing.T) {
	_, featureDir, _, pythonBin := setupUpstreamAndForkRemoteHEADFeatureRepo(t)

	output := runRepoScriptTextAtDir(t, featureDir, pythonBin, "../../scripts/worktree_cleanup_report.py")
	if !strings.Contains(output, "- Base ref: `upstream/develop`") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestWorktreeAuditPrefersUpstreamMainOverOtherNonOriginMainRefs(t *testing.T) {
	_, featureDir, _, pythonBin := setupUpstreamAndForkRemoteMainFeatureRepo(t)

	output := runRepoScriptTextAtDir(t, featureDir, pythonBin, "../../scripts/worktree_audit.py")
	if !strings.Contains(output, "# Base ref: upstream/main") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestWorktreePrunePlanPrefersUpstreamMainOverOtherNonOriginMainRefs(t *testing.T) {
	_, featureDir, _, pythonBin := setupUpstreamAndForkRemoteMainFeatureRepo(t)

	output := runRepoScriptTextAtDir(t, featureDir, pythonBin, "../../scripts/worktree_prune_plan.py")
	if !strings.Contains(output, "# Base ref: upstream/main") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestWorktreeCleanupReportPrefersUpstreamMainOverOtherNonOriginMainRefs(t *testing.T) {
	_, featureDir, _, pythonBin := setupUpstreamAndForkRemoteMainFeatureRepo(t)

	output := runRepoScriptTextAtDir(t, featureDir, pythonBin, "../../scripts/worktree_cleanup_report.py")
	if !strings.Contains(output, "- Base ref: `upstream/main`") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestWorktreeAuditPrefersUpstreamMasterOverForkMainAcrossBranchTypes(t *testing.T) {
	_, featureDir, _, pythonBin := setupUpstreamMasterAndForkMainFeatureRepo(t)

	output := runRepoScriptTextAtDir(t, featureDir, pythonBin, "../../scripts/worktree_audit.py")
	if !strings.Contains(output, "# Base ref: upstream/master") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestWorktreePrunePlanPrefersUpstreamMasterOverForkMainAcrossBranchTypes(t *testing.T) {
	_, featureDir, _, pythonBin := setupUpstreamMasterAndForkMainFeatureRepo(t)

	output := runRepoScriptTextAtDir(t, featureDir, pythonBin, "../../scripts/worktree_prune_plan.py")
	if !strings.Contains(output, "# Base ref: upstream/master") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestWorktreeCleanupReportPrefersUpstreamMasterOverForkMainAcrossBranchTypes(t *testing.T) {
	_, featureDir, _, pythonBin := setupUpstreamMasterAndForkMainFeatureRepo(t)

	output := runRepoScriptTextAtDir(t, featureDir, pythonBin, "../../scripts/worktree_cleanup_report.py")
	if !strings.Contains(output, "- Base ref: `upstream/master`") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestWorktreeAuditPrefersOriginHeadTargetOverOriginMainWhenBothExist(t *testing.T) {
	_, featureDir, _, pythonBin := setupRemoteHeadOverridesOriginMainRepo(t)

	output := runRepoScriptTextAtDir(t, featureDir, pythonBin, "../../scripts/worktree_audit.py")
	if !strings.Contains(output, "# Base ref: origin/develop") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestWorktreePrunePlanPrefersOriginHeadTargetOverOriginMainWhenBothExist(t *testing.T) {
	_, featureDir, _, pythonBin := setupRemoteHeadOverridesOriginMainRepo(t)

	output := runRepoScriptTextAtDir(t, featureDir, pythonBin, "../../scripts/worktree_prune_plan.py")
	if !strings.Contains(output, "# Base ref: origin/develop") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestWorktreeCleanupReportPrefersOriginHeadTargetOverOriginMainWhenBothExist(t *testing.T) {
	_, featureDir, _, pythonBin := setupRemoteHeadOverridesOriginMainRepo(t)

	output := runRepoScriptTextAtDir(t, featureDir, pythonBin, "../../scripts/worktree_cleanup_report.py")
	if !strings.Contains(output, "- Base ref: `origin/develop`") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestWorktreeAuditDefaultsToHEADWhenDetachedAndNoStandardBaseRefExists(t *testing.T) {
	repoDir, _, pythonBin := setupDetachedTrunkRepo(t)

	output := runRepoScriptText(t, repoDir, pythonBin, "../../scripts/worktree_audit.py")
	if !strings.Contains(output, "# Base ref: HEAD") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestWorktreePrunePlanDefaultsToHEADWhenDetachedAndNoStandardBaseRefExists(t *testing.T) {
	repoDir, _, pythonBin := setupDetachedTrunkRepo(t)

	output := runRepoScriptText(t, repoDir, pythonBin, "../../scripts/worktree_prune_plan.py")
	if !strings.Contains(output, "# Base ref: HEAD") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestWorktreeCleanupReportDefaultsToHEADWhenDetachedAndNoStandardBaseRefExists(t *testing.T) {
	repoDir, _, pythonBin := setupDetachedTrunkRepo(t)

	output := runRepoScriptText(t, repoDir, pythonBin, "../../scripts/worktree_cleanup_report.py")
	if !strings.Contains(output, "- Base ref: `HEAD`") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestWorktreeAuditWritesTextWithExplicitUTF8Stdout(t *testing.T) {
	repoDir, _, pythonBin := setupUnicodePatchEquivalentRepo(t)

	output := runRepoScriptTextWithEnv(
		t,
		repoDir,
		pythonBin,
		"../../scripts/worktree_audit.py",
		[]string{
			"PYTHONDONTWRITEBYTECODE=1",
			"LC_ALL=C",
			"LANG=C",
			"PYTHONCOERCECLOCALE=0",
			"PYTHONUTF8=0",
		},
		"main",
	)
	if !strings.Contains(output, "特性") {
		t.Fatalf("unexpected output: %s", output)
	}
	if !strings.Contains(output, "工作树") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestWorktreeCleanupReportWritesTextWithExplicitUTF8Stdout(t *testing.T) {
	repoDir, _, pythonBin := setupUnicodePatchEquivalentRepo(t)

	output := runRepoScriptTextWithEnv(
		t,
		repoDir,
		pythonBin,
		"../../scripts/worktree_cleanup_report.py",
		[]string{
			"PYTHONDONTWRITEBYTECODE=1",
			"LC_ALL=C",
			"LANG=C",
			"PYTHONCOERCECLOCALE=0",
			"PYTHONUTF8=0",
		},
		"main",
	)
	if !strings.Contains(output, "# Worktree Cleanup Report") {
		t.Fatalf("unexpected output: %s", output)
	}
	if !strings.Contains(output, "特性") {
		t.Fatalf("unexpected output: %s", output)
	}
	if !strings.Contains(output, "工作树") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestWorktreeAuditWritesJSONWithExplicitUTF8Characters(t *testing.T) {
	repoDir, _, pythonBin := setupUnicodePatchEquivalentRepo(t)

	output := runRepoScriptTextWithEnv(
		t,
		repoDir,
		pythonBin,
		"../../scripts/worktree_audit.py",
		[]string{
			"PYTHONDONTWRITEBYTECODE=1",
			"LC_ALL=C",
			"LANG=C",
			"PYTHONCOERCECLOCALE=0",
			"PYTHONUTF8=0",
		},
		"--json",
		"main",
	)
	if !strings.Contains(output, "特性") || !strings.Contains(output, "工作树") {
		t.Fatalf("unexpected output: %s", output)
	}
	if strings.Contains(output, "\\u7279\\u6027") || strings.Contains(output, "\\u5de5\\u4f5c\\u6811") {
		t.Fatalf("unexpected escaped unicode: %s", output)
	}
}

func TestWorktreePrunePlanWritesJSONWithExplicitUTF8Characters(t *testing.T) {
	repoDir, _, pythonBin := setupUnicodePatchEquivalentRepo(t)

	output := runRepoScriptTextWithEnv(
		t,
		repoDir,
		pythonBin,
		"../../scripts/worktree_prune_plan.py",
		[]string{
			"PYTHONDONTWRITEBYTECODE=1",
			"LC_ALL=C",
			"LANG=C",
			"PYTHONCOERCECLOCALE=0",
			"PYTHONUTF8=0",
		},
		"--json",
		"main",
	)
	if !strings.Contains(output, "特性") || !strings.Contains(output, "工作树") {
		t.Fatalf("unexpected output: %s", output)
	}
	if strings.Contains(output, "\\u7279\\u6027") || strings.Contains(output, "\\u5de5\\u4f5c\\u6811") {
		t.Fatalf("unexpected escaped unicode: %s", output)
	}
}

func TestWorktreeCleanupReportWritesJSONWithExplicitUTF8Characters(t *testing.T) {
	repoDir, _, pythonBin := setupUnicodePatchEquivalentRepo(t)

	output := runRepoScriptTextWithEnv(
		t,
		repoDir,
		pythonBin,
		"../../scripts/worktree_cleanup_report.py",
		[]string{
			"PYTHONDONTWRITEBYTECODE=1",
			"LC_ALL=C",
			"LANG=C",
			"PYTHONCOERCECLOCALE=0",
			"PYTHONUTF8=0",
		},
		"--json",
		"main",
	)
	if !strings.Contains(output, "特性") || !strings.Contains(output, "工作树") {
		t.Fatalf("unexpected output: %s", output)
	}
	if strings.Contains(output, "\\u7279\\u6027") || strings.Contains(output, "\\u5de5\\u4f5c\\u6811") {
		t.Fatalf("unexpected escaped unicode: %s", output)
	}
}

func TestWorktreeAuditRejectsUnknownOptionGracefully(t *testing.T) {
	repoDir, _, pythonBin := setupPatchEquivalentRepo(t)

	output := runRepoScriptFailureText(t, repoDir, pythonBin, "../../scripts/worktree_audit.py", "--bogus")
	if strings.Contains(output, "usage: git rev-list") {
		t.Fatalf("unexpected git usage leakage: %s", output)
	}
	if !strings.Contains(output, "unknown option: --bogus") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestWorktreeAuditRejectsDuplicateOptionGracefully(t *testing.T) {
	repoDir, _, pythonBin := setupPatchEquivalentRepo(t)

	output := runRepoScriptFailureText(t, repoDir, pythonBin, "../../scripts/worktree_audit.py", "--json", "--json")
	if strings.Contains(output, "usage: git rev-list") {
		t.Fatalf("unexpected git usage leakage: %s", output)
	}
	if !strings.Contains(output, "duplicate option: --json") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestWorktreeAuditRejectsUnknownBaseRefGracefully(t *testing.T) {
	repoDir, _, pythonBin := setupPatchEquivalentRepo(t)

	output := runRepoScriptFailureText(t, repoDir, pythonBin, "../../scripts/worktree_audit.py", "does-not-exist")
	if strings.Contains(output, "usage: git rev-list") {
		t.Fatalf("unexpected git usage leakage: %s", output)
	}
	if !strings.Contains(output, "unknown base ref: does-not-exist") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestWorktreeAuditRejectsMultipleViewFiltersGracefully(t *testing.T) {
	repoDir, _, pythonBin := setupPatchEquivalentRepo(t)

	output := runRepoScriptFailureText(t, repoDir, pythonBin, "../../scripts/worktree_audit.py", "--merged-only", "--integrated-only")
	if strings.Contains(output, "usage: git rev-list") {
		t.Fatalf("unexpected git usage leakage: %s", output)
	}
	if !strings.Contains(output, "choose at most one of --merged-only, --integrated-only, --prune-only, or --stale-only") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestWorktreePrunePlanRejectsUnknownOptionGracefully(t *testing.T) {
	repoDir, _, pythonBin := setupPatchEquivalentRepo(t)

	output := runRepoScriptFailureText(t, repoDir, pythonBin, "../../scripts/worktree_prune_plan.py", "--bogus")
	if strings.Contains(output, "usage: git rev-list") {
		t.Fatalf("unexpected git usage leakage: %s", output)
	}
	if !strings.Contains(output, "unknown option: --bogus") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestWorktreePrunePlanRejectsDuplicateOptionGracefully(t *testing.T) {
	repoDir, _, pythonBin := setupPatchEquivalentRepo(t)

	output := runRepoScriptFailureText(t, repoDir, pythonBin, "../../scripts/worktree_prune_plan.py", "--json", "--json")
	if strings.Contains(output, "usage: git rev-list") {
		t.Fatalf("unexpected git usage leakage: %s", output)
	}
	if !strings.Contains(output, "duplicate option: --json") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestWorktreePrunePlanRejectsUnknownBaseRefGracefully(t *testing.T) {
	repoDir, _, pythonBin := setupPatchEquivalentRepo(t)

	output := runRepoScriptFailureText(t, repoDir, pythonBin, "../../scripts/worktree_prune_plan.py", "does-not-exist")
	if strings.Contains(output, "usage: git rev-list") {
		t.Fatalf("unexpected git usage leakage: %s", output)
	}
	if !strings.Contains(output, "unknown base ref: does-not-exist") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestWorktreeCleanupReportCreatesMissingParentDir(t *testing.T) {
	repoDir, _, pythonBin := setupPatchEquivalentRepo(t)

	outputPath := filepath.Join(t.TempDir(), "nested", "reports", "cleanup.md")
	scriptPath, err := filepath.Abs(filepath.Clean("../../scripts/worktree_cleanup_report.py"))
	if err != nil {
		t.Fatalf("abs script path: %v", err)
	}
	cmd := exec.Command(pythonBin, scriptPath, "main", outputPath)
	cmd.Dir = repoDir
	cmd.Env = append(os.Environ(), "PYTHONDONTWRITEBYTECODE=1")
	output := runCommandSuccess(t, cmd)

	reportData := mustReadFile(t, outputPath)
	if !strings.Contains(string(output), "# Worktree Cleanup Report") {
		t.Fatalf("unexpected stdout: %s", output)
	}
	if !strings.Contains(string(reportData), "# Worktree Cleanup Report") {
		t.Fatalf("unexpected report file: %s", reportData)
	}
}

func TestWorktreeCleanupReportRejectsDirectoryOutputPathGracefully(t *testing.T) {
	repoDir, _, pythonBin := setupPatchEquivalentRepo(t)

	outputDir := t.TempDir()
	output := runRepoScriptFailureText(t, repoDir, pythonBin, "../../scripts/worktree_cleanup_report.py", "main", outputDir)
	if strings.Contains(output, "Traceback") {
		t.Fatalf("unexpected traceback: %s", output)
	}
	if !strings.Contains(output, "unable to write worktree cleanup report") {
		t.Fatalf("unexpected output: %s", output)
	}
	if !strings.Contains(output, outputDir) {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestWorktreeCleanupReportRejectsUncreatableParentDirGracefully(t *testing.T) {
	repoDir, _, pythonBin := setupPatchEquivalentRepo(t)

	baseDir := t.TempDir()
	blocked := filepath.Join(baseDir, "blocked")
	if err := os.WriteFile(blocked, []byte("file\n"), 0o644); err != nil {
		t.Fatalf("write blocker file: %v", err)
	}
	outputPath := filepath.Join(blocked, "cleanup.md")
	output := runRepoScriptFailureText(t, repoDir, pythonBin, "../../scripts/worktree_cleanup_report.py", "main", outputPath)
	if strings.Contains(output, "Traceback") {
		t.Fatalf("unexpected traceback: %s", output)
	}
	if !strings.Contains(output, "unable to prepare worktree cleanup report parent dir") {
		t.Fatalf("unexpected output: %s", output)
	}
	if !strings.Contains(output, blocked) {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestWorktreeCleanupReportRendersNonASCIIOutputPathGracefully(t *testing.T) {
	repoDir, _, pythonBin := setupPatchEquivalentRepo(t)

	baseDir := t.TempDir()
	blocked := filepath.Join(baseDir, "阻止")
	if err := os.WriteFile(blocked, []byte("file\n"), 0o644); err != nil {
		t.Fatalf("write blocker file: %v", err)
	}
	outputPath := filepath.Join(blocked, "报告.md")
	output := runRepoScriptFailureTextWithEnv(
		t,
		repoDir,
		pythonBin,
		"../../scripts/worktree_cleanup_report.py",
		[]string{
			"PYTHONDONTWRITEBYTECODE=1",
			"LC_ALL=C",
			"LANG=C",
			"PYTHONCOERCECLOCALE=0",
			"PYTHONUTF8=0",
		},
		"main",
		outputPath,
	)
	if strings.Contains(output, "Traceback") {
		t.Fatalf("unexpected traceback: %s", output)
	}
	if strings.Contains(output, "\\udc") {
		t.Fatalf("unexpected surrogate escapes: %s", output)
	}
	if !strings.Contains(output, "unable to prepare worktree cleanup report parent dir") {
		t.Fatalf("unexpected output: %s", output)
	}
	if !strings.Contains(output, "阻止") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestWorktreeCleanupReportRejectsUnknownBaseRefGracefully(t *testing.T) {
	repoDir, _, pythonBin := setupPatchEquivalentRepo(t)

	output := runRepoScriptFailureText(t, repoDir, pythonBin, "../../scripts/worktree_cleanup_report.py", "does-not-exist")
	if strings.Contains(output, "Traceback") {
		t.Fatalf("unexpected traceback: %s", output)
	}
	if !strings.Contains(output, "unknown base ref: does-not-exist") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func TestWorktreeCleanupReportRejectsNonASCIIMissingConfiguredPythonGracefully(t *testing.T) {
	repoDir, _, pythonBin := setupPatchEquivalentRepo(t)

	output := runRepoScriptFailureTextWithEnv(
		t,
		repoDir,
		pythonBin,
		"../../scripts/worktree_cleanup_report.py",
		[]string{
			"PYTHONDONTWRITEBYTECODE=1",
			"LC_ALL=C",
			"LANG=C",
			"PYTHONCOERCECLOCALE=0",
			"PYTHONUTF8=0",
			"PYTHON=/缺失/python3",
		},
		"--json",
		"main",
	)
	if strings.Contains(output, "Traceback") {
		t.Fatalf("unexpected traceback: %s", output)
	}
	if !strings.Contains(output, "unable to execute worktree_audit.py via configured python command") {
		t.Fatalf("unexpected output: %s", output)
	}
	if !strings.Contains(output, "/缺失/python3") {
		t.Fatalf("unexpected output: %s", output)
	}
}

func runCmd(t *testing.T, dir, bin string, args ...string) {
	t.Helper()
	cmd := exec.Command(bin, args...)
	cmd.Dir = dir
	runCommandSuccess(t, cmd)
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

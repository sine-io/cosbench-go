#!/usr/bin/env python3

import json
import os
import subprocess
import sys
import tempfile
from datetime import datetime, timezone


REQUIRED_SECRETS = [
    "COSBENCH_SMOKE_ENDPOINT",
    "COSBENCH_SMOKE_ACCESS_KEY",
    "COSBENCH_SMOKE_SECRET_KEY",
]
WORKFLOW_NAMES = [
    "Smoke Local",
    "Smoke S3",
    "Smoke S3 Matrix",
    "Legacy Live Compare",
    "Legacy Live Compare Matrix",
    "Remote Smoke Local",
    "Remote Smoke Matrix",
    "Remote Smoke Recovery",
    "Remote Smoke Recovery Matrix",
]
DEFAULT_REPO = "sine-io/cosbench-go"
SMOKE_S3_WORKFLOW = "Smoke S3"
SMOKE_S3_MATRIX_WORKFLOW = "Smoke S3 Matrix"
LEGACY_LIVE_WORKFLOW = "Legacy Live Compare"
LEGACY_LIVE_MATRIX_WORKFLOW = "Legacy Live Compare Matrix"
REMOTE_SMOKE_LOCAL_WORKFLOW = "Remote Smoke Local"
REMOTE_SMOKE_MATRIX_WORKFLOW = "Remote Smoke Matrix"
REMOTE_SMOKE_RECOVERY_WORKFLOW = "Remote Smoke Recovery"
REMOTE_SMOKE_RECOVERY_MATRIX_WORKFLOW = "Remote Smoke Recovery Matrix"
LEGACY_STEP_NAME = "Run legacy live compare"


def generated_at():
    return datetime.now(timezone.utc).isoformat().replace("+00:00", "Z")


def parse_name_list(raw):
    items = []
    for part in raw.replace("\n", ",").split(","):
        item = part.strip()
        if item:
            items.append(item)
    return items


def run(*args):
    return subprocess.run(args, check=False, text=True, capture_output=True)


def resolve_repo():
    if "SMOKE_READY_REPO" in os.environ:
        return os.environ["SMOKE_READY_REPO"], None
    return DEFAULT_REPO, None


def load_repo_secret_names(repo):
    if "SMOKE_READY_MOCK_REPO_SECRETS" in os.environ:
        return set(parse_name_list(os.environ["SMOKE_READY_MOCK_REPO_SECRETS"])), True, ""
    proc = run("gh", "secret", "list", "--repo", repo)
    if proc.returncode != 0:
        error = (proc.stderr or proc.stdout).strip()
        return set(), False, error
    names = set()
    for line in proc.stdout.splitlines():
        line = line.strip()
        if not line:
            continue
        names.add(line.split("\t", 1)[0].split(None, 1)[0])
    return names, True, ""


def load_workflow_names(repo):
    if "SMOKE_READY_MOCK_WORKFLOWS" in os.environ:
        return set(parse_name_list(os.environ["SMOKE_READY_MOCK_WORKFLOWS"])), True, ""
    proc = run("gh", "workflow", "list", "--repo", repo)
    if proc.returncode != 0:
        error = (proc.stderr or proc.stdout).strip()
        return set(), False, error
    names = set()
    for line in proc.stdout.splitlines():
        line = line.strip()
        if not line:
            continue
        names.add(line.split("\t", 1)[0].split("  ", 1)[0].strip())
    return names, True, ""


def load_workflow_latest_runs(repo):
    if "SMOKE_READY_MOCK_WORKFLOW_RUNS_JSON" in os.environ:
        return json.loads(os.environ["SMOKE_READY_MOCK_WORKFLOW_RUNS_JSON"]), True, ""
    latest = {}
    for name in WORKFLOW_NAMES:
        proc = run(
            "gh",
            "run",
            "list",
            "--repo",
            repo,
            "--workflow",
            name,
            "--limit",
            "1",
            "--json",
            "databaseId,status,conclusion,createdAt,url",
        )
        if proc.returncode != 0:
            error = (proc.stderr or proc.stdout).strip()
            return {}, False, error
        rows = json.loads(proc.stdout or "[]")
        if rows:
            row = rows[0]
            latest[name] = {
                "database_id": row.get("databaseId"),
                "status": row.get("status", ""),
                "conclusion": row.get("conclusion", ""),
                "created_at": row.get("createdAt", ""),
                "url": row.get("url", ""),
            }
        else:
            latest[name] = None
    return latest, True, ""


def load_legacy_workflow_details(repo, workflow_latest):
    if "SMOKE_READY_MOCK_WORKFLOW_DETAILS_JSON" in os.environ:
        return json.loads(os.environ["SMOKE_READY_MOCK_WORKFLOW_DETAILS_JSON"]), True, ""
    details = {}
    latest = workflow_latest.get(LEGACY_LIVE_WORKFLOW) or {}
    run_id = latest.get("database_id")
    if run_id:
        with tempfile.TemporaryDirectory(prefix="smoke-ready-legacy-live-") as tmpdir:
            proc = run("gh", "run", "download", str(run_id), "--repo", repo, "-n", "legacy-live-compare-output", "-D", tmpdir)
            if proc.returncode == 0:
                result_path = os.path.join(tmpdir, "result.json")
                if not os.path.exists(result_path):
                    result_path = os.path.join(tmpdir, ".artifacts", "legacy-live-compare", "result.json")
                payload = {}
                if os.path.exists(result_path):
                    with open(result_path, "r", encoding="utf-8") as f:
                        payload["result"] = json.load(f)
                if payload:
                    details[LEGACY_LIVE_WORKFLOW] = payload
            elif latest.get("status") == "completed":
                return {}, False, (proc.stderr or proc.stdout).strip()
    if LEGACY_LIVE_WORKFLOW not in details:
        details[LEGACY_LIVE_WORKFLOW] = None

    latest = workflow_latest.get(LEGACY_LIVE_MATRIX_WORKFLOW) or {}
    run_id = latest.get("database_id")
    if run_id:
        with tempfile.TemporaryDirectory(prefix="smoke-ready-legacy-live-matrix-") as tmpdir:
            proc = run(
                "gh",
                "run",
                "download",
                str(run_id),
                "--repo",
                repo,
                "-n",
                "legacy-live-compare-matrix-aggregate",
                "-D",
                tmpdir,
            )
            if proc.returncode == 0:
                summary_path = os.path.join(tmpdir, "summary.json")
                if os.path.exists(summary_path):
                    with open(summary_path, "r", encoding="utf-8") as f:
                        details[LEGACY_LIVE_MATRIX_WORKFLOW] = json.load(f)
            elif latest.get("status") == "completed":
                return {}, False, (proc.stderr or proc.stdout).strip()
    if LEGACY_LIVE_MATRIX_WORKFLOW not in details:
        details[LEGACY_LIVE_MATRIX_WORKFLOW] = None

    for name in (LEGACY_LIVE_WORKFLOW, LEGACY_LIVE_MATRIX_WORKFLOW):
        if details.get(name) is None:
            latest = workflow_latest.get(name) or {}
            run_id = latest.get("database_id")
            if not run_id:
                continue
            proc = run("gh", "run", "view", str(run_id), "--repo", repo, "--json", "jobs")
            if proc.returncode != 0:
                error = (proc.stderr or proc.stdout).strip()
                return {}, False, error
            details[name] = json.loads(proc.stdout or "{}")
    return details, True, ""


def load_real_endpoint_details(repo, workflow_latest):
    if "SMOKE_READY_MOCK_WORKFLOW_DETAILS_JSON" in os.environ:
        details = json.loads(os.environ["SMOKE_READY_MOCK_WORKFLOW_DETAILS_JSON"])
        return {
            SMOKE_S3_WORKFLOW: details.get(SMOKE_S3_WORKFLOW),
            SMOKE_S3_MATRIX_WORKFLOW: details.get(SMOKE_S3_MATRIX_WORKFLOW),
        }, True, ""

    details = {}
    smoke_s3 = workflow_latest.get(SMOKE_S3_WORKFLOW) or {}
    smoke_s3_id = smoke_s3.get("database_id")
    if smoke_s3_id:
        with tempfile.TemporaryDirectory(prefix="smoke-ready-smoke-s3-") as tmpdir:
            proc = run("gh", "run", "download", str(smoke_s3_id), "--repo", repo, "-n", "smoke-s3-output", "-D", tmpdir)
            if proc.returncode == 0:
                output_path = os.path.join(tmpdir, "smoke-s3-output.txt")
                summary_path = os.path.join(tmpdir, "summary.json")
                if not os.path.exists(summary_path):
                    summary_path = os.path.join(tmpdir, ".artifacts", "smoke-s3-summary", "summary.json")
                payload = {}
                if os.path.exists(output_path):
                    with open(output_path, "r", encoding="utf-8") as f:
                        payload["output"] = f.read()
                if os.path.exists(summary_path):
                    with open(summary_path, "r", encoding="utf-8") as f:
                        payload["summary"] = json.load(f)
                if payload:
                    details[SMOKE_S3_WORKFLOW] = payload
            elif smoke_s3.get("status") == "completed":
                return {}, False, (proc.stderr or proc.stdout).strip()
    else:
        details[SMOKE_S3_WORKFLOW] = None

    smoke_s3_matrix = workflow_latest.get(SMOKE_S3_MATRIX_WORKFLOW) or {}
    smoke_s3_matrix_id = smoke_s3_matrix.get("database_id")
    if smoke_s3_matrix_id:
        with tempfile.TemporaryDirectory(prefix="smoke-ready-smoke-s3-matrix-") as tmpdir:
            proc = run(
                "gh",
                "run",
                "download",
                str(smoke_s3_matrix_id),
                "--repo",
                repo,
                "-n",
                "smoke-s3-matrix-aggregate",
                "-D",
                tmpdir,
            )
            if proc.returncode == 0:
                summary_path = os.path.join(tmpdir, "summary.json")
                if os.path.exists(summary_path):
                    with open(summary_path, "r", encoding="utf-8") as f:
                        details[SMOKE_S3_MATRIX_WORKFLOW] = json.load(f)
            elif smoke_s3_matrix.get("status") == "completed":
                return {}, False, (proc.stderr or proc.stdout).strip()
    else:
        details[SMOKE_S3_MATRIX_WORKFLOW] = None

    return details, True, ""


def load_remote_workflow_details(repo, workflow_latest):
    if "SMOKE_READY_MOCK_WORKFLOW_DETAILS_JSON" in os.environ:
        details = json.loads(os.environ["SMOKE_READY_MOCK_WORKFLOW_DETAILS_JSON"])
        return {
            REMOTE_SMOKE_LOCAL_WORKFLOW: details.get(REMOTE_SMOKE_LOCAL_WORKFLOW),
            REMOTE_SMOKE_MATRIX_WORKFLOW: details.get(REMOTE_SMOKE_MATRIX_WORKFLOW),
            REMOTE_SMOKE_RECOVERY_WORKFLOW: details.get(REMOTE_SMOKE_RECOVERY_WORKFLOW),
            REMOTE_SMOKE_RECOVERY_MATRIX_WORKFLOW: details.get(REMOTE_SMOKE_RECOVERY_MATRIX_WORKFLOW),
        }, True, ""

    details = {}

    single_targets = [
        (REMOTE_SMOKE_LOCAL_WORKFLOW, "remote-smoke-output"),
        (REMOTE_SMOKE_RECOVERY_WORKFLOW, "remote-smoke-recovery-summary"),
    ]
    for workflow_name, artifact_name in single_targets:
        latest = workflow_latest.get(workflow_name) or {}
        run_id = latest.get("database_id")
        if not run_id:
            details[workflow_name] = None
            continue
        with tempfile.TemporaryDirectory(prefix="smoke-ready-remote-smoke-") as tmpdir:
            proc = run("gh", "run", "download", str(run_id), "--repo", repo, "-n", artifact_name, "-D", tmpdir)
            if proc.returncode == 0:
                summary_path = os.path.join(tmpdir, "summary.json")
                if not os.path.exists(summary_path):
                    summary_path = os.path.join(tmpdir, "remote-smoke", "summary.json")
                if os.path.exists(summary_path):
                    with open(summary_path, "r", encoding="utf-8") as f:
                        details[workflow_name] = json.load(f)
            elif latest.get("status") == "completed":
                return {}, False, (proc.stderr or proc.stdout).strip()
        if workflow_name not in details:
            details[workflow_name] = None

    matrix_targets = [
        (REMOTE_SMOKE_MATRIX_WORKFLOW, "remote-smoke-matrix-aggregate"),
        (REMOTE_SMOKE_RECOVERY_MATRIX_WORKFLOW, "remote-smoke-recovery-matrix-aggregate"),
    ]
    for workflow_name, artifact_name in matrix_targets:
        latest = workflow_latest.get(workflow_name) or {}
        run_id = latest.get("database_id")
        if not run_id:
            details[workflow_name] = None
            continue
        with tempfile.TemporaryDirectory(prefix="smoke-ready-remote-smoke-matrix-") as tmpdir:
            proc = run("gh", "run", "download", str(run_id), "--repo", repo, "-n", artifact_name, "-D", tmpdir)
            if proc.returncode == 0:
                summary_path = os.path.join(tmpdir, "summary.json")
                if os.path.exists(summary_path):
                    with open(summary_path, "r", encoding="utf-8") as f:
                        details[workflow_name] = json.load(f)
            elif latest.get("status") == "completed":
                return {}, False, (proc.stderr or proc.stdout).strip()
        if workflow_name not in details:
            details[workflow_name] = None

    return details, True, ""


def step_state(job, step_name):
    for step in job.get("steps", []):
        if step.get("name") == step_name:
            return {
                "status": step.get("status", ""),
                "conclusion": step.get("conclusion", ""),
            }
    return {"status": "", "conclusion": ""}


def classify_single_legacy_result(latest, detail):
    if not latest:
        return "none"
    if latest.get("status") != "completed":
        return "pending"
    if not detail:
        return "failed"
    result = (detail.get("result") or {}).get("result", "")
    if result in {"executed", "skipped", "failed"}:
        return result
    for job in detail.get("jobs", []):
        step = step_state(job, LEGACY_STEP_NAME)
        if step["status"] or step["conclusion"]:
            if step["status"] != "completed":
                return "pending"
            if step["conclusion"] == "success":
                return "executed"
            if step["conclusion"] == "skipped":
                return "skipped"
            return "failed"
    return "failed"


def classify_matrix_legacy_result(latest, detail):
    if not latest:
        return "none"
    if latest.get("status") != "completed":
        return "pending"
    if not detail:
        return "failed"
    if "rows" in detail:
        row_results = []
        for row in detail.get("rows", []):
            row_status = row.get("status", "")
            if row_status in {"executed", "skipped", "failed"}:
                row_results.append(row_status)
            else:
                row_results.append("failed")
        if not row_results:
            return "failed"
        unique = set(row_results)
        if unique == {"executed"}:
            return "executed"
        if unique == {"skipped"}:
            return "skipped"
        if unique == {"failed"}:
            return "failed"
        if "pending" in unique:
            return "pending"
        return "partial"
    row_results = []
    for job in detail.get("jobs", []):
        if not str(job.get("name", "")).startswith("legacy-live-compare-matrix ("):
            continue
        step = step_state(job, LEGACY_STEP_NAME)
        if step["status"] != "completed":
            row_results.append("pending")
            continue
        if step["conclusion"] == "success":
            row_results.append("executed")
        elif step["conclusion"] == "skipped":
            row_results.append("skipped")
        else:
            row_results.append("failed")
    if not row_results:
        return "failed"
    unique = set(row_results)
    if unique == {"executed"}:
        return "executed"
    if unique == {"skipped"}:
        return "skipped"
    if unique == {"failed"}:
        return "failed"
    if "pending" in unique:
        return "pending"
    return "partial"


def smoke_output_result(latest, output):
    if not latest:
        return "none"
    if latest.get("status") != "completed":
        return "pending"
    if not output:
        return "failed"
    if "--- SKIP:" in output and "PASS" in output and "t.Fatalf" not in output:
        return "skipped"
    if "PASS" in output and "--- SKIP:" not in output:
        return "executed"
    return "failed"


def smoke_summary_result(latest, detail):
    if not latest:
        return "none"
    if latest.get("status") != "completed":
        return "pending"
    if not detail:
        return "failed"
    summary = detail.get("summary") or {}
    result = summary.get("result", "")
    if result in {"executed", "skipped", "failed"}:
        return result
    return smoke_output_result(latest, detail.get("output", ""))


def smoke_matrix_result(latest, detail):
    if not latest:
        return "none"
    if latest.get("status") != "completed":
        return "pending"
    if not detail:
        return "failed"
    row_results = []
    for row in detail.get("rows", []):
        output = row.get("output", "")
        row_status = row.get("status", "")
        if row_status in {"executed", "skipped", "failed"}:
            row_results.append(row_status)
            continue
        if row_status != "present":
            row_results.append("failed")
            continue
        row_results.append(smoke_output_result({"status": "completed"}, output))
    if not row_results:
        return "failed"
    unique = set(row_results)
    if unique == {"executed"}:
        return "executed"
    if unique == {"skipped"}:
        return "skipped"
    if unique == {"failed"}:
        return "failed"
    if "pending" in unique:
        return "pending"
    return "partial"


def remote_single_result(latest, detail):
    if not latest:
        return "none"
    if latest.get("status") != "completed":
        return "pending"
    if not detail:
        return "failed"
    if detail.get("overall") == "pass":
        return "executed"
    if detail.get("overall") == "fail":
        return "failed"
    return "failed"


def remote_matrix_result(latest, detail):
    if not latest:
        return "none"
    if latest.get("status") != "completed":
        return "pending"
    if not detail:
        return "failed"
    overall = detail.get("overall", "")
    if overall == "pass":
        return "executed"
    if overall == "partial":
        return "partial"
    if overall == "fail":
        return "failed"
    return "failed"


def pick_latest_result(workflow_latest, names):
    latest_name = None
    latest_created_at = ""
    for name in names:
        row = workflow_latest.get(name) or {}
        created_at = row.get("created_at", "")
        if created_at >= latest_created_at:
            latest_name = name
            latest_created_at = created_at
    return latest_name


def latest_url(workflow_latest, workflow_name):
    if not workflow_name:
        return ""
    return (workflow_latest.get(workflow_name) or {}).get("url", "")


def latest_created_at(workflow_latest, workflow_name):
    if not workflow_name:
        return ""
    return (workflow_latest.get(workflow_name) or {}).get("created_at", "")


def build_payload():
    repo, repo_error = resolve_repo()
    local_env = {name: bool(os.getenv(name, "").strip()) for name in REQUIRED_SECRETS}
    local_ready = all(local_env.values())

    repo_secret_names, repo_secrets_accessible, repo_secrets_error = load_repo_secret_names(repo)
    workflow_names, workflows_accessible, workflows_error = load_workflow_names(repo)
    workflow_latest, workflow_runs_accessible, workflow_runs_error = load_workflow_latest_runs(repo)
    legacy_details, legacy_details_accessible, legacy_details_error = load_legacy_workflow_details(repo, workflow_latest)
    real_endpoint_details, real_endpoint_details_accessible, real_endpoint_details_error = load_real_endpoint_details(repo, workflow_latest)
    remote_details, remote_details_accessible, remote_details_error = load_remote_workflow_details(repo, workflow_latest)

    repo_secret_presence = {name: name in repo_secret_names for name in REQUIRED_SECRETS}
    workflow_presence = {name: name in workflow_names for name in WORKFLOW_NAMES}
    local_workflow_ready = workflows_accessible and workflow_presence["Smoke Local"]
    real_endpoint_ready = workflows_accessible and workflow_presence["Smoke S3"]
    real_endpoint_matrix_ready = workflows_accessible and workflow_presence["Smoke S3 Matrix"]
    legacy_live_ready = workflows_accessible and workflow_presence["Legacy Live Compare"]
    legacy_live_matrix_ready = workflows_accessible and workflow_presence["Legacy Live Compare Matrix"]
    remote_happy_ready = workflows_accessible and workflow_presence["Remote Smoke Local"] and workflow_presence["Remote Smoke Matrix"]
    remote_recovery_ready = workflows_accessible and workflow_presence["Remote Smoke Recovery"] and workflow_presence["Remote Smoke Recovery Matrix"]
    real_endpoint_latest_result = smoke_summary_result(workflow_latest.get(SMOKE_S3_WORKFLOW), real_endpoint_details.get(SMOKE_S3_WORKFLOW))
    real_endpoint_matrix_latest_result = smoke_matrix_result(workflow_latest.get(SMOKE_S3_MATRIX_WORKFLOW), real_endpoint_details.get(SMOKE_S3_MATRIX_WORKFLOW))
    legacy_live_latest_result = classify_single_legacy_result(workflow_latest.get(LEGACY_LIVE_WORKFLOW), legacy_details.get(LEGACY_LIVE_WORKFLOW))
    legacy_live_matrix_latest_result = classify_matrix_legacy_result(workflow_latest.get(LEGACY_LIVE_MATRIX_WORKFLOW), legacy_details.get(LEGACY_LIVE_MATRIX_WORKFLOW))
    remote_happy_latest_name = pick_latest_result(workflow_latest, [REMOTE_SMOKE_LOCAL_WORKFLOW, REMOTE_SMOKE_MATRIX_WORKFLOW])
    remote_recovery_latest_name = pick_latest_result(workflow_latest, [REMOTE_SMOKE_RECOVERY_WORKFLOW, REMOTE_SMOKE_RECOVERY_MATRIX_WORKFLOW])
    remote_happy_latest_result = (
        remote_single_result(workflow_latest.get(REMOTE_SMOKE_LOCAL_WORKFLOW), remote_details.get(REMOTE_SMOKE_LOCAL_WORKFLOW))
        if remote_happy_latest_name == REMOTE_SMOKE_LOCAL_WORKFLOW
        else remote_matrix_result(workflow_latest.get(REMOTE_SMOKE_MATRIX_WORKFLOW), remote_details.get(REMOTE_SMOKE_MATRIX_WORKFLOW))
    )
    remote_recovery_latest_result = (
        remote_single_result(workflow_latest.get(REMOTE_SMOKE_RECOVERY_WORKFLOW), remote_details.get(REMOTE_SMOKE_RECOVERY_WORKFLOW))
        if remote_recovery_latest_name == REMOTE_SMOKE_RECOVERY_WORKFLOW
        else remote_matrix_result(workflow_latest.get(REMOTE_SMOKE_RECOVERY_MATRIX_WORKFLOW), remote_details.get(REMOTE_SMOKE_RECOVERY_MATRIX_WORKFLOW))
    )
    real_endpoint_latest_success = real_endpoint_latest_result == "executed"
    real_endpoint_matrix_latest_success = real_endpoint_matrix_latest_result == "executed"
    legacy_live_latest_success = legacy_live_latest_result == "executed"
    legacy_live_matrix_latest_success = legacy_live_matrix_latest_result == "executed"
    remote_happy_latest_success = remote_happy_latest_result == "executed"
    remote_recovery_latest_success = remote_recovery_latest_result == "executed"
    ready = local_ready or local_workflow_ready

    blockers = []
    if repo_error:
        blockers.append(f"unable to resolve repo: {repo_error}")
    if not ready:
        missing_local = [name for name, present in local_env.items() if not present]
        if missing_local:
            blockers.append(f"missing local smoke env: {', '.join(missing_local)}")
        if not workflows_accessible and workflows_error:
            blockers.append(f"unable to query workflows: {workflows_error}")
        elif not workflow_presence["Smoke Local"]:
            blockers.append("required workflow missing: Smoke Local")
    if workflows_accessible and not workflow_runs_accessible and workflow_runs_error:
        blockers.append(f"unable to query workflow runs: {workflow_runs_error}")
    if workflow_runs_accessible and not real_endpoint_details_accessible and real_endpoint_details_error:
        blockers.append(f"unable to query real-endpoint workflow details: {real_endpoint_details_error}")
    if workflow_runs_accessible and not remote_details_accessible and remote_details_error:
        blockers.append(f"unable to query remote workflow details: {remote_details_error}")
    if workflow_runs_accessible and not legacy_details_accessible and legacy_details_error:
        blockers.append(f"unable to query legacy workflow details: {legacy_details_error}")

    return {
        "generated_at": generated_at(),
        "repo": repo,
        "required": REQUIRED_SECRETS,
        "local_env": local_env,
        "repo_secrets": {
            "accessible": repo_secrets_accessible,
            "error": repo_secrets_error,
            "present": repo_secret_presence,
        },
        "workflows": {
            "accessible": workflows_accessible,
            "error": workflows_error,
            "present": workflow_presence,
            "latest_accessible": workflow_runs_accessible,
            "latest_error": workflow_runs_error,
            "latest": workflow_latest,
        },
        "summary": {
            "local_env_ready": local_ready,
            "local_workflow_ready": local_workflow_ready,
            "real_endpoint_ready": real_endpoint_ready,
            "real_endpoint_matrix_ready": real_endpoint_matrix_ready,
            "legacy_live_ready": legacy_live_ready,
            "legacy_live_matrix_ready": legacy_live_matrix_ready,
            "remote_happy_ready": remote_happy_ready,
            "remote_recovery_ready": remote_recovery_ready,
            "real_endpoint_latest_success": real_endpoint_latest_success,
            "real_endpoint_matrix_latest_success": real_endpoint_matrix_latest_success,
            "real_endpoint_latest_result": real_endpoint_latest_result,
            "real_endpoint_matrix_latest_result": real_endpoint_matrix_latest_result,
            "real_endpoint_latest_source": SMOKE_S3_WORKFLOW,
            "real_endpoint_matrix_latest_source": SMOKE_S3_MATRIX_WORKFLOW,
            "real_endpoint_latest_url": latest_url(workflow_latest, SMOKE_S3_WORKFLOW),
            "real_endpoint_matrix_latest_url": latest_url(workflow_latest, SMOKE_S3_MATRIX_WORKFLOW),
            "real_endpoint_latest_created_at": latest_created_at(workflow_latest, SMOKE_S3_WORKFLOW),
            "real_endpoint_matrix_latest_created_at": latest_created_at(workflow_latest, SMOKE_S3_MATRIX_WORKFLOW),
            "legacy_live_latest_success": legacy_live_latest_success,
            "legacy_live_matrix_latest_success": legacy_live_matrix_latest_success,
            "legacy_live_latest_result": legacy_live_latest_result,
            "legacy_live_matrix_latest_result": legacy_live_matrix_latest_result,
            "legacy_live_latest_source": LEGACY_LIVE_WORKFLOW,
            "legacy_live_matrix_latest_source": LEGACY_LIVE_MATRIX_WORKFLOW,
            "legacy_live_latest_url": latest_url(workflow_latest, LEGACY_LIVE_WORKFLOW),
            "legacy_live_matrix_latest_url": latest_url(workflow_latest, LEGACY_LIVE_MATRIX_WORKFLOW),
            "legacy_live_latest_created_at": latest_created_at(workflow_latest, LEGACY_LIVE_WORKFLOW),
            "legacy_live_matrix_latest_created_at": latest_created_at(workflow_latest, LEGACY_LIVE_MATRIX_WORKFLOW),
            "remote_happy_latest_success": remote_happy_latest_success,
            "remote_recovery_latest_success": remote_recovery_latest_success,
            "remote_happy_latest_result": remote_happy_latest_result,
            "remote_recovery_latest_result": remote_recovery_latest_result,
            "remote_happy_latest_source": remote_happy_latest_name or "none",
            "remote_recovery_latest_source": remote_recovery_latest_name or "none",
            "remote_happy_latest_url": latest_url(workflow_latest, remote_happy_latest_name),
            "remote_recovery_latest_url": latest_url(workflow_latest, remote_recovery_latest_name),
            "remote_happy_latest_created_at": latest_created_at(workflow_latest, remote_happy_latest_name),
            "remote_recovery_latest_created_at": latest_created_at(workflow_latest, remote_recovery_latest_name),
            "ready": ready,
        },
        "blockers": blockers,
    }


def yes_no(value):
    return "yes" if value else "no"


def set_unset(value):
    return "set" if value else "unset"


def available_missing(value):
    return "available" if value else "missing"


def latest_display(value):
    if not value:
        return "none"
    status = value.get("status", "") or "unknown"
    conclusion = value.get("conclusion", "") or "unknown"
    created_at = value.get("created_at", "") or "unknown"
    return f"{status}/{conclusion} @ {created_at}"


def print_text(payload):
    print("# Smoke Ready")
    print()
    print(f"Repository: `{payload['repo']}`")
    print(f"Generated at: `{payload['generated_at']}`")
    print()
    print("## Local Env")
    print()
    for name in REQUIRED_SECRETS:
        print(f"- {name}: `{set_unset(payload['local_env'][name])}`")
    print()
    print("## Repository Secrets")
    print()
    if not payload["repo_secrets"]["accessible"] and payload["repo_secrets"]["error"]:
        print(f"- query: `error`")
        print(f"- error: `{payload['repo_secrets']['error']}`")
    else:
        for name in REQUIRED_SECRETS:
            print(f"- {name}: `{set_unset(payload['repo_secrets']['present'][name])}`")
    print()
    print("## Workflows")
    print()
    if not payload["workflows"]["accessible"] and payload["workflows"]["error"]:
        print(f"- query: `error`")
        print(f"- error: `{payload['workflows']['error']}`")
    else:
        for name in WORKFLOW_NAMES:
            print(f"- {name}: `{available_missing(payload['workflows']['present'][name])}`")
    print()
    print("## Latest Runs")
    print()
    if not payload["workflows"]["latest_accessible"] and payload["workflows"]["latest_error"]:
        print(f"- query: `error`")
        print(f"- error: `{payload['workflows']['latest_error']}`")
    else:
        for name in WORKFLOW_NAMES:
            print(f"- {name}: `{latest_display(payload['workflows']['latest'].get(name))}`")
    print()
    print("## Summary")
    print()
    print(f"- Local Env Ready: `{yes_no(payload['summary']['local_env_ready'])}`")
    print(f"- Local Workflow Ready: `{yes_no(payload['summary']['local_workflow_ready'])}`")
    print(f"- Real Endpoint Ready: `{yes_no(payload['summary']['real_endpoint_ready'])}`")
    print(f"- Real Endpoint Matrix Ready: `{yes_no(payload['summary']['real_endpoint_matrix_ready'])}`")
    print(f"- Legacy Live Ready: `{yes_no(payload['summary']['legacy_live_ready'])}`")
    print(f"- Legacy Live Matrix Ready: `{yes_no(payload['summary']['legacy_live_matrix_ready'])}`")
    print(f"- Remote Happy Ready: `{yes_no(payload['summary']['remote_happy_ready'])}`")
    print(f"- Remote Recovery Ready: `{yes_no(payload['summary']['remote_recovery_ready'])}`")
    print(f"- Real Endpoint Latest Success: `{yes_no(payload['summary']['real_endpoint_latest_success'])}`")
    print(f"- Real Endpoint Matrix Latest Success: `{yes_no(payload['summary']['real_endpoint_matrix_latest_success'])}`")
    print(f"- Real Endpoint Latest Result: `{payload['summary']['real_endpoint_latest_result']}`")
    print(f"- Real Endpoint Matrix Latest Result: `{payload['summary']['real_endpoint_matrix_latest_result']}`")
    print(f"- Real Endpoint Latest Source: `{payload['summary']['real_endpoint_latest_source']}`")
    print(f"- Real Endpoint Matrix Latest Source: `{payload['summary']['real_endpoint_matrix_latest_source']}`")
    print(f"- Real Endpoint Latest URL: `{payload['summary']['real_endpoint_latest_url']}`")
    print(f"- Real Endpoint Matrix Latest URL: `{payload['summary']['real_endpoint_matrix_latest_url']}`")
    print(f"- Real Endpoint Latest Created At: `{payload['summary']['real_endpoint_latest_created_at']}`")
    print(f"- Real Endpoint Matrix Latest Created At: `{payload['summary']['real_endpoint_matrix_latest_created_at']}`")
    print(f"- Legacy Live Latest Success: `{yes_no(payload['summary']['legacy_live_latest_success'])}`")
    print(f"- Legacy Live Matrix Latest Success: `{yes_no(payload['summary']['legacy_live_matrix_latest_success'])}`")
    print(f"- Legacy Live Latest Result: `{payload['summary']['legacy_live_latest_result']}`")
    print(f"- Legacy Live Matrix Latest Result: `{payload['summary']['legacy_live_matrix_latest_result']}`")
    print(f"- Legacy Live Latest Source: `{payload['summary']['legacy_live_latest_source']}`")
    print(f"- Legacy Live Matrix Latest Source: `{payload['summary']['legacy_live_matrix_latest_source']}`")
    print(f"- Legacy Live Latest URL: `{payload['summary']['legacy_live_latest_url']}`")
    print(f"- Legacy Live Matrix Latest URL: `{payload['summary']['legacy_live_matrix_latest_url']}`")
    print(f"- Legacy Live Latest Created At: `{payload['summary']['legacy_live_latest_created_at']}`")
    print(f"- Legacy Live Matrix Latest Created At: `{payload['summary']['legacy_live_matrix_latest_created_at']}`")
    print(f"- Remote Happy Latest Success: `{yes_no(payload['summary']['remote_happy_latest_success'])}`")
    print(f"- Remote Recovery Latest Success: `{yes_no(payload['summary']['remote_recovery_latest_success'])}`")
    print(f"- Remote Happy Latest Result: `{payload['summary']['remote_happy_latest_result']}`")
    print(f"- Remote Recovery Latest Result: `{payload['summary']['remote_recovery_latest_result']}`")
    print(f"- Remote Happy Latest Source: `{payload['summary']['remote_happy_latest_source']}`")
    print(f"- Remote Recovery Latest Source: `{payload['summary']['remote_recovery_latest_source']}`")
    print(f"- Remote Happy Latest URL: `{payload['summary']['remote_happy_latest_url']}`")
    print(f"- Remote Recovery Latest URL: `{payload['summary']['remote_recovery_latest_url']}`")
    print(f"- Remote Happy Latest Created At: `{payload['summary']['remote_happy_latest_created_at']}`")
    print(f"- Remote Recovery Latest Created At: `{payload['summary']['remote_recovery_latest_created_at']}`")
    print(f"- Overall ready: `{yes_no(payload['summary']['ready'])}`")
    print()
    print("## Blockers")
    print()
    if payload["blockers"]:
        for blocker in payload["blockers"]:
            print(f"- {blocker}")
    else:
        print("- none")


def main():
    json_mode = "--json" in sys.argv[1:]
    payload = build_payload()
    if json_mode:
        print(json.dumps(payload, indent=2, ensure_ascii=False))
        return
    print_text(payload)


if __name__ == "__main__":
    main()

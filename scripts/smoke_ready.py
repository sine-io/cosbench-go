#!/usr/bin/env python3

import json
import os
import subprocess
import sys
from datetime import datetime, timezone


REQUIRED_SECRETS = [
    "COSBENCH_SMOKE_ENDPOINT",
    "COSBENCH_SMOKE_ACCESS_KEY",
    "COSBENCH_SMOKE_SECRET_KEY",
]
WORKFLOW_NAMES = [
    "Smoke Local",
    "Remote Smoke Local",
    "Remote Smoke Matrix",
    "Remote Smoke Recovery",
    "Remote Smoke Recovery Matrix",
]
DEFAULT_REPO = "sine-io/cosbench-go"


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
            "status,conclusion,createdAt,url",
        )
        if proc.returncode != 0:
            error = (proc.stderr or proc.stdout).strip()
            return {}, False, error
        rows = json.loads(proc.stdout or "[]")
        if rows:
            row = rows[0]
            latest[name] = {
                "status": row.get("status", ""),
                "conclusion": row.get("conclusion", ""),
                "created_at": row.get("createdAt", ""),
                "url": row.get("url", ""),
            }
        else:
            latest[name] = None
    return latest, True, ""


def build_payload():
    repo, repo_error = resolve_repo()
    local_env = {name: bool(os.getenv(name, "").strip()) for name in REQUIRED_SECRETS}
    local_ready = all(local_env.values())

    repo_secret_names, repo_secrets_accessible, repo_secrets_error = load_repo_secret_names(repo)
    workflow_names, workflows_accessible, workflows_error = load_workflow_names(repo)
    workflow_latest, workflow_runs_accessible, workflow_runs_error = load_workflow_latest_runs(repo)

    repo_secret_presence = {name: name in repo_secret_names for name in REQUIRED_SECRETS}
    workflow_presence = {name: name in workflow_names for name in WORKFLOW_NAMES}
    local_workflow_ready = workflows_accessible and workflow_presence["Smoke Local"]
    remote_happy_ready = workflows_accessible and workflow_presence["Remote Smoke Local"] and workflow_presence["Remote Smoke Matrix"]
    remote_recovery_ready = workflows_accessible and workflow_presence["Remote Smoke Recovery"] and workflow_presence["Remote Smoke Recovery Matrix"]
    remote_happy_latest_success = any(
        (workflow_latest.get(name) or {}).get("conclusion") == "success"
        for name in ("Remote Smoke Local", "Remote Smoke Matrix")
    )
    remote_recovery_latest_success = any(
        (workflow_latest.get(name) or {}).get("conclusion") == "success"
        for name in ("Remote Smoke Recovery", "Remote Smoke Recovery Matrix")
    )
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
            "remote_happy_ready": remote_happy_ready,
            "remote_recovery_ready": remote_recovery_ready,
            "remote_happy_latest_success": remote_happy_latest_success,
            "remote_recovery_latest_success": remote_recovery_latest_success,
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
    print(f"- Remote Happy Ready: `{yes_no(payload['summary']['remote_happy_ready'])}`")
    print(f"- Remote Recovery Ready: `{yes_no(payload['summary']['remote_recovery_ready'])}`")
    print(f"- Remote Happy Latest Success: `{yes_no(payload['summary']['remote_happy_latest_success'])}`")
    print(f"- Remote Recovery Latest Success: `{yes_no(payload['summary']['remote_recovery_latest_success'])}`")
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

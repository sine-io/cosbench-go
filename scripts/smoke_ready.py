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
SMOKE_WORKFLOW_NAME = "Smoke S3"
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


def build_payload():
    repo, repo_error = resolve_repo()
    local_env = {name: bool(os.getenv(name, "").strip()) for name in REQUIRED_SECRETS}
    local_ready = all(local_env.values())

    repo_secret_names, repo_secrets_accessible, repo_secrets_error = load_repo_secret_names(repo)
    workflow_names, workflows_accessible, workflows_error = load_workflow_names(repo)

    repo_secret_presence = {name: name in repo_secret_names for name in REQUIRED_SECRETS}
    repo_secrets_ready = repo_secrets_accessible and all(repo_secret_presence.values())
    workflow_presence = {SMOKE_WORKFLOW_NAME: SMOKE_WORKFLOW_NAME in workflow_names}
    workflow_ready = workflows_accessible and workflow_presence[SMOKE_WORKFLOW_NAME] and repo_secrets_ready
    ready = local_ready or workflow_ready

    blockers = []
    if repo_error:
        blockers.append(f"unable to resolve repo: {repo_error}")
    if not ready:
        missing_local = [name for name, present in local_env.items() if not present]
        if missing_local:
            blockers.append(f"missing local smoke env: {', '.join(missing_local)}")
        if not repo_secrets_accessible and repo_secrets_error:
            blockers.append(f"unable to query repo secrets: {repo_secrets_error}")
        elif not repo_secrets_ready:
            missing_repo = [name for name, present in repo_secret_presence.items() if not present]
            if missing_repo:
                blockers.append(f"missing repo smoke secrets: {', '.join(missing_repo)}")
        if not workflows_accessible and workflows_error:
            blockers.append(f"unable to query workflows: {workflows_error}")
        elif not workflow_presence[SMOKE_WORKFLOW_NAME]:
            blockers.append(f"required workflow missing: {SMOKE_WORKFLOW_NAME}")

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
        },
        "summary": {
            "local_ready": local_ready,
            "workflow_ready": workflow_ready,
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
        print(f"- {SMOKE_WORKFLOW_NAME}: `{available_missing(payload['workflows']['present'][SMOKE_WORKFLOW_NAME])}`")
    print()
    print("## Summary")
    print()
    print(f"- Local ready: `{yes_no(payload['summary']['local_ready'])}`")
    print(f"- Workflow ready: `{yes_no(payload['summary']['workflow_ready'])}`")
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

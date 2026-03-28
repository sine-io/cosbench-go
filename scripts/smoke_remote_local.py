#!/usr/bin/env python3

from datetime import datetime
import json
import os
import signal
import socket
import subprocess
import sys
import tempfile
import time
import urllib.error
import urllib.parse
import urllib.request
import uuid
from pathlib import Path


ROOT = Path(__file__).resolve().parents[1]
ARTIFACT_DIR = ROOT / ".artifacts" / "remote-smoke"
S3_FIXTURE = ROOT / "testdata" / "workloads" / "remote-smoke-s3-two-driver.xml"
SIO_FIXTURE = ROOT / "testdata" / "workloads" / "remote-smoke-sio-two-driver.xml"
S3_MULTISTAGE_FIXTURE = ROOT / "testdata" / "workloads" / "remote-smoke-s3-multistage-two-driver.xml"
SIO_MULTISTAGE_FIXTURE = ROOT / "testdata" / "workloads" / "remote-smoke-sio-multistage-two-driver.xml"
DEFAULT_MINIO = ROOT / ".artifacts" / "minio" / "bin" / "minio"
DEFAULT_MINIO_DATA = ROOT / ".artifacts" / "remote-smoke" / "minio-data"
DEFAULT_MINIO_URL = "https://dl.min.io/server/minio/release/linux-amd64/minio"
DEFAULT_HOST = "127.0.0.1"
DEFAULT_GO = os.environ.get("GO", "/snap/bin/go")
DEFAULT_MINIO_ACCESS_KEY = "minioadmin"
DEFAULT_MINIO_SECRET_KEY = "minioadmin"


def build_summary(
    *,
    backend,
    scenario,
    controller_url,
    driver_urls,
    job_id,
    job_status,
    drivers_seen,
    units_claimed,
    drivers_participated,
    operation_count,
    byte_count,
    stage_names,
    stages_seen,
    checks,
):
    overall = "pass" if all(value == "pass" for value in checks.values()) else "fail"
    return {
        "backend": backend,
        "scenario": scenario,
        "controller_url": controller_url,
        "driver_urls": driver_urls,
        "job_id": job_id,
        "job_status": job_status,
        "drivers_seen": drivers_seen,
        "units_claimed": units_claimed,
        "drivers_participated": drivers_participated,
        "operation_count": operation_count,
        "byte_count": byte_count,
        "stage_names": stage_names,
        "stages_seen": stages_seen,
        "checks": checks,
        "overall": overall,
    }


def build_failure_summary(failed_at, error):
    return {
        "failed_at": failed_at,
        "error": error,
        "overall": "fail",
    }


def render_summary_md(summary):
    lines = ["# Remote Smoke", ""]
    for key in [
        "backend",
        "scenario",
        "controller_url",
        "driver_urls",
        "job_id",
        "job_status",
        "drivers_seen",
        "units_claimed",
        "drivers_participated",
        "operation_count",
        "byte_count",
        "stage_names",
        "stages_seen",
        "overall",
    ]:
        if key in summary:
            lines.append(f"- `{key}`: `{summary[key]}`")
    if "checks" in summary:
        lines.append("")
        lines.append("## Checks")
        for name, status in summary["checks"].items():
            lines.append(f"- `{name}`: `{status}`")
    if "error" in summary:
        lines.append("")
        lines.append("## Error")
        lines.append(summary["error"])
    return "\n".join(lines) + "\n"


def write_summary(summary):
    ARTIFACT_DIR.mkdir(parents=True, exist_ok=True)
    json_path = ARTIFACT_DIR / "summary.json"
    md_path = ARTIFACT_DIR / "summary.md"
    json_path.write_text(json.dumps(summary, indent=2) + "\n", encoding="utf-8")
    md_path.write_text(render_summary_md(summary), encoding="utf-8")


def run_mock():
    mode = os.environ.get("SMOKE_REMOTE_LOCAL_MOCK", "").strip()
    if not mode:
        return None
    if mode == "success":
        summary = build_summary(
            backend="s3",
            scenario="single",
            controller_url="http://127.0.0.1:19088",
            driver_urls=["http://127.0.0.1:18081", "http://127.0.0.1:18082"],
            job_id="job-1",
            job_status="succeeded",
            drivers_seen=2,
            units_claimed=2,
            drivers_participated=2,
            operation_count=2,
            byte_count=2000,
            stage_names=["main"],
            stages_seen=1,
            checks={
                "process_ready": "pass",
                "drivers_healthy": "pass",
                "units_distributed": "pass",
                "job_succeeded": "pass",
                "visibility": "pass",
            },
        )
        write_summary(summary)
        sys.stdout.write(render_summary_md(summary))
        return 0
    summary = build_failure_summary("controller", "mocked remote smoke failure")
    write_summary(summary)
    sys.stdout.write(render_summary_md(summary))
    return 1


def find_free_port():
    with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as sock:
        sock.bind((DEFAULT_HOST, 0))
        return sock.getsockname()[1]


def ensure_minio():
    if DEFAULT_MINIO.exists() and DEFAULT_MINIO.stat().st_size > 0:
        DEFAULT_MINIO.chmod(0o755)
        return str(DEFAULT_MINIO)
    DEFAULT_MINIO.parent.mkdir(parents=True, exist_ok=True)
    with urllib.request.urlopen(DEFAULT_MINIO_URL) as response, DEFAULT_MINIO.open("wb") as target:
        target.write(response.read())
    DEFAULT_MINIO.chmod(0o755)
    return str(DEFAULT_MINIO)


def wait_for_socket(host, port, timeout_seconds=20):
    deadline = time.time() + timeout_seconds
    while time.time() < deadline:
        with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as sock:
            sock.settimeout(0.5)
            try:
                sock.connect((host, port))
                return
            except OSError:
                time.sleep(0.1)
    raise RuntimeError(f"socket did not become ready on {host}:{port}")


def fixture_for_backend(backend):
    backend = backend.strip().lower()
    if backend == "s3":
        return S3_FIXTURE
    if backend == "sio":
        return SIO_FIXTURE
    raise ValueError(f"unsupported backend: {backend}")


def normalize_scenario(scenario):
    normalized = scenario.strip().lower()
    if not normalized:
        return "single"
    if normalized in {"single", "multistage"}:
        return normalized
    raise ValueError(f"unsupported scenario: {scenario}")


def fixture_for_selection(backend, scenario):
    backend = backend.strip().lower()
    scenario = normalize_scenario(scenario)
    if scenario == "single":
        return fixture_for_backend(backend)
    if backend == "s3":
        return S3_MULTISTAGE_FIXTURE
    if backend == "sio":
        return SIO_MULTISTAGE_FIXTURE
    raise ValueError(f"unsupported remote smoke scenario: backend={backend} scenario={scenario}")


def parse_time(value):
    if not value:
        return None
    return datetime.fromisoformat(value.replace("Z", "+00:00"))


def multistage_barrier_holds(stages):
    if len(stages) < 2:
        return False
    for index in range(len(stages) - 1):
        current_stage = stages[index]
        next_stage = stages[index + 1]
        current_finished = parse_time(current_stage.get("finished_at"))
        next_started = parse_time(next_stage.get("started_at"))
        if current_finished is None or next_started is None:
            return False
        if current_finished > next_started:
            return False
    return True


def wait_for_http(url, timeout_seconds=20):
    deadline = time.time() + timeout_seconds
    while time.time() < deadline:
        try:
            with urllib.request.urlopen(url, timeout=2) as response:
                if response.status < 500:
                    return
        except Exception:
            time.sleep(0.2)
    raise RuntimeError(f"http endpoint did not become ready: {url}")


def start_process(cmd, env, log_path):
    log_path.parent.mkdir(parents=True, exist_ok=True)
    log_file = log_path.open("w", encoding="utf-8")
    proc = subprocess.Popen(
        cmd,
        cwd=ROOT,
        env=env,
        stdout=log_file,
        stderr=subprocess.STDOUT,
        preexec_fn=os.setsid,
        text=True,
    )
    return proc, log_file


def stop_process(proc, log_file):
    try:
        os.killpg(os.getpgid(proc.pid), signal.SIGINT)
    except OSError:
        pass
    try:
        proc.wait(timeout=5)
    except Exception:
        try:
            os.killpg(os.getpgid(proc.pid), signal.SIGKILL)
        except OSError:
            pass
    log_file.close()


def make_multipart_request(url, field_name, file_path):
    boundary = "----codex-" + uuid.uuid4().hex
    file_name = file_path.name
    body = (
        f"--{boundary}\r\n"
        f'Content-Disposition: form-data; name="{field_name}"; filename="{file_name}"\r\n'
        "Content-Type: application/xml\r\n\r\n"
    ).encode("utf-8") + file_path.read_bytes() + f"\r\n--{boundary}--\r\n".encode("utf-8")
    req = urllib.request.Request(url, data=body, method="POST")
    req.add_header("Content-Type", f"multipart/form-data; boundary={boundary}")
    return req


class NoRedirect(urllib.request.HTTPRedirectHandler):
    def redirect_request(self, req, fp, code, msg, headers, newurl):
        return None


def submit_workload(controller_url, fixture_path):
    opener = urllib.request.build_opener(NoRedirect)
    req = make_multipart_request(controller_url.rstrip("/") + "/workloads", "workload", fixture_path)
    try:
        opener.open(req, timeout=5)
    except urllib.error.HTTPError as err:
        if err.code not in (302, 303):
            raise
        location = err.headers.get("Location", "")
        if not location.startswith("/jobs/"):
            raise RuntimeError(f"unexpected workload redirect: {location}")
        return location.split("/")[-1]
    raise RuntimeError("workload submission did not redirect to job detail")


def post_empty(url):
    req = urllib.request.Request(url, data=b"", method="POST")
    with urllib.request.urlopen(req, timeout=5) as response:
        return response.status


def fetch_json(url):
    with urllib.request.urlopen(url, timeout=5) as response:
        return json.loads(response.read().decode("utf-8"))


def wait_for_job(controller_url, job_id, timeout_seconds=30):
    deadline = time.time() + timeout_seconds
    url = controller_url.rstrip("/") + f"/api/controller/jobs/{job_id}"
    while time.time() < deadline:
        payload = fetch_json(url)
        status = payload["job"]["status"]
        if status in {"succeeded", "failed", "cancelled"}:
            return payload
        time.sleep(0.25)
    raise RuntimeError(f"job {job_id} did not finish before timeout")


def load_json_rows(directory):
    rows = []
    path = Path(directory)
    if not path.exists():
        return rows
    for item in sorted(path.glob("*.json")):
        rows.append(json.loads(item.read_text(encoding="utf-8")))
    return rows


def run_real():
    backend = os.environ.get("SMOKE_REMOTE_LOCAL_BACKEND", "s3").strip().lower()
    scenario = normalize_scenario(os.environ.get("SMOKE_REMOTE_LOCAL_SCENARIO", "single"))
    fixture = fixture_for_selection(backend, scenario)
    ARTIFACT_DIR.mkdir(parents=True, exist_ok=True)
    controller_port = find_free_port()
    driver1_port = find_free_port()
    driver2_port = find_free_port()
    minio_port = find_free_port()
    controller_url = f"http://{DEFAULT_HOST}:{controller_port}"
    driver1_url = f"http://{DEFAULT_HOST}:{driver1_port}"
    driver2_url = f"http://{DEFAULT_HOST}:{driver2_port}"
    minio_url = f"http://{DEFAULT_HOST}:{minio_port}"
    shared_token = "remote-smoke-token"

    controller_data = Path(tempfile.mkdtemp(prefix="cosbench-controller-", dir=ARTIFACT_DIR))
    driver1_data = Path(tempfile.mkdtemp(prefix="cosbench-driver1-", dir=ARTIFACT_DIR))
    driver2_data = Path(tempfile.mkdtemp(prefix="cosbench-driver2-", dir=ARTIFACT_DIR))
    DEFAULT_MINIO_DATA.mkdir(parents=True, exist_ok=True)

    minio_bin = ensure_minio()

    processes = []
    try:
        minio_proc, minio_log = start_process(
            [minio_bin, "server", str(DEFAULT_MINIO_DATA), "--address", f"{DEFAULT_HOST}:{minio_port}"],
            {
                **os.environ,
                "MINIO_ROOT_USER": DEFAULT_MINIO_ACCESS_KEY,
                "MINIO_ROOT_PASSWORD": DEFAULT_MINIO_SECRET_KEY,
            },
            ARTIFACT_DIR / "minio.log",
        )
        processes.append((minio_proc, minio_log))
        wait_for_socket(DEFAULT_HOST, minio_port)

        controller_proc, controller_log = start_process(
            [
                DEFAULT_GO, "run", "./cmd/server",
                "-listen", f"{DEFAULT_HOST}:{controller_port}",
                "-mode", "controller-only",
                "-data-dir", str(controller_data),
                "-view-dir", "web/templates",
                "-driver-shared-token", shared_token,
            ],
            os.environ.copy(),
            ARTIFACT_DIR / "controller.log",
        )
        processes.append((controller_proc, controller_log))
        wait_for_http(controller_url + "/")

        common_driver_env = {
            **os.environ,
            "COSBENCH_SMOKE_ENDPOINT": minio_url,
            "COSBENCH_SMOKE_ACCESS_KEY": DEFAULT_MINIO_ACCESS_KEY,
            "COSBENCH_SMOKE_SECRET_KEY": DEFAULT_MINIO_SECRET_KEY,
        }
        driver1_proc, driver1_log = start_process(
            [
                DEFAULT_GO, "run", "./cmd/server",
                "-listen", f"{DEFAULT_HOST}:{driver1_port}",
                "-mode", "driver-only",
                "-data-dir", str(driver1_data),
                "-view-dir", "web/templates",
                "-driver-shared-token", shared_token,
                "-controller-url", controller_url,
                "-driver-name", "driver-1",
            ],
            common_driver_env,
            ARTIFACT_DIR / "driver1.log",
        )
        processes.append((driver1_proc, driver1_log))
        wait_for_http(driver1_url + "/")

        driver2_proc, driver2_log = start_process(
            [
                DEFAULT_GO, "run", "./cmd/server",
                "-listen", f"{DEFAULT_HOST}:{driver2_port}",
                "-mode", "driver-only",
                "-data-dir", str(driver2_data),
                "-view-dir", "web/templates",
                "-driver-shared-token", shared_token,
                "-controller-url", controller_url,
                "-driver-name", "driver-2",
            ],
            common_driver_env,
            ARTIFACT_DIR / "driver2.log",
        )
        processes.append((driver2_proc, driver2_log))
        wait_for_http(driver2_url + "/")

        job_id = submit_workload(controller_url, fixture)
        status = post_empty(controller_url + f"/jobs/{job_id}/start")
        if status not in (200, 303):
            raise RuntimeError(f"unexpected start status: {status}")
        payload = wait_for_job(controller_url, job_id)

        controller_drivers = load_json_rows(controller_data / "drivers")
        controller_missions = load_json_rows(controller_data / "missions")
        job_stages = payload["job"].get("stages", [])
        stage_names = [item.get("name") for item in job_stages if item.get("name")]
        checks = {
            "process_ready": "pass",
            "drivers_healthy": "pass" if len(controller_drivers) == 2 and all(item.get("status") == "healthy" for item in controller_drivers) else "fail",
            "units_distributed": "pass" if len({item.get("lease", {}).get("driver_id") for item in controller_missions if item.get("lease") and item.get("status") in {"claimed", "running", "succeeded", "failed"}}) >= 2 else "fail",
            "job_succeeded": "pass" if payload["job"]["status"] == "succeeded" else "fail",
            "visibility": "pass",
        }
        if scenario == "multistage":
            mission_stage_names = {item.get("stage_name") for item in controller_missions if item.get("stage_name")}
            stage_totals = payload["result"].get("stage_totals", [])
            checks.update({
                "stages_present": "pass" if len(job_stages) >= 2 and all(item.get("status") == "succeeded" for item in job_stages) else "fail",
                "stage_coverage": "pass" if len(stage_names) >= 2 and set(stage_names).issubset(mission_stage_names) else "fail",
                "stage_barrier": "pass" if multistage_barrier_holds(job_stages) else "fail",
                "stage_aggregation": "pass" if len(stage_totals) >= 2 else "fail",
            })

        driver_ids = [item["id"] for item in controller_drivers]
        # prove driver APIs are live
        if driver_ids:
            fetch_json(driver1_url + f"/api/driver/self?driver_id={urllib.parse.quote(driver_ids[0])}")
            fetch_json(driver2_url + f"/api/driver/self?driver_id={urllib.parse.quote(driver_ids[-1])}")
        fetch_json(controller_url + f"/api/controller/jobs/{job_id}/timeline")
        fetch_json(controller_url + "/api/controller/jobs")
        if scenario == "multistage":
            for stage_name in stage_names:
                fetch_json(controller_url + f"/api/controller/jobs/{job_id}/stages/{urllib.parse.quote(stage_name)}")

        drivers_participated = len({item.get("lease", {}).get("driver_id") for item in controller_missions if item.get("lease") and item.get("status") in {"claimed", "running", "succeeded", "failed"}})
        units_claimed = len({item.get("work_unit_id") for item in controller_missions if item.get("lease") and item.get("status") in {"claimed", "running", "succeeded", "failed"}})
        summary = build_summary(
            backend=backend,
            scenario=scenario,
            controller_url=controller_url,
            driver_urls=[driver1_url, driver2_url],
            job_id=job_id,
            job_status=payload["job"]["status"],
            drivers_seen=len(controller_drivers),
            units_claimed=units_claimed,
            drivers_participated=drivers_participated,
            operation_count=payload["result"]["metrics"]["operation_count"],
            byte_count=payload["result"]["metrics"]["byte_count"],
            stage_names=stage_names,
            stages_seen=len(stage_names),
            checks=checks,
        )
        write_summary(summary)
        sys.stdout.write(render_summary_md(summary))
        raise SystemExit(0 if summary["overall"] == "pass" else 1)
    except Exception as err:
        summary = build_failure_summary("orchestration", str(err))
        write_summary(summary)
        sys.stdout.write(render_summary_md(summary))
        raise SystemExit(1)
    finally:
        for proc, log_file in reversed(processes):
            stop_process(proc, log_file)


def main():
    mock = run_mock()
    if mock is not None:
        raise SystemExit(mock)
    run_real()


if __name__ == "__main__":
    main()

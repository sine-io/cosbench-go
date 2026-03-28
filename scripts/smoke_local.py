#!/usr/bin/env python3

import os
import signal
import socket
import subprocess
import sys
import time
import urllib.request
from pathlib import Path


ROOT = Path(__file__).resolve().parents[1]
DEFAULT_MINIO = ROOT / ".artifacts" / "minio" / "bin" / "minio"
DEFAULT_MINIO_DATA = ROOT / ".artifacts" / "minio" / "data"
DEFAULT_MINIO_URL = "https://dl.min.io/server/minio/release/linux-amd64/minio"
DEFAULT_HOST = "127.0.0.1"
DEFAULT_GO = os.environ.get("GO", "/snap/bin/go")
DEFAULT_MINIO_ACCESS_KEY = "minioadmin"
DEFAULT_MINIO_SECRET_KEY = "minioadmin"


def print_summary(endpoint, s3_status, sio_status):
    overall = "pass" if s3_status == "pass" and sio_status == "pass" else "fail"
    print("# Smoke Local")
    print()
    print("Provider: `minio`")
    print(f"Endpoint: `{endpoint}`")
    print(f"S3 smoke: `{s3_status}`")
    print(f"SIO smoke: `{sio_status}`")
    print(f"Overall: `{overall}`")
    return overall


def run_mock():
    mode = os.environ.get("SMOKE_LOCAL_MOCK", "").strip()
    if mode == "success":
        print_summary("http://127.0.0.1:9000", "pass", "pass")
        return 0
    if mode == "failure":
        print_summary("http://127.0.0.1:9000", "pass", "fail")
        return 1
    return None


def find_free_port():
    with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as sock:
        sock.bind((DEFAULT_HOST, 0))
        return sock.getsockname()[1]


def ensure_minio():
    if DEFAULT_MINIO.exists():
        return str(DEFAULT_MINIO)

    DEFAULT_MINIO.parent.mkdir(parents=True, exist_ok=True)
    with urllib.request.urlopen(DEFAULT_MINIO_URL) as response, DEFAULT_MINIO.open("wb") as target:
        target.write(response.read())
    DEFAULT_MINIO.chmod(0o755)
    return str(DEFAULT_MINIO)


def wait_for_socket(host, port, timeout_seconds=10):
    deadline = time.time() + timeout_seconds
    while time.time() < deadline:
        with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as sock:
            sock.settimeout(0.5)
            try:
                sock.connect((host, port))
                return
            except OSError:
                time.sleep(0.1)
    raise RuntimeError(f"minio server did not become ready on {host}:{port}")


def run_smoke(go_bin, endpoint, backend, extra_env=None):
    env = os.environ.copy()
    env["GO"] = go_bin
    env["COSBENCH_SMOKE_ENDPOINT"] = endpoint
    env["COSBENCH_SMOKE_ACCESS_KEY"] = DEFAULT_MINIO_ACCESS_KEY
    env["COSBENCH_SMOKE_SECRET_KEY"] = DEFAULT_MINIO_SECRET_KEY
    env["COSBENCH_SMOKE_BACKEND"] = backend
    env.setdefault("COSBENCH_SMOKE_BUCKET_PREFIX", "cosbench-go-local-smoke")
    if extra_env:
        env.update(extra_env)
    proc = subprocess.run(
        ["make", "smoke-s3"],
        cwd=ROOT,
        env=env,
        text=True,
        capture_output=True,
        check=False,
    )
    return proc


def main():
    mock_result = run_mock()
    if mock_result is not None:
        raise SystemExit(mock_result)

    minio_server = ensure_minio()
    port = find_free_port()
    endpoint = f"http://{DEFAULT_HOST}:{port}"
    DEFAULT_MINIO_DATA.mkdir(parents=True, exist_ok=True)
    server = subprocess.Popen(
        [minio_server, "server", str(DEFAULT_MINIO_DATA), "--address", f"{DEFAULT_HOST}:{port}"],
        cwd=ROOT,
        env={
            **os.environ,
            "MINIO_ROOT_USER": DEFAULT_MINIO_ACCESS_KEY,
            "MINIO_ROOT_PASSWORD": DEFAULT_MINIO_SECRET_KEY,
        },
        stdout=subprocess.DEVNULL,
        stderr=subprocess.DEVNULL,
        preexec_fn=os.setsid,
    )
    try:
        wait_for_socket(DEFAULT_HOST, port)
        s3_proc = run_smoke(DEFAULT_GO, endpoint, "s3")
        sio_proc = run_smoke(DEFAULT_GO, endpoint, "sio", {"COSBENCH_SMOKE_PATH_STYLE": "true"})
        s3_status = "pass" if s3_proc.returncode == 0 else "fail"
        sio_status = "pass" if sio_proc.returncode == 0 else "fail"
        overall = print_summary(endpoint, s3_status, sio_status)
        if s3_proc.returncode != 0:
            sys.stdout.write("\n## S3 Output\n\n")
            sys.stdout.write(s3_proc.stdout)
            sys.stdout.write(s3_proc.stderr)
        if sio_proc.returncode != 0:
            sys.stdout.write("\n## SIO Output\n\n")
            sys.stdout.write(sio_proc.stdout)
            sys.stdout.write(sio_proc.stderr)
        raise SystemExit(0 if overall == "pass" else 1)
    finally:
        try:
            os.killpg(os.getpgid(server.pid), signal.SIGINT)
        except OSError:
            pass
        server.wait(timeout=5)


if __name__ == "__main__":
    main()

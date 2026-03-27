#!/usr/bin/env python3

import os
import shutil
import signal
import socket
import subprocess
import sys
import time
from pathlib import Path


ROOT = Path(__file__).resolve().parents[1]
DEFAULT_VENV = ROOT / ".artifacts" / "moto-venv"
DEFAULT_HOST = "127.0.0.1"
DEFAULT_GO = os.environ.get("GO", "/snap/bin/go")


def print_summary(endpoint, s3_status, sio_status):
    overall = "pass" if s3_status == "pass" and sio_status == "pass" else "fail"
    print("# Smoke Local")
    print()
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


def ensure_moto_server():
    moto_server = DEFAULT_VENV / "bin" / "moto_server"
    if moto_server.exists():
        return str(moto_server)

    python_bin = shutil.which("python3")
    if not python_bin:
        raise RuntimeError("python3 not found for smoke-local setup")

    subprocess.run([python_bin, "-m", "venv", str(DEFAULT_VENV)], check=True)
    pip_bin = DEFAULT_VENV / "bin" / "python"
    subprocess.run([str(pip_bin), "-m", "pip", "install", "--upgrade", "pip"], check=True)
    subprocess.run([str(DEFAULT_VENV / "bin" / "pip"), "install", "moto[server]"], check=True)
    return str(moto_server)


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
    raise RuntimeError(f"moto server did not become ready on {host}:{port}")


def run_smoke(go_bin, endpoint, backend, extra_env=None):
    env = os.environ.copy()
    env["GO"] = go_bin
    env["COSBENCH_SMOKE_ENDPOINT"] = endpoint
    env["COSBENCH_SMOKE_ACCESS_KEY"] = "test"
    env["COSBENCH_SMOKE_SECRET_KEY"] = "test"
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

    moto_server = ensure_moto_server()
    port = find_free_port()
    endpoint = f"http://{DEFAULT_HOST}:{port}"
    server = subprocess.Popen(
        [moto_server, "-H", DEFAULT_HOST, "-p", str(port)],
        cwd=ROOT,
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

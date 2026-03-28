import importlib.util
from pathlib import Path


def load_module():
    path = Path("scripts/summarize_smoke_s3_output.py")
    spec = importlib.util.spec_from_file_location("summarize_smoke_s3_output", path)
    module = importlib.util.module_from_spec(spec)
    assert spec.loader is not None
    spec.loader.exec_module(module)
    return module


def test_classify_smoke_output_skipped():
    module = load_module()
    text = """go test ./internal/driver/s3 -run Smoke -v
=== RUN   TestSmokeObjectLifecycle
--- SKIP: TestSmokeObjectLifecycle (0.00s)
=== RUN   TestSmokeSIOMultipartLifecycle
--- SKIP: TestSmokeSIOMultipartLifecycle (0.00s)
PASS
"""
    assert module.classify_output(text) == "skipped"


def test_classify_smoke_output_executed():
    module = load_module()
    text = """go test ./internal/driver/s3 -run Smoke -v
=== RUN   TestSmokeObjectLifecycle
--- PASS: TestSmokeObjectLifecycle (0.10s)
=== RUN   TestSmokeSIOMultipartLifecycle
--- PASS: TestSmokeSIOMultipartLifecycle (0.11s)
PASS
"""
    assert module.classify_output(text) == "executed"


def test_classify_smoke_output_failed():
    module = load_module()
    text = """go test ./internal/driver/s3 -run Smoke -v
=== RUN   TestSmokeObjectLifecycle
--- FAIL: TestSmokeObjectLifecycle (0.10s)
FAIL
"""
    assert module.classify_output(text) == "failed"

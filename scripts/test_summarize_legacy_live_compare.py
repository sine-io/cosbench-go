import importlib.util
from pathlib import Path


def load_module():
    path = Path("scripts/summarize_legacy_live_compare.py")
    spec = importlib.util.spec_from_file_location("summarize_legacy_live_compare", path)
    module = importlib.util.module_from_spec(spec)
    assert spec.loader is not None
    spec.loader.exec_module(module)
    return module


def test_classify_legacy_summary_skipped():
    module = load_module()
    payload = {"status": "skipped", "reason": "missing secrets"}
    assert module.classify_summary(payload) == "skipped"


def test_classify_legacy_summary_executed():
    module = load_module()
    payload = {"workload": "xml-sample", "stages": 2, "works": 2, "samples": 4, "errors": 0}
    assert module.classify_summary(payload) == "executed"


def test_classify_legacy_summary_failed():
    module = load_module()
    payload = {}
    assert module.classify_summary(payload) == "failed"

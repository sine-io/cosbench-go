import subprocess
from pathlib import Path


def run_render(tmp_path, fixture, backend):
    output = tmp_path / "rendered.xml"
    subprocess.run(
        [
            "python3",
            "scripts/render_legacy_live_compare_workload.py",
            fixture,
            str(output),
            backend,
            "http://example.test",
            "ak",
            "sk",
            "us-east-1",
            "true",
        ],
        check=True,
        cwd=Path.cwd(),
    )
    return output.read_text(encoding="utf-8")


def test_render_s3_legacy_fixture_replaces_placeholders(tmp_path):
    text = run_render(tmp_path, "testdata/legacy/s3-config-sample.xml", "s3")
    assert "<accesskey>" not in text
    assert "<scretkey>" not in text
    assert "<endpoint>" not in text
    assert 'type="s3"' in text
    assert "accesskey=ak" in text
    assert "secretkey=sk" in text
    assert "endpoint=http://example.test" in text


def test_render_sio_legacy_fixture_replaces_placeholders(tmp_path):
    text = run_render(tmp_path, "testdata/legacy/sio-config-sample.xml", "sio")
    assert "<accesskey>" not in text
    assert "<scretkey>" not in text
    assert "<endpoint>" not in text
    assert 'type="sio"' in text
    assert "accesskey=ak" in text
    assert "secretkey=sk" in text
    assert "endpoint=http://example.test" in text

from pathlib import Path
import re


def test_all_workflow_upload_artifact_actions_use_v7():
    workflows = sorted(Path(".github/workflows").glob("*.yml"))
    assert workflows, "expected workflow files"
    mismatches = []
    pattern = re.compile(r"actions/upload-artifact@([^\s]+)")
    for workflow in workflows:
        text = workflow.read_text(encoding="utf-8")
        for version in pattern.findall(text):
            if version != "v7.0.0":
                mismatches.append((workflow.name, version))
    assert not mismatches, f"unexpected upload-artifact versions: {mismatches}"

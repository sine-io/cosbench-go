from pathlib import Path


def read_manifest(manifest_path: str):
    fixtures = []
    for raw_line in Path(manifest_path).read_text().splitlines():
        line = raw_line.strip()
        if not line or line.startswith("#"):
            continue
        name, workload = line.split()
        fixtures.append({"name": name, "workload": workload})
    return fixtures


def parse_filter(raw_filter: str):
    items = [item for item in raw_filter.split(",") if item]
    if items == ["all"]:
        return []
    return items


def validate_filter(fixtures, raw_filter: str):
    selected = parse_filter(raw_filter)
    if not selected and raw_filter in ("", "all"):
        return
    known = {fixture["name"] for fixture in fixtures}
    for name in selected:
        if name not in known:
            raise ValueError(name)

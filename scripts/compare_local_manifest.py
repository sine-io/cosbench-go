from pathlib import Path


class ManifestFormatError(ValueError):
    pass


def read_manifest(manifest_path: str):
    fixtures = []
    for line_no, raw_line in enumerate(Path(manifest_path).read_text().splitlines(), start=1):
        line = raw_line.strip()
        if not line or line.startswith("#"):
            continue
        fields = line.split()
        if len(fields) != 2:
            raise ManifestFormatError(
                f"invalid compare-local manifest line {line_no} in {manifest_path}: {line!r}"
            )
        name, workload = fields
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

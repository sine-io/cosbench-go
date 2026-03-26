from pathlib import Path


class ManifestError(ValueError):
    pass


class ManifestFormatError(ManifestError):
    pass


class ManifestReadError(ManifestError):
    pass


def read_manifest(manifest_path: str):
    fixtures = []
    seen_names = {}
    try:
        lines = Path(manifest_path).read_text().splitlines()
    except FileNotFoundError:
        raise ManifestReadError(f"compare-local manifest not found: {manifest_path}")

    for line_no, raw_line in enumerate(lines, start=1):
        line = raw_line.strip()
        if not line or line.startswith("#"):
            continue
        fields = line.split()
        if len(fields) != 2:
            raise ManifestFormatError(
                f"invalid compare-local manifest line {line_no} in {manifest_path}: {line!r}"
            )
        name, workload = fields
        if name in seen_names:
            raise ManifestFormatError(
                f"duplicate compare-local fixture name {name!r} on line {line_no} in {manifest_path}; first seen on line {seen_names[name]}"
            )
        seen_names[name] = line_no
        fixtures.append({"name": name, "workload": workload})
    return fixtures


def parse_filter(raw_filter: str):
    items = [item.strip() for item in raw_filter.split(",") if item.strip()]
    if items == ["all"]:
        return []
    return items


def normalize_filter(raw_filter: str):
    selected = parse_filter(raw_filter)
    if not selected:
        return "all"
    return ",".join(selected)


def select_fixtures(fixtures, raw_filter: str):
    selected = set(parse_filter(raw_filter))
    if not selected:
        return fixtures
    return [fixture for fixture in fixtures if fixture["name"] in selected]


def validate_filter(fixtures, raw_filter: str):
    selected = parse_filter(raw_filter)
    if not selected and raw_filter in ("", "all"):
        return
    known = {fixture["name"] for fixture in fixtures}
    for name in selected:
        if name not in known:
            raise ValueError(name)

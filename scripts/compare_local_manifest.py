from pathlib import Path


class ManifestError(ValueError):
    pass


class ManifestFormatError(ManifestError):
    pass


class ManifestReadError(ManifestError):
    pass


class FilterError(ValueError):
    pass


class InvalidFilterError(FilterError):
    pass


class UnknownFixtureError(FilterError):
    pass


def read_manifest(manifest_path: str):
    fixtures = []
    seen_names = {}
    try:
        lines = Path(manifest_path).read_text().splitlines()
    except FileNotFoundError:
        raise ManifestReadError(f"compare-local manifest not found: {manifest_path}")
    except UnicodeDecodeError as err:
        raise ManifestReadError(f"unable to decode compare-local manifest {manifest_path}: {err}")
    except OSError as err:
        raise ManifestReadError(f"unable to read compare-local manifest {manifest_path}: {err}")

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
    items = []
    seen = set()
    for raw_item in raw_filter.split(","):
        item = raw_item.strip()
        if not item or item in seen:
            continue
        seen.add(item)
        items.append(item)
    if items == ["all"]:
        return []
    return items


def normalize_filter(raw_filter: str):
    selected = parse_filter(raw_filter)
    if not selected:
        return "all"
    return ",".join(selected)


def select_fixtures(fixtures, raw_filter: str):
    selected = parse_filter(raw_filter)
    if not selected:
        return fixtures
    fixture_map = {fixture["name"]: fixture for fixture in fixtures}
    return [fixture_map[name] for name in selected if name in fixture_map]


def format_filter_error(fixtures, err: FilterError):
    if isinstance(err, InvalidFilterError):
        return f"invalid compare-local filter: {err}"
    names = "".join(f"  - {fixture['name']}\n" for fixture in fixtures)
    return f"unknown compare-local fixture: {err}\nknown fixtures:\n{names}"


def validate_filter(fixtures, raw_filter: str):
    selected = parse_filter(raw_filter)
    if not selected and raw_filter in ("", "all"):
        return
    if "all" in selected and len(selected) > 1:
        raise InvalidFilterError("'all' cannot be combined with specific fixtures")
    known = {fixture["name"] for fixture in fixtures}
    for name in selected:
        if name not in known:
            raise UnknownFixtureError(name)

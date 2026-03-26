from pathlib import Path, PurePosixPath, PureWindowsPath
import sys


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


def display_text(value: str):
    return value.encode("utf-8", "surrogateescape").decode("utf-8", "replace")


def format_os_error(err: OSError):
    parts = []
    if getattr(err, "errno", None) is not None:
        parts.append(f"[Errno {err.errno}]")
    if getattr(err, "strerror", None):
        parts.append(display_text(str(err.strerror)))
    elif str(err):
        parts.append(display_text(str(err)))
    return " ".join(parts) or err.__class__.__name__


def configure_utf8_stdio():
    for stream in (sys.stdout, sys.stderr):
        if hasattr(stream, "reconfigure"):
            stream.reconfigure(encoding="utf-8")


def validate_fixture_name(name: str):
    if name in (".", "..") or "/" in name or "\\" in name:
        raise ManifestFormatError(
            f"invalid compare-local fixture name {name!r}: must not contain path separators or dot-path segments"
        )
    if name == "all":
        raise ManifestFormatError(
            f"invalid compare-local fixture name {name!r}: reserved for the all-fixtures selector"
        )
    if name.startswith("--"):
        raise ManifestFormatError(
            f"invalid compare-local fixture name {name!r}: must not start with --"
        )


def validate_workload_path(workload: str):
    posix_path = PurePosixPath(workload)
    windows_path = PureWindowsPath(workload)
    if "\\" in workload:
        raise ManifestFormatError(
            f"invalid compare-local workload path {workload!r}: must use forward slashes instead of backslashes"
        )
    if posix_path.is_absolute() or windows_path.is_absolute():
        raise ManifestFormatError(
            f"invalid compare-local workload path {workload!r}: must not be absolute"
        )
    if any(part == ".." for part in posix_path.parts) or any(part == ".." for part in windows_path.parts):
        raise ManifestFormatError(
            f"invalid compare-local workload path {workload!r}: must be repo-relative without '..' segments"
        )
def read_manifest(manifest_path: str):
    fixtures = []
    seen_names = {}
    manifest_display = display_text(manifest_path)
    try:
        lines = Path(manifest_path).read_text(encoding="utf-8-sig").splitlines()
    except FileNotFoundError:
        raise ManifestReadError(f"compare-local manifest not found: {manifest_display}")
    except UnicodeDecodeError as err:
        raise ManifestReadError(f"unable to decode compare-local manifest {manifest_display}: {err}")
    except OSError as err:
        raise ManifestReadError(f"unable to read compare-local manifest {manifest_display}: {format_os_error(err)}")

    for line_no, raw_line in enumerate(lines, start=1):
        line = raw_line.strip()
        if not line or line.startswith("#"):
            continue
        fields = line.split()
        if len(fields) != 2:
            raise ManifestFormatError(
                f"invalid compare-local manifest line {line_no} in {manifest_display}: {line!r}"
            )
        name, workload = fields
        validate_fixture_name(name)
        validate_workload_path(workload)
        if name in seen_names:
            raise ManifestFormatError(
                f"duplicate compare-local fixture name {name!r} on line {line_no} in {manifest_display}; first seen on line {seen_names[name]}"
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
    if not selected:
        stripped = raw_filter.strip()
        non_empty_tokens = [item.strip() for item in raw_filter.split(",") if item.strip()]
        if stripped == "" or non_empty_tokens == ["all"]:
            return
        raise InvalidFilterError("filter did not include any fixture names")
    if "all" in selected and len(selected) > 1:
        raise InvalidFilterError("'all' cannot be combined with specific fixtures")
    known = {fixture["name"] for fixture in fixtures}
    for name in selected:
        if name not in known:
            raise UnknownFixtureError(name)
